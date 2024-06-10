// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"mime/multipart"

	"github.com/danielgtaylor/huma/v2"
)

type UploadData struct {
	UploadFile multipart.File `form-data:"file" required:"true"`
}

type UploadInput struct {
	KID      string `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
	RawBody  huma.MultipartFormFiles[UploadData]
}

type UploadOutput struct {
	Body struct {
		KID          string            `json:"file_id" example:"1" doc:"karaoke ID"`
		CheckResults CheckS3FileOutput `json:"check_results"`
	}
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	resp := &UploadOutput{}

	file := input.RawBody.Form.File["file"][0]
	fd, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	err = SaveFileToS3(ctx, fd, input.KID, input.FileType, file.Size)
	if err != nil {
		return nil, err
	}

	res, err := CheckKara(ctx, input.KID)
	if err != nil {
		return nil, err
	}

	resp.Body.CheckResults = *res
	resp.Body.KID = input.KID

	return resp, nil
}
