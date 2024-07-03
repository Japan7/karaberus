// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
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
		api.Post("/api/tags/author",
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

	path := fmt.Sprintf("/api/tags/author/%d", data.Body.author.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestFindAuthorTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/author",
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

	resp = assertRespCode(t, api.Get("/api/tags/author?name=author_name"), 200)

	path := fmt.Sprintf("/api/tags/author/%d", data.Body.author.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestArtistTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/artist",
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

	path := fmt.Sprintf("/api/tags/artist/%d", data.Body.Artist.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestFindArtistTag(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/artist",
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

	resp = assertRespCode(t, api.Get("/api/tags/artist?name=artist_name"), 200)

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
					"additional_names": []string{},
				},
			),
			200,
		)

		data := MediaOutput{}
		dec := json.NewDecoder(resp.Body)
		dec.Decode(&data.Body)

		path := fmt.Sprintf("/api/tags/media/%d", data.Body.Media.ID)
		assertRespCode(t, api.Delete(path), 204)
	}

}

func TestFindMedia(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/tags/media",
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

	resp = assertRespCode(t, api.Get("/api/tags/media?name=media_name"), 200)

	path := fmt.Sprintf("/api/tags/media/%d", data.Body.Media.ID)
	assertRespCode(t, api.Delete(path), 204)
}

func TestCreateKara(t *testing.T) {
	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/kara",
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
				"language":      "",
			}),
		200,
	)

	data := KaraOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	path := fmt.Sprintf("/api/kara/%d", data.Body.Kara.ID)
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
	fwriter, err := multipart_writer.CreateFormFile("file", "karaberus_test.mkv")
	if err != nil {
		panic("failed to create multipart file")
	}

	tmpbuf := make([]byte, 1024*4)
	for {
		n, err := f.Read(tmpbuf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}

		fwriter.Write(tmpbuf[:n])
	}

	multipart_writer.Close()

	path := fmt.Sprintf("/api/kara/%d/upload/%s", kid, filetype)
	headers := "Content-Type: multipart/form-data; boundary=" + multipart_writer.Boundary()
	resp := assertRespCode(t,
		api.Put(path, headers, buf),
		200,
	)

	data_upload := UploadOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data_upload.Body)
	return data_upload
}

func TestUploadKara(t *testing.T) {
	skipCI(t)

	api := getTestAPI(t)

	resp := assertRespCode(t,
		api.Post("/api/kara",
			map[string]any{
				"title":         "kara_upload_title",
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
				"language":      "",
			}),
		200,
	)

	data := KaraOutput{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data.Body)

	mkv_test_file := path.Join(CONFIG.GENERATED_TEST_DIR, "karaberus_test.mkv")
	video_upload := uploadFile(t, api, data.Body.Kara.ID, mkv_test_file, "video")
	if !video_upload.Body.CheckResults.Video.Passed {
		t.Log("Video did not pass checks.")
		t.Fail()
	}
	if video_upload.Body.CheckResults.Instrumental != nil {
		t.Log("Instrumental should be uploaded yet.")
		t.Fail()
	}
	if video_upload.Body.CheckResults.Subtitles != nil {
		t.Log("Subtitles should be uploaded yet.")
		t.Fail()
	}

	ass_test_file := path.Join(CONFIG.TEST_DIR, "test.ass")
	sub_upload := uploadFile(t, api, data.Body.Kara.ID, ass_test_file, "sub")
	if !sub_upload.Body.CheckResults.Video.Passed {
		t.Log("Video did not pass checks.")
		t.Fail()
	}
	if !sub_upload.Body.CheckResults.Subtitles.Passed {
		t.Log("Subtitles did not pass checks.")
		t.Fail()
	}
	if sub_upload.Body.CheckResults.Instrumental != nil {
		t.Log("Instrumental should be uploaded yet.")
		t.Fail()
	}

	inst_test_file := path.Join(CONFIG.GENERATED_TEST_DIR, "karaberus_test.opus")
	inst_upload := uploadFile(t, api, data.Body.Kara.ID, inst_test_file, "inst")
	if !inst_upload.Body.CheckResults.Video.Passed {
		t.Log("Video did not pass checks.")
		t.Fail()
	}
	if !inst_upload.Body.CheckResults.Subtitles.Passed {
		t.Log("Subtitles did not pass checks.")
		t.Fail()
	}
	if !inst_upload.Body.CheckResults.Instrumental.Passed {
		t.Log("Instrumental did not pass checks.")
		t.Fail()
	}
}
