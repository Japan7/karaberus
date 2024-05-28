// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
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

func routes(api huma.API) {
	// Create a registry and register a type.
	registry := huma.NewMapRegistry("#/karaberus", huma.DefaultSchemaNamer)

	// Register POST /upload
	huma.Register(api, huma.Operation{
		OperationID: "upload",
		Summary:     "Upload karaoke file",
		Method:      http.MethodPost,
		Path:        "/upload/{kid}",
		RequestBody: &huma.RequestBody{
			Content: map[string]*huma.MediaType{
				"multipart/form-data": {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(UploadInputDef{})),
				},
			},
		},
	}, UploadKaraFile)

	huma.Post(api, "/kara", CreateKara)

	huma.Get(api, "/tags/audio", GetAudioTags)
	huma.Get(api, "/tags/video", GetVideoTags)

	huma.Get(api, "/tags/author", FindAuthor)
	huma.Get(api, "/tags/author/{id}", GetAuthor)
	huma.Delete(api, "/tags/author/{id}", DeleteAuthor)
	huma.Post(api, "/tags/author", CreateAuthor)

	huma.Get(api, "/tags/artist", FindArtist)
	huma.Get(api, "/tags/artist/{id}", GetArtist)
	huma.Delete(api, "/tags/artist/{id}", DeleteArtist)
	huma.Post(api, "/tags/artist", CreateArtist)

	huma.Get(api, "/tags/media", FindMedia)
	huma.Get(api, "/tags/media/{id}", GetMedia)
	huma.Delete(api, "/tags/media/{id}", DeleteMedia)
	huma.Post(api, "/tags/media", CreateMedia)
}

func RunKaraberus() {
	// Create a new router & API
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024, // 1GiB
	})
	api := humafiber.New(app, huma.DefaultConfig("My API", "1.0.0"))

	routes(api)

	log.Printf("Starting server at %s...\n", LISTEN_ADDR)
	log.Fatal(app.Listen(LISTEN_ADDR))
}
