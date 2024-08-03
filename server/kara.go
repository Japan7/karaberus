// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
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
	Language string `json:"language" example:"FR"`
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

func (info KaraInfo) to_KaraInfoDB(tx *gorm.DB) (KaraInfoDB, error) {
	tags, err := makeTags(tx, info)
	if err != nil {
		return KaraInfoDB{}, err
	}

	kara_info := KaraInfoDB{
		VideoTags:   tags.Video,
		AudioTags:   tags.Audio,
		Authors:     tags.Authors,
		Artists:     tags.Artists,
		Medias:      tags.Media,
		Title:       info.Title,
		ExtraTitles: makeExtraTitles(info),
		Comment:     info.Comment,
		Version:     info.Version,
		SongOrder:   info.SongOrder,
		UploadInfo:  NewUploadInfo(),
	}

	if info.SourceMedia > 0 {
		kara_info.SourceMedia, err = getMediaByID(tx, info.SourceMedia)
	}

	return kara_info, err
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
		kara, err := input.Body.to_KaraInfoDB(tx)
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

func SetKaraUploadTime(ctx context.Context, input *SetKaraUploadTimeInput) (*KaraOutput, error) {
	db := GetDB(ctx)
	out := &KaraOutput{}
	err := db.Transaction(func(tx *gorm.DB) error {
		kara, err := GetKaraByID(tx, input.Id)
		if err != nil {
			return err
		}

		kara.KaraokeCreationTime = time.Unix(input.Body.CreationDate, 0)
		err = tx.Save(&kara).Error
		if err != nil {
			return err
		}
		out.Body.Kara = kara
		return err
	})
	return out, err
}

type UpdateKaraInput struct {
	Id   uint `path:"id"`
	Body struct {
		KaraInfo
		KaraokeCreationDate int64 `json:"karaoke_creation_time" example:"42"`
		IsHardsub           bool  `json:"is_hardsub" example:"false"`
	}
}

func UpdateKara(ctx context.Context, input *UpdateKaraInput) (*KaraOutput, error) {
	user := getCurrentUser(ctx)

	db := GetDB(ctx)
	current_kara := KaraInfoDB{}
	err := db.First(&current_kara, input.Id).Error
	if err != nil {
		return nil, err
	}
	new_kara, err := input.Body.to_KaraInfoDB(db)
	if err != nil {
		return nil, err
	}

	upload_info := current_kara.UploadInfo

	if user.Admin {
		upload_info.Hardsubbed = input.Body.IsHardsub
		upload_info.KaraokeCreationTime = time.Unix(input.Body.KaraokeCreationDate, 0)
	}

	new_kara.ID = input.Id
	new_kara.UploadInfo = upload_info

	err = db.Save(&new_kara).Error
	if err != nil {
		return nil, err
	}

	out := &KaraOutput{}
	out.Body.Kara = new_kara

	return out, nil
}

type GetKaraInput struct {
	Id uint `path:"id"`
}

func GetKara(ctx context.Context, input *GetKaraInput) (*KaraOutput, error) {
	db := GetDB(ctx)

	kara_output := &KaraOutput{}
	err := db.Preload(clause.Associations).First(&kara_output.Body.Kara, input.Id).Error
	return kara_output, DBErrToHumaErr(err)
}

type DeleteKaraResponse struct {
	Status int
}

func DeleteKara(ctx context.Context, input *GetKaraInput) (*DeleteKaraResponse, error) {
	db := GetDB(ctx)
	err := db.Delete(&TimingAuthor{}, input.Id).Error
	return &DeleteKaraResponse{204}, DBErrToHumaErr(err)
}

func GetKaraByID(db *gorm.DB, kara_id uint) (KaraInfoDB, error) {
	kara := &KaraInfoDB{}
	err := db.Preload(clause.Associations).First(kara, kara_id).Error
	return *kara, err
}

type GetAllKarasOutput struct {
	Body struct {
		Karas []KaraInfoDB
	}
}

func GetAllKaras(ctx context.Context, input *struct{}) (*GetAllKarasOutput, error) {
	out := &GetAllKarasOutput{}
	db := GetDB(ctx)
	err := db.Preload(clause.Associations).Find(&out.Body.Karas).Error
	return out, DBErrToHumaErr(err)
}
