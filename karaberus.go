// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package main

// #cgo pkg-config: dakara_check
// #include <stdlib.h>
// #include <dakara_check.h>
//
// void karaberus_dakara_check_results_free(struct dakara_check_results *res) {
//   dakara_check_results_free(res);
// }
//
import "C"

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var FILES_DIR = getEnvDefault("FILES_DIR", "files")
var LISTEN_ADDR = getEnvDefault("LISTEN_ADDR", "127.0.0.1:8888")

type UploadInput struct {
	FileID uuid.UUID
}

func getEnvDefault(name string, defaultValue string) string {
	envVar := os.Getenv("KARABERUS_" + name)
	if envVar != "" {
		return envVar
	}

	return defaultValue
}

func getFilePathStr(fileIDStr string) string {
	return filepath.Join(FILES_DIR, fileIDStr)
}

func getFilePath(fileID uuid.UUID) string {
	return getFilePathStr(fileID.String())
}

func saveFile(fd multipart.File) (*uuid.UUID, error) {
	buf := make([]byte, 4*1024*1024) // 4MiB
	filesdir := FILES_DIR

	err := os.MkdirAll(filesdir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	fileUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	filepath := filepath.Join(filesdir, fileUUID.String())

	wfd, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer wfd.Close()

	for {
		n, err := fd.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			goto cleanup
		}

		if n < len(buf) {
			buf = buf[:n]
		}
		_, err = wfd.Write(buf)
		if err != nil {
			goto cleanup
		}
	}

	return &fileUUID, nil

cleanup:
	log.Println("cleaning up " + filepath)
	err = os.Remove(filepath)
	if err != nil {
		log.Println(err)
	}
	return nil, err
}

func (m *UploadInput) Resolve(ctx huma.Context) []error {
	multipartForm, err := ctx.GetMultipartForm()

	if err != nil {
		log.Fatal(err)
		return []error{err}
	}

	file := multipartForm.File["file"][0]

	fd, err := file.Open()
	if err != nil {
		return []error{err}
	}
	defer fd.Close()

	fileid, err := saveFile(fd)
	if err != nil {
		return []error{err}
	}

	m.FileID = *fileid

	return nil
}

// Ensure UploadInput implements huma.Resolver
var _ huma.Resolver = (*UploadInput)(nil)

// only used to create the OpenAPI schema
type UploadInputDef struct {
	File string `json:"file" format:"binary" example:"@file.mkv"`
}

type UploadOutput struct {
	Body struct {
		FileID string `json:"file_id" example:"79d97afe-b2db-4b55-af82-f16b60d8ae77" doc:"file ID"`
		Passed bool   `json:"passed" example:"true" doc:"true if file passed all checks"`
	}
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
	}, func(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
		resp := &UploadOutput{}
		cfilepath := C.CString(getFilePath(input.FileID))
		defer C.free(unsafe.Pointer(cfilepath))
		dakara_check_results := C.dakara_check(cfilepath, 0)
		defer C.karaberus_dakara_check_results_free(dakara_check_results)

		resp.Body.Passed = bool(dakara_check_results.passed)
		resp.Body.FileID = input.FileID.String()

		return resp, nil
	})

}

func main() {
	// Create a new router & API
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024, // 1GiB
	})
	api := humafiber.New(app, huma.DefaultConfig("My API", "1.0.0"))

	routes(api)

	log.Printf("Starting server at %s...\n", LISTEN_ADDR)
	log.Fatal(app.Listen(LISTEN_ADDR))
}
