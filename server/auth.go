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

type Time int64

type KaraberusClaims struct {
	Subject       string `json:"sub,omitempty"`
	IssuedAt      Time   `json:"iat,omitempty"`
	Expiration    Time   `json:"exp,omitempty"`
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
	_, signed, err := CreateTokenForUser(r.Context(), sub, &expiresAt)
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

func CreateTokenForUser(ctx context.Context, sub string, expiresAt *time.Time) (*jwt.Token, string, error) {
	user, err := getOrCreateUser(ctx, sub)
	if err != nil {
		return nil, "", err
	}
	claims := KaraberusClaims{
		Subject:  user.ID,
		IssuedAt: Time(time.Now().Unix()),
	}
	if expiresAt != nil {
		claims.Expiration = Time(expiresAt.Unix())
	}
	return createToken(claims)
}

func createToken(claims KaraberusClaims) (*jwt.Token, string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := jwtToken.SignedString([]byte(CONFIG.OIDC.JwtSignKey))
	return jwtToken, signed, err
}

func authMiddleware(ctx huma.Context, next func(huma.Context)) {
	token, err := getRequestToken(ctx)
	if err != nil {
		getLogger().Print(err)
		ctx.SetStatus(http.StatusUnauthorized)
		return
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(CONFIG.OIDC.JwtSignKey), nil
	})
	if err != nil || !jwtToken.Valid {
		getLogger().Printf("Invalid token, %v", err)
		ctx.SetStatus(http.StatusUnauthorized)
		return
	}

	sub, err := jwtToken.Claims.GetSubject()
	if err != nil {
		getLogger().Print(err)
		ctx.SetStatus(http.StatusUnauthorized)
		return
	}

	db := GetDB(ctx.Context())
	user := User{ID: sub}
	if err = db.First(&user, sub).Error; err != nil {
		getLogger().Print(err)
		ctx.SetStatus(http.StatusUnauthorized)
		return
	}

	ctx = huma.WithValue(ctx, currentUserCtxKey, user)
	next(ctx)
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

func getCurrentUser(ctx context.Context) User {
	return ctx.Value(currentUserCtxKey).(User)
}
