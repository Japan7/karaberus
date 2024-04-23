// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func getTestAPI(t *testing.T) humatest.TestAPI {
	_, api := humatest.New(t)

	routes(api)
	return api
}

func assertRespCode(t *testing.T, resp *httptest.ResponseRecorder, expected_code int) *httptest.ResponseRecorder {
	if resp.Code != expected_code {
		t.Fatal("returned an invalid status code", resp.Code)
	}
	return resp
}

func TestGenericTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/tags/generic/artist",
			map[string]any{
				"name":             "artist_name",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := TagOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	path := fmt.Sprintf("/tags/generic/%d", data.Body.tag.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestFindGenericTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/tags/generic/artist",
			map[string]any{
				"name":             "artist_name",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := TagOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	resp = assertRespCode(t, api.Get("/tags/generic?name=artist_name&type=artist"), 200)

	path := fmt.Sprintf("/tags/generic/%d", data.Body.tag.ID)
	assertRespCode(t, api.Delete(path), 204)
}
