package server

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type UploadFontInput struct {
	RawBody huma.MultipartFormFiles[UploadData]
}

type UploadFontOutput struct {
	Body struct {
		Font Font `json:"font"`
	}
}

func createFont(ctx context.Context, name string) (Font, error) {
	db := GetDB(ctx)
	font := Font{Name: name}

	err := db.Transaction(
		func(tx *gorm.DB) error {
			err := tx.Create(&font).Error
			return DBErrToHumaErr(err)
		})

	return font, err
}

func UploadFont(ctx context.Context, input *UploadFontInput) (*UploadFontOutput, error) {
	out := &UploadFontOutput{}

	file := input.RawBody.Form.File["file"][0]

	fd, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	font, err := createFont(ctx, file.Filename)
	if err != nil {
		return nil, err
	}

	err = SaveFontToS3(ctx, fd, font.ID, file.Size)
	if err != nil {
		return nil, err
	}

	font.UpdatedAt = time.Now()
	err = GetDB(ctx).Save(&font).Error
	out.Body.Font = font

	return out, err
}

type DownloadFontInput struct {
	ID uint `path:"id" example:"1"`
}

func DownloadFont(ctx context.Context, input *DownloadFontInput) (*huma.StreamResponse, error) {
	db := GetDB(ctx)
	id := input.ID

	font := &Font{}
	err := db.First(font, id).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}

	obj, err := GetFontObject(ctx, id)
	if err != nil {
		return nil, err
	}

	return serveObject(obj)
}
