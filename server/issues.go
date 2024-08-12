// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import "context"

type GetIssuesOutput struct {
	Body struct {
		Issues []KaraIssue `json:"issues"`
	}
}

func GetIssues(ctx context.Context, input *struct{}) (*GetIssuesOutput, error) {
	db := GetDB(ctx)
	out := &GetIssuesOutput{}
	err := db.Find(&out.Body.Issues).Error
	return out, err
}
