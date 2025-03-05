// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

type KaraberusTestConfig struct {
	GeneratedDirectory string `envkey:"DIR_GENERATED"`
	Directory          string `envkey:"DIR"`
}

func getKaraberusTestConfig() KaraberusTestConfig {
	config := KaraberusTestConfig{}

	config_value := reflect.ValueOf(&config).Elem()
	config_type := reflect.TypeOf(config)

	setConfigValue(config_value, config_type, "KARABERUS_TEST_")

	return config
}

var TEST_CONFIG = getKaraberusTestConfig()

func testUserMiddleware(ctx huma.Context, next func(huma.Context)) {
	ctx = huma.WithValue(ctx, currentUserCtxKey, User{
		ID:    "test_user",
		Admin: false,
	})
	next(ctx)
}

func getTestAPI(t *testing.T) humatest.TestAPI {
	_, api := humatest.New(t)

	init_db(context.Background())
	api.UseMiddleware(testUserMiddleware)
	addRoutes(api)
	return api
}

// func testAdminMiddleware(ctx huma.Context, next func(huma.Context)) {
// 	ctx = huma.WithValue(ctx, currentUserCtxKey, User{
// 		ID:    "test_admin",
// 		Admin: true,
// 	})
// 	next(ctx)
// }
//
// func getTestAPIAdmin(t *testing.T) humatest.TestAPI {
// 	_, api := humatest.New(t)
//
// 	api.UseMiddleware(testAdminMiddleware)
// 	addRoutes(api)
// 	return api
// }

func assertRespCode(t *testing.T, resp *httptest.ResponseRecorder, expected_code int) *httptest.ResponseRecorder {
	if resp.Code != expected_code {
		t.Fatal("returned an invalid status code", resp.Code)
	}
	return resp
}

func TestAuthorTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/author",
			map[string]any{
				"name": "author_name",
			},
		),
		200,
	)

	data := AuthorOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/api/tags/author/%d", data.Body.Author.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestFindAuthorTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/author",
			map[string]any{
				"name": "author_name_find_test",
			},
		),
		200,
	)

	data := AuthorOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertRespCode(t, api.Get("/api/tags/author?name=author_name_find_test"), 200)

	path := fmt.Sprintf("/api/tags/author/%d", data.Body.Author.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestArtistTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/artist",
			map[string]any{
				"name":             "artist_name",
				"additional_names": []string{"additional_artist_name"},
			},
		),
		200,
	)

	data := ArtistOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/api/tags/artist/%d", data.Body.Artist.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestFindArtistTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/artist",
			map[string]any{
				"name":             "artist_name_find_test",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := ArtistOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertRespCode(t, api.Get("/api/tags/artist?name=artist_name_find_test"), 200)

	path := fmt.Sprintf("/api/tags/artist/%d", data.Body.Artist.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestMediaTag(t *testing.T) {
	api := getTestAPI(t)

	for _, v := range MediaTypes {
		resp := assertRespCode(t,
			api.Post("/api/tags/media",
				map[string]any{
					"name":             "media_name",
					"media_type":       v.ID,
					"additional_names": []string{"additional_media_name"},
				},
			),
			200,
		)

		data := MediaOutput{}
		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&data.Body)
		if err != nil {
			t.Fatal(err)
		}

		if data.Body.Media.AdditionalNames[0].Name != "additional_media_name" {
			t.Log("Failed to set media additional name")
			t.Fail()
		}

		path := fmt.Sprintf("/api/tags/media/%d", data.Body.Media.ID)
		assertRespCode(t, api.Delete(path), 204)
	}

}

func TestFindMedia(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/media",
			map[string]any{
				"name":             "media_name_find_test",
				"media_type":       "ANIME",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := MediaOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertRespCode(t, api.Get("/api/tags/media?name=media_name_find_test"), 200)

	path := fmt.Sprintf("/api/tags/media/%d", data.Body.Media.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestCreateKara(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/kara",
			map[string]any{
				"title":         "kara_title",
				"title_aliases": []string{"kara_title_alias"},
				"authors":       []uint{},
				"artists":       []uint{},
				"source_media":  0,
				"song_order":    0,
				"medias":        []uint{},
				"audio_tags":    []string{},
				"video_tags":    []string{},
				"comment":       "",
				"version":       "",
				"language":      "",
			}),
		200,
	)

	data := KaraOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	if data.Body.Kara.ExtraTitles[0].Name != "kara_title_alias" {
		t.Log("failed to set extra title to karaoke")
		t.Fail()
	}

	path := fmt.Sprintf("/api/kara/%d", data.Body.Kara.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestUpdateKara(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/kara",
			map[string]any{
				"title":         "kara_title_pre_update",
				"title_aliases": []string{"kara_update_title_alias"},
				"authors":       []uint{},
				"artists":       []uint{},
				"source_media":  0,
				"song_order":    0,
				"medias":        []uint{},
				"audio_tags":    []string{},
				"video_tags":    []string{},
				"comment":       "",
				"version":       "",
				"language":      "",
			}),
		200,
	)

	data := KaraOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/api/kara/%d", data.Body.Kara.ID)

	resp = assertRespCode(t,
		api.Patch(path,
			map[string]any{
				"title":                 "kara_title_post_update",
				"title_aliases":         []string{"kara_update_title_alias"},
				"authors":               []uint{},
				"artists":               []uint{},
				"source_media":          0,
				"song_order":            0,
				"medias":                []uint{},
				"audio_tags":            []string{},
				"video_tags":            []string{},
				"comment":               "",
				"version":               "",
				"language":              "",
				"is_hardsub":            false,
				"karaoke_creation_time": 0,
			}),
		200,
	)

	prev_ID := data.Body.Kara.ID

	data = KaraOutput{}
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	if data.Body.Kara.ID != prev_ID {
		t.Fatal("Karaoke ID changed after PATCH")
	}

	if data.Body.Kara.Title != "kara_title_post_update" {
		t.Fatal("Failed to update karaoke")
	}

	history_path := fmt.Sprintf("/api/kara/%d/history", data.Body.Kara.ID)
	resp = assertRespCode(t, api.Get(history_path), 200)

	history_data := GetKaraHistoryOutput{}
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&history_data.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(history_data.Body.History) != 1 {
		t.Fatalf("wrong number of history entries: %d", len(history_data.Body.History))
	}

	assertRespCode(t, api.Delete(path), 204)
}

func TestUpdateAuthor(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/author",
			map[string]any{
				"name": "author_name_update_test",
			},
		),
		200,
	)

	data := AuthorOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	author_path := fmt.Sprintf("/api/tags/author/%d", data.Body.Author.ID)

	resp = assertRespCode(t,
		api.Patch(author_path,
			map[string]any{
				"name": "author_name_update_test2",
			},
		),
		200,
	)

	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}
	if data.Body.Author.Name != "author_name_update_test2" {
		t.Fatal("failed to update author name")
	}

	resp = assertRespCode(t,
		api.Patch(author_path,
			map[string]any{
				"name": "author_name_update_test3",
			},
		),
		200,
	)

	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}
	if data.Body.Author.Name != "author_name_update_test3" {
		t.Fatal("failed to update author name a second time")
	}
}

func TestUpdateMedia(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/media",
			map[string]any{
				"name":             "media_name_update_test",
				"media_type":       "ANIME",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := MediaOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	media_path := fmt.Sprintf("/api/tags/media/%d", data.Body.Media.ID)

	resp = assertRespCode(t,
		api.Patch(media_path,
			map[string]any{
				"name":             "media_name_update_test2",
				"media_type":       "LIVE",
				"additional_names": []string{},
			},
		),
		200,
	)

	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}
	if data.Body.Media.Name != "media_name_update_test2" {
		t.Fatal("failed to update media name")
	}
	if data.Body.Media.Type != "LIVE" {
		t.Fatal("failed to update media type")
	}

	resp = assertRespCode(t,
		api.Patch(media_path,
			map[string]any{
				"name":             "media_name_update_test3",
				"media_type":       "ANIME",
				"additional_names": []string{},
			},
		),
		200,
	)

	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}
	if data.Body.Media.Name != "media_name_update_test3" {
		t.Fatal("failed to update media name a second time")
	}
	if data.Body.Media.Type != "ANIME" {
		t.Fatal("failed to update media type a second time")
	}
}

func TestUpdateArtist(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/artist",
			map[string]any{
				"name":             "artist_name_update_test",
				"additional_names": []string{},
			},
		),
		200,
	)

	data := ArtistOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	artist_path := fmt.Sprintf("/api/tags/artist/%d", data.Body.Artist.ID)

	resp = assertRespCode(t,
		api.Patch(artist_path,
			map[string]any{
				"name":             "artist_name_update_test2",
				"additional_names": []string{},
			},
		),
		200,
	)

	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}
	if data.Body.Artist.Name != "artist_name_update_test2" {
		t.Fatal("failed to update artist name")
	}

	resp = assertRespCode(t,
		api.Patch(artist_path,
			map[string]any{
				"name":             "artist_name_update_test3",
				"additional_names": []string{},
			},
		),
		200,
	)

	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}
	if data.Body.Artist.Name != "artist_name_update_test3" {
		t.Fatal("failed to update artist name a second time")
	}
}
