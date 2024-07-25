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
		fmt.Sprintf("%v/api/oidc/callback", CONFIG.Listen.BaseUrl),
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
	return createJwt(claims)
}

func createJwt(claims KaraberusClaims) (*jwt.Token, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(CONFIG.OIDC.JwtSignKey))
	return token, signed, err
}

func authMiddleware(ctx huma.Context, next func(huma.Context)) {
	var (
		token string
		user  *User
		err   error
	)

	// Check for a token in the request.
	token, err = getRequestToken(ctx)
	// If we have a token, try to get the user.
	if err == nil {
		user, err = getUserFromApiToken(ctx.Context(), token)
		if err != nil {
			user, err = getUserFromJwt(ctx.Context(), token)
		}
	}
	// If we have a user, add it to the context.
	if err == nil {
		ctx = huma.WithValue(ctx, currentUserCtxKey, user)
	}

	if ok := checkOperationSecurity(ctx, user); ok {
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

func getUserFromApiToken(ctx context.Context, token string) (*User, error) {
	db := GetDB(ctx)
	apiToken := Token{ID: token}
	if err := db.First(&apiToken).Error; err != nil {
		return nil, err
	}
	return &apiToken.User, nil
}

func getUserFromJwt(ctx context.Context, token string) (*User, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(CONFIG.OIDC.JwtSignKey), nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, &KaraberusError{"invalid token"}
	}

	sub, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	db := GetDB(ctx)
	user := User{ID: sub}
	if err := db.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func checkOperationSecurity(ctx huma.Context, user *User) bool {
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

	if user.Admin {
		return true
	}

	for _, v := range opScopes {
		if !user.Scopes.HasScope(v) {
			return false
		}
	}

	return true
}
