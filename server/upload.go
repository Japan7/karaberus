// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"log"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UploadInput struct {
	FileID uuid.UUID `path:"kid" example:"9da21d68-6c4d-493b-bd1f-ab1d08c1234f"`
}

func getFilePathStr(fileIDStr string) string {
	return filepath.Join(FILES_DIR, fileIDStr)
}

func getFilePath(fileID uuid.UUID) string {
	return getFilePathStr(fileID.String())
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

	err = SaveFileToS3(ctx.Context(), fd, m.FileID, "video")
	if err != nil {
		return []error{err}
	}

	return nil
}

// Ensure UploadInput implements huma.Resolver
var _ huma.Resolver = (*UploadInput)(nil)

// only used to create the OpenAPI schema
type UploadInputDef struct {
	FileID uuid.UUID `path:"kid" example:"9da21d68-6c4d-493b-bd1f-ab1d08c1234f"`
	File string `json:"file" format:"binary" example:"@file.mkv"`
}

type UploadOutput struct {
	Body struct {
		FileID string `json:"file_id" example:"79d97afe-b2db-4b55-af82-f16b60d8ae77" doc:"file ID"`
		CheckResults CheckS3FileOutput
	}
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	resp := &UploadOutput{}

	res, err := CheckS3File(ctx, input.FileID)
	if (err != nil) {
		return nil, err
	}

	resp.Body.CheckResults = *res
	resp.Body.FileID = input.FileID.String()

	return resp, nil
}
