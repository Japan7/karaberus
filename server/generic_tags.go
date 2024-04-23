// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type GetTagInput struct {
	Id uint `path:"id" example:"1"`
}

type TagOutput struct {
	Body struct {
		tag Tag
	}
}

func GetTag(Ctx context.Context, input *GetTagInput) (*TagOutput, error) {
	db := GetDB()

	tag_output := &TagOutput{}
	tx := db.First(&tag_output.Body.tag, input.Id)
	if tx.Error != nil {
		return nil, huma.Error404NotFound("tag not found", tx.Error)
	}

	return tag_output, nil
}

type CreateTagInput struct {
	TagTypeName string `path:"tag_type" example:"artist"`
	Body        struct {
		Name            string   `json:"name"`
		AdditionalNames []string `json:"additional_names"`
	}
}

func createTag(name string, tag_type TagType) (*Tag, error) {
	tag := Tag{Name: name, Type: tag_type.Value}
	tx := GetDB().Create(&tag)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &tag, nil
}

func getTagType(TagTypeName string) (*TagType, error) {
	for _, tag := range TagTypes {
		if tag.Type == TagTypeName {
			return &tag, nil
		}
	}
	return nil, &KaraberusError{"unknown tag type:" + TagTypeName}
}

func CreateTag(Ctx context.Context, input *CreateTagInput) (*TagOutput, error) {
	tag_output := &TagOutput{}
	tag_type, err := getTagType(input.TagTypeName)
	if err != nil {
		return nil, err
	}
	tag, err := createTag(input.Body.Name, *tag_type)
	if err != nil {
		return nil, err
	}
	tag_output.Body.tag = *tag

	return tag_output, nil
}

type DeleteTagResponse struct {
	Status int
}

func DeleteTag(Ctx context.Context, input *GetTagInput) (*DeleteTagResponse, error) {
	tx := GetDB().Delete(&Tag{}, input.Id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	return &DeleteTagResponse{204}, nil
}

type FindTagInput struct {
	Name string `query:"name"`
	Type string `query:"type"`
}

func FindTag(Ctx context.Context, input *FindTagInput) (*TagOutput, error) {
	tag_type, err := getTagType(input.Type)
	if err != nil {
		return nil, err
	}
	tag := Tag{}
	tx := GetDB().Where(&Tag{Name: input.Name, Type: tag_type.Value}).First(&tag)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("tag not found")
		}
		return nil, tx.Error
	}

	out := &TagOutput{}
	out.Body.tag = tag

	return out, nil
}
