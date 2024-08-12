// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
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
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	oidc_security := []map[string][]string{{"oidc": []string{""}}}
	kara_security := []map[string][]string{{"oidc": []string{""}, "scopes": []string{"kara"}}}
	kara_ro_security := []map[string][]string{{"oidc": []string{""}, "scopes": []string{"kara_ro"}}}

	huma.Get(api, "/api/kara", GetAllKaras, setSecurity(kara_ro_security))
	huma.Get(api, "/api/kara/{id}", GetKara, setSecurity(kara_ro_security))
	huma.Get(api, "/api/kara/{id}/history", GetKaraHistory, setSecurity(kara_ro_security))
	huma.Get(api, "/api/kara/{id}/issues", GetKaraIssues, setSecurity(kara_ro_security))
	huma.Post(api, "/api/kara/{id}/issues", CreateKaraIssue, setSecurity(kara_security))
	huma.Delete(api, "/api/kara/{id}", DeleteKara, setSecurity(kara_security))
	huma.Patch(api, "/api/kara/{id}", UpdateKara, setSecurity(kara_security))
	huma.Post(api, "/api/kara", CreateKara, setSecurity(kara_security))
	huma.Put(api, "/api/kara/{id}/upload/{filetype}", UploadKaraFile, setSecurity(kara_security))
	huma.Register(api, huma.Operation{
		OperationID: "kara-download-head",
		Method:      http.MethodHead,
		Path:        "/api/kara/{id}/download/{filetype}",
		Security:    kara_ro_security,
	}, DownloadHead)
	huma.Get(api, "/api/kara/{id}/download/{filetype}", DownloadFile, setSecurity(kara_ro_security))

	huma.Get(api, "api/issues", GetIssues, setSecurity(kara_security))

	huma.Get(api, "/api/font", GetAllFonts, setSecurity(kara_ro_security))
	huma.Post(api, "/api/font", UploadFont, setSecurity(kara_security))
	huma.Get(api, "/api/font/{id}/download", DownloadFont, setSecurity(kara_ro_security))

	huma.Get(api, "/api/tags/audio", GetAudioTags, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/video", GetVideoTags, setSecurity(kara_ro_security))

	huma.Get(api, "/api/tags/author", GetAllAuthors, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/author/search", FindAuthor, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/author/{id}", GetAuthor, setSecurity(kara_ro_security))
	huma.Delete(api, "/api/tags/author/{id}", DeleteAuthor, setSecurity(kara_security))
	huma.Post(api, "/api/tags/author", CreateAuthor, setSecurity(kara_security))

	huma.Get(api, "/api/tags/artist", GetAllArtists, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/artist/search", FindArtist, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/artist/{id}", GetArtist, setSecurity(kara_ro_security))
	huma.Delete(api, "/api/tags/artist/{id}", DeleteArtist, setSecurity(kara_security))
	huma.Post(api, "/api/tags/artist", CreateArtist, setSecurity(kara_security))

	huma.Get(api, "/api/tags/media", GetAllMedias, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/media/types", GetAllMediaTypes, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/media/search", FindMedia, setSecurity(kara_ro_security))
	huma.Get(api, "/api/tags/media/{id}", GetMedia, setSecurity(kara_ro_security))
	huma.Delete(api, "/api/tags/media/{id}", DeleteMedia, setSecurity(kara_security))
	huma.Post(api, "/api/tags/media", CreateMedia, setSecurity(kara_security))

	huma.Post(api, "/api/mugen", ImportMugenKara, setSecurity(kara_security))
	huma.Post(api, "/api/mugen/refresh", RefreshMugen, setSecurity(kara_security))
	huma.Get(api, "/api/mugen", GetMugenImports, setSecurity(kara_ro_security))
	huma.Delete(api, "/api/mugen/{id}", DeleteMugenImport, setSecurity(kara_security))

	huma.Post(api, "/api/dakara/sync", StartDakaraSync, setSecurity(kara_security))

	huma.Get(api, "/api/token", GetAllUserTokens, setSecurity(oidc_security))
	huma.Post(api, "/api/token", CreateToken, setSecurity(oidc_security))
	huma.Delete(api, "/api/token/{token}", DeleteToken, setSecurity(oidc_security))
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
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
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

	addOidcRoutes(app)

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

func RunKaraberus(app *fiber.App, api huma.API) {
	err := CONFIG.OIDC.Validate()
	if err != nil {
		panic(err)
	}

	db := GetDB(context.Background())
	init_model(db)

	if CONFIG.Dakara.BaseURL != "" {
		go SyncDakaraLoop(context.Background())
	}

	go SyncMugen(context.Background())

	listen_addr := CONFIG.Listen.Addr()
	getLogger().Printf("Starting server on %s...\n", listen_addr)
	getLogger().Fatal(app.Listen(listen_addr))
}
