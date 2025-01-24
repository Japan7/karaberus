// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ironsmile/nedomi/utils/httputils"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
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
	Name  string
	CRC32 uint32
}

type UploadInput struct {
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" enum:"video,sub,inst" example:"video"`
	File     UploadTempFile
}

func CreateTempFile(ctx context.Context, tempfile *UploadTempFile, reader io.Reader) error {
	fd, err := os.CreateTemp("", "karaberus-*")
	if err != nil {
		return err
	}

	hasher := crc32.NewIEEE()
	// roughly io.Copy but with a small buffer
	// don't change mindlessly
	buf := make([]byte, 1024*8)
	for {
		n, err := reader.Read(buf)
		if errors.Is(err, io.EOF) {
			if n == 0 {
				break
			}
		} else if err != nil {
			return err
		}
		_, err = hasher.Write(buf[:n])
		if err != nil {
			return err
		}
		_, err = fd.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	tempfile.CRC32 = hasher.Sum32()
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

	return nil
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
	reader := ctx.BodyReader()

	var bodyReader io.Reader = nil
	if strings.HasPrefix(content_type, "multipart/") {
		boundary := params["boundary"]
		mr := multipart.NewReader(reader, boundary)
		for {
			part, err := mr.NextPart()
			if err != nil {
				return err
			}
			if part.FormName() == "file" {
				bodyReader = part
				tempfile.Name = part.FileName()
				break
			} else {
				Closer(part)
			}
		}
	} else if strings.HasPrefix(content_type, "application/octet-stream") {
		bodyReader = reader
		tempfile.Name = ctx.Header("Filename")
	}

	if bodyReader == nil {
		return huma.Error422UnprocessableEntity("content type: " + content_type)
	}

	return CreateTempFile(ctx.Context(), tempfile, bodyReader)
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

func UploadKaraFile(ctx context.Context, input *UploadInput) (*UploadOutput, error) {
	db := GetDB(ctx)
	var err error
	defer func() {
		err := os.Remove(input.File.Fd.Name())
		if err != nil {
			getLogger().Println(err)
		}
	}()
	defer Closer(input.File.Fd)

	kid := input.KID
	kara, err := GetKaraByID(db, kid)
	if err != nil {
		return nil, err
	}

	switch input.FileType {
	case "video":
		res := CheckS3Video(ctx, input.File.Fd, input.File.Size)
		if !res.Passed {
			return nil, res.Error()
		}
	case "inst":
		res := CheckS3Inst(ctx, input.File.Fd, input.File.Size)
		if !res.Passed {
			return nil, res.Error()
		}
	case "sub":
		res, err := CheckS3Ass(ctx, input.File.Fd, input.File.Size)
		if err != nil {
			return nil, err
		}
		if !res.Passed {
			return nil, huma.Error422UnprocessableEntity("Uploaded subtitles file cannot be read")
		}
	}

	_, err = input.File.Fd.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	resp := &UploadOutput{}
	err = db.Transaction(func(tx *gorm.DB) error {
		res, err := SaveTempFileToS3(ctx, tx, input.File, &kara, input.FileType)
		if err != nil {
			return err
		}

		resp.Body.CheckResults = *res
		resp.Body.KID = input.KID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, err
}

type DownloadInput struct {
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" enum:"video,sub,inst" example:"video"`
	Range    string `header:"Range"`
}

type FileSender struct {
	// Reader should be already at the Range.Start location
	Fd        io.ReadCloser
	Range     httputils.Range
	BytesRead uint64
}

func (f *FileSender) Read(buf []byte) (int, error) {
	toread := f.Range.Length - f.BytesRead
	if toread < uint64(len(buf)) {
		buf = buf[:toread]
	}
	return f.Fd.Read(buf)
}

func (f *FileSender) Close() error {
	return f.Fd.Close()
}

func serveObject(obj_file string, range_header string, filename string) (*huma.StreamResponse, error) {

	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			fiber_ctx := ctx.BodyWriter().(*fasthttp.RequestCtx)

			obj, err := GetObject(context.Background(), obj_file)
			if err != nil {
				ctx.SetStatus(500)
				return
			}

			defer func() {
				r := recover()
				if r != nil {
					// unlikely, but close object on panic just in case it happens
					err = obj.Close()
					if err != nil {
						panic(err)
					}
					panic(r)
				}
			}()

			stat, err := obj.Stat()

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

			ctx.SetHeader("Accept-Range", "bytes")

			var reqRange httputils.Range
			if range_header == "" {
				reqRange = httputils.Range{Start: 0, Length: uint64(stat.Size)}
			} else {
				ranges, err := httputils.ParseRequestRange(range_header, uint64(stat.Size))
				if err != nil {
					ctx.SetStatus(416)
					ctx.SetHeader("Content-Range", fmt.Sprintf("bytes */%d", stat.Size))
					getLogger().Println(err.Error())
					return
				}
				reqRange = ranges[0]
				ctx.SetStatus(206)
				ctx.SetHeader("Content-Range", reqRange.ContentRange(uint64(stat.Size)))
			}

			ctx.SetHeader("Content-Type", "application/octet-stream")
			ctx.SetHeader("Content-Length", strconv.FormatUint(reqRange.Length, 10))
			encoded_filename := url.PathEscape(filename)
			ctx.SetHeader(
				"Content-Disposition",
				fmt.Sprintf("attachment; filename*=UTF-8''%s", encoded_filename),
			)

			_, err = obj.Seek(int64(reqRange.Start), 0)
			if err != nil {
				getLogger().Println(err.Error())
				return
			}

			filesender := FileSender{obj, reqRange, 0}

			fiber_ctx.SetBodyStream(&filesender, int(reqRange.Length))
		},
	}, nil
}

type DownloadHeadOutput struct {
	AcceptRange        string `header:"Accept-Range"`
	ContentLength      int64  `header:"Content-Length"`
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
}

func ContentDispositionHeaderContent(kara KaraInfoDB, filetype string) string {
	encoded_filename := url.PathEscape(kara.FriendlyName()) + FileTypeExtension(filetype)
	return fmt.Sprintf("attachment; filename*=UTF-8''%s", encoded_filename)
}

func FileTypeExtension(filetype string) string {
	switch filetype {
	case "video":
		return ".mkv"
	case "inst":
		return ".mka"
	case "sub":
		return ".ass"
	}

	return ""
}

func DownloadHead(ctx context.Context, input *DownloadInput) (*DownloadHeadOutput, error) {
	db := GetDB(ctx)
	kara, err := GetKaraByID(db, input.KID)
	if err != nil {
		return nil, err
	}

	obj, err := GetKaraObject(ctx, kara, input.FileType)
	if err != nil {
		return nil, err
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	return &DownloadHeadOutput{
		AcceptRange:        "bytes",
		ContentLength:      stat.Size,
		ContentType:        "application/octet-stream",
		ContentDisposition: ContentDispositionHeaderContent(kara, input.FileType),
	}, nil
}

func DownloadFile(ctx context.Context, input *DownloadInput) (*huma.StreamResponse, error) {
	db := GetDB(ctx)
	kid := input.KID

	kara, err := GetKaraByID(db, kid)
	if err != nil {
		return nil, err
	}

	if kara.Private && getCurrentUser(ctx) == nil {
		// return forbidden response for private karas for external users
		return nil, huma.Error403Forbidden("private kara")
	}

	obj, err := getKaraObjectFilename(kara, input.FileType)
	if err != nil {
		return nil, err
	}

	return serveObject(obj, input.Range, kara.FriendlyName()+FileTypeExtension(input.FileType))
}

type DeleteInput struct {
	KID      uint   `path:"id" example:"1"`
	FileType string `path:"filetype" enum:"video,sub,inst" example:"video"`
}

type DeleteOutput struct {
	Body struct {
		Deleted string `json:"deleted"`
	}
}

func DeleteKaraFile(ctx context.Context, input *DownloadInput) (*DeleteOutput, error) {
	db := GetDB(ctx)
	kid := input.KID

	kara, err := GetKaraByID(db, kid)
	if err != nil {
		return nil, err
	}

	obj, err := getKaraObjectFilename(kara, input.FileType)
	if err != nil {
		return nil, err
	}

	err = deleteFile(ctx, obj)
	if err != nil {
		return nil, err
	}

	switch input.FileType {
	case "video":
		err = db.Model(&kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{VideoUploaded: false}}).Error
	case "inst":
		err = db.Model(&kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{InstrumentalUploaded: false}}).Error
	case "sub":
		err = db.Model(&kara).Updates(&KaraInfoDB{UploadInfo: UploadInfo{InstrumentalUploaded: false}}).Error
	}
	if err != nil {
		return nil, err
	}

	out := DeleteOutput{}
	out.Body.Deleted = "file deleted"
	return &out, nil
}
