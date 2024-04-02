// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

// #cgo pkg-config: dakara_check
// #include <stdlib.h>
// #include <dakara_check.h>
import "C"

import (
	"context"
	"log"
	"mime/multipart"
	"path/filepath"
	"unsafe"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UploadInput struct {
	FileID uuid.UUID `path:"kid"`
}

func getFilePathStr(fileIDStr string) string {
	return filepath.Join(FILES_DIR, fileIDStr)
}

func getFilePath(fileID uuid.UUID) string {
	return getFilePathStr(fileID.String())
}

func saveFile(fd multipart.File, kid uuid.UUID, type_directory string) (error) {
	filename := filepath.Join(type_directory, "/", kid.String())
	return UploadToS3(fd, BUCKET_NAME, filename)
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

	err = saveFile(fd, m.FileID, "video")
	if err != nil {
		return []error{err}
	}

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

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	resp := &UploadOutput{}
	cfilepath := C.CString(getFilePath(input.FileID))
	defer C.free(unsafe.Pointer(cfilepath))
	dakara_check_results := C.dakara_check(cfilepath, 0)
	defer C.dakara_check_results_free(dakara_check_results)

	resp.Body.Passed = bool(dakara_check_results.passed)
	resp.Body.FileID = input.FileID.String()

	return resp, nil
}
