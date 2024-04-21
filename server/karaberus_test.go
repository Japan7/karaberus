// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func getTestAPI(t *testing.T) humatest.TestAPI {
	_, api := humatest.New(t)

	routes(api)
	return api
}

func assertRespCode(t *testing.T, resp *httptest.ResponseRecorder, expected_code int) {
	if resp.Code != expected_code {
		t.Fatal("returned an invalid status code", resp.Code)
	}
}

func TestTags(t *testing.T) {
	api := getTestAPI(t)

	assertRespCode(t, api.Get("/tags/generic/1"), 404)

	assertRespCode(t,
		api.Post("/tags/generic/artist",
			map[string]any{
				"name":             "artist_name",
				"additional_names": []string{},
			},
		),
		200,
	)

	assertRespCode(t, api.Delete("/tags/generic/1"), 204)
}
