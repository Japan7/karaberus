// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"

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

func GetArtistByID(tx *gorm.DB, Id uint) (*Artist, error) {
	artist := &Artist{}
	err := tx.First(artist, Id).Error
	return artist, DBErrToHumaErr(err)
}

func GetArtist(ctx context.Context, input *GetArtistInput) (*ArtistOutput, error) {
	tx := GetDB(ctx)

	artist_output := &ArtistOutput{}
	artist, err := GetArtistByID(tx, input.Id)
	if err != nil {
		return nil, err
	}
	artist_output.Body.Artist = *artist
	return artist_output, nil
}

type CreateArtistInput struct {
	Body struct {
		Name            string   `json:"name"`
		AdditionalNames []string `json:"additional_names"`
	}
}

func createArtist(tx *gorm.DB, name string, additional_names []string, artist *Artist) error {
	return tx.Transaction(
		func(tx *gorm.DB) error {
			additional_names_db := createAdditionalNames(additional_names)
			artist.Name = name
			artist.AdditionalNames = additional_names_db
			err := tx.Create(artist).Error
			return DBErrToHumaErr(err)
		})
}

func CreateArtist(ctx context.Context, input *CreateArtistInput) (*ArtistOutput, error) {
	artist_output := &ArtistOutput{}

	db := GetDB(ctx)
	err := createArtist(db, input.Body.Name, input.Body.AdditionalNames, &artist_output.Body.Artist)
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}

	return artist_output, nil
}

type DeleteArtistResponse struct {
	Status int
}

func DeleteArtist(ctx context.Context, input *GetArtistInput) (*DeleteArtistResponse, error) {
	tx := GetDB(ctx)
	err := tx.Delete(&Artist{}, input.Id).Error
	return &DeleteArtistResponse{204}, DBErrToHumaErr(err)
}

type FindArtistInput struct {
	Name string `query:"name"`
}

func findArtist(tx *gorm.DB, names []string, artist *Artist) error {
	err := tx.Where("Name in ?", names).First(&artist).Error
	return err
}

func FindArtist(ctx context.Context, input *FindArtistInput) (*ArtistOutput, error) {
	out := &ArtistOutput{}
	err := findArtist(GetDB(ctx), []string{input.Name}, &out.Body.Artist)
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}

	return out, nil
}

type AllArtistsOutput struct {
	Body []Artist `json:"artists"`
}

func GetAllArtists(ctx context.Context, input *struct{}) (*AllArtistsOutput, error) {
	db := GetDB(ctx)
	out := &AllArtistsOutput{}
	err := db.Preload("AdditionalNames").Find(&out.Body).Error
	return out, DBErrToHumaErr(err)
}
