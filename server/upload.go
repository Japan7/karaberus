// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"io"
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

func parseID(kid_str string) (uint, error) {
	kid_int, err := strconv.Atoi(kid_str)
	if kid_int < 0 {
		return 0, errors.New("Kara ID cannot be negative")
	}
	return uint(kid_int), err
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	db := GetDB(ctx)

	kid, err := parseID(input.KID)
	if err != nil {
		return nil, err
	}

	kara, err := GetKaraByID(db, kid)

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

	err = updateKaraokeAfterUpload(db, &kara, input.FileType)
	if err != nil {
		return nil, err
	}

	res, err := CheckKara(ctx, kara)
	if err != nil {
		return nil, err
	}

	resp := &UploadOutput{}
	resp.Body.CheckResults = *res
	resp.Body.KID = input.KID

	return resp, nil
}

type DownloadInput struct {
	KID      string `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
}

func DownloadFile(ctx context.Context, input *DownloadInput) (*huma.StreamResponse, error) {
	db := GetDB(ctx)
	kid, err := parseID(input.KID)
	if err != nil {
		return nil, err
	}

	kara, err := GetKaraByID(db, kid)
	if err != nil {
		return nil, err
	}

	obj, err := GetKaraObject(ctx, kara, input.FileType)
	if err != nil {
		return nil, err
	}

	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			writer := ctx.BodyWriter()
			var n int
			for {
				buf := make([]byte, 1024*1024)
				n, err = obj.Read(buf)
				writer.Write(buf[:n])
				if err != nil {
					if errors.Is(err, io.EOF) {
						err = nil
					}
					break
				}
			}
		},
	}, err
}
