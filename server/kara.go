// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Media struct {
	Name      string `json:"name" example:"Shinseiki Evangelion"`
	MediaType string `json:"media_type" example:"Anime"`
}

type KaraInfo struct {
	// Main name of the karaoke
	Title string `json:"title" example:"Zankoku na Tenshi no These"`
	// More names relating to this karaoke
	ExtraTitles []string `json:"title_aliases" example:"[\"A Cruel Angel's Thesis\"]"`
	// Karaoke authors
	Authors []uint `json:"authors" example:"[1]"`
	// Artists of the original song
	Artists []uint `json:"artists" example:"[1]"`
	// Name of the Media
	SourceMedia uint `json:"source_media" example:"1"`
	// Number of the track related to the media.
	SongOrder uint `json:"song_order" example:"0"`
	// Medias related to the karaoke
	Medias []uint `json:"medias"`
	// Audio tags
	AudioTags []string `json:"audio_tags" example:"[\"Opening\"]"`
	// Video tags
	VideoTags []string `json:"video_tags" example:"[\"Opening\"]"`
	// Generic comment
	Comment string `json:"comment" example:"From https://youtu.be/dQw4w9WgXcQ"`
	// Version (8-bit, Episode 12, ...)
	Version string `json:"version" example:"iykyk"`
	// Language (FR, EN, ...)
	Language            string `json:"language" example:"FR"`
	KaraokeCreationDate *int64 `json:"karaoke_creation_time,omitempty" example:"42"`
	IsHardsub           *bool  `json:"is_hardsub,omitempty" example:"false"`
	Private             bool   `json:"private,omitempty" example:"false"`
}

type AllTags struct {
	Authors []TimingAuthor
	Artists []Artist
	Video   []VideoTagDB
	Audio   []AudioTagDB
	Media   []MediaDB
}

func makeTags(tx *gorm.DB, info KaraInfo) (AllTags, error) {
	tags := AllTags{}
	tags.Authors = make([]TimingAuthor, len(info.Authors))

	for i, author := range info.Authors {
		author, err := GetAuthorById(tx, author)
		if err != nil {
			return tags, err
		}
		tags.Authors[i] = *author
	}

	tags.Artists = make([]Artist, len(info.Artists))

	for i, artist := range info.Artists {
		artist, err := GetArtistByID(tx, artist)
		if err != nil {
			return tags, err
		}
		tags.Artists[i] = *artist
	}

	tags.Media = make([]MediaDB, len(info.Medias))
	for i, media := range info.Medias {
		media, err := getMediaByID(tx, media)
		if err != nil {
			return tags, err
		}
		tags.Media[i] = media
	}

	tags.Video = make([]VideoTagDB, len(info.VideoTags))
	for i, video_type := range info.VideoTags {
		video_tag, err := getVideoTag(video_type)
		if err != nil {
			return tags, err
		}
		tags.Video[i] = VideoTagDB{video_tag.ID}
	}

	tags.Audio = make([]AudioTagDB, len(info.AudioTags))
	for i, audio_type := range info.AudioTags {
		audio_tag, err := getAudioTag(audio_type)
		if err != nil {
			return tags, err
		}
		tags.Audio[i] = AudioTagDB{audio_tag.ID}
	}

	return tags, nil
}

func makeExtraTitles(info KaraInfo) []AdditionalName {
	extra_titles := make([]AdditionalName, len(info.ExtraTitles))

	for i, title := range info.ExtraTitles {
		extra_titles[i] = AdditionalName{Name: title}
	}

	return extra_titles
}

func (info KaraInfo) to_KaraInfoDB(ctx context.Context, tx *gorm.DB, kara_info *KaraInfoDB) error {
	tags, err := makeTags(tx, info)
	if err != nil {
		return err
	}

	kara_info.VideoTags = tags.Video
	kara_info.AudioTags = tags.Audio
	kara_info.Authors = tags.Authors
	kara_info.Artists = tags.Artists
	kara_info.Medias = tags.Media
	kara_info.Title = info.Title
	kara_info.ExtraTitles = makeExtraTitles(info)
	kara_info.Comment = info.Comment
	kara_info.Version = info.Version
	kara_info.SongOrder = info.SongOrder
	kara_info.Language = info.Language
	kara_info.Private = info.Private

	user, err := getCurrentUser(ctx)
	if err != nil {
		return err
	}
	if user.Admin {
		if info.IsHardsub != nil {
			kara_info.Hardsubbed = *info.IsHardsub
		}
		if info.KaraokeCreationDate != nil {
			kara_info.KaraokeCreationTime = time.Unix(*info.KaraokeCreationDate, 0)
		}
	}

	if info.SourceMedia > 0 {
		source_media, err := getMediaByID(tx, info.SourceMedia)
		if err != nil {
			return err
		}
		kara_info.SourceMedia = &source_media
	} else {
		kara_info.SourceMediaID = nil
	}

	return err
}

type CreateKaraInput struct {
	Body KaraInfo
}

type KaraOutput struct {
	Body struct {
		Kara KaraInfoDB `json:"kara"`
	}
}

func CreateKara(ctx context.Context, input *CreateKaraInput) (*KaraOutput, error) {
	db := GetDB(ctx)
	output := KaraOutput{}

	err := db.Transaction(func(tx *gorm.DB) error {
		kara := KaraInfoDB{}
		err := input.Body.to_KaraInfoDB(ctx, tx, &kara)
		if err != nil {
			return err
		}
		output.Body.Kara = kara

		err = tx.Create(&output.Body.Kara).Error
		return err
	})

	return &output, err
}

type SetKaraUploadTimeInput struct {
	Id   uint `path:"id"`
	Body struct {
		CreationDate int64 `json:"creation_time" example:"42"`
	}
}

type UpdateKaraInput struct {
	Id   uint `path:"id"`
	Body KaraInfo
}

func updateKara(tx *gorm.DB, kara *KaraInfoDB) error {
	err := tx.Model(&kara).Select("*").Updates(&kara).Error
	if err != nil {
		return err
	}
	prev_context := tx.Statement.Context
	tx = WithAssociationsUpdate(tx)
	defer tx.WithContext(prev_context)
	err = tx.Model(&kara).Association("AudioTags").Replace(&kara.AudioTags)
	if err != nil {
		return err
	}
	err = tx.Model(&kara).Association("VideoTags").Replace(&kara.VideoTags)
	if err != nil {
		return err
	}
	err = tx.Model(&kara).Association("Medias").Replace(&kara.Medias)
	if err != nil {
		return err
	}
	err = tx.Model(&kara).Association("Authors").Replace(&kara.Authors)
	if err != nil {
		return err
	}
	err = tx.Model(&kara).Association("ExtraTitles").Replace(&kara.ExtraTitles)
	if err != nil {
		return err
	}
	return nil
}

func UpdateKara(ctx context.Context, input *UpdateKaraInput) (*KaraOutput, error) {
	db := GetDB(ctx)
	kara := KaraInfoDB{}
	err := db.Scopes(CurrentKaras).First(&kara, input.Id).Error
	if err != nil {
		return nil, err
	}
	err = input.Body.to_KaraInfoDB(ctx, db, &kara)
	if err != nil {
		return nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return updateKara(tx, &kara)
	})
	if err != nil {
		return nil, err
	}

	out := &KaraOutput{}
	out.Body.Kara = kara

	return out, nil
}

type GetKaraInput struct {
	Id uint `path:"id"`
}

func GetKara(ctx context.Context, input *GetKaraInput) (*KaraOutput, error) {
	db := GetDB(ctx)

	kara_output := &KaraOutput{}
	err := db.Scopes(CurrentKaras).Preload(clause.Associations).First(&kara_output.Body.Kara, input.Id).Error
	return kara_output, DBErrToHumaErr(err)
}

type DeleteKaraResponse struct {
	Status int
}

func DeleteKara(ctx context.Context, input *GetKaraInput) (*DeleteKaraResponse, error) {
	db := GetDB(ctx)
	err := db.Scopes(CurrentKaras).Delete(&KaraInfoDB{}, input.Id).Error
	return &DeleteKaraResponse{204}, DBErrToHumaErr(err)
}

func GetKaraByID(db *gorm.DB, kara_id uint) (KaraInfoDB, error) {
	kara := &KaraInfoDB{}
	err := db.Preload(clause.Associations).First(kara, kara_id).Error
	return *kara, err
}

type GetAllKarasInput struct {
	IfNoneMatch string `header:"If-None-Match"`
}

type GetAllKarasBody struct {
}

type GetAllKarasOutput struct {
	ETag   string `header:"ETag"`
	Status int
	Body   struct {
		Karas []KaraInfoDB
	}
}

func getKaraETag(tx *gorm.DB, etag_header *string) error {
	etag := uint(0)

	last_kara := KaraInfoDB{}
	err := tx.Last(&last_kara).Error
	err = extendETag(last_kara.ID, err, &etag)
	if err != nil {
		return err
	}

	// need to check artists since they are included in the response
	last_artist := Artist{}
	err = tx.Last(&last_artist).Error
	err = extendETag(last_artist.ID, err, &etag)
	if err != nil {
		return err
	}

	// need to check medias since they are included in the response
	last_media := Artist{}
	err = tx.Last(&last_media).Error
	err = extendETag(last_media.ID, err, &etag)
	if err != nil {
		return err
	}

	// need to check authors since they are included in the response
	// they are not versioned so it uses the last update time
	last_author := TimingAuthor{}
	err = tx.Order("updated_at DESC").Take(&last_author).Error
	err = extendETag(uint(last_author.UpdatedAt.Unix()), err, &etag)
	if err != nil {
		return err
	}

	*etag_header = fmt.Sprint(etag)
	return nil
}

func GetAllKaras(ctx context.Context, input *GetAllKarasInput) (*GetAllKarasOutput, error) {
	out := &GetAllKarasOutput{}
	db := GetDB(ctx)

	err := getKaraETag(db, &out.ETag)
	if err != nil {
		return nil, err
	}

	if out.ETag == input.IfNoneMatch {
		out.Status = 304
	} else {
		out.Status = 200
		err = db.Preload(clause.Associations).Scopes(CurrentKaras).Find(&out.Body.Karas).Error
	}
	return out, DBErrToHumaErr(err)
}

type GetKaraHistoryOutput struct {
	Body struct {
		History []KaraInfoDB `json:"history"`
	}
}

func GetKaraHistory(ctx context.Context, input *GetKaraInput) (*GetKaraHistoryOutput, error) {
	out := &GetKaraHistoryOutput{}
	db := GetDB(ctx)
	err := db.Preload(clause.Associations).Where(&KaraInfoDB{CurrentKaraInfoID: &input.Id}).Find(&out.Body.History).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}
