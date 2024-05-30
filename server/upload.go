// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UploadInput struct {
	KID      string `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
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
		panic(err)
	}

	file := multipartForm.File["file"][0]

	fd, err := file.Open()
	if err != nil {
		return []error{err}
	}
	defer fd.Close()

	err = SaveFileToS3(ctx.Context(), fd, m.KID, m.FileType, file.Size)
	if err != nil {
		return []error{err}
	}

	return nil
}

// Ensure UploadInput implements huma.Resolver
var _ huma.Resolver = (*UploadInput)(nil)

// only used to create the OpenAPI schema
type UploadInputDef struct {
	KID      uuid.UUID `path:"kid" example:"9da21d68-6c4d-493b-bd1f-ab1d08c1234f"`
	FileType string    `path:"filetype" example:"video"`
	File     string    `json:"file" format:"binary" example:"@file.mkv"`
}

type UploadOutput struct {
	Body struct {
		KID          string            `json:"file_id" example:"1" doc:"karaoke ID"`
		CheckResults CheckS3FileOutput `json:"check_results"`
	}
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	resp := &UploadOutput{}

	res, err := CheckKara(ctx, input.KID)
	if err != nil {
		return nil, err
	}

	resp.Body.CheckResults = *res
	resp.Body.KID = input.KID

	return resp, nil
}
