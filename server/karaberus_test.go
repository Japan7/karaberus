// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
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

func testAdminMiddleware(ctx huma.Context, next func(huma.Context)) {
	ctx = huma.WithValue(ctx, currentUserCtxKey, User{
		ID:    "test_admin",
		Admin: true,
	})
	next(ctx)
}

func getTestAPI(t *testing.T) humatest.TestAPI {
	_, api := humatest.New(t)

	db := GetDB(context.TODO())
	init_model(db)

	api.UseMiddleware(testUserMiddleware)
	addRoutes(api)
	return api
}

func getTestAPIAdmin(t *testing.T) humatest.TestAPI {
	_, api := humatest.New(t)

	db := GetDB(context.TODO())
	init_model(db)

	api.UseMiddleware(testAdminMiddleware)
	addRoutes(api)
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

	assertRespCode(t, api.Delete(path), 204)
}

func skipCI(t *testing.T) {
	if os.Getenv("SKIP_S3_TESTS") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func uploadFile(t *testing.T, api humatest.TestAPI, kid uint, filepath string, filetype string) UploadOutput {
	f, err := os.Open(filepath)
	if err != nil {
		panic("failed to open " + filepath)
	}

	buf := new(bytes.Buffer)
	multipart_writer := multipart.NewWriter(buf)
	fwriter, err := multipart_writer.CreateFormFile("file", filepath)
	if err != nil {
		panic("failed to create multipart file")
	}

	_, err = io.Copy(fwriter, f)
	if err != nil {
		t.Fatal(err)
	}

	err = multipart_writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/api/kara/%d/upload/%s", kid, filetype)
	headers := "Content-Type: multipart/form-data; boundary=" + multipart_writer.Boundary()
	resp := assertRespCode(t,
		api.Put(path, headers, buf),
		200,
	)

	data_upload := UploadOutput{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&data_upload.Body)
	if err != nil {
		t.Fatal(err)
	}
	return data_upload
}

func uploadFont(t *testing.T, api humatest.TestAPI, filepath string) UploadFontOutput {
	f, err := os.Open(filepath)
	if err != nil {
		panic("failed to open " + filepath)
	}

	buf := new(bytes.Buffer)
	multipart_writer := multipart.NewWriter(buf)
	fwriter, err := multipart_writer.CreateFormFile("file", filepath)
	if err != nil {
		panic("failed to create multipart file")
	}

	_, err = io.Copy(fwriter, f)
	if err != nil {
		t.Fatal(err)
	}

	err = multipart_writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	path := "/api/font"
	headers := "Content-Type: multipart/form-data; boundary=" + multipart_writer.Boundary()
	resp := assertRespCode(t,
		api.Post(path, headers, buf),
		200,
	)

	data_upload := UploadFontOutput{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&data_upload.Body)
	if err != nil {
		t.Fatal(err)
	}
	return data_upload
}

func CompareDownloadedFile(t *testing.T, api humatest.TestAPI, original_file string, kid uint, filetype string) {
	path := fmt.Sprintf("/api/kara/%d/download/%s", kid, filetype)
	resp := assertRespCode(t,
		api.Get(path),
		200,
	)

	orig, err := os.Open(original_file)
	if err != nil {
		t.Fatalf("failed to open %s", original_file)
	}

	for {
		orig_buf := make([]byte, 1024*1024)
		dl_buf := make([]byte, 1024*1024)
		_, oerr := orig.Read(orig_buf)
		_, dlerr := resp.Result().Body.Read(dl_buf)

		if oerr != dlerr {
			t.Fatalf("%s %s: oerr=%v != dlerr=%v", original_file, path, oerr, dlerr)
		}

		if len(orig_buf) != len(dl_buf) {
			t.Fatalf("buf sizes differ %s", original_file)
		}

		for i := 0; i < len(orig_buf); i++ {
			if orig_buf[i] != dl_buf[i] {
				t.Fatalf("downloaded file is different from original file %s", original_file)
			}
		}

		if errors.Is(oerr, io.EOF) {
			break
		}
	}
}

func CompareDownloadedFont(t *testing.T, api humatest.TestAPI, original_file string, id uint) {
	path := fmt.Sprintf("/api/font/%d/download", id)
	resp := assertRespCode(t, api.Get(path), 200)

	orig, err := os.Open(original_file)
	if err != nil {
		t.Fatalf("failed to open %s", original_file)
	}

	for {
		orig_buf := make([]byte, 1024*1024)
		dl_buf := make([]byte, 1024*1024)
		_, oerr := orig.Read(orig_buf)
		_, dlerr := resp.Result().Body.Read(dl_buf)

		if oerr != dlerr {
			t.Fatalf("%s %s: oerr=%v != dlerr=%v", original_file, path, oerr, dlerr)
		}

		if len(orig_buf) != len(dl_buf) {
			t.Fatalf("buf sizes differ %s", original_file)
		}

		for i := 0; i < len(orig_buf); i++ {
			if orig_buf[i] != dl_buf[i] {
				t.Fatalf("downloaded file is different from original file %s", original_file)
			}
		}

		if errors.Is(oerr, io.EOF) {
			break
		}
	}
}

func TestUploadKara(t *testing.T) {
	skipCI(t)

	api := getTestAPIAdmin(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/media",
			map[string]any{
				"name":             "test_media_karaupload",
				"media_type":       "ANIME",
				"additional_names": []string{},
			}),
		200,
	)

	media_data := MediaOutput{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&media_data.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp = assertRespCode(t,
		api.Post("/api/tags/artist",
			map[string]any{
				"name":             "test_artist_karaupload",
				"additional_names": []string{},
			}),
		200,
	)
	artist_data := ArtistOutput{}
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&artist_data.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp = assertRespCode(t,
		api.Post("/api/kara",
			map[string]any{
				"title":         "kara_upload_title",
				"title_aliases": []string{},
				"authors":       []uint{},
				"artists":       []uint{artist_data.Body.Artist.ID},
				"source_media":  media_data.Body.Media.ID,
				"song_order":    0,
				"medias":        []uint{media_data.Body.Media.ID},
				"audio_tags":    []string{},
				"video_tags":    []string{},
				"comment":       "",
				"version":       "",
				"language":      "",
			}),
		200,
	)

	data := KaraOutput{}
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	if data.Body.Kara.SourceMedia.ID != media_data.Body.Media.ID {
		t.Fatal("Kara source media is not set")
	}

	if data.Body.Kara.Medias[0].ID != media_data.Body.Media.ID {
		t.Fatal("Kara medias are badly set")
	}

	if data.Body.Kara.Artists[0].ID != artist_data.Body.Artist.ID {
		t.Fatal("Kara artists are badly set")
	}

	mkv_test_file := path.Join(TEST_CONFIG.GeneratedDirectory, "karaberus_test.mkv")
	video_upload := uploadFile(t, api, data.Body.Kara.ID, mkv_test_file, "video")
	if !video_upload.Body.CheckResults.Video.Passed {
		t.Fatal("Video did not pass checks.")
	}
	if video_upload.Body.CheckResults.Instrumental != nil {
		t.Fatal("Instrumental should be uploaded yet.")
	}
	if video_upload.Body.CheckResults.Subtitles != nil {
		t.Fatal("Subtitles should be uploaded yet.")
	}

	ass_test_file := path.Join(TEST_CONFIG.Directory, "test.ass")
	sub_upload := uploadFile(t, api, data.Body.Kara.ID, ass_test_file, "sub")
	if !sub_upload.Body.CheckResults.Video.Passed {
		t.Fatal("Video did not pass checks.")
	}
	if !sub_upload.Body.CheckResults.Subtitles.Passed {
		t.Fatal("Subtitles did not pass checks.")
	}
	if sub_upload.Body.CheckResults.Instrumental != nil {
		t.Fatal("Instrumental should be uploaded yet.")
	}

	inst_test_file := path.Join(TEST_CONFIG.GeneratedDirectory, "karaberus_test.opus")
	inst_upload := uploadFile(t, api, data.Body.Kara.ID, inst_test_file, "inst")
	if !inst_upload.Body.CheckResults.Video.Passed {
		t.Fatal("Video did not pass checks.")
	}
	if !inst_upload.Body.CheckResults.Subtitles.Passed {
		t.Fatal("Subtitles did not pass checks.")
	}
	if !inst_upload.Body.CheckResults.Instrumental.Passed {
		t.Fatal("Instrumental did not pass checks.")
	}

	// can't do that from humatest because we stream the files using fasthttp's Response
	// which obviously humatest is not aware of
	//
	// CompareDownloadedFile(t, api, mkv_test_file, data.Body.Kara.ID, "video")
	// CompareDownloadedFile(t, api, ass_test_file, data.Body.Kara.ID, "sub")
	// CompareDownloadedFile(t, api, inst_test_file, data.Body.Kara.ID, "inst")

	kara_path := fmt.Sprintf("/api/kara/%d", data.Body.Kara.ID)
	resp = assertRespCode(t, api.Get(kara_path), 200)
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	kara_path_creation_time := fmt.Sprintf("/api/kara/%d", data.Body.Kara.ID)
	newCreationTime := int64(3600)
	resp = assertRespCode(t,
		api.Patch(kara_path_creation_time, map[string]any{
			"title":                 "kara_upload_title",
			"title_aliases":         []string{},
			"authors":               []uint{},
			"artists":               []uint{artist_data.Body.Artist.ID},
			"source_media":          media_data.Body.Media.ID,
			"song_order":            0,
			"medias":                []uint{media_data.Body.Media.ID},
			"audio_tags":            []string{},
			"video_tags":            []string{},
			"comment":               "",
			"version":               "",
			"language":              "",
			"karaoke_creation_time": newCreationTime,
			"is_hardsub":            false,
		}),
		200,
	)
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&data.Body)
	if err != nil {
		t.Fatal(err)
	}

	if data.Body.Kara.KaraokeCreationTime.Unix() != newCreationTime {
		t.Fatal("failed to set karaoke creation date", data.Body.Kara.KaraokeCreationTime)
	}
}

func TestUploadFont(t *testing.T) {
	skipCI(t)

	api := getTestAPI(t)

	font_test_file := path.Join(TEST_CONFIG.Directory, "KaraberusTestFont.ttf")
	uploadFont(t, api, font_test_file)

	// can't do that from humatest because we stream the files using fasthttp's Response
	// which obviously humatest is not aware of
	//
	// CompareDownloadedFont(t, api, font_test_file, resp.Body.Font.ID)
}
