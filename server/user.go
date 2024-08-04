package server

import (
	"context"
	"errors"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"gorm.io/gorm"
)

func getCurrentUser(ctx context.Context) User {
	return ctx.Value(currentUserCtxKey).(User)
}

func getOrCreateUser(ctx context.Context, sub string, info *oidc.UserInfo) (*User, error) {
	db := GetDB(ctx)
	user := User{ID: sub}
	if err := db.First(&user).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			// The user doesn't exist yet
			if err = db.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
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
