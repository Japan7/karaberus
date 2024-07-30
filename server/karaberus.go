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
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	huma.Get(api, "/api/kara", GetAllKaras, setSecurity(kara_security))
	huma.Get(api, "/api/kara/{id}", GetKara, setSecurity(kara_security))
	huma.Delete(api, "/api/kara/{id}", DeleteKara, setSecurity(kara_security))
	// TODO: should be reserved to admins
	huma.Patch(api, "/api/kara/{id}/creation_time", SetKaraUploadTime, setSecurity(kara_security))
	huma.Post(api, "/api/kara", CreateKara, setSecurity(kara_security))
	huma.Put(api, "/api/kara/{id}/upload/{filetype}", UploadKaraFile, setSecurity(kara_security))
	huma.Get(api, "/api/kara/{id}/download/{filetype}", DownloadFile, setSecurity(kara_security))

	huma.Get(api, "/api/font", GetAllFonts, setSecurity(kara_security))
	huma.Post(api, "/api/font", UploadFont, setSecurity(kara_security))
	huma.Get(api, "/api/font/{id}", DownloadFont, setSecurity(kara_security))

	huma.Get(api, "/api/tags/audio", GetAudioTags, setSecurity(kara_security))
	huma.Get(api, "/api/tags/video", GetVideoTags, setSecurity(kara_security))

	huma.Get(api, "/api/tags/author", GetAllAuthors, setSecurity(kara_security))
	huma.Get(api, "/api/tags/author/search", FindAuthor, setSecurity(kara_security))
	huma.Get(api, "/api/tags/author/{id}", GetAuthor, setSecurity(kara_security))
	huma.Delete(api, "/api/tags/author/{id}", DeleteAuthor, setSecurity(kara_security))
	huma.Post(api, "/api/tags/author", CreateAuthor, setSecurity(kara_security))

	huma.Get(api, "/api/tags/artist", GetAllArtists, setSecurity(kara_security))
	huma.Get(api, "/api/tags/artist/search", FindArtist, setSecurity(kara_security))
	huma.Get(api, "/api/tags/artist/{id}", GetArtist, setSecurity(kara_security))
	huma.Delete(api, "/api/tags/artist/{id}", DeleteArtist, setSecurity(kara_security))
	huma.Post(api, "/api/tags/artist", CreateArtist, setSecurity(kara_security))

	huma.Get(api, "/api/tags/media", GetAllMedias, setSecurity(kara_security))
	huma.Get(api, "/api/tags/media/types", GetAllMediaTypes, setSecurity(kara_security))
	huma.Get(api, "/api/tags/media/search", FindMedia, setSecurity(kara_security))
	huma.Get(api, "/api/tags/media/{id}", GetMedia, setSecurity(kara_security))
	huma.Delete(api, "/api/tags/media/{id}", DeleteMedia, setSecurity(kara_security))
	huma.Post(api, "/api/tags/media", CreateMedia, setSecurity(kara_security))

	huma.Post(api, "/api/mugen", ImportMugenKara, setSecurity(kara_security))

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
		BodyLimit: 1024 * 1024 * 1024, // 1GiB
	})

	app.Use(logger.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

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

	ctx := context.Background()
	db := GetDB(ctx)
	init_model(db)

	if CONFIG.Dakara.BaseURL != "" {
		go SyncDakara(ctx)
	}

	go SyncMugen(ctx)

	listen_addr := CONFIG.Listen.Addr()
	getLogger().Printf("Starting server on %s...\n", listen_addr)
	getLogger().Fatal(app.Listen(listen_addr))
}
