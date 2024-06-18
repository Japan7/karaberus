// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type GetArtistInput struct {
	Id uint `path:"id" example:"1"`
}

type ArtistOutput struct {
	Body struct {
		Artist Artist `json:"artist"`
	}
}

func GetArtistByID(Id uint) Artist {
	db := GetDB()

	artist := Artist{}
	tx := db.First(&artist, Id)
	if tx.Error != nil {
		panic(tx.Error.Error())
	}

	return artist
}

func GetArtist(Ctx context.Context, input *GetArtistInput) (*ArtistOutput, error) {
	db := GetDB()

	artist_output := &ArtistOutput{}
	tx := db.First(&artist_output.Body.Artist, input.Id)
	if tx.Error != nil {
		return nil, huma.Error404NotFound("tag not found", tx.Error)
	}

	return artist_output, nil
}

type CreateArtistInput struct {
	Body struct {
		Name            string   `json:"name"`
		AdditionalNames []string `json:"additional_names"`
	}
}

func createArtist(name string) (*Artist, error) {
	artist := Artist{Name: name}
	tx := GetDB().Create(&artist)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &artist, nil
}

func CreateArtist(Ctx context.Context, input *CreateArtistInput) (*ArtistOutput, error) {
	artist_output := &ArtistOutput{}
	artist, err := createArtist(input.Body.Name)
	if err != nil {
		return nil, err
	}
	artist_output.Body.Artist = *artist

	return artist_output, nil
}

type DeleteArtistResponse struct {
	Status int
}

func DeleteArtist(Ctx context.Context, input *GetArtistInput) (*DeleteArtistResponse, error) {
	tx := GetDB().Delete(&Artist{}, input.Id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	return &DeleteArtistResponse{204}, nil
}

type FindArtistInput struct {
	Name string `query:"name"`
}

func FindArtist(Ctx context.Context, input *FindArtistInput) (*ArtistOutput, error) {
	artist := Artist{}
	tx := GetDB().Where(&Artist{Name: input.Name}).First(&artist)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	out := &ArtistOutput{}
	out.Body.Artist = artist

	return out, nil
}
