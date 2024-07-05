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
		KID          string          `json:"file_id" example:"1" doc:"karaoke ID"`
		CheckResults CheckKaraOutput `json:"check_results"`
	}
}

func updateKaraokeAfterUpload(tx *gorm.DB, kara *KaraInfoDB, filetype string) error {
	switch filetype {
	case "video":
		err := tx.Model(kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{VideoUploaded: true}}).Error
		return DBErrToHumaErr(err)
	case "inst":
		err := tx.Model(kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{InstrumentalUploaded: true}}).Error
		return DBErrToHumaErr(err)
	case "sub":
		err := tx.Model(kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{SubtitlesUploaded: true}}).Error
		return DBErrToHumaErr(err)
	}
	return errors.New("Unknown file type " + filetype)
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	db := GetDB(ctx)

	kid, err := strconv.Atoi(input.KID)
	if err != nil {
		return nil, err
	}

	kara := &KaraInfoDB{}
	err = db.First(kara, kid).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
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

	err = updateKaraokeAfterUpload(db, kara, input.FileType)
	if err != nil {
		return nil, err
	}

	res, err := CheckKara(ctx, *kara)
	if err != nil {
		return nil, err
	}

	resp := &UploadOutput{}
	resp.Body.CheckResults = *res
	resp.Body.KID = input.KID

	return resp, nil
}
