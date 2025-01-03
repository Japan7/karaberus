// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 Japan7

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func dakaraApiEndpoint(path string) string {
	return fmt.Sprintf("%s%s", CONFIG.Dakara.BaseURL, path)
}

type DakaraOutput struct {
	Status int
	Body   []byte
}

func dakaraSendRequest(ctx context.Context, method string, path string, bodyData any) (*http.Response, error) {
	endpoint := dakaraApiEndpoint(path)
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", CONFIG.Dakara.Token))
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
			getLogger().Printf("Failed to read body of dakara response: %s\n%s %s\nbody: %+v", buf[:n], method, path, bodyData)
			return nil, err
		}
		getLogger().Printf("dakara response: %+v\n%s", resp, buf[:n])
		return nil, fmt.Errorf("dakara responded with status code %d", resp.StatusCode)
	}

	return resp, err
}

func dakaraPost(ctx context.Context, path string, bodyData any) (*http.Response, error) {
	return dakaraSendRequest(ctx, http.MethodPost, path, bodyData)
}

func dakaraPut(ctx context.Context, path string, bodyData any) (*http.Response, error) {
	return dakaraSendRequest(ctx, http.MethodPut, path, bodyData)
}

func dakaraDelete(ctx context.Context, path string) (*http.Response, error) {
	return dakaraSendRequest(ctx, http.MethodDelete, path, struct{}{})
}

func dakaraGet(ctx context.Context, path string, page int) (*http.Response, error) {
	if page > 0 {
		path = fmt.Sprintf("%s?page=%d", path, page)
	}
	return dakaraSendRequest(ctx, http.MethodGet, path, struct{}{})
}

type DakaraPagination struct {
	Current int `json:"current" example:"1"`
	Last    int `json:"last" example:"10"`
}

type DakaraPaginatedResponse struct {
	Pagination DakaraPagination `json:"pagination"`
	Count      int              `json:"count" example:"99"`
}

type DakaraArtist struct {
	DakaraArtistBody
	ID        int `json:"id"`
	SongCount int `json:"song_count"`
}

type DakaraGetArtistsResponse struct {
	DakaraPaginatedResponse
	Results []DakaraArtist `json:"results"`
}

func dakaraGetArtists(ctx context.Context) (map[string]*DakaraArtist, error) {
	page := 1
	artists := map[string]*DakaraArtist{}

	for {
		resp, err := dakaraGet(ctx, "/api/library/artists/", page)
		if err != nil {
			return nil, err
		}
		defer Closer(resp.Body)

		dec := json.NewDecoder(resp.Body)
		data := DakaraGetArtistsResponse{}
		err = dec.Decode(&data)
		if err != nil {
			return nil, err
		}
		for _, artist := range data.Results {
			artists[artist.Name] = &artist
		}

		if data.Pagination.Current == data.Pagination.Last {
			break
		}
		page++
	}

	return artists, nil
}

type DakaraGetWorkTypesResponse struct {
	DakaraPaginatedResponse
	Results []DakaraWorkType `json:"results"`
}

func dakaraGetWorkTypes(ctx context.Context) (map[string]*DakaraWorkType, error) {
	page := 1
	worktypes := map[string]*DakaraWorkType{}

	for {
		resp, err := dakaraGet(ctx, "/api/library/work-types/", page)
		if err != nil {
			return nil, err
		}
		defer Closer(resp.Body)

		dec := json.NewDecoder(resp.Body)
		data := DakaraGetWorkTypesResponse{}
		err = dec.Decode(&data)
		if err != nil {
			return nil, err
		}
		for _, worktype := range data.Results {
			worktypes[worktype.QueryName] = &worktype
		}

		if data.Pagination.Current == data.Pagination.Last {
			break
		}
		page++
	}

	return worktypes, nil
}

type DakaraWork struct {
	ID                int                `json:"id"`
	Title             string             `json:"title"`
	Subtitle          string             `json:"subtitle"` // we ignore subtitles
	AlternativeTitles []string           `json:"AlternativeTitles"`
	WorkType          DakaraWorkTypeBody `json:"work_type"`
	SongCount         int                `json:"song_count"`
}

type DakaraGetWorksResponse struct {
	DakaraPaginatedResponse
	Results []DakaraWork `json:"results"`
}

func dakaraGetWorks(ctx context.Context) (map[string]map[string]*DakaraWork, error) {
	page := 1
	worktypes := map[string]map[string]*DakaraWork{}
	for _, media_type := range MediaTypes {
		worktypes[strings.ToLower(media_type.ID)] = map[string]*DakaraWork{}
	}

	for {
		resp, err := dakaraGet(ctx, "/api/library/works/", page)
		if err != nil {
			return nil, err
		}
		defer Closer(resp.Body)

		dec := json.NewDecoder(resp.Body)
		data := DakaraGetWorksResponse{}
		err = dec.Decode(&data)
		if err != nil {
			return nil, err
		}
		for _, work := range data.Results {
			worktype_map := worktypes[work.WorkType.QueryName]
			if worktype_map == nil {
				continue
			}
			worktype_map[work.Title] = &work
		}

		if data.Pagination.Current == data.Pagination.Last {
			break
		}
		page++
	}

	return worktypes, nil
}

type DakaraTag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ColorHue int    `json:"color_hue"`
	Disabled bool   `json:"disabled"`
}

type DakaraGetTagsResponse struct {
	DakaraPaginatedResponse
	Results []DakaraTag `json:"results"`
}

func dakaraGetTags(ctx context.Context) (map[string]*DakaraTag, error) {
	page := 1
	worktypes := map[string]*DakaraTag{}

	for {
		resp, err := dakaraGet(ctx, "/api/library/song-tags/", page)
		if err != nil {
			return nil, err
		}
		defer Closer(resp.Body)

		dec := json.NewDecoder(resp.Body)
		data := DakaraGetTagsResponse{}
		err = dec.Decode(&data)
		if err != nil {
			return nil, err
		}
		for _, worktype := range data.Results {
			worktypes[worktype.Name] = &worktype
		}

		if data.Pagination.Current == data.Pagination.Last {
			break
		}
		page++
	}

	return worktypes, nil
}

type DakaraLyricsPreview struct {
	Text      string `json:"text"`
	Truncated bool   `json:"truncated"`
}

type DakaraSong struct {
	ID              int                 `json:"id"`
	Title           string              `json:"title"`
	Filename        string              `json:"filename"` // basically our ID
	Duration        int32               `json:"duration"`
	Directory       string              `json:"directory"`
	Version         string              `json:"version"`
	Detail          string              `json:"detail"`
	DetailVideo     string              `json:"detail_video"`
	Tags            []DakaraTag         `json:"tags"`
	Artists         []DakaraArtist      `json:"artists"`
	Works           []DakaraSongWork    `json:"works"`
	LyricsPreview   DakaraLyricsPreview `json:"lyrics_preview"`
	HasInstrumental bool                `json:"has_instrumental"`
	DateCreated     time.Time           `json:"date_created"`
	DateUpdated     time.Time           `json:"date_updated"`
}

type DakaraGetSongsResponse struct {
	DakaraPaginatedResponse
	Results []DakaraSong `json:"results"`
}

func dakaraGetSongs(ctx context.Context) (map[string]*DakaraSong, error) {
	page := 1
	worktypes := map[string]*DakaraSong{}

	for {
		resp, err := dakaraGet(ctx, "/api/library/songs/", page)
		if err != nil {
			return nil, err
		}
		defer Closer(resp.Body)

		dec := json.NewDecoder(resp.Body)
		data := DakaraGetSongsResponse{}
		err = dec.Decode(&data)
		if err != nil {
			return nil, err
		}
		for _, worktype := range data.Results {
			worktypes[worktype.Filename] = &worktype
		}

		if data.Pagination.Current == data.Pagination.Last {
			break
		}
		page++
	}

	return worktypes, nil
}

type DakaraArtistBody struct {
	Name string `json:"name"`
}

func dakaraAddArtist(ctx context.Context, artist Artist) error {
	artist_body := DakaraArtistBody{Name: artist.Name}

	resp, err := dakaraPost(ctx, "/api/library/artists/", artist_body)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	return nil
}

type DakaraWorkTypeBody struct {
	QueryName  string `json:"query_name"`
	Name       string `json:"name"`
	NamePlural string `json:"name_plural"`
	IconName   string `json:"icon_name"`
}

type DakaraWorkType struct {
	DakaraWorkTypeBody
	// not yet implemented in dakara-server
	// ID int `json:"id"`
}

func dakaraWorkType(media_type MediaType) DakaraWorkTypeBody {
	return DakaraWorkTypeBody{
		QueryName:  strings.ToLower(media_type.ID),
		Name:       media_type.Name,
		NamePlural: media_type.Name + "s", // works for now
		IconName:   media_type.IconName,
	}
}

func dakaraAddWorkType(ctx context.Context, media_type MediaType) error {
	worktype := dakaraWorkType(media_type)

	resp, err := dakaraPost(ctx, "/api/library/work-types/", worktype)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	return nil
}

type DakaraWorkBody struct {
	Title    string             `json:"title"`
	WorkType DakaraWorkTypeBody `json:"work_type"`
}

func dakaraAddWork(ctx context.Context, work DakaraWorkBody) error {
	resp, err := dakaraPost(ctx, "/api/library/works/", work)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	return nil
}

type DakaraTagBody struct {
	Name string `json:"name"`
	Hue  uint   `json:"color_hue"`
}

func dakaraAddTag(ctx context.Context, tag TagInterface) error {
	body := DakaraTagBody{
		Name: tag.getName(),
		Hue:  tag.getHue(),
	}

	resp, err := dakaraPost(ctx, "/api/library/song-tags/", body)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	return nil
}

func dakaraPutTag(ctx context.Context, id int, tag TagInterface) error {
	body := DakaraTagBody{
		Name: tag.getName(),
		Hue:  tag.getHue(),
	}

	path := fmt.Sprintf("/api/library/song-tags/%d/", id)
	resp, err := dakaraPut(ctx, path, body)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	return nil
}

type DakaraSongWork struct {
	Work           DakaraWork `json:"work"`
	LinkType       string     `json:"link_type"`
	LinkTypeNumber *uint      `json:"link_type_number"`
}

type DakaraSongBody struct {
	Title           string           `json:"title"`
	Filename        string           `json:"filename"` // basically our ID
	Duration        int32            `json:"duration"`
	Directory       string           `json:"directory"`
	Version         string           `json:"version"`
	Detail          string           `json:"detail"`
	DetailVideo     string           `json:"detail_video"`
	Tags            []DakaraTag      `json:"tags"`
	Artists         []DakaraArtist   `json:"artists"`
	Works           []DakaraSongWork `json:"works"`
	Lyrics          string           `json:"lyrics"`
	HasInstrumental bool             `json:"has_instrumental"`
}

func getLinkTypeNumber(n *uint) uint {
	if n == nil {
		return 0
	} else {
		return *n
	}
}

func (body DakaraSongBody) HasChanged(ref DakaraSong) bool {
	// lyrics are truncated in the response from dakara
	// so compare to our truncated text even if that means that we could miss an update
	possibly_truncated_lyrics := body.Lyrics
	if ref.LyricsPreview.Truncated {
		if len(body.Lyrics) > len(ref.LyricsPreview.Text) {
			possibly_truncated_lyrics = trimWhitespace(body.Lyrics[:len(ref.LyricsPreview.Text)])
		}
	}

	tags_equal := len(body.Tags) == len(ref.Tags)
	if tags_equal {
		for _, tag := range body.Tags {
			found := false
			for _, ref_tag := range ref.Tags {
				if ref_tag.ID == tag.ID {
					found = true
				}
			}
			if !found {
				tags_equal = false
			}
		}
	}

	artists_equal := len(body.Artists) == len(ref.Artists)
	if artists_equal {
		for _, artist := range body.Artists {
			found := false
			for _, ref_artist := range ref.Artists {
				if ref_artist.ID == artist.ID {
					found = true
				}
			}
			if !found {
				artists_equal = false
			}
		}
	}

	works_equal := len(body.Works) == len(ref.Works)
	if works_equal {
		for _, work := range body.Works {
			found := false
			for _, ref_work := range ref.Works {
				if ref_work.Work.ID == work.Work.ID && ref_work.LinkType == work.LinkType && getLinkTypeNumber(ref_work.LinkTypeNumber) == getLinkTypeNumber(work.LinkTypeNumber) {
					found = true
				}
			}
			if !found {
				works_equal = false
			}
		}
	}

	return body.Title != ref.Title ||
		body.Filename != ref.Filename ||
		body.Duration != ref.Duration ||
		body.Directory != ref.Directory ||
		body.Version != ref.Version ||
		body.Detail != ref.Detail ||
		body.DetailVideo != ref.DetailVideo ||
		possibly_truncated_lyrics != ref.LyricsPreview.Text ||
		!tags_equal ||
		!artists_equal ||
		!works_equal
}

func dakaraAddSong(ctx context.Context, song *DakaraSongBody) error {
	resp, err := dakaraPost(ctx, "/api/library/songs/", song)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	return nil
}

// func worktypeShouldExist(worktype string) bool {
// 	for _, media_type := range MediaTypes {
// 		if strings.ToLower(media_type.ID) == worktype {
// 			return true
// 		}
// 	}
// 	return false
// }

// func cleanUpWorkTypes(worktypes []DakaraWorkType) error {
// 	for _, worktype := range worktypes {
// 		if !worktypeShouldExist(worktype.QueryName) {
// 			// ID is missing from the work type struct
// 			// return DakaraDeleteWorkType(ctx, worktype)
// 		}
// 	}
// 	return nil
// }

func dakaraFilterAudioTags(audio_tags []AudioTag) []AudioTag {
	out := []AudioTag{}
	for _, tag := range audio_tags {
		switch tag.ID {
		case "OP", "ED", "INS", "IS":
			continue
		default:
			out = append(out, tag)
		}
	}

	return out
}

func UploadedKaras(db *gorm.DB) *gorm.DB {
	return db.Where("video_uploaded AND (subtitles_uploaded OR hardsubbed)")
}

var syncDakaraNotifyMutex = sync.Mutex{}
var syncDakaraNotifyChannel = make(chan bool)

// Notify the sync worker (non blocking)
func SyncDakaraNotify() {
	go func() {
		if syncDakaraNotifyMutex.TryLock() {
			defer syncDakaraNotifyMutex.Unlock()
			getLogger().Println("notifying dakara worker")
			syncDakaraNotifyChannel <- true
			getLogger().Println("notified dakara worker")
		}
	}()
}

func SyncDakaraLoop(ctx context.Context) {
	for {
		SyncDakara(ctx)
		<-syncDakaraNotifyChannel
	}
}

func SyncDakara(ctx context.Context) {
	logger := getLogger()
	defer logger.Println("Done syncing dakara")
	defer func() {
		r := recover()
		if r != nil {
			getLogger().Printf("recovered from panic in SyncDakara %s\n%s\n", r, string(debug.Stack()))
		}
	}()

	db := GetDB(ctx)
	// sync media types / work types
	worktypes, err := dakaraGetWorkTypes(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	for _, media_type := range MediaTypes {
		if worktypes[strings.ToLower(media_type.ID)] == nil {
			err = dakaraAddWorkType(ctx, media_type)
			if err != nil {
				getLogger().Println(err)
				return
			}
		}
	}

	all_karas := []KaraInfoDB{}
	err = db.Preload(clause.Associations).Scopes(UploadedKaras, CurrentKaras).Find(&all_karas).Error
	if err != nil {
		getLogger().Println(err)
		return
	}
	logger.Printf("Syncing %d karas to Dakara", len(all_karas))

	// sync media / works
	works, err := dakaraGetWorks(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	all_medias := map[uint]MediaDB{}
	all_artists := map[uint]Artist{}

	for _, kara := range all_karas {
		if kara.SourceMedia != nil {
			all_medias[kara.SourceMedia.ID] = *kara.SourceMedia
		}
		for _, media := range kara.Medias {
			all_medias[media.ID] = media
		}
		for _, artist := range kara.Artists {
			all_artists[artist.ID] = artist
		}
	}

	new_works := 0
	for _, media := range all_medias {
		if works[strings.ToLower(media.Type)][media.Name] == nil {
			media_type := getMediaType(media.Type)
			err = dakaraAddWork(ctx, DakaraWorkBody{
				Title:    media.Name,
				WorkType: dakaraWorkType(media_type),
			})
			if err != nil {
				getLogger().Println(err)
				return
			}
			new_works++
		}
	}
	logger.Printf("Added %d new works to dakara\n", new_works)

	// sync artists
	dakara_artists, err := dakaraGetArtists(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	new_artists := 0
	for _, artist := range all_artists {
		if dakara_artists[artist.Name] == nil {
			err = dakaraAddArtist(ctx, artist)
			if err != nil {
				getLogger().Println(err)
				return
			}
			new_artists++
		}
	}
	logger.Printf("Added %d new artists to dakara\n", new_artists)

	// sync tags
	dakara_tags, err := dakaraGetTags(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	// sync audio tags
	for _, audio_tag := range dakaraFilterAudioTags(AudioTags) {
		dakara_tag := dakara_tags[audio_tag.Name]
		if dakara_tag == nil {
			err = dakaraAddTag(ctx, audio_tag)
			if err != nil {
				getLogger().Println(err)
				return
			}
		} else {
			err = dakaraPutTag(ctx, dakara_tag.ID, audio_tag)
			if err != nil {
				getLogger().Println(err)
				return
			}
		}
	}

	// sync video tags
	for _, video_tag := range VideoTags {
		dakara_tag := dakara_tags[video_tag.Name]
		if dakara_tag == nil {
			err = dakaraAddTag(ctx, video_tag)
			if err != nil {
				getLogger().Println(err)
				return
			}
		} else {
			err = dakaraPutTag(ctx, dakara_tag.ID, video_tag)
			if err != nil {
				getLogger().Println(err)
				return
			}
		}
	}
	logger.Println("Synced tags with dakara")

	// sync karas / songs

	songs, err := dakaraGetSongs(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	dakara_tags, err = dakaraGetTags(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	dakara_artists, err = dakaraGetArtists(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	works, err = dakaraGetWorks(ctx)
	if err != nil {
		getLogger().Println(err)
		return
	}

	logger.Println("syncing new songs")

	new_songs := 0
	updated_songs := 0
	for _, kara := range all_karas {
		song_body, err := createDakaraSongBody(ctx, kara, dakara_tags, dakara_artists, works)
		if err != nil {
			getLogger().Printf("Kara %d: %s\n", kara.ID, err)
			continue
		}
		dakara_song := songs[kara.VideoFilename()]
		if dakara_song == nil {
			err = dakaraAddSong(ctx, song_body)
			if err != nil {
				getLogger().Println(err)
				continue
			}
			new_songs++
		} else if song_body.HasChanged(*dakara_song) {
			// overly spammy in practice
			// logger.Printf("kara changed\n%+v\n%+v\n", dakara_song, song_body)
			err = dakaraUpdateSong(ctx, dakara_song, song_body)
			if err != nil {
				getLogger().Println(err)
				continue
			}
			updated_songs++
		}
	}
	logger.Printf("Created %d songs, Updated %d songs\n", new_songs, updated_songs)

	err = cleanUpDakaraSongs(ctx, songs)
	if err != nil {
		getLogger().Println(err)
	}
	err = cleanUpDakaraWorks(ctx)
	if err != nil {
		getLogger().Println(err)
	}
	err = cleanUpDakaraArtists(ctx)
	if err != nil {
		getLogger().Println(err)
	}

	// cleanUpWorkTypes(worktypes)
}

func dakaraSongEndpoint(dakara_song_id int) string {
	return fmt.Sprintf("/api/library/songs/%d/", dakara_song_id)
}

func dakaraUpdateSong(ctx context.Context, dakara_song *DakaraSong, song_body *DakaraSongBody) error {
	// TODO: try to skip the update if unnecessary

	path := dakaraSongEndpoint(dakara_song.ID)
	resp, err := dakaraPut(ctx, path, song_body)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	return nil
}

func cleanUpDakaraSongs(ctx context.Context, songs map[string]*DakaraSong) error {
	db := GetDB(ctx)

	deleted_songs := 0
	for _, song := range songs {
		id_str, _, _ := strings.Cut(song.Filename, ".")
		id, err := strconv.Atoi(id_str)
		if err != nil {
			// not our song, probably
			err = deleteDakaraSong(ctx, song)
			if err != nil {
				return err
			}
			deleted_songs++
			continue
		}

		kara := &KaraInfoDB{}
		err = db.Scopes(UploadedKaras).First(kara, id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// We don't know this karaoke (deleted or never existed)
				err = deleteDakaraSong(ctx, song)
				if err != nil {
					return err
				}
				deleted_songs++
			} else {
				return err
			}
		}
	}
	getLogger().Printf("deleted %d songs in dakara\n", deleted_songs)

	return nil
}

func deleteDakaraSong(ctx context.Context, song *DakaraSong) error {
	path := dakaraSongEndpoint(song.ID)
	resp, err := dakaraDelete(ctx, path)
	defer Closer(resp.Body)
	return err
}

func cleanUpDakaraWorks(ctx context.Context) error {
	path := "/api/library/works/prune/"
	resp, err := dakaraDelete(ctx, path)
	defer Closer(resp.Body)
	return err
}

func cleanUpDakaraArtists(ctx context.Context) error {
	path := "/api/library/artists/prune/"
	resp, err := dakaraDelete(ctx, path)
	defer Closer(resp.Body)
	return err
}

func createDakaraSongBody(ctx context.Context, kara KaraInfoDB, dakara_tags map[string]*DakaraTag, dakara_artists map[string]*DakaraArtist, dakara_works map[string]map[string]*DakaraWork) (*DakaraSongBody, error) {
	audio_tags, err := kara.getAudioTags()
	if err != nil {
		return nil, err
	}
	audio_tags = dakaraFilterAudioTags(audio_tags)

	video_tags, err := kara.getVideoTags()
	if err != nil {
		return nil, err
	}

	n_tags := len(audio_tags) + len(video_tags)
	tags := make([]DakaraTag, n_tags)

	for i, tag := range audio_tags {
		tags[i] = *dakara_tags[tag.Name]
	}

	for i, tag := range video_tags {
		tags[i+len(audio_tags)] = *dakara_tags[tag.Name]
	}

	artists := make([]DakaraArtist, len(kara.Artists))
	for i, artist := range kara.Artists {
		artists[i] = *dakara_artists[artist.Name]
	}

	works := make([]DakaraSongWork, 0)

	if kara.SourceMedia != nil {
		dakara_worktype := dakara_works[strings.ToLower(kara.SourceMedia.Type)]
		dakara_work := dakara_worktype[kara.SourceMedia.Name]

		if dakara_work == nil {
			return nil, fmt.Errorf("could not find source media for: %+v", kara)
		}
		link_type_number := &kara.SongOrder
		if kara.SongOrder == 0 {
			link_type_number = nil
		}
		linktype := getWorkLinkType(kara)
		if linktype != "" {
			works = append(works, DakaraSongWork{
				Work:           *dakara_work,
				LinkType:       linktype,
				LinkTypeNumber: link_type_number,
			})
		}
	}
	// NOTE: kara.Media is not usable because we don't know what link_type should be

	comment := kara.Comment
	if len(kara.Comment) > 255 {
		getLogger().Printf("kara %d: comment is %d chars long", kara.ID, len(kara.Comment))
		comment = kara.Comment[:255]
	}

	lyrics, err := GetKaraLyrics(ctx, kara)
	if err != nil {
		return nil, err
	}

	var version string
	if kara.Version != "" && kara.Language != "" {
		version = fmt.Sprintf("%s, %s", strings.Trim(kara.Version, " \n"), strings.Trim(kara.Language, " \n"))
	} else if kara.Language != "" {
		version = strings.Trim(kara.Language, " \n")
	} else {
		version = strings.Trim(kara.Version, " \n")
	}

	return &DakaraSongBody{
		Title:           kara.Title,
		Filename:        kara.VideoFilename(),
		Duration:        kara.Duration,
		Directory:       "",
		Version:         version,
		Detail:          strings.Trim(comment, " \n"),
		DetailVideo:     "",
		Tags:            tags,
		Artists:         artists,
		Works:           works,
		Lyrics:          lyrics,
		HasInstrumental: false,
	}, nil
}

func getWorkLinkType(kara KaraInfoDB) string {
	for _, audio_tag := range kara.AudioTags {
		switch audio_tag.ID {
		case "OP":
			return "OP"
		case "ED":
			return "ED"
		case "IS":
			return "IS"
		case "INS":
			return "IN"
		}
	}

	return ""
}

func StartDakaraSync(ctx context.Context, input *struct{}) (*struct{}, error) {
	if !getCurrentUser(ctx).Admin {
		return nil, huma.Error403Forbidden("endpoint reserved to adminitrators")
	}

	SyncDakaraNotify()
	return &struct{}{}, nil
}
