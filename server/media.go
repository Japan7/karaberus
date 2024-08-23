package server

import (
	"context"

	"gorm.io/gorm"
)

type MediaInfo struct {
	Name            string   `json:"name" example:"Shinseiki Evangelion"`
	MediaType       string   `json:"media_type" example:"ANIME"`
	AdditionalNames []string `json:"additional_names" example:"[]"`
}

type CreateMediaInput struct {
	Body MediaInfo
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

func createMedia(db *gorm.DB, media *MediaDB, info *MediaInfo) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := info.to_MediaDB(media); err != nil {
			return err
		}
		err := tx.Create(media).Error
		return DBErrToHumaErr(err)
	})
}

func CreateMedia(ctx context.Context, input *CreateMediaInput) (*MediaOutput, error) {
	db := GetDB(ctx)
	output := MediaOutput{}

	err := db.Transaction(func(tx *gorm.DB) error {
		media := MediaDB{}
		if err := createMedia(tx, &media, &input.Body); err != nil {
			return err
		}
		output.Body.Media = media
		return nil
	})

	return &output, err
}

type UpdateMediaInput struct {
	Id   uint `path:"id"`
	Body MediaInfo
}

func updateMedia(tx *gorm.DB, media *MediaDB) error {
	err := tx.Model(&media).Select("*").Updates(&media).Error
	if err != nil {
		return err
	}
	prev_context := tx.Statement.Context
	tx = WithAssociationsUpdate(tx)
	defer tx.WithContext(prev_context)
	err = tx.Model(&media).Association("AdditionalName").Replace(&media.AdditionalNames)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMedia(ctx context.Context, input *UpdateMediaInput) (*MediaOutput, error) {
	db := GetDB(ctx)
	media := MediaDB{}
	err := db.First(&media, input.Id).Error
	if err != nil {
		return nil, err
	}
	err = input.Body.to_MediaDB(&media)
	if err != nil {
		return nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return updateMedia(tx, &media)
	})
	if err != nil {
		return nil, err
	}

	out := &MediaOutput{}
	out.Body.Media = media

	return out, nil
}

func (info MediaInfo) to_MediaDB(media *MediaDB) error {
	media.Name = info.Name
	media.Type = getMediaType(info.MediaType).ID
	media.AdditionalNames = createAdditionalNames(info.AdditionalNames)
	return nil
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
