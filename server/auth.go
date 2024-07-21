package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"gorm.io/gorm"
)

var sessionCookieName = "karaberus_session"

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

func authMiddleware(ctx huma.Context, next func(huma.Context)) {
	token, err := getRequestToken(ctx)
	if err != nil {
		getLogger().Print(err)
		ctx.SetStatus(http.StatusUnauthorized)
		return
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(CONFIG.OIDC.JwtSignKey), nil
	})
	if err != nil || !jwtToken.Valid {
		getLogger().Printf("Invalid token, %v", err)
		ctx.SetStatus(http.StatusUnauthorized)
		return
	}
	claims := jwtToken.Claims.(jwt.MapClaims)

	db := GetDB(ctx.Context())
	user_id := claims["sub"].(string)
	user := User{ID: user_id}
	err = db.First(&user, user_id).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			// The user doesn't exist yet
			err = db.Create(&user).Error
			if err != nil {
				getLogger().Print(err)
				ctx.SetStatus(http.StatusInternalServerError)
				return
			}
		} else {
			getLogger().Print(err)
			ctx.SetStatus(http.StatusInternalServerError)
			return
		}
	}

	ctx = huma.WithValue(ctx, "current_user", user)
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
	return ctx.Value("current_user").(User)
}
