package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

func generateToken() (string, error) {
	token_bytes := make([]byte, 64)
	_, err := rand.Reader.Read(token_bytes)
	if err != nil {
		return "", err
	}
	token_str := hex.EncodeToString(token_bytes)
	return token_str, nil
}

func CreateSystemToken() (string, error) {
	token_id, err := generateToken()
	if err != nil {
		return "", err
	}

	token := Token{
		ID:       token_id,
		Admin:    true,
		ReadOnly: false,
		Scopes:   AllScopes,
	}
	err = GetDB(context.TODO()).Create(&token).Error
	return token_id, DBErrToHumaErr(err)
}

func CreateToken(ctx context.Context, input *CreateTokenInput) (*CreateTokenOutput, error) {
	db := GetDB(ctx)

	token_id, err := generateToken()
	if err != nil {
		return nil, err
	}

	out := &CreateTokenOutput{}
	out.Body.Token = token_id

	token := Token{
		ID:       token_id,
		Admin:    input.Body.Admin,
		ReadOnly: input.Body.User,
		Scopes:   input.Body.Scopes,
		User:     getCurrentUser(ctx),
	}

	err = db.Create(&token).Error
	return out, DBErrToHumaErr(err)
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
