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
	err := tx.Scopes(CurrentArtists).First(artist, Id).Error
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

type ArtistInfo struct {
	Name            string   `json:"name"`
	AdditionalNames []string `json:"additional_names"`
}

type CreateArtistInput struct {
	Body ArtistInfo
}

func createArtist(db *gorm.DB, artist *Artist, info *ArtistInfo) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := info.to_Artist(artist); err != nil {
			return err
		}
		err := tx.Create(artist).Error
		return DBErrToHumaErr(err)
	})
}

func CreateArtist(ctx context.Context, input *CreateArtistInput) (*ArtistOutput, error) {
	db := GetDB(ctx)
	output := ArtistOutput{}

	err := db.Transaction(func(tx *gorm.DB) error {
		artist := Artist{}
		if err := createArtist(tx, &artist, &input.Body); err != nil {
			return err
		}
		output.Body.Artist = artist
		return nil
	})

	return &output, err
}

type UpdateArtistInput struct {
	Id   uint `path:"id"`
	Body ArtistInfo
}

func updateArtist(tx *gorm.DB, artist *Artist) error {
	err := tx.Model(&artist).Select("*").Updates(&artist).Error
	if err != nil {
		return err
	}
	prev_context := tx.Statement.Context
	tx = WithAssociationsUpdate(tx)
	defer tx.WithContext(prev_context)
	err = tx.Model(&artist).Association("AdditionalNames").Replace(&artist.AdditionalNames)
	if err != nil {
		return err
	}
	return nil
}

func UpdateArtist(ctx context.Context, input *UpdateArtistInput) (*ArtistOutput, error) {
	db := GetDB(ctx)
	artist := Artist{}
	err := db.First(&artist, input.Id).Error
	if err != nil {
		return nil, err
	}
	err = input.Body.to_Artist(&artist)
	if err != nil {
		return nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return updateArtist(tx, &artist)
	})
	if err != nil {
		return nil, err
	}

	out := &ArtistOutput{}
	out.Body.Artist = artist

	return out, nil
}

func (info ArtistInfo) to_Artist(artist *Artist) error {
	artist.Name = info.Name
	artist.AdditionalNames = createAdditionalNames(info.AdditionalNames)
	return nil
}

type DeleteArtistResponse struct {
	Status int
}

func DeleteArtist(ctx context.Context, input *GetArtistInput) (*DeleteArtistResponse, error) {
	tx := GetDB(ctx)
	err := tx.Scopes(CurrentArtists).Delete(&Artist{}, input.Id).Error
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

type AllArtistsInput struct {
	IfNoneMatch string `header:"If-None-Match"`
}

type AllArtistsOutput struct {
	ETag   string `header:"ETag"`
	Status int
	Body   []Artist `json:"artists"`
}

func GetAllArtists(ctx context.Context, input *AllArtistsInput) (*AllArtistsOutput, error) {
	db := GetDB(ctx)
	out := &AllArtistsOutput{}

	last_artist := Artist{}
	err := db.Last(&last_artist).Error
	err = setETag(last_artist.ID, err, &out.ETag)
	if err != nil {
		return nil, err
	}

	if out.ETag == input.IfNoneMatch {
		out.Status = 304
	} else {
		out.Status = 200
		err = db.Preload("AdditionalNames").Scopes(CurrentArtists).Find(&out.Body).Error
	}
	return out, DBErrToHumaErr(err)
}
