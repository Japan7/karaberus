// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type UploadData struct {
	UploadFile multipart.File `form-data:"file" required:"true"`
}

type UploadInputDefinition struct {
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
	RawBody  huma.MultipartFormFiles[UploadData]
}

type UploadTempFile struct {
	Fd   *os.File
	Size int64
	// original file name
	Name string
}

type UploadInput struct {
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" example:"video"`
	File     UploadTempFile
}

// write uploaded file to temporary file so fasthttp can't store it in memory
// which could easily lead to OOMs.
// minio PutObject wants a io.Seeker or will write file to a buffer which is
// another reason for a temporary file (otherwise we could just stream it).
func createTempFile(ctx huma.Context, tempfile *UploadTempFile) error {
	content_type, params, err := mime.ParseMediaType(ctx.Header("Content-Type"))
	if err != nil {
		return err
	}
	if !strings.HasPrefix(content_type, "multipart/") {
		return errors.New("not a multipart request")
	}
	boundary := params["boundary"]

	reader := ctx.BodyReader()
	mr := multipart.NewReader(reader, boundary)

	if err != nil {
		return err
	}

	for {
		part, err := mr.NextPart()
		if err != nil {
			return err
		}
		defer part.Close()

		if part.FormName() == "file" {
			tempfile.Name = part.FileName()

			fd, err := os.CreateTemp("", "karaberus-*")
			if err != nil {
				return err
			}

			// roughly io.Copy but with a small buffer
			// don't change mindlessly
			buf := make([]byte, 1024*8)
			for {
				n, err := part.Read(buf)
				if errors.Is(err, io.EOF) {
					if n == 0 {
						break
					}
				} else if err != nil {
					return err
				}
				_, err = fd.Write(buf[:n])
				if err != nil {
					return err
				}
			}

			tempfile.Fd = fd

			stat, err := fd.Stat()
			if err != nil {
				return err
			}
			tempfile.Size = stat.Size()

			_, err = fd.Seek(0, 0)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (i *UploadInput) Resolve(ctx huma.Context) []error {
	err := createTempFile(ctx, &i.File)
	if err != nil {
		return []error{err}
	}
	return nil
}

var _ huma.Resolver = (*UploadInput)(nil)

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
	defer os.Remove(input.File.Fd.Name())
	defer input.File.Fd.Close()

	kid := input.KID
	kara, err := GetKaraByID(db, kid)

	resp := &UploadOutput{}
	err = db.Transaction(func(tx *gorm.DB) error {
		res, err := SaveFileToS3(ctx, tx, input.File.Fd, &kara, input.FileType, input.File.Size)
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
