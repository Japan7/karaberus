// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type UploadData struct {
	UploadFile multipart.File `form-data:"file" required:"true"`
}

type UploadInput struct {
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
	RawBody  huma.MultipartFormFiles[UploadData]
}

type UploadOutput struct {
	Body struct {
		KID          uint            `json:"file_id" example:"1" doc:"karaoke ID"`
		CheckResults CheckKaraOutput `json:"check_results"`
	}
}

func updateKaraokeAfterUpload(tx *gorm.DB, kara *KaraInfoDB, filetype string) error {
	currentTime := time.Now().UTC()
	switch filetype {
	case "video":
		err := tx.Model(kara).Updates(&KaraInfoDB{
			UploadInfo: UploadInfo{
				VideoUploaded: true,
				VideoModTime:  currentTime,
			}}).Error
		return DBErrToHumaErr(err)
	case "inst":
		err := tx.Model(kara).Updates(&KaraInfoDB{
			UploadInfo: UploadInfo{
				InstrumentalUploaded: true,
				InstrumentalModTime:  currentTime,
			}}).Error
		return DBErrToHumaErr(err)
	case "sub":
		if kara.KaraokeCreationTime.Unix() == 0 {
			err := tx.Model(kara).Updates(&KaraInfoDB{
				UploadInfo: UploadInfo{
					SubtitlesUploaded:   true,
					SubtitlesModTime:    currentTime,
					KaraokeCreationTime: currentTime,
				}}).Error
			return DBErrToHumaErr(err)
		} else {
			err := tx.Model(kara).Updates(&KaraInfoDB{
				UploadInfo: UploadInfo{
					SubtitlesUploaded: true,
					SubtitlesModTime:  currentTime,
				}}).Error
			return DBErrToHumaErr(err)
		}
	}
	return errors.New("Unknown file type " + filetype)
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	db := GetDB(ctx)

	kid := input.KID
	kara, err := GetKaraByID(db, kid)

	file := input.RawBody.Form.File["file"][0]
	fd, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	err = SaveFileToS3(ctx, fd, kid, input.FileType, file.Size)
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
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
}

func serveObject(obj *minio.Object) (*huma.StreamResponse, error) {
	var err error
	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			defer obj.Close()

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

func DownloadFile(ctx context.Context, input *DownloadInput) (*huma.StreamResponse, error) {
	db := GetDB(ctx)
	kid := input.KID

	kara, err := GetKaraByID(db, kid)
	if err != nil {
		return nil, err
	}

	obj, err := GetKaraObject(ctx, kara, input.FileType)
	if err != nil {
		return nil, err
	}

	return serveObject(obj)
}
