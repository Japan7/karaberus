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

type GetUserInput struct {
	ID string `path:"id"`
}

type GetUserOutput struct {
	Body User
}

func GetUser(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
	return getUser(ctx, input.ID)
}

func GetMe(ctx context.Context, input *struct{}) (*GetUserOutput, error) {
	user := getCurrentUser(ctx)
	return &GetUserOutput{Body: *user}, nil
}

func getUser(ctx context.Context, sub string) (*GetUserOutput, error) {
	out := &GetUserOutput{}
	db := GetDB(ctx)
	user := User{}
	err := db.Where(&User{ID: sub}).First(&user).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}
	out.Body = user
	return out, nil
}

type UpdateUserAuthorInput struct {
	ID string `path:"id"`
	UpdateMeAuthorInput
}

type UpdateMeAuthorInput struct {
	Body struct {
		Id *uint `json:"id"`
	}
}

type UpdateMeAuthorOutput struct {
	Status int
}

func UpdateUserAuthor(ctx context.Context, input *UpdateUserAuthorInput) (*UpdateMeAuthorOutput, error) {
	user := &User{}
	err := GetDB(ctx).Where(&User{ID: input.ID}).First(user).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}
	return updateUserAuthor(ctx, user, input.Body.Id)
}

func UpdateMeAuthor(ctx context.Context, input *UpdateMeAuthorInput) (*UpdateMeAuthorOutput, error) {
	user := getCurrentUser(ctx)
	return updateUserAuthor(ctx, user, input.Body.Id)
}

func updateUserAuthor(ctx context.Context, user *User, authorId *uint) (*UpdateMeAuthorOutput, error) {
	tx := GetDB(ctx)
	if authorId != nil {
		if _, err := GetAuthorById(tx, *authorId); err != nil {
			return nil, DBErrToHumaErr(err)
		}
	}
	user.TimingProfileID = authorId
	err := tx.Model(&user).Select("timing_profile_id").Updates(&user).Error
	return &UpdateMeAuthorOutput{Status: 204}, DBErrToHumaErr(err)
}
