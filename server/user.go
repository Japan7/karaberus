package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
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

// pass *info to upsert an user
func getOrUpsertUser(ctx context.Context, sub string, info *oidc.UserInfo) (*User, error) {
	db := GetDB(ctx)
	user := User{ID: sub}

	var err error
	if info != nil {
		err = db.FirstOrCreate(&user).Error
	} else {
		err = db.First(&user).Error
	}
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

func GetUserEndpoint(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
	return getUserEndpoint(ctx, &input.ID)
}

func GetMeEndpoint(ctx context.Context, input *struct{}) (*GetUserOutput, error) {
	return getUserEndpoint(ctx, nil)
}

func getUserEndpoint(ctx context.Context, sub *string) (*GetUserOutput, error) {
	out := &GetUserOutput{}
	if sub != nil {
		user, err := getOrUpsertUser(ctx, *sub, nil)
		if err != nil {
			return nil, DBErrToHumaErr(err)
		}
		out.Body = *user
	} else {
		user := getCurrentUser(ctx)
		out.Body = *user
	}
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

func UpdateUserAuthorEndpoint(ctx context.Context, input *UpdateUserAuthorInput) (*UpdateMeAuthorOutput, error) {
	return updateUserAuthorEndpoint(ctx, &input.ID, input.Body.Id)
}

func UpdateMeAuthorEndpoint(ctx context.Context, input *UpdateMeAuthorInput) (*UpdateMeAuthorOutput, error) {
	return updateUserAuthorEndpoint(ctx, nil, input.Body.Id)
}

func updateUserAuthorEndpoint(ctx context.Context, sub *string, authorId *uint) (*UpdateMeAuthorOutput, error) {
	user := getCurrentUser(ctx)
	if sub != nil {
		if !user.Admin {
			return nil, huma.Error403Forbidden("Only admins can update other users")
		}
		maybeUser, err := getOrUpsertUser(ctx, *sub, nil)
		if err != nil {
			return nil, DBErrToHumaErr(err)
		}
		user = maybeUser
	}
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
