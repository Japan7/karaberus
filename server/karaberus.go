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

var FILES_DIR = getEnvDefault("FILES_DIR", "files")
var LISTEN_ADDR = getEnvDefault("LISTEN_ADDR", "127.0.0.1:8888")

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
		Path:        "/upload",
		RequestBody: &huma.RequestBody{
			Content: map[string]*huma.MediaType{
				"multipart/form-data": {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(UploadInputDef{})),
				},
			},
		},
	}, UploadKaraFile)

	huma.Register(api, huma.Operation{
		OperationID: "create_kara",
		Summary:     "Create karaoke",
		Method:      http.MethodPost,
		Path:        "/kara",
	}, CreateKara)

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
