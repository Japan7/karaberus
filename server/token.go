package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type APIToken struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Scopes    Scopes    `gorm:"embedded" json:"scopes"`
}

type GetAllTokensOutput struct {
	Body []APIToken
}

func GetAllUserTokens(ctx context.Context, input *struct{}) (*GetAllTokensOutput, error) {
	db := GetDB(ctx)
	user := *getCurrentUser(ctx)
	out := &GetAllTokensOutput{}
	err := db.Model(&TokenV2{}).Where(&TokenV2{User: user}).Find(&out.Body).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}
	return out, nil
}

type CreateTokenInput struct {
	Body struct {
		Name   string `json:"name"`
		Scopes Scopes `json:"scopes"`
	}
}

type CreateTokenOutput struct {
	Body struct {
		Token string `json:"token"`
	}
}

func CreateToken(ctx context.Context, input *CreateTokenInput) (*CreateTokenOutput, error) {
	token, err := createTokenForUser(ctx, *getCurrentUser(ctx), input.Body.Name, input.Body.Scopes)
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}
	out := &CreateTokenOutput{}
	out.Body.Token = token.Token
	return out, nil
}

func createTokenForUser(ctx context.Context, user User, name string, scopes Scopes) (*TokenV2, error) {
	token_str, err := generateToken()
	if err != nil {
		return nil, err
	}
	db := GetDB(ctx)
	token := TokenV2{
		Token:  token_str,
		User:   user,
		Name:   name,
		Scopes: scopes,
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
	TokenID uint `path:"token"`
}

type DeleteTokenOutput struct {
	Body struct {
		Message string `json:"message" example:"Token 123 deleted."`
	}
}

func DeleteToken(ctx context.Context, input *DeleteTokenInput) (*DeleteTokenOutput, error) {
	db := GetDB(ctx)
	user := *getCurrentUser(ctx)

	var err error
	if user.Admin {
		err = db.Delete(&TokenV2{}, input.TokenID).Error
	} else {
		err = db.Where(&TokenV2{User: user}).Delete(&TokenV2{}, input.TokenID).Error
	}
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}

	out := &DeleteTokenOutput{}
	out.Body.Message = fmt.Sprintf("Token %d deleted.", input.TokenID)

	return out, nil
}
