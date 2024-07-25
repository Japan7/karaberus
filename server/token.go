package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type CreateTokenOutput struct {
	Body struct {
		Token string `json:"token"`
	}
}

func CreateToken(ctx context.Context, input *struct{}) (*CreateTokenOutput, error) {
	token, err := createTokenForUser(ctx, getCurrentUser(ctx))
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}
	out := &CreateTokenOutput{}
	out.Body.Token = token.ID
	return out, nil
}

func createTokenForUser(ctx context.Context, user User) (*Token, error) {
	token_id, err := generateToken()
	if err != nil {
		return nil, err
	}
	db := GetDB(ctx)
	token := Token{
		ID:        token_id,
		User:      user,
		CreatedAt: time.Now(),
	}
	if err = db.Create(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func generateToken() (string, error) {
	token_bytes := make([]byte, 64)
	_, err := rand.Reader.Read(token_bytes)
	if err != nil {
		return "", err
	}
	token_str := hex.EncodeToString(token_bytes)
	return token_str, nil
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
	db := GetDB(ctx)
	user := getCurrentUser(ctx)

	var err error
	if user.Admin {
		err = db.Delete(&Token{}, input.TokenID).Error
	} else {
		err = db.Where(&Token{ID: input.TokenID, User: user}).Delete(&Token{}).Error
	}
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}

	out := &DeleteTokenOutput{}
	out.Body.Message = fmt.Sprintf("Token %s deleted.", input.TokenID)

	return out, nil
}
