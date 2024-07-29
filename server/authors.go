// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"

	"gorm.io/gorm"
)

type GetAuthorInput struct {
	Id uint `path:"id" example:"1"`
}

type AuthorOutput struct {
	Body struct {
		Author TimingAuthor `json:"author"`
	}
}

func GetAuthorById(tx *gorm.DB, Id uint) (*TimingAuthor, error) {
	author := &TimingAuthor{}
	err := tx.First(author, Id).Error
	return author, DBErrToHumaErr(err)
}

func GetAuthor(ctx context.Context, input *GetAuthorInput) (*AuthorOutput, error) {
	tx := GetDB(ctx)

	author_output := &AuthorOutput{}
	author, err := GetAuthorById(tx, input.Id)
	if err != nil {
		return nil, err
	}
	author_output.Body.Author = *author
	return author_output, nil
}

type CreateAuthorInput struct {
	Body struct {
		Name string `json:"name"`
	}
}

func CreateAuthor(ctx context.Context, input *CreateAuthorInput) (*AuthorOutput, error) {
	db := GetDB(ctx)
	author_output := &AuthorOutput{}

	err := db.Transaction(
		func(tx *gorm.DB) error {
			author_output.Body.Author = TimingAuthor{Name: input.Body.Name}
			err := tx.Create(&author_output.Body.Author).Error
			return DBErrToHumaErr(err)
		})

	return author_output, err
}

type DeleteAuthorResponse struct {
	Status int
}

func DeleteAuthor(ctx context.Context, input *GetArtistInput) (*DeleteAuthorResponse, error) {
	db := GetDB(ctx)
	err := db.Delete(&TimingAuthor{}, input.Id).Error
	return &DeleteAuthorResponse{204}, DBErrToHumaErr(err)
}

type FindAuthorInput struct {
	Name string `query:"name"`
}

func FindAuthor(ctx context.Context, input *FindAuthorInput) (*AuthorOutput, error) {
	db := GetDB(ctx)
	out := &AuthorOutput{}
	err := db.Where(&TimingAuthor{Name: input.Name}).First(&out.Body.Author).Error
	return out, DBErrToHumaErr(err)
}

type AllAuthorOutput struct {
	Body []TimingAuthor `json:"authors"`
}

func GetAllAuthors(ctx context.Context, input *struct{}) (*AllAuthorOutput, error) {
	db := GetDB(ctx)
	out := &AllAuthorOutput{}
	err := db.Find(&out.Body).Error
	return out, DBErrToHumaErr(err)
}
