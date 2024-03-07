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
	Title       string   `json:"title" example:"Zankoku na Tenshi no These"`
	ExtraTitles []string `json:"title_aliases" example:"[\"A Cruel Angel's Thesis\"]"`
	Authors     []string `json:"authors" example:"[\"odrling\"]"`
	Artists     []string `json:"artists" example:"[\"Yoko Takahashi\"]"`
	Medias      []Media  `json:"medias"`
	AudioTags   []string `json:"audio_tags" example:"[\"Opening\"]"`
	VideoTags   []string `json:"video_tags" example:"[\"Opening\"]"`
	Comment     string   `json:"comment" example:"From https://youtu.be/dQw4w9WgXcQ"`
	Version     string   `json:"version" example:"iykyk"`
	SongOrder   int      `json:"song_order" example:"0"`
}

func (info KaraInfo) count_tags() int {
	tags := 0
	tags += len(info.Artists)

	return tags
}

type AllTags struct {
	Generic []Tag
	Authors []TimingAuthor
	Video   []VideoTagDB
	Audio   []AudioTagDB
	Media   []MediaDB
}

func makeTags(info KaraInfo) AllTags {
	authors := make([]TimingAuthor, len(info.Authors))
	auth_i := 0

	for _, author_name := range info.Authors {
		authors[auth_i] = getAuthor(author_name)
		auth_i++
	}

	tags := make([]Tag, info.count_tags())
	tag_i := 0

	for _, artist_name := range info.Artists {
		tags[tag_i] = getArtist(artist_name)
		tag_i++
	}

	medias := make([]MediaDB, len(info.Medias))
	media_i := 0
	for _, media := range info.Medias {
		medias[media_i] = getMedia(media.Name, media.MediaType)
		media_i++
	}

	video_tags := make([]VideoTagDB, len(info.VideoTags))
	vt_i := 0
	for _, video_type := range info.VideoTags {
		video_tags[vt_i] = getVideoTag(video_type)
		vt_i++
	}

	audio_tags := make([]AudioTagDB, len(info.AudioTags))
	at_i := 0
	for _, audio_type := range info.AudioTags {
		audio_tags[at_i] = getAudioTag(audio_type)
		at_i++
	}

	return AllTags{
		Generic: tags,
		Authors: authors,
		Video:   video_tags,
		Audio:   audio_tags,
		Media:   medias,
	}
}

func makeExtraTitles(info KaraInfo) []AdditionalName {
	extra_titles := make([]AdditionalName, len(info.ExtraTitles))
	i := 0

	for _, title := range info.ExtraTitles {
		extra_titles[i] = AdditionalName{Name: title}
		i++
	}

	return extra_titles
}

func (info KaraInfo) to_KaraInfoDB() KaraInfoDB {
	tags := makeTags(info)
	return KaraInfoDB{
		Tags:        tags.Generic,
		VideoTags:   tags.Video,
		AudioTags:   tags.Audio,
		Authors:     tags.Authors,
		Medias:      tags.Media,
		Title:       info.Title,
		ExtraTitles: makeExtraTitles(info),
		Comment:     info.Comment,
		SongOrder:   info.SongOrder,
	}
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
