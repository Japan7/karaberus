package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/Japan7/karaberus/server/clients/mugen"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var VIDEO_GAME_MUGEN_TAG_ID = "dbedd6b3-d125-4cd8-aa32-c4175e4ca3a3"
var ANIMATION_MUGEN_TAG_ID = "0377db02-3af6-43b8-9b08-c759df3d25c3"
var WEST_MUGEN_TAG_ID = "efe171c0-e8a1-4d03-98c0-60ecf741ad52"
var COVER_MUGEN_TAG_ID = "03e1e1d2-8641-47b7-bbcb-39a3df9ff21c"

var ANIME_TYPE = getMediaType("ANIME")
var GAME_TYPE = getMediaType("GAME")
var CARTOON_TYPE = getMediaType("CARTOON")
var LIVE_TYPE = getMediaType("LIVE")

func getMugenMedia(tx *gorm.DB, tag mugen.MugenTag, origins []mugen.MugenTag, collections []mugen.MugenTag, media *MediaDB) error {
	var media_type *MediaType = nil
	// we have to find the media type
	// anime is the easiest to find because we can find the anilist ID
	if tag.ExternalDatabaseIDs.Anilist != nil {
		media_type = &ANIME_TYPE
	}
	// if origins contains video game tag, then we guess that it is a video game
	is_animation := false
	if media_type == nil {
		for _, origin := range origins {
			if origin.TID.String() == VIDEO_GAME_MUGEN_TAG_ID {
				media_type = &GAME_TYPE
			}
			if origin.TID.String() == ANIMATION_MUGEN_TAG_ID {
				is_animation = true
			}
		}
	}
	// if origins contains animation tag and kara is in the "West" collection assume cartoon
	if is_animation && media_type == nil {
		for _, collection := range collections {
			if collection.TID.String() == WEST_MUGEN_TAG_ID {
				media_type = &CARTOON_TYPE
			}
		}
	}
	// if we still didn't find it and it's not animated could it be live action, perchance?
	if media_type == nil && !is_animation {
		media_type = &LIVE_TYPE
	}
	if media_type == nil {
		return errors.New("could not guess media type for media " + tag.Name)
	}

	additional_names := []string{}
	err := findMedia(tx, []string{tag.Name}, media)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = createMedia(tx, media, &MediaInfo{tag.Name, media_type.ID, additional_names})
	}
	if err != nil {
		return err
	}
	return nil
}

func getMugenArtist(tx *gorm.DB, mugen_artist mugen.MugenTag, karaberus_artist *Artist) error {
	artistNames := []string{mugen_artist.Name}
	artistNames = append(artistNames, mugen_artist.Aliases...)

	err := findArtist(tx, artistNames, karaberus_artist)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = createArtist(tx, karaberus_artist, &ArtistInfo{mugen_artist.Name, mugen_artist.Aliases})
	}
	return err
}

func getMugenTimingAuthor(tx *gorm.DB, mugen_author mugen.MugenTag, author *TimingAuthor) error {
	err := tx.Where(&TimingAuthor{MugenID: &mugen_author.TID}).First(author).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = tx.Where(&TimingAuthor{Name: mugen_author.Name}).First(author).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			author.Name = mugen_author.Name
			author.MugenID = &mugen_author.TID
			err = tx.Create(author).Error
		}
	}
	return err
}

func mugenKaraToKaraInfoDB(tx *gorm.DB, k mugen.Kara, kara_info *KaraInfoDB) error {
	kara_info.Title = k.Titles[k.TitleDefaultLanguage]

	n_titles := len(k.Titles) + len(k.TitleAliases) - 1
	titles := make([]AdditionalName, n_titles)

	i := 0
	for lang, title := range k.Titles {
		if lang == k.TitleDefaultLanguage {
			continue
		}
		titles[i] = AdditionalName{Name: title}
		i++
	}

	for _, title := range k.TitleAliases {
		titles[i] = AdditionalName{Name: title}
		i++
	}
	kara_info.ExtraTitles = titles

	if k.SongOrder == nil {
		kara_info.SongOrder = 0
	} else {
		kara_info.SongOrder = *k.SongOrder
	}

	if len(k.Languages) == 1 {
		kara_info.Language = k.Languages[0].Name
	}

	kara_info.Comment = k.Comment

	if len(k.Series) > 0 {
		if len(k.Series) == 1 {
			source_media := MediaDB{}
			err := getMugenMedia(tx, k.Series[0], k.Origins, k.Collections, &source_media)
			if err != nil {
				return err
			}
			kara_info.SourceMedia = &source_media
		} else {
			kara_info.Medias = make([]MediaDB, len(k.Series))
			for i, series := range k.Series {
				err := getMugenMedia(tx, series, k.Origins, k.Collections, &kara_info.Medias[i])
				if err != nil {
					return err
				}
			}
		}
	}

	// authors
	kara_info.Authors = make([]TimingAuthor, len(k.Authors))

	for i, author := range k.Authors {
		err := getMugenTimingAuthor(tx, author, &kara_info.Authors[i])
		if err != nil {
			return err
		}
	}

	// artists
	artistTags := make([]mugen.MugenTag, 0)
	artistTags = append(artistTags, k.SingerGroups...)
	artistTags = append(artistTags, k.Singers...)
	artistTags = append(artistTags, k.SongWriters...)

	kara_info.Artists = make([]Artist, len(artistTags))

	for i, artist := range artistTags {
		err := getMugenArtist(tx, artist, &kara_info.Artists[i])
		if err != nil {
			return err
		}
	}

	// videotags
	mugenTags := make([]mugen.MugenTag, 0)
	mugenTags = append(mugenTags, k.SongTypes...)
	mugenTags = append(mugenTags, k.Warnings...)

	for _, mugen_tag := range mugenTags {
		for _, video_tag := range VideoTags {
			for _, mapped_tag := range video_tag.MugenTags {
				if mapped_tag == mugen_tag.TID.String() {
					kara_info.VideoTags = append(kara_info.VideoTags, VideoTagDB{ID: video_tag.ID})
				}
			}
		}
	}
	// audiotags
	for _, audio_tag := range AudioTags {
		for _, mugen_tag := range mugenTags {
			for _, mapped_tag := range audio_tag.MugenTags {
				if mapped_tag == mugen_tag.TID.String() {
					kara_info.AudioTags = append(kara_info.AudioTags, AudioTagDB{ID: audio_tag.ID})
				}
			}
		}
	}
	// Version
	versions := []string{}

	for _, version_tag := range k.Versions {
		if version_tag.TID.String() == COVER_MUGEN_TAG_ID {
			continue
		}
		versions = append(versions, version_tag.Name)
	}

	kara_info.Version = strings.Join(versions, " ")

	return nil
}

func reimportMugenKara(ctx context.Context, mugen_import *MugenImport) error {
	getLogger().Println("reimporting ", mugen_import.MugenKID)
	client := mugen.GetClient()
	kara, err := client.GetKara(ctx, mugen_import.MugenKID)
	if err != nil {
		return err
	}

	err = GetDB(context.Background()).Transaction(func(tx *gorm.DB) error {
		kara_info := &mugen_import.Kara
		err = mugenKaraToKaraInfoDB(tx, *kara, kara_info)
		if err != nil {
			return err
		}

		return updateKara(tx, kara_info)
	})

	return err
}

func importMugenKara(ctx context.Context, kid uuid.UUID, mugen_import *MugenImport) error {
	client := mugen.GetClient()
	kara, err := client.GetKara(ctx, kid)
	if err != nil {
		return err
	}

	db_ctx := context.Background()
	db := GetDB(db_ctx)
	getLogger().Printf("Importing kid %s for %s\n", kid, getCurrentUser(ctx).ID)
	err = db.Transaction(func(tx *gorm.DB) error {
		kara_info := KaraInfoDB{}
		err = mugenKaraToKaraInfoDB(tx, *kara, &kara_info)
		if err != nil {
			return err
		}

		err = tx.Create(&kara_info).Error
		if err != nil {
			return err
		}

		mugen_import.MugenKID = kid
		mugen_import.Kara = kara_info
		err = tx.Create(mugen_import).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	go MugenDownload(db_ctx, db, *mugen_import)

	return nil
}

type ImportMugenKaraInput struct {
	Body struct {
		MugenKID uuid.UUID `json:"mugen_kid"`
	}
}

type ImportMugenKaraOutput struct {
	Status int
	Body   struct {
		Import MugenImport `json:"import"`
	}
}

func ImportMugenKara(ctx context.Context, input *ImportMugenKaraInput) (*ImportMugenKaraOutput, error) {
	out := &ImportMugenKaraOutput{}
	err := importMugenKara(ctx, input.Body.MugenKID, &out.Body.Import)
	if err == nil {
		out.Status = 200
	} else {
		err_select := GetDB(ctx).First(&out.Body.Import, input.Body.MugenKID).Error
		if err_select == nil {
			err = nil
			out.Status = 204
		} else {
			out.Status = 500
		}
	}
	return out, err
}

func RefreshMugenImports(ctx context.Context) error {
	mugen_imports := make([]MugenImport, 0)
	db := GetDB(ctx)
	err := db.Preload(clause.Associations).Find(&mugen_imports).Error
	if err != nil {
		return err
	}

	for _, mugen_import := range mugen_imports {
		// kara was deleted ignore
		if mugen_import.Kara.ID == 0 {
			getLogger().Printf("Not updating %s because the kara is not initialized", mugen_import.MugenKID)
			continue
		}
		// karaoke was edited, don't refresh and we don't need to query
		if mugen_import.Kara.EditorUserID != nil {
			getLogger().Printf("Not updating %d because the editor is not NULL", mugen_import.Kara.ID)
			continue
		}

		kara := KaraInfoDB{}
		err = db.Where("editor_user_id IS NOT NULL").Where(&KaraInfoDB{CurrentKaraInfoID: &mugen_import.KaraID}).First(&kara).Error
		if err == nil {
			getLogger().Printf("Not updating %d because it was updated by %s", mugen_import.Kara.ID, *kara.EditorUserID)
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		err = reimportMugenKara(ctx, &mugen_import)
		if err != nil {
			return err
		}
	}

	return nil
}

func RefreshMugen(ctx context.Context, input *struct{}) (*struct{}, error) {
	user := *getCurrentUser(ctx)
	if !user.Admin {
		return nil, huma.Error403Forbidden("This endpoint is reserved to administrators")
	}

	err := RefreshMugenImports(ctx)
	return &struct{}{}, err
}

func SaveMugenResponseToS3(ctx context.Context, tx *gorm.DB, resp *http.Response, kara MugenImport, type_directory string, user_metadata map[string]string) (*CheckKaraOutput, error) {
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d: failed to download, received code %d", kara.MugenKID, resp.StatusCode)
	}
	content_length, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}

	return SaveFileToS3WithMetadata(ctx, tx, resp.Body, &kara.Kara, type_directory, content_length, user_metadata)
}

var MugenDownloadSemaphore = semaphore.NewWeighted(5)

func mugenDownload(ctx context.Context, tx *gorm.DB, mugen_import MugenImport) error {
	if mugen_import.Kara.ID == 0 {
		return errors.New("trying to download a karaoke that has ID 0")
	}
	err := MugenDownloadSemaphore.Acquire(ctx, 1)
	if err != nil {
		return err
	}

	defer func() {
		r := recover()
		if r != nil {
			getLogger().Printf("recovered from panic in SyncMugen: %s\n%s\n", r, string(debug.Stack()))
		}
		MugenDownloadSemaphore.Release(1)
	}()

	mugen_client := mugen.GetClient()
	mugen_kara, err := mugen_client.GetKara(ctx, mugen_import.MugenKID)
	if err != nil {
		return err
	}

	// video
	obj, err := GetKaraObject(ctx, mugen_import.Kara, "video")
	if err != nil {
		return err
	}
	defer Closer(obj)

	should_download_video := !mugen_import.Kara.VideoUploaded

	if !should_download_video {
		stat, err := obj.Stat()
		if err != nil {
			resp := minio.ToErrorResponse(err)
			if resp.Code == "NoSuchKey" {
				should_download_video = true
			} else {
				return resp
			}
		} else {
			// afaik file size is the only possible check (other than downloading on
			// any update of the metadata)
			should_download_video = stat.Size != int64(mugen_kara.MediaSize)
		}
	}

	if should_download_video {
		getLogger().Printf("Downloading %s (%s)", mugen_kara.MediaFile, mugen_kara.KID)
		resp, err := mugen_client.DownloadMedia(ctx, mugen_kara.MediaFile)
		if err != nil {
			return err
		}
		defer Closer(resp.Body)
		_, err = SaveMugenResponseToS3(ctx, tx, resp, mugen_import, "video", nil)
		if err != nil {
			return err
		}
	}

	// sub

	obj, err = GetKaraObject(ctx, mugen_import.Kara, "sub")
	if err != nil {
		return err
	}
	defer Closer(obj)

	should_download_sub := !mugen_import.Kara.SubtitlesUploaded

	if !should_download_sub {
		stat, err := obj.Stat()
		if err != nil {
			resp := minio.ToErrorResponse(err)
			if resp.Code == "NoSuchKey" {
				should_download_sub = true
			} else {
				return resp
			}
		} else {
			should_download_sub = stat.UserMetadata["Mugenchecksum"] != mugen_kara.SubChecksum
		}
	}

	if should_download_sub {
		getLogger().Printf("Downloading %s (%s)", mugen_kara.SubFile, mugen_kara.KID)
		resp, err := mugen_client.DownloadLyrics(ctx, mugen_kara.SubFile)
		if err != nil {
			return err
		}
		defer Closer(resp.Body)
		// we're essentially using the checksum as a version
		user_metadata := map[string]string{"Mugenchecksum": mugen_kara.SubChecksum}
		_, err = SaveMugenResponseToS3(ctx, tx, resp, mugen_import, "sub", user_metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func MugenDownload(ctx context.Context, tx *gorm.DB, mugen_import MugenImport) {
	err := mugenDownload(ctx, tx, mugen_import)
	if err != nil {
		getLogger().Println(err)
	}
}

func SyncMugen(ctx context.Context) {
	mugen_imports := []MugenImport{}
	db := GetDB(ctx)
	err := db.Preload(clause.Associations).Find(&mugen_imports).Error
	if err != nil {
		getLogger().Println(err)
		return
	}

	getLogger().Printf("Syncing %d karaokes from Mugen", len(mugen_imports))

	for _, mugen_import := range mugen_imports {
		err = mugenDownload(ctx, db, mugen_import)
		if err != nil {
			getLogger().Println(err)
		}
	}
}

type GetMugenImportsOutput struct {
	Body struct {
		Imports []MugenImport `json:"imports"`
	}
}

func GetMugenImports(ctx context.Context, input *struct{}) (*GetMugenImportsOutput, error) {
	out := &GetMugenImportsOutput{}
	err := GetDB(ctx).Preload(clause.Associations).Find(&out.Body.Imports).Error
	return out, err
}

type DeleteMugenImportInput struct {
	ID uuid.UUID `path:"id"`
}

type DeleteMugenImportOutput struct {
	Status int
}

func DeleteMugenImport(ctx context.Context, input *DeleteMugenImportInput) (*DeleteMugenImportOutput, error) {
	user := getCurrentUser(ctx)
	if !user.Admin {
		return nil, huma.Error403Forbidden("You must be an administrator to use this endpoint")
	}

	db := GetDB(ctx)
	err := db.Delete(&MugenImport{}, input.ID).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}
	return &DeleteMugenImportOutput{Status: 204}, nil
}
