// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"mime/multipart"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
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

func updateKaraokeAfterUpload(kara *KaraInfoDB, filetype string) error {
	switch filetype {
	case "video":
		tx := GetDB().Model(kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{VideoUploaded: true}})
		return tx.Error
	case "inst":
		tx := GetDB().Model(kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{InstrumentalUploaded: true}})
		return tx.Error
	case "sub":
		tx := GetDB().Model(kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{SubtitlesUploaded: true}})
		return tx.Error
	}
	return errors.New("Unknown file type " + filetype)
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	resp := &UploadOutput{}

	kid, err := strconv.Atoi(input.KID)
	if err != nil {
		return nil, err
	}

	kara := &KaraInfoDB{}
	tx := GetDB().First(kara, kid)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("Karaoke not found", tx.Error)
		}
		return nil, tx.Error
	}

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

	updateKaraokeAfterUpload(kara, input.FileType)

	res, err := CheckKara(ctx, *kara)
	if err != nil {
		return nil, err
	}

	resp.Body.CheckResults = *res
	resp.Body.KID = input.KID

	return resp, nil
}
