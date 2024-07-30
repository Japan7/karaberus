package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type KaraberusClaims struct {
	oidc.IDTokenClaims
	Scopes
	jwt.MapClaims `json:"-"`
}

var (
	sessionCookieName = "karaberus_session"
	currentUserCtxKey = "current_user"
)

func addOidcRoutes(app *fiber.App) {
	provider, err := rp.NewRelyingPartyOIDC(
		context.TODO(),
		CONFIG.OIDC.Issuer,
		CONFIG.OIDC.ClientID,
		CONFIG.OIDC.ClientSecret,
		fmt.Sprintf("%v/api/oidc/callback", CONFIG.Listen.BaseURL),
		strings.Split(CONFIG.OIDC.Scopes, " "))
	if err != nil {
		getLogger().Print(err)
	}
	state := func() string {
		return uuid.New().String()
	}
	app.Get("/api/oidc/login",
		adaptor.HTTPHandler(rp.AuthURLHandler(state, provider)))
	app.Get("/api/oidc/callback",
		adaptor.HTTPHandler(rp.CodeExchangeHandler(rp.UserinfoCallback(callbackHandler), provider)))
}

func callbackHandler(
	w http.ResponseWriter,
	r *http.Request,
	tokens *oidc.Tokens[*oidc.IDTokenClaims],
	state string,
	rp rp.RelyingParty,
	info *oidc.UserInfo) {
	sub := info.Subject
	if CONFIG.OIDC.IDClaim != "" {
		sub = info.Claims[CONFIG.OIDC.IDClaim].(string)
	}

	expiresAt := time.Now().Add(time.Hour)
	_, signed, err := CreateJwtForUser(r.Context(), sub, &expiresAt, info)
	if err != nil {
		getLogger().Print(err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    sessionCookieName,
		Value:   signed,
		Path:    "/",
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func CreateJwtForUser(
	ctx context.Context,
	sub string,
	expiresAt *time.Time,
	info *oidc.UserInfo) (*jwt.Token, string, error) {
	user, err := getOrCreateUser(ctx, sub)
	if err != nil {
		return nil, "", err
	}
	claims := KaraberusClaims{}
	if info != nil {
		claims.SetUserInfo(info)
	}
	claims.Subject = user.ID
	claims.IssuedAt = oidc.Time(time.Now().Unix())
	if expiresAt != nil {
		claims.IDTokenClaims.TokenClaims.Expiration = oidc.Time(expiresAt.Unix())
	}
	claims.Scopes = AllScopes
	return createJwt(claims)
}

func createJwt(claims KaraberusClaims) (*jwt.Token, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(CONFIG.OIDC.JwtSignKey))
	return token, signed, err
}

func authMiddleware(ctx huma.Context, next func(huma.Context)) {
	var (
		token  string
		user   *User
		scopes *Scopes
		err    error
	)

	// Check for a token in the request.
	token, err = getRequestToken(ctx)
	// If we have a token, try to get the user.
	if err == nil {
		user, scopes, err = getUserScopesFromApiToken(ctx.Context(), token)
		if err != nil {
			user, scopes, err = getUserScopesFromJwt(ctx.Context(), token)
		}
	}
	// If we have a user, add it to the context.
	if err == nil {
		ctx = huma.WithValue(ctx, currentUserCtxKey, *user)
	}

	if ok := checkOperationSecurity(ctx, user, scopes); ok {
		next(ctx)
		return
	}

	if err != nil {
		getLogger().Print(err)
		ctx.SetStatus(http.StatusUnauthorized)
	} else {
		ctx.SetStatus(http.StatusForbidden)
	}
}

func getRequestToken(ctx huma.Context) (string, error) {
	authHeader := ctx.Header("authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, oidc.BearerToken) {
			return strings.TrimPrefix(authHeader, oidc.PrefixBearer), nil
		}
		return "", errors.New("invalid authorization header")
	} else {
		cookie, err := huma.ReadCookie(ctx, sessionCookieName)
		if err != nil {
			return "", err
		}
		return cookie.Value, nil
	}
}

func getUserScopesFromApiToken(ctx context.Context, token string) (*User, *Scopes, error) {
	db := GetDB(ctx)
	apiToken := Token{ID: token}
	if err := db.First(&apiToken).Error; err != nil {
		return nil, nil, err
	}
	return &apiToken.User, &apiToken.Scopes, nil
}

func getUserScopesFromJwt(ctx context.Context, token string) (*User, *Scopes, error) {
	// claims := KaraberusClaims{} // FIXME: deserialization is not working with Scopes
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(CONFIG.OIDC.JwtSignKey), nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, nil, &KaraberusError{"invalid token"}
	}

	sub, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return nil, nil, err
	}

	db := GetDB(ctx)
	user := User{ID: sub}
	if err := db.First(&user).Error; err != nil {
		return nil, nil, err
	}

	// FIXME: replace with proper deserialization
	scopes := Scopes{}
	for key, val := range claims {
		if key == "kara" {
			scopes.Kara = val.(bool)
		}
		if key == "user" {
			scopes.User = val.(bool)
		}
	}

	return &user, &scopes, nil
}

func checkOperationSecurity(ctx huma.Context, user *User, scopes *Scopes) bool {
	var authRequired bool
	var opScopes []string = []string{}
	for _, opScheme := range ctx.Operation().Security {
		var ok bool
		if _, ok = opScheme["oidc"]; ok {
			authRequired = true
		}
		if opScopes, ok = opScheme["scopes"]; ok {
			break
		}
	}

	if authRequired && user == nil {
		return false
	}

	for _, v := range opScopes {
		if !scopes.HasScope(v) {
			return false
		}
	}

	return true
}
