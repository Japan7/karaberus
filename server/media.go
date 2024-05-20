package server

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type CreateMediaInput struct {
	Body struct {
		Name            string   `json:"name" example:"Shinseiki Evangelion"`
		MediaType       string   `json:"media_type" example:"anime"`
		AdditionalNames []string `json:"additional_names" example:"[]"`
	}
}

type MediaOutput struct {
	Body struct {
		Media MediaDB `json:"media"`
	}
}

func createMedia(name string, media_type MediaType) (*MediaDB, error) {
	media_tag := MediaDB{Name: name, Type: media_type.Value}
	tx := GetDB().Create(&media_tag)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &media_tag, nil
}

func getMediaType(media_type_name string) MediaType {
	for _, v := range MediaTypes {
		if v.Type == media_type_name {
			return v
		}
	}

	// TODO: make huma check the input
	panic("unknown media type " + media_type_name)
}

func getMedia(name string, media_type_str string) MediaDB {
	media_type := getMediaType(media_type_str)
	media := MediaDB{}
	tx := GetDB().Where(&MediaDB{Name: name, Type: media_type.Value}).FirstOrCreate(&media)
	if tx.Error != nil {
		panic(tx.Error.Error())
	}

	return media
}

func CreateMedia(Ctx context.Context, input *CreateMediaInput) (*MediaOutput, error) {
	media_output := &MediaOutput{}
	tag_type := getMediaType(input.Body.MediaType)
	media, err := createMedia(input.Body.Name, tag_type)
	if err != nil {
		return nil, err
	}
	media_output.Body.Media = *media

	return media_output, nil
}

type DeleteMediaResponse struct {
	Status int
}

type GetMediaInput struct {
	Id uint `path:"id" example:"1"`
}

func DeleteMedia(Ctx context.Context, input *GetMediaInput) (*DeleteMediaResponse, error) {
	tx := GetDB().Delete(&MediaDB{}, input.Id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	return &DeleteMediaResponse{204}, nil
}

func GetMedia(Ctx context.Context, input *GetMediaInput) (*MediaOutput, error) {
	db := GetDB()

	media_output := &MediaOutput{}
	tx := db.First(&media_output.Body.Media, input.Id)
	if tx.Error != nil {
		return nil, huma.Error404NotFound("tag not found", tx.Error)
	}

	return media_output, nil
}

type FindMediaInput struct {
	Name string `query:"name"`
}

func FindMedia(Ctx context.Context, input *FindMediaInput) (*MediaOutput, error) {
	media := MediaDB{}
	tx := GetDB().Where(&MediaDB{Name: input.Name}).First(&media)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	out := &MediaOutput{}
	out.Body.Media = media

	return out, nil
}
