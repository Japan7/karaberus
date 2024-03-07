// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import "context"

type GetAuthorInput struct {
	Id uint `path:"id" example:"1"`
}

type AuthorOutput struct {
	Body struct {
		Author Tag
	}
}

func GetAuthor(Ctx context.Context, input *GetAuthorInput) (*AuthorOutput, error) {
	db := GetDB()

	author_output := &AuthorOutput{}
	tx := db.First(&author_output.Body.Author, input.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return author_output, nil
}

type CreateAuthorInput struct {
	Body struct {
		Name            string
		AdditionalNames []string
	}
}

func CreateAuthor(Ctx context.Context, input *CreateAuthorInput) (*AuthorOutput, error) {
	author_output := &AuthorOutput{}
	author_output.Body.Author = getAuthor(input.Body.Name)

	return author_output, nil
}
