package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

var sessionCookieName = "karaberus_session"

func oidcRoutes(app *fiber.App) {
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
		adaptor.HTTPHandler(rp.CodeExchangeHandler(rp.UserinfoCallback(setSessionCookie), provider)))
}

func setSessionCookie(
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
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":          sub,
		"access_token": tokens.AccessToken,
	})
	signed, err := jwtToken.SignedString([]byte(CONFIG.OIDC.JwtSignKey))
	if err != nil {
		getLogger().Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    sessionCookieName,
		Value:   signed,
		Path:    "/",
		Expires: tokens.Expiry,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
