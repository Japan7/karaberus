// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

type KaraberusError struct {
	Message string
}

func (m *KaraberusError) Error() string {
	return m.Message
}

func addMiddlewares(api huma.API) {
	api.UseMiddleware(authMiddleware)
}

func addRoutes(api huma.API) {
	oidc := RouteSecurity{OIDC: true}.toSecurity()
	oidc_admin := RouteSecurity{OIDC: true, Admin: true}.toSecurity()
	kara := RouteSecurity{OIDC: true, Scopes: Scopes{Kara: true}}.toSecurity()
	kara_admin := RouteSecurity{OIDC: true, Admin: true, Scopes: Scopes{Kara: true}}.toSecurity()
	kara_ro := RouteSecurity{OIDC: true, Scopes: Scopes{KaraRO: true}}.toSecurity()
	kara_ro_basic := RouteSecurity{OIDC: true, Basic: true, Scopes: Scopes{KaraRO: true}}.toSecurity()
	user := RouteSecurity{OIDC: true, Scopes: Scopes{User: true}}.toSecurity()
	user_admin := RouteSecurity{OIDC: true, Admin: true, Scopes: Scopes{User: true}}.toSecurity()

	huma.Get(api, "/api/kara", GetAllKaras, setSecurity(kara_ro))
	huma.Get(api, "/api/kara/{id}", GetKara, setSecurity(kara_ro))
	huma.Get(api, "/api/kara/{id}/history", GetKaraHistory, setSecurity(kara_ro))
	huma.Delete(api, "/api/kara/{id}", DeleteKara, setSecurity(kara))
	huma.Patch(api, "/api/kara/{id}", UpdateKara, setSecurity(kara))
	huma.Post(api, "/api/kara", CreateKara, setSecurity(kara))
	huma.Put(api, "/api/kara/{id}/upload/{filetype}", UploadKaraFile, setSecurity(kara))
	huma.Delete(api, "/api/kara/{id}/{filetype}", DeleteKaraFile, setSecurity(kara_admin))
	huma.Register(api, huma.Operation{
		OperationID: "kara-download-head",
		Method:      http.MethodHead,
		Path:        "/api/kara/{id}/download/{filetype}",
		Security:    kara_ro,
	}, DownloadHead)
	huma.Get(api, "/api/kara/{id}/download/{filetype}", DownloadFile, setSecurity(kara_ro_basic))
	huma.Get(api, "/api/kara/{id}/mugen/export", MugenExportKara, setSecurity(kara_admin))

	huma.Get(api, "/api/font", GetAllFonts, setSecurity(kara_ro))
	huma.Post(api, "/api/font", UploadFont, setSecurity(kara))
	huma.Get(api, "/api/font/{id}/download", DownloadFont, setSecurity(kara_ro))

	huma.Get(api, "/api/tags/audio", GetAudioTags, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/video", GetVideoTags, setSecurity(kara_ro))

	huma.Get(api, "/api/tags/author", GetAllAuthors, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/author/search", FindAuthor, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/author/{id}", GetAuthor, setSecurity(kara_ro))
	huma.Delete(api, "/api/tags/author/{id}", DeleteAuthor, setSecurity(kara))
	huma.Patch(api, "/api/tags/author/{id}", UpdateAuthor, setSecurity(kara))
	huma.Post(api, "/api/tags/author", CreateAuthor, setSecurity(kara))

	huma.Get(api, "/api/tags/artist", GetAllArtists, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/artist/search", FindArtist, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/artist/{id}", GetArtist, setSecurity(kara_ro))
	huma.Delete(api, "/api/tags/artist/{id}", DeleteArtist, setSecurity(kara))
	huma.Patch(api, "/api/tags/artist/{id}", UpdateArtist, setSecurity(kara))
	huma.Post(api, "/api/tags/artist", CreateArtist, setSecurity(kara))

	huma.Get(api, "/api/tags/media", GetAllMedias, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/media/types", GetAllMediaTypes, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/media/search", FindMedia, setSecurity(kara_ro))
	huma.Get(api, "/api/tags/media/{id}", GetMedia, setSecurity(kara_ro))
	huma.Delete(api, "/api/tags/media/{id}", DeleteMedia, setSecurity(kara))
	huma.Patch(api, "/api/tags/media/{id}", UpdateMedia, setSecurity(kara))
	huma.Post(api, "/api/tags/media", CreateMedia, setSecurity(kara))

	huma.Post(api, "/api/mugen", ImportMugenKara, setSecurity(kara))
	huma.Post(api, "/api/mugen/refresh", RefreshMugen, setSecurity(kara_admin))
	huma.Get(api, "/api/mugen", GetMugenImports, setSecurity(kara_ro))
	huma.Delete(api, "/api/mugen/{id}", DeleteMugenImport, setSecurity(kara_admin))

	huma.Post(api, "/api/dakara/sync", StartDakaraSync, setSecurity(kara))

	huma.Get(api, "/api/token", GetAllUserTokens, setSecurity(oidc))
	huma.Post(api, "/api/token", CreateToken, setSecurity(oidc))
	huma.Delete(api, "/api/token/{token}", DeleteToken, setSecurity(oidc))

	huma.Get(api, "/api/gitlab/authorize", GitlabAuth, setSecurity(oidc_admin))
	huma.Get(api, "/api/gitlab/callback", GitlabCallback, setSecurity(oidc_admin))

	huma.Get(api, "/api/user/{id}", GetUser, setSecurity(user))
	huma.Get(api, "/api/me", GetMe, setSecurity(user))
	huma.Put(api, "/api/user/{id}/author", UpdateUserAuthor, setSecurity(user_admin))
	huma.Put(api, "/api/me/author", UpdateMeAuthor, setSecurity(user))
}

type RouteSecurity struct {
	OIDC  bool
	Admin bool
	Basic bool
	Scopes
}

func (security RouteSecurity) toSecurity() []map[string][]string {
	security_mapping := map[string][]string{}

	if security.OIDC {
		security_mapping["oidc"] = []string{""}
	}

	if security.Admin {
		security_mapping["admin"] = []string{""}
	}

	if security.Basic {
		security_mapping["basic"] = []string{""}
	}

	// use the json names of the scopes, could be done through reflection instead
	json_scopes, err := json.Marshal(security.Scopes)
	if err != nil {
		panic(err)
	}
	scopes_map := map[string]bool{}
	err = json.Unmarshal(json_scopes, &scopes_map)
	if err != nil {
		panic(err)
	}

	scopes_list := []string{}
	for k, v := range scopes_map {
		if v {
			scopes_list = append(scopes_list, k)
		}
	}

	if len(scopes_list) > 0 {
		security_mapping["scopes"] = scopes_list
	}

	return []map[string][]string{security_mapping}
}

func setSecurity(security []map[string][]string) func(o *huma.Operation) {
	return func(o *huma.Operation) {
		o.Security = security
	}
}

func setupKaraberus() (*fiber.App, huma.API) {
	// Create a new router & API
	app := fiber.New(fiber.Config{
		BodyLimit:                    1024 * 1024 * 1024, // 1GiB
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
	})

	app.Use(logger.New())
	app.Use(healthcheck.New())
	app.Use(compress.New())
	if CONFIG.Listen.Profiling {
		app.Use(pprof.New())
	}

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

	addMiddlewares(api)
	addRoutes(api)

	// sec := huma.SecurityScheme{
	// 	Type: "openIdConnect",
	// 	Name: "oidc",
	// 	In: "header",
	// 	Scheme: "bearer",
	// }

	return app, api
}

type KaraberusInit struct{}

func isKaraberusInit(ctx context.Context) bool {
	return ctx.Value(KaraberusInit{}) != nil
}

func RunKaraberus(app *fiber.App, api huma.API) {
	err := CONFIG.OIDC.Validate()
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), KaraberusInit{}, true)
	addOidcRoutes(ctx, app)
	initS3Clients(ctx)
	init_db(ctx)

	if CONFIG.Dakara.BaseURL != "" {
		go SyncDakaraLoop(context.Background())
	}

	go SyncMugen(context.Background())

	listen_addr := CONFIG.Listen.Addr()
	getLogger().Printf("Starting server on %s...\n", listen_addr)
	getLogger().Fatal(app.Listen(listen_addr))
}
