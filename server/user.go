package server

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func getCurrentUser(ctx context.Context) *User {
	val := ctx.Value(currentUserCtxKey)
	if val == nil {
		return nil
	} else {
		return val.(*User)
	}
}

func getOrCreateUser(ctx context.Context, sub string, info *oidc.UserInfo) (*User, error) {
	db := GetDB(ctx)
	user := User{ID: sub}
	err := db.FirstOrCreate(&user).Error
	if err != nil {
		return nil, err
	}
	if info != nil {
		user.Admin = false
		var groups = info.Claims[CONFIG.OIDC.GroupsClaim].([]any)
		for _, group := range groups {
			if group == CONFIG.OIDC.AdminGroup {
				user.Admin = true
				break
			}
		}
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

type GetMeOutput struct {
	Body User
}

func GetMe(ctx context.Context, input *struct{}) (*GetMeOutput, error) {
	user := getCurrentUser(ctx)
	return &GetMeOutput{Body: *user}, nil
}

type UpdateMeAuthorInput struct {
	Body struct {
		Id *uint `json:"id"`
	}
}

type UpdateMeAuthorOutput struct {
	Status int
}

func UpdateMeAuthor(ctx context.Context, input *UpdateMeAuthorInput) (*UpdateMeAuthorOutput, error) {
	tx := GetDB(ctx)
	user := getCurrentUser(ctx)
	if input.Body.Id != nil {
		if _, err := GetAuthorById(tx, *input.Body.Id); err != nil {
			return nil, DBErrToHumaErr(err)
		}
	}
	user.TimingProfileID = input.Body.Id
	err := tx.Model(&user).Select("timing_profile_id").Updates(&user).Error
	return &UpdateMeAuthorOutput{Status: 204}, DBErrToHumaErr(err)
}
