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

func TestAuthorTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/tags/author",
			map[string]any{
				"name":             "author_name",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := AuthorOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	path := fmt.Sprintf("/tags/author/%d", data.Body.author.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestFindAuthorTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/tags/author",
			map[string]any{
				"name":             "author_name",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := AuthorOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	resp = assertRespCode(t, api.Get("/tags/author?name=author_name"), 200)

	path := fmt.Sprintf("/tags/author/%d", data.Body.author.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestArtistTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/tags/artist",
			map[string]any{
				"name":             "artist_name",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := ArtistOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	path := fmt.Sprintf("/tags/artist/%d", data.Body.artist.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestFindArtistTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/tags/artist",
			map[string]any{
				"name":             "artist_name",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := ArtistOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	resp = assertRespCode(t, api.Get("/tags/artist?name=artist_name"), 200)

	path := fmt.Sprintf("/tags/artist/%d", data.Body.artist.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestMediaTag(t *testing.T) {
	api := getTestAPI(t)

	for _, v := range MediaTypes {
		resp := assertRespCode(t,
			api.Post("/tags/media",
				map[string]any{
					"name":             "media_name",
					"media_type":       v.ID,
					"additional_names": []string{},
				},
			),
			200,
		)

		data := MediaOutput{}
		dec := json.NewDecoder(resp.Body)
		dec.Decode(&data.Body)

		path := fmt.Sprintf("/tags/media/%d", data.Body.Media.ID)
		assertRespCode(t, api.Delete(path), 204)
	}

}

func TestFindMedia(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/tags/media",
			map[string]any{
				"name":             "media_name",
				"media_type":       "ANIME",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := MediaOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	resp = assertRespCode(t, api.Get("/tags/media?name=media_name"), 200)

	path := fmt.Sprintf("/tags/media/%d", data.Body.Media.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestCreateKara(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/kara",
			map[string]any{
				"title":         "kara_title",
				"title_aliases": []string{},
				"authors":       []uint{},
				"artists":       []uint{},
				"source_media":  0,
				"song_order":    0,
				"medias":        []uint{},
				"audio_tags":    []string{},
				"video_tags":    []string{},
				"comment":       "",
				"version":       "",
			}),
		200,
	)

	data := KaraOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	path := fmt.Sprintf("/kara/%d", data.Body.Kara.ID)
	assertRespCode(t, api.Delete(path), 204)
}
