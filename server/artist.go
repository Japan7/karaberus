// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import "context"

type GetArtistInput struct {
	Id uint `path:"id" example:"1"`
}

type ArtistOutput struct {
	Body struct {
		Artist Tag
	}
}

func GetArtist(Ctx context.Context, input *GetArtistInput) (*ArtistOutput, error) {
	db := GetDB()

	artist_output := &ArtistOutput{}
	tx := db.First(&artist_output.Body.Artist, input.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return artist_output, nil
}

type CreateArtistInput struct {
	Body struct {
		Name            string
		AdditionalNames []string
	}
}

func getArtist(artist_name string) Tag {
	return getTag(artist_name, KaraTagArtist)
}

func CreateArtist(Ctx context.Context, input *CreateArtistInput) (*ArtistOutput, error) {
	artist_output := &ArtistOutput{}
	artist_output.Body.Artist = getArtist(input.Body.Name)

	return artist_output, nil
}
