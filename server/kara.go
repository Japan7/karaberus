// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
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
}

type AllTags struct {
	Authors []TimingAuthor
	Artists []Artist
	Video   []VideoTagDB
	Audio   []AudioTagDB
	Media   []MediaDB
}

func makeTags(info KaraInfo) AllTags {
	authors := make([]TimingAuthor, len(info.Authors))

	for i, author := range info.Authors {
		authors[i] = GetAuthorById(author)
	}

	artists := make([]Artist, len(info.Artists))

	for i, artist := range info.Artists {
		artists[i] = GetArtistByID(artist)
	}

	medias := make([]MediaDB, len(info.Medias))
	for i, media := range info.Medias {
		medias[i] = getMediaByID(media)
	}

	video_tags := make([]VideoTagDB, len(info.VideoTags))
	for i, video_type := range info.VideoTags {
		video_tags[i] = VideoTagDB{getVideoTag(video_type).ID}
	}

	audio_tags := make([]AudioTagDB, len(info.AudioTags))
	for i, audio_type := range info.AudioTags {
		audio_tags[i] = AudioTagDB{getAudioTag(audio_type).ID}
	}

	return AllTags{
		Authors: authors,
		Artists: artists,
		Video:   video_tags,
		Audio:   audio_tags,
		Media:   medias,
	}
}

func makeExtraTitles(info KaraInfo) []AdditionalName {
	extra_titles := make([]AdditionalName, len(info.ExtraTitles))

	for i, title := range info.ExtraTitles {
		extra_titles[i] = AdditionalName{Name: title}
	}

	return extra_titles
}

func (info KaraInfo) to_KaraInfoDB() KaraInfoDB {
	tags := makeTags(info)

	kara_info := KaraInfoDB{
		VideoTags:   tags.Video,
		AudioTags:   tags.Audio,
		Authors:     tags.Authors,
		Artists:     tags.Artists,
		Medias:      tags.Media,
		Title:       info.Title,
		ExtraTitles: makeExtraTitles(info),
		Comment:     info.Comment,
		SongOrder:   info.SongOrder,
	}

	if info.SourceMedia > 0 {
		kara_info.SourceMedia = getMediaByID(info.SourceMedia)
	}

	return kara_info
}

type CreateKaraInput struct {
	Body KaraInfo
}

type KaraOutput struct {
	Body struct {
		Kara KaraInfoDB
	}
}

func CreateKara(ctx context.Context, input *CreateKaraInput) (*KaraOutput, error) {
	kara := input.Body.to_KaraInfoDB()

	db := GetDB()

	result := db.Create(&kara)
	if result.Error != nil {
		return nil, result.Error
	}

	output := KaraOutput{}
	output.Body.Kara = kara
	return &output, nil
}

type GetKaraInput struct {
	Id uint `path:"id"`
}

func GetKara(Ctx context.Context, input *GetKaraInput) (*KaraOutput, error) {
	db := GetDB()

	kara_output := &KaraOutput{}
	tx := db.First(&kara_output.Body.Kara, input.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return kara_output, nil
}
