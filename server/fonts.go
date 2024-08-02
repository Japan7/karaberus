package server

import (
	"context"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type GetAllFontsOutput struct {
	Body struct {
		Fonts []Font
	}
}

func GetAllFonts(ctx context.Context, input *struct{}) (*GetAllFontsOutput, error) {
	db := GetDB(ctx)
	out := &GetAllFontsOutput{}
	err := db.Find(&out.Body.Fonts).Error
	return out, DBErrToHumaErr(err)
}

type UploadFontInputDefinition struct {
	RawBody huma.MultipartFormFiles[UploadData]
}

type UploadFontInput struct {
	File UploadTempFile
}

func (i *UploadFontInput) Resolve(ctx huma.Context) []error {
	err := createTempFile(ctx, &i.File)
	if err != nil {
		return []error{err}
	}
	return nil
}

var _ huma.Resolver = (*UploadFontInput)(nil)

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
	defer os.Remove(input.File.Fd.Name())
	defer input.File.Fd.Close()

	font, err := createFont(ctx, input.File.Name)
	if err != nil {
		return nil, err
	}

	err = SaveFontToS3(ctx, input.File.Fd, font.ID, input.File.Size)
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
