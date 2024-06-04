// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type KaraberusError struct {
	Message string
}

func (m *KaraberusError) Error() string {
	return m.Message
}

var FILES_DIR = getEnvDefault("FILES_DIR", "files")
var LISTEN_ADDR = getEnvDefault("LISTEN_ADDR", "127.0.0.1:8888")
var BUCKET_NAME = getEnvDefault("BUCKET_NAME", "karaberus")
var S3_ENDPOINT = getEnvDefault("S3_ENDPOINT", "")
var S3_KEYID = getEnvDefault("S3_KEYID", "")
var S3_SECRET = getEnvDefault("S3_SECRET", "")

func getEnvDefault(name string, defaultValue string) string {
	envVar := os.Getenv("KARABERUS_" + name)
	if envVar != "" {
		return envVar
	}

	return defaultValue
}

func init() {
	init_db()
	init_model()
}

func getBearerToken(auth string) string {
	if auth == "" {
		panic("Authorization header is missing.")
	}
	if strings.HasPrefix(auth, oidc.BearerToken) {
		return strings.TrimPrefix(auth, oidc.PrefixBearer)
	}
	return ""
}

func setSecurity(security []map[string][]string) func(o *huma.Operation) {
	return func(o *huma.Operation) {
		o.Security = security
	}
}

func routes(api huma.API) {
	oidc_security := []map[string][]string{{"oidc": []string{""}}}

	huma.Get(api, "/kara/{id}", GetKara, setSecurity(oidc_security))
	huma.Delete(api, "/kara/{id}", DeleteKara, setSecurity(oidc_security))
	huma.Post(api, "/kara", CreateKara, setSecurity(oidc_security))
	huma.Put(api, "/kara/{id}/upload/{filetype}", UploadKaraFile, setSecurity(oidc_security))

	huma.Get(api, "/tags/audio", GetAudioTags, setSecurity(oidc_security))
	huma.Get(api, "/tags/video", GetVideoTags, setSecurity(oidc_security))

	huma.Get(api, "/tags/author", FindAuthor, setSecurity(oidc_security))
	huma.Get(api, "/tags/author/{id}", GetAuthor, setSecurity(oidc_security))
	huma.Delete(api, "/tags/author/{id}", DeleteAuthor, setSecurity(oidc_security))
	huma.Post(api, "/tags/author", CreateAuthor, setSecurity(oidc_security))

	huma.Get(api, "/tags/artist", FindArtist, setSecurity(oidc_security))
	huma.Get(api, "/tags/artist/{id}", GetArtist, setSecurity(oidc_security))
	huma.Delete(api, "/tags/artist/{id}", DeleteArtist, setSecurity(oidc_security))
	huma.Post(api, "/tags/artist", CreateArtist, setSecurity(oidc_security))

	huma.Get(api, "/tags/media", FindMedia, setSecurity(oidc_security))
	huma.Get(api, "/tags/media/{id}", GetMedia, setSecurity(oidc_security))
	huma.Delete(api, "/tags/media/{id}", DeleteMedia, setSecurity(oidc_security))
	huma.Post(api, "/tags/media", CreateMedia, setSecurity(oidc_security))
}

var OIDC_ISSUER = getEnvDefault("OIDC_ISSUER", "")
var OIDC_KEY_PATH = getEnvDefault("OIDC_KEY", "")
var OIDC_ID_CLAIM = getEnvDefault("OIDC_ID_CLAIM", "")

func middlewares(api huma.API) {
	if OIDC_ISSUER == "" {
		panic("OIDC issuer is not set")
	}
	if OIDC_KEY_PATH == "" {
		panic("OIDC key is not set")
	}

	provider, err := rs.NewResourceServerFromKeyFile(context.TODO(), OIDC_ISSUER, OIDC_KEY_PATH)
	if err != nil {
		panic(err)
	}

	// OIDC/Auth middleware
	api.UseMiddleware(
		func(ctx huma.Context, next func(huma.Context)) {
			auth := ctx.Header("authorization")
			bearer_token := getBearerToken(auth)

			for _, sec := range ctx.Operation().Security {
				if len(sec["oidc"]) > 0 {
					if bearer_token == "" {
						continue
					}
					resp, err := rs.Introspect[*oidc.IntrospectionResponse](ctx.Context(), provider, bearer_token)
					if err != nil {
						huma.WriteErr(api, ctx, 403, "Forbidden", err)
						return
					}
					if !resp.Active {
						huma.WriteErr(api, ctx, 403, "Forbidden: Inactive account")
						return
					}

					var user_id string
					if OIDC_ID_CLAIM == "" {
						user_id = resp.Subject
					} else {
						user_id = fmt.Sprintf("%v", resp.Claims[OIDC_ID_CLAIM])
					}

					user := User{ID: user_id}
					tx := GetDB().First(&user, resp.Subject)
					if tx.Error != nil {
						if errors.Is(gorm.ErrRecordNotFound, tx.Error) {
							tx = GetDB().Create(&user)
							if tx.Error != nil {
								huma.WriteErr(api, ctx, 500, "Failed to create user account")
								return
							}
						} else {
							huma.WriteErr(api, ctx, 500, "Failed to find user account")
							return
						}
					}
					ctx = huma.WithValue(ctx, "current_user", user)
					next(ctx)
				}
			}

			huma.WriteErr(api, ctx, 403, "Forbidden")
		},
	)
}

func RunKaraberus() {
	// Create a new router & API
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024, // 1GiB
	})
	api := humafiber.New(app, huma.DefaultConfig("My API", "1.0.0"))

	// sec := huma.SecurityScheme{
	// 	Type: "openIdConnect",
	// 	Name: "oidc",
	// 	In: "header",
	// 	Scheme: "bearer",
	// }

	middlewares(api)
	routes(api)

	log.Printf("Starting server at %s...\n", LISTEN_ADDR)
	log.Fatal(app.Listen(LISTEN_ADDR))
}
