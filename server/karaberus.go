// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type KaraberusError struct {
	Message string
}

func (m *KaraberusError) Error() string {
	return m.Message
}

func getCurrentUser(ctx context.Context) User {
	return ctx.Value("current_user").(User)
}

func getBearerToken(auth string) (string, error) {
	if strings.HasPrefix(auth, oidc.BearerToken) {
		return strings.TrimPrefix(auth, oidc.PrefixBearer), nil
	}
	return "", &KaraberusError{"Authorization header is missing."}
}

func setSecurity(security []map[string][]string) func(o *huma.Operation) {
	return func(o *huma.Operation) {
		o.Security = security
	}
}

func routes(api huma.API) {
	oidc_security := []map[string][]string{{"oidc": []string{""}}}
	kara_security := []map[string][]string{{"oidc": []string{""}, "scopes": []string{"kara"}}}

	huma.Get(api, "/api/kara", GetAllKaras, setSecurity(kara_security))
	huma.Get(api, "/api/kara/{id}", GetKara, setSecurity(kara_security))
	huma.Delete(api, "/api/kara/{id}", DeleteKara, setSecurity(kara_security))
	// TODO: should be reserved to admins
	huma.Patch(api, "/api/kara/{id}/creation_time", SetKaraUploadTime, setSecurity(kara_security))
	huma.Post(api, "/api/kara", CreateKara, setSecurity(kara_security))
	huma.Put(api, "/api/kara/{id}/upload/{filetype}", UploadKaraFile, setSecurity(kara_security))
	huma.Get(api, "/api/kara/{id}/download/{filetype}", DownloadFile, setSecurity(kara_security))

	huma.Get(api, "/api/tags/audio", GetAudioTags, setSecurity(kara_security))
	huma.Get(api, "/api/tags/video", GetVideoTags, setSecurity(kara_security))

	huma.Get(api, "/api/tags/author", FindAuthor, setSecurity(kara_security))
	huma.Get(api, "/api/tags/author/{id}", GetAuthor, setSecurity(kara_security))
	huma.Delete(api, "/api/tags/author/{id}", DeleteAuthor, setSecurity(kara_security))
	huma.Post(api, "/api/tags/author", CreateAuthor, setSecurity(kara_security))

	huma.Get(api, "/api/tags/artist", FindArtist, setSecurity(kara_security))
	huma.Get(api, "/api/tags/artist/{id}", GetArtist, setSecurity(kara_security))
	huma.Delete(api, "/api/tags/artist/{id}", DeleteArtist, setSecurity(kara_security))
	huma.Post(api, "/api/tags/artist", CreateArtist, setSecurity(kara_security))

	huma.Get(api, "/api/tags/media", FindMedia, setSecurity(kara_security))
	huma.Get(api, "/api/tags/media/{id}", GetMedia, setSecurity(kara_security))
	huma.Delete(api, "/api/tags/media/{id}", DeleteMedia, setSecurity(kara_security))
	huma.Post(api, "/api/tags/media", CreateMedia, setSecurity(kara_security))

	huma.Post(api, "/api/token", CreateToken, setSecurity(oidc_security))
	huma.Delete(api, "/api/token/{token}", DeleteToken, setSecurity(oidc_security))

	huma.Get(api, "/api/oidc_discovery", getOIDCDiscovery)
}

func checkToken(ctx huma.Context, bearer_token string, operation_security []map[string][]string) (huma.Context, error) {
	provider, err := rs.NewResourceServerJWTProfile(ctx.Context(), CONFIG.OIDC.Issuer, CONFIG.OIDC.ClientID, CONFIG.OIDC.KeyID, []byte(CONFIG.OIDC.Key))
	if err != nil {
		getLogger().Print(err)
	}

	if bearer_token == "" {
		return ctx, errors.New("No bearer token")
	}

	db := GetDB(ctx.Context())

	for _, sec := range operation_security {
		if len(sec["scopes"]) > 0 {
			scope := sec["scopes"][0]
			db_token := &Token{}
			err := db.Where(Token{ID: bearer_token}).First(db_token).Error
			if err == nil {
				if db_token.HasScope(scope) {
					ctx = huma.WithValue(ctx, "current_user", db_token.User)
					return ctx, nil
				} else {
					return ctx, errors.New(fmt.Sprintf("Token doesn't have the %s API scope", scope))
				}
			} else {
				getLogger().Printf(DBErrToHumaErr(err).Error())
			}
		}

		if len(sec["oidc"]) > 0 {
			resp, err := rs.Introspect[*oidc.IntrospectionResponse](ctx.Context(), provider, bearer_token)
			if err != nil {
				return ctx, err
			}
			if !resp.Active {
				return ctx, errors.New("Forbidden: Inactive account")
			}

			var user_id string
			if CONFIG.OIDC.IDClaim == "" {
				user_id = resp.Subject
			} else {
				user_id = fmt.Sprintf("%v", resp.Claims[CONFIG.OIDC.IDClaim])
			}

			user := User{ID: user_id}
			err = db.First(&user, resp.Subject).Error
			if err != nil {
				if errors.Is(gorm.ErrRecordNotFound, err) {
					// The user doesn't exist yet
					err = db.Create(&user).Error
					if err != nil {
						return ctx, DBErrToHumaErr(err)
					}
				} else {
					return ctx, DBErrToHumaErr(err)
				}
			}
			ctx = huma.WithValue(ctx, "current_user", user)
			return ctx, nil
		}
	}

	return ctx, errors.New("Forbidden")
}

func middlewares(api huma.API) {
	// OIDC/Auth middleware
	api.UseMiddleware(
		func(ctx huma.Context, next func(huma.Context)) {
			auth := ctx.Header("authorization")
			bearer_token, _ := getBearerToken(auth)
			// error value is not needed

			operation_security := ctx.Operation().Security

			if len(operation_security) == 0 {
				next(ctx)
				return
			}

			ctx, err := checkToken(ctx, bearer_token, operation_security)
			if err != nil {
				huma.WriteErr(api, ctx, 403, "Forbidden", err)
				return
			}
			next(ctx)
		},
	)
}

func setupKaraberus() (*fiber.App, huma.API) {
	// Create a new router & API
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024, // 1GiB
	})

	app.Use(logger.New())

	app.Use(filesystem.New(filesystem.Config{
		Root:         http.Dir(CONFIG.UIDistDir),
		Index:        "index.html",
		NotFoundFile: "index.html",
		MaxAge:       3600,
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api")
		},
	}))

	api := humafiber.New(app, huma.DefaultConfig("My API", "1.0.0"))

	// sec := huma.SecurityScheme{
	// 	Type: "openIdConnect",
	// 	Name: "oidc",
	// 	In: "header",
	// 	Scheme: "bearer",
	// }

	middlewares(api)
	routes(api)

	return app, api
}

func RunKaraberus(app *fiber.App, api huma.API) {
	err := CONFIG.OIDC.Validate()
	if err != nil {
		panic(err)
	}

	db := GetDB(context.TODO())
	init_model(db)

	listen_addr := CONFIG.Listen.Addr()
	getLogger().Printf("Starting server on %s...\n", listen_addr)
	getLogger().Fatal(app.Listen(listen_addr))
}
