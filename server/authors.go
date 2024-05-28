// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type GetAuthorInput struct {
	Id uint `path:"id" example:"1"`
}

type AuthorOutput struct {
	Body struct {
		author TimingAuthor
	}
}

func GetAuthorById(Id uint) TimingAuthor {
	db := GetDB()

	author := TimingAuthor{}
	tx := db.First(&author, Id)
	if tx.Error != nil {
		panic(tx.Error.Error())
	}

	return author
}

func GetAuthor(Ctx context.Context, input *GetAuthorInput) (*AuthorOutput, error) {
	db := GetDB()

	author_output := &AuthorOutput{}
	tx := db.First(&author_output.Body.author, input.Id)
	if tx.Error != nil {
		return nil, huma.Error404NotFound("tag not found", tx.Error)
	}

	return author_output, nil
}

type CreateAuthorInput struct {
	Body struct {
		Name            string   `json:"name"`
		AdditionalNames []string `json:"additional_names"`
	}
}

func createAuthor(name string) (*TimingAuthor, error) {
	author := TimingAuthor{Name: name}
	tx := GetDB().Create(&author)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &author, nil
}

func CreateAuthor(Ctx context.Context, input *CreateAuthorInput) (*AuthorOutput, error) {
	author_output := &AuthorOutput{}
	author, err := createAuthor(input.Body.Name)
	if err != nil {
		return nil, err
	}
	author_output.Body.author = *author

	return author_output, nil
}

type DeleteAuthorResponse struct {
	Status int
}

func DeleteAuthor(Ctx context.Context, input *GetArtistInput) (*DeleteAuthorResponse, error) {
	tx := GetDB().Delete(&TimingAuthor{}, input.Id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	return &DeleteAuthorResponse{204}, nil
}

type FindAuthorInput struct {
	Name string `query:"name"`
}

func FindAuthor(Ctx context.Context, input *FindAuthorInput) (*AuthorOutput, error) {
	author := TimingAuthor{}
	tx := GetDB().Where(&TimingAuthor{Name: input.Name}).First(&author)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	out := &AuthorOutput{}
	out.Body.author = author

	return out, nil
}
