package server

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func getCurrentUser(ctx context.Context) User {
	return ctx.Value(currentUserCtxKey).(User)
}

func getOrCreateUser(ctx context.Context, sub string) (*User, error) {
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
	return &user, nil
}
