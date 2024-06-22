package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type CreateTokenInput struct {
	Body struct {
		Scopes
		ReadOnly bool `json:"read_only"`
		Admin    bool `json:"admin"`
	}
}

type CreateTokenOutput struct {
	Body struct {
		Token string `json:"token"`
	}
}

func generateToken() (*string, error) {
	token_bytes := make([]byte, 64)
	_, err := rand.Reader.Read(token_bytes)
	if err != nil {
		return nil, err
	}
	token_str := hex.EncodeToString(token_bytes)
	return &token_str, nil
}

func CreateSystemToken() (string, error) {
	token_id, err := generateToken()
	if err != nil {
		return "", err
	}

	token := Token{
		ID:       *token_id,
		Admin:    true,
		ReadOnly: false,
		Scopes:   AllScopes,
	}
	GetDB().Create(&token)

	return *token_id, nil
}

func CreateToken(ctx context.Context, input *CreateTokenInput) (*CreateTokenOutput, error) {
	token_id, err := generateToken()
	if err != nil {
		return nil, err
	}

	token := Token{
		ID:       *token_id,
		Admin:    input.Body.Admin,
		ReadOnly: input.Body.User,
		Scopes:   input.Body.Scopes,
		User:     getCurrentUser(ctx),
	}
	tx := GetDB().Create(&token)
	if tx.Error != nil {
		return nil, tx.Error
	}

	out := &CreateTokenOutput{}
	out.Body.Token = *token_id

	return out, nil
}

type DeleteTokenInput struct {
	TokenID string `path:"token"`
}

type DeleteTokenOutput struct {
	Body struct {
		Message string `json:"message" example:"Token 123 deleted."`
	}
}

func DeleteToken(ctx context.Context, input *DeleteTokenInput) (*DeleteTokenOutput, error) {
	user := getCurrentUser(ctx)
	var err error
	if user.Admin {
		tx := GetDB().Delete(&Token{}, input.TokenID)
		err = tx.Error
	} else {
		tx := GetDB().Where(&Token{ID: input.TokenID, User: user}).Delete(&Token{})
		err = tx.Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("Token not found")
		}
	}

	out := &DeleteTokenOutput{}
	out.Body.Message = fmt.Sprintf("Token %s deleted.", input.TokenID)

	return out, nil
}
