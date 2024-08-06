// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 Japan7
package mugen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ExternalDatabaseIDs struct {
	Anilist *int `json:"anilist"`
}

type MugenTag struct {
	TID                 uuid.UUID           `json:"tid"`
	Name                string              `json:"name"`
	Short               string              `json:"short"`
	I18n                map[string]string   `json:"i18n"`
	Aliases             []string            `json:"aliases"`
	ExternalDatabaseIDs ExternalDatabaseIDs `json:"external_database_ids"`
}

type Kara struct {
	KID                  uuid.UUID         `json:"kid"`
	Titles               map[string]string `json:"titles"`
	TitleAliases         []string          `json:"titles_aliases"`
	TitleDefaultLanguage string            `json:"titles_default_language"`
	MediaFile            string            `json:"mediafile"`
	MediaSize            uint64            `json:"mediasize"`
	SubFile              string            `json:"subfile"`
	SubChecksum          string            `json:"subchecksum"`
	Duration             int               `json:"duration"`
	SongOrder            *uint             `json:"songorder"`
	CreatedAt            time.Time         `json:"created_at"`
	ModifiedAt           time.Time         `json:"modified_at"`
	Series               []MugenTag        `json:"series"`
	Singers              []MugenTag        `json:"singers"`
	SongTypes            []MugenTag        `json:"songtypes"`
	Creators             []MugenTag        `json:"creators"`
	Languages            []MugenTag        `json:"langs"`
	Authors              []MugenTag        `json:"authors"`
	Misc                 []MugenTag        `json:"misc"`
	SongWriters          []MugenTag        `json:"songwriters"`
	Families             []MugenTag        `json:"families"`
	Origins              []MugenTag        `json:"origins"`
	Genres               []MugenTag        `json:"genres"`
	Platforms            []MugenTag        `json:"platforms"`
	Versions             []MugenTag        `json:"versions"`
	Warnings             []MugenTag        `json:"warnings"`
	Collections          []MugenTag        `json:"Collections"`
	SingerGroups         []MugenTag        `json:"singergroups"`
	Franchises           []MugenTag        `json:"franchises"`
	Comment              string            `json:"comment"`
}

type Client interface {
	GetKara(ctx context.Context, kid uuid.UUID) (*Kara, error)
	DownloadMedia(ctx context.Context, mediafile string) (*http.Response, error)
	DownloadLyrics(ctx context.Context, karafile string) (*http.Response, error)
}

type MugenClient struct {
	Server      string
	MediaServer string
	HTTPClient  *http.Client
}

func (c MugenClient) GetEndpoint(path string) string {
	server_base := strings.TrimLeft(c.Server, "/")
	return fmt.Sprintf("%s%s", server_base, path)
}

func Closer(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (c MugenClient) SendRequest(ctx context.Context, method string, path string, bodyData any) (*http.Response, error) {
	endpoint := c.GetEndpoint(path)
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		defer Closer(resp.Body)
		buf := make([]byte, resp.ContentLength)
		n, err := resp.Body.Read(buf)
		if err != nil {
			fmt.Printf("Failed to read body of mugen response: %s\n%s %s\nbody: %+v", buf[:n], method, path, bodyData)
			return nil, err
		}
		fmt.Printf("mugen response: %+v\n%s", resp, buf[:n])
		return nil, fmt.Errorf("dakara responded with status code %d", resp.StatusCode)
	}

	return resp, err
}

func (c MugenClient) GetKara(ctx context.Context, kid uuid.UUID) (*Kara, error) {
	path := fmt.Sprintf("/karas/%s/", kid)
	resp, err := c.SendRequest(ctx, http.MethodGet, path, struct{}{})
	if err != nil {
		return nil, err
	}
	defer Closer(resp.Body)

	dec := json.NewDecoder(resp.Body)
	data := &Kara{}
	err = dec.Decode(data)

	return data, err
}

func (c MugenClient) DownloadMedia(ctx context.Context, mediafile string) (*http.Response, error) {
	url := fmt.Sprintf("%s/medias/%s", c.MediaServer, mediafile)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.HTTPClient.Do(req)
}

func (c MugenClient) DownloadLyrics(ctx context.Context, karafile string) (*http.Response, error) {
	url := fmt.Sprintf("%s/lyrics/%s", c.MediaServer, karafile)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.HTTPClient.Do(req)
}

var MUGEN_SERVER = "https://kara.moe/api/"
var MUGEN_MEDIA_SERVER = "https://kara.moe/downloads/"
var MUGEN_CLIENT_INST Client = nil

func GetClient() Client {
	if MUGEN_CLIENT_INST == nil {
		http_client := &http.Client{}
		MUGEN_CLIENT_INST = MugenClient{
			Server:      MUGEN_SERVER,
			MediaServer: MUGEN_MEDIA_SERVER,
			HTTPClient:  http_client,
		}
	}

	return MUGEN_CLIENT_INST
}
