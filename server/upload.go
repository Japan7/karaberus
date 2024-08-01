// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"fmt"
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
		kara.VideoUploaded = true
		kara.VideoModTime = currentTime
		return nil
	case "inst":
		kara.InstrumentalUploaded = true
		kara.InstrumentalModTime = currentTime
		return nil
	case "sub":
		kara.SubtitlesUploaded = true
		kara.SubtitlesModTime = currentTime
		if kara.KaraokeCreationTime.Unix() == 0 {
			kara.KaraokeCreationTime = currentTime
		}
		return nil
	}
	return errors.New("Unknown file type " + filetype)
}

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	db := GetDB(ctx)

	defer input.RawBody.Form.RemoveAll()

	kid := input.KID
	kara, err := GetKaraByID(db, kid)

	file := input.RawBody.Form.File["file"][0]
	fd, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	resp := &UploadOutput{}
	err = db.Transaction(func(tx *gorm.DB) error {
		res, err := SaveFileToS3(ctx, tx, fd, &kara, input.FileType, file.Size)
		if err != nil {
			return err
		}

		resp.Body.CheckResults = *res
		resp.Body.KID = input.KID
		return nil
	})

	return resp, err
}

type DownloadInput struct {
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
}

func serveObject(obj *minio.Object) (*huma.StreamResponse, error) {
	stat, err := obj.Stat()

	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			defer obj.Close()

			if err != nil {
				resp := minio.ToErrorResponse(err)
				if resp.Code == "NoSuchKey" {
					ctx.SetStatus(404)
				} else {
					ctx.SetStatus(500)
					getLogger().Printf("%+v\n", resp)
				}
				return
			}

			ctx.SetHeader("Content-Length", fmt.Sprintf("%d", stat.Size))

			writer := ctx.BodyWriter()

			buf := make([]byte, 1024*1024)
			var n int
			for {
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
