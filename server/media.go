package server

import (
	"context"

	"gorm.io/gorm"
)

type CreateMediaInput struct {
	Body struct {
		Name            string   `json:"name" example:"Shinseiki Evangelion"`
		MediaType       string   `json:"media_type" example:"ANIME"`
		AdditionalNames []string `json:"additional_names" example:"[]"`
	}
}

type MediaOutput struct {
	Body struct {
		Media MediaDB `json:"media"`
	}
}

func getMediaType(media_type_id string) MediaType {
	for _, v := range MediaTypes {
		if v.ID == media_type_id {
			return v
		}
	}

	// TODO: make huma check the input
	panic("unknown media type " + media_type_id)
}

func getMediaByID(tx *gorm.DB, Id uint) (MediaDB, error) {
	media := MediaDB{}
	err := tx.First(&media, Id).Error
	return media, DBErrToHumaErr(err)
}

func getMedia(tx *gorm.DB, name string, media_type_str string) (MediaDB, error) {
	media_type := getMediaType(media_type_str)
	media := MediaDB{}
	err := tx.Where(&MediaDB{Name: name, Type: media_type.ID}).FirstOrCreate(&media).Error
	return media, DBErrToHumaErr(err)
}

func CreateMedia(ctx context.Context, input *CreateMediaInput) (*MediaOutput, error) {
	db := GetDB(ctx)
	media_output := &MediaOutput{}
	media_type := getMediaType(input.Body.MediaType)

	err := db.Transaction(func(tx *gorm.DB) error {
		media_output.Body.Media = MediaDB{Name: input.Body.Name, Type: media_type.ID}

		additional_names := createAdditionalNames(input.Body.AdditionalNames)
		media_output.Body.Media.AdditionalNames = additional_names

		err := tx.Create(&media_output.Body.Media).Error
		return DBErrToHumaErr(err)
	})

	return media_output, err
}

type DeleteMediaResponse struct {
	Status int
}

type GetMediaInput struct {
	Id uint `path:"id" example:"1"`
}

func DeleteMedia(ctx context.Context, input *GetMediaInput) (*DeleteMediaResponse, error) {
	db := GetDB(ctx)
	err := db.Delete(&MediaDB{}, input.Id).Error
	return &DeleteMediaResponse{204}, DBErrToHumaErr(err)
}

func GetMedia(ctx context.Context, input *GetMediaInput) (*MediaOutput, error) {
	db := GetDB(ctx)
	media_output := &MediaOutput{}
	err := db.First(&media_output.Body.Media, input.Id).Error
	return media_output, DBErrToHumaErr(err)
}

type FindMediaInput struct {
	Name string `query:"name"`
}

func FindMedia(ctx context.Context, input *FindMediaInput) (*MediaOutput, error) {
	db := GetDB(ctx)
	out := &MediaOutput{}
	err := db.Where(&MediaDB{Name: input.Name}).First(&out.Body.Media).Error
	return out, DBErrToHumaErr(err)
}
