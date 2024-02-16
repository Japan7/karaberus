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
	Types       []string `json:"types" example:"[\"OP\"]"`
	Comment     string   `json:"comment" example:"From https://youtu.be/dQw4w9WgXcQ"`
	SongOrder   int      `json:"song_order" example:"0"`
}

func (info KaraInfo) count_tags() int {
	tags := 0
	tags += len(info.Authors)
	tags += len(info.Artists)

	return tags
}

func makeTags(info KaraInfo) []Tag {
	tags := make([]Tag, info.count_tags())
	tag_i := 0

	for _, author_name := range info.Authors {
		tags[tag_i] = getAuthor(author_name)
		tag_i++
	}

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

	kara_types := make([]KaraType, len(info.Types))
	kt_i := 0
	for _, kara_type := range info.Types {
		kara_types[kt_i] = getKaraType(kara_type)
		kt_i++
	}

	return tags
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
	return KaraInfoDB{
		Tags:        makeTags(info),
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
