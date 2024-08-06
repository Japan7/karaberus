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
	err := tx.Scopes(CurrentMedias).First(&media, Id).Error
	return media, DBErrToHumaErr(err)
}

// func getMedia(tx *gorm.DB, name string, media_type_str string) (MediaDB, error) {
// 	media_type := getMediaType(media_type_str)
// 	media := MediaDB{}
// 	err := tx.Where(&MediaDB{Name: name, Type: media_type.ID}).FirstOrCreate(&media).Error
// 	return media, DBErrToHumaErr(err)
// }

func createMedia(tx *gorm.DB, name string, media_type MediaType, additional_names []string, media *MediaDB) error {
	err := tx.Transaction(func(tx *gorm.DB) error {
		media.Name = name
		media.Type = media_type.ID
		media.AdditionalNames = createAdditionalNames(additional_names)

		err := tx.Create(&media).Error
		return DBErrToHumaErr(err)
	})

	return err
}

func CreateMedia(ctx context.Context, input *CreateMediaInput) (*MediaOutput, error) {
	media_output := &MediaOutput{}
	media_type := getMediaType(input.Body.MediaType)

	db := GetDB(ctx)
	err := createMedia(db, input.Body.Name, media_type, input.Body.AdditionalNames, &media_output.Body.Media)

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
	err := db.Scopes(CurrentMedias).Delete(&MediaDB{}, input.Id).Error
	return &DeleteMediaResponse{204}, DBErrToHumaErr(err)
}

func GetMedia(ctx context.Context, input *GetMediaInput) (*MediaOutput, error) {
	db := GetDB(ctx)
	media_output := &MediaOutput{}
	err := db.Scopes(CurrentMedias).First(&media_output.Body.Media, input.Id).Error
	return media_output, DBErrToHumaErr(err)
}

type FindMediaInput struct {
	Name string `query:"name"`
}

func findMedia(tx *gorm.DB, names []string, media *MediaDB) error {
	err := tx.Scopes(CurrentMedias).Where("Name in ?", names).First(&media).Error
	return err
}

func FindMedia(ctx context.Context, input *FindMediaInput) (*MediaOutput, error) {
	out := &MediaOutput{}
	err := findMedia(GetDB(ctx), []string{input.Name}, &out.Body.Media)
	return out, DBErrToHumaErr(err)
}

type AllMediasOutput struct {
	Body []MediaDB `json:"medias"`
}

func GetAllMedias(ctx context.Context, input *struct{}) (*AllMediasOutput, error) {
	db := GetDB(ctx)
	out := &AllMediasOutput{}
	err := db.Preload("AdditionalNames").Scopes(CurrentMedias).Find(&out.Body).Error
	return out, DBErrToHumaErr(err)
}

type AllMediaTypesOutput struct {
	Body []MediaType `json:"media_types"`
}

func GetAllMediaTypes(ctx context.Context, input *struct{}) (*AllMediaTypesOutput, error) {
	return &AllMediaTypesOutput{MediaTypes}, nil
}
