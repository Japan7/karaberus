package server

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaType struct {
	ID       string // used in the database/API
	Name     string // user visible name
	IconName string // font-awesome icon name
}

type HasName interface {
	getName() string // user visible name
}

type TagInterface interface {
	getID() string   // used in the database/API
	getName() string // user visible name
	// Hue in deg
	getHue() uint
}

type TagType struct {
	ID   string // used in the database/API
	Name string // user visible name
	// Hue in deg
	Hue uint
	// Mugen tag ID (optional)
	MugenTags []string
}

type AudioTag struct {
	ID   string // used in the database/API
	Name string // user visible name
	// Hue in deg
	Hue uint
	// Mugen tag ID (optional)
	MugenTags []string
	// true if this type can have a song order
	HasSongOrder bool
}

type VideoTag TagType

func (t AudioTag) getID() string {
	return t.ID
}

func (t AudioTag) getName() string {
	return t.Name
}

func (t AudioTag) getHue() uint {
	return t.Hue
}

func (t VideoTag) getID() string {
	return t.ID
}

func (t VideoTag) getName() string {
	return t.Name
}

func (t VideoTag) getHue() uint {
	return t.Hue
}

// Users
type User struct {
	ID              string         `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Admin           bool           `json:"admin"`
	TimingProfileID *uint          `json:"timing_profile_id"`
	TimingProfile   *TimingAuthor  `gorm:"foreignKey:TimingProfileID;references:ID" json:"timing_profile"`
}

type TimingAuthor struct {
	gorm.Model
	Name    string     `gorm:"uniqueIndex:idx_timing_author_name"`
	MugenID *uuid.UUID `gorm:"uniqueIndex:idx_timing_author_mugen_id"`
}

func (name *TimingAuthor) BeforeSave(tx *gorm.DB) error {
	name.Name = trimWhitespace(name.Name)
	return nil
}

type Scopes struct {
	Kara   bool `json:"kara"`
	KaraRO bool `json:"kara_ro"`
	User   bool `json:"user"`
}

var AllScopes = Scopes{Kara: true, KaraRO: true, User: true}

func (scopes Scopes) HasScope(scope string) bool {
	if scope == "kara" {
		return scopes.Kara
	}
	if scope == "kara_ro" {
		return scopes.Kara || scopes.KaraRO
	}
	if scope == "user" {
		return scopes.User
	}

	panic("unknown scope " + scope)
}

type TokenV2 struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `gorm:"uniqueIndex:idx_token" json:"token"`
	UserID    string    `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Name      string    `json:"name"`
	Scopes    Scopes    `gorm:"embedded" json:"scopes"`
}

func (name *TokenV2) BeforeSave(tx *gorm.DB) error {
	name.Name = trimWhitespace(name.Name)
	return nil
}

// Artists

type Artist struct {
	gorm.Model
	Name            string           `gorm:"uniqueIndex:idx_artist_name_v2,where:current_artist_id IS NULL AND deleted_at IS NULL"`
	AdditionalNames []AdditionalName `gorm:"many2many:artists_additional_name"`
	CurrentArtistID *uint
	CurrentArtist   *Artist
	Editor
}

func CurrentArtists(tx *gorm.DB) *gorm.DB {
	return tx.Where("current_artist_id IS NULL")
}

func (a *Artist) AfterUpdate(tx *gorm.DB) error {
	if a.CurrentArtistID == nil {
		SyncDakaraNotify()
	}
	return nil
}

func (a *Artist) BeforeUpdate(tx *gorm.DB) error {
	if isAssociationsUpdate(tx) {
		return nil
	}
	orig_artist := &Artist{}
	err := tx.First(orig_artist, a.ID).Error
	if err != nil {
		return err
	}

	// create historic entry with the previous value
	orig_artist.ID = 0
	orig_artist.CurrentArtist = a
	err = tx.Create(orig_artist).Error

	return err
}

func (a *Artist) BeforeSave(tx *gorm.DB) error {
	a.Name = trimWhitespace(a.Name)

	// set editor for this new update
	a.EditorUser = getCurrentUser(tx.Statement.Context)
	if a.EditorUser == nil {
		a.EditorUserID = nil
	}
	return nil
}

// Media types
var MediaTypes []MediaType = []MediaType{
	{ID: "ANIME", Name: "Anime", IconName: "tv"},
	{ID: "GAME", Name: "Game", IconName: "gamepad"},
	{ID: "LIVE", Name: "Live action", IconName: "film"},
	{ID: "CARTOON", Name: "Cartoon", IconName: "globe"},
}

type MediaDB struct {
	gorm.Model
	Name            string           `json:"name" example:"Shinseiki Evangelion" gorm:"uniqueIndex:idx_media_name_type_v2,where:current_media_id IS NULL AND deleted_at IS NULL"`
	Type            string           `json:"media_type" example:"ANIME" gorm:"uniqueIndex:idx_media_name_type_v2,where:current_media_id IS NULL AND deleted_at IS NULL"`
	AdditionalNames []AdditionalName `json:"additional_name" gorm:"many2many:media_additional_name"`
	CurrentMediaID  *uint
	CurrentMedia    *MediaDB
	Editor
}

func CurrentMedias(tx *gorm.DB) *gorm.DB {
	return tx.Where("current_media_id IS NULL")
}

func (m *MediaDB) AfterUpdate(tx *gorm.DB) error {
	if m.CurrentMediaID == nil {
		SyncDakaraNotify()
	}
	return nil
}

func (m *MediaDB) BeforeUpdate(tx *gorm.DB) error {
	if isAssociationsUpdate(tx) {
		return nil
	}
	orig_media := &MediaDB{}
	err := tx.First(orig_media, m.ID).Error
	if err != nil {
		return err
	}

	// create historic entry with the previous value
	orig_media.ID = 0
	orig_media.CurrentMedia = m
	err = tx.Create(orig_media).Error

	return err
}

func (m *MediaDB) BeforeSave(tx *gorm.DB) error {
	m.Name = trimWhitespace(m.Name)

	// set editor for this new update
	m.EditorUser = getCurrentUser(tx.Statement.Context)
	if m.EditorUser == nil {
		m.EditorUserID = nil
	}
	return nil
}

// Video tags
var VideoTags = []VideoTag{
	// Mugen tag is AMV, not exhaustive but best we can do
	{ID: "FANMADE", Name: "Fanmade", Hue: 140, MugenTags: []string{"a6c79ce5-89ee-4d50-afe8-3abd7317f6c2"}},
	{ID: "STREAM", Name: "Stream", Hue: 160, MugenTags: []string{"55ce3d79-dcc2-453c-b00a-60ce0c1eba1c"}},
	{ID: "CONCERT", Name: "Concert", Hue: 260, MugenTags: []string{"a0167949-580c-4de3-bf13-497e462e02f3"}},
	{ID: "AD", Name: "Advertisement", Hue: 120, MugenTags: []string{"2ddb5358-e674-46fa-a6e1-7f5c5d56f8fa"}},
	{ID: "NSFW", Name: "Not Safe For Work", Hue: 0, MugenTags: []string{"e82ce681-6d7b-4fb6-abe4-daa8aaa9bbf9"}},
	{ID: "SPOILER", Name: "Spoiler", Hue: 20, MugenTags: []string{"24371984-5e4c-4485-a937-fb0c480ca23b"}},
	{ID: "EPILEPSY", Name: "Epilepsy", Hue: 0, MugenTags: []string{"51288600-29e0-4e41-a42b-77f0498e5691"}},
	{ID: "MV", Name: "Music Video", Hue: 120, MugenTags: []string{"7be1b15c-cff8-4b37-a649-5c90f3d569a9"}},
}

// Audio tags
var AudioTags = []AudioTag{
	{ID: "OP", Name: "Opening", Hue: 280, MugenTags: []string{"f02ad9b3-0bd9-4aad-85b3-9976739ba0e4"}, HasSongOrder: true},
	{ID: "ED", Name: "Ending", Hue: 280, MugenTags: []string{"38c77c56-2b95-4040-b676-0994a8cb0597"}, HasSongOrder: true},
	{ID: "INS", Name: "Insert", Hue: 280, MugenTags: []string{"5e5250d9-351a-4a82-98eb-55db50ad8962"}, HasSongOrder: true},
	{ID: "IS", Name: "Image Song", Hue: 280, MugenTags: []string{"10a1ad3e-a05c-4f5c-84b6-f491e3e3a92e"}, HasSongOrder: true},
	// Mugen tags are Concert and Streaming
	{ID: "LIVE", Name: "Live performance", Hue: 240, MugenTags: []string{"a0167949-580c-4de3-bf13-497e462e02f3", "55ce3d79-dcc2-453c-b00a-60ce0c1eba1c"}},
	// Mugen tags are version tags: Cover, Metal
	{ID: "REMIX", Name: "Remix/Cover", Hue: 220, MugenTags: []string{"03e1e1d2-8641-47b7-bbcb-39a3df9ff21c", "188a5c46-63ff-4e9f-89e4-763468b6ea4a"}},
}

type AdditionalName struct {
	gorm.Model
	Name string
}

func trimWhitespace(s string) string {
	return strings.Trim(s, " \t\n")
}

func (name *AdditionalName) BeforeSave(tx *gorm.DB) error {
	name.Name = trimWhitespace(name.Name)
	return nil
}

type VideoTagDB struct {
	ID string
}

type AudioTagDB struct {
	ID string
}

type UploadInfo struct {
	VideoUploaded        bool
	VideoModTime         time.Time
	VideoSize            int64
	VideoCRC32           uint32
	InstrumentalUploaded bool
	InstrumentalModTime  time.Time
	InstrumentalSize     int64
	InstrumentalCRC32    uint32
	SubtitlesUploaded    bool
	SubtitlesModTime     time.Time
	SubtitlesSize        int64
	SubtitlesCRC32       uint32
	Hardsubbed           bool
	Duration             int32
	// date of the first upload of the sub file
	KaraokeCreationTime time.Time
}

type Editor struct {
	EditorUserID *string
	EditorUser   *User `gorm:"foreignKey:EditorUserID;references:ID"`
}

type KaraInfoDB struct {
	gorm.Model
	Authors       []TimingAuthor `gorm:"many2many:kara_authors_tags"`
	Artists       []Artist       `gorm:"many2many:kara_artist_tags"`
	VideoTags     []VideoTagDB   `gorm:"many2many:kara_video_tags"`
	AudioTags     []AudioTagDB   `gorm:"many2many:kara_audio_tags"`
	SourceMediaID *uint
	SourceMedia   *MediaDB  `gorm:"foreignKey:SourceMediaID;references:ID"`
	Medias        []MediaDB `gorm:"many2many:kara_media_tags"`
	Title         string
	ExtraTitles   []AdditionalName `gorm:"many2many:kara_info_additional_name"`
	Private       bool
	Version       string
	Comment       string
	SongOrder     uint
	Language      string
	UploadInfo
	// Can't be set by users
	CurrentKaraInfoID *uint
	CurrentKaraInfo   *KaraInfoDB
	Editor
}

// try not to go over 255 chars
func (k KaraInfoDB) FriendlyName() string {
	parts := []string{}

	if k.SourceMedia != nil {
		parts = append(parts, k.SourceMedia.Name)
	}

	if len(k.Artists) > 0 {
		artists := k.Artists[0].Name
		for _, artist := range k.Artists[1:] {
			// try not to use a name that is too long
			if len(artists)+len(artist.Name) > 100 {
				break
			}
			artists += ", " + artist.Name
		}
		parts = append(parts, artists)
	}

	parts = append(parts, k.Title)

	return strings.Join(parts, " – ")
}

// Filter out historic entries
func CurrentKaras(tx *gorm.DB) *gorm.DB {
	return tx.Where("current_kara_info_id IS NULL")
}

func isNewKaraUpdate(tx *gorm.DB, kara *KaraInfoDB) bool {
	// check for unix time 0 is for older karaokes, because we also used
	// that at some point
	return kara.SubtitlesUploaded && kara.VideoUploaded &&
		kara.KaraokeCreationTime.Before(time.Unix(1, 0)) &&
		tx.Statement.Context.Value(PossiblyNewKaraUpdate{}) != nil
}

type PossiblyNewKaraUpdate struct{}

func WithPossiblyNewKaraUpdate(tx *gorm.DB) *gorm.DB {
	return tx.WithContext(context.WithValue(tx.Statement.Context, PossiblyNewKaraUpdate{}, true))
}

type UpdateAssociations struct{}

func WithAssociationsUpdate(tx *gorm.DB) *gorm.DB {
	return tx.WithContext(context.WithValue(tx.Statement.Context, UpdateAssociations{}, true))
}

func isAssociationsUpdate(tx *gorm.DB) bool {
	return tx.Statement.Context.Value(UpdateAssociations{}) != nil
}

var gitlabUpdateMutex = sync.Mutex{}

func UploadHookGitlab(tx *gorm.DB, ki *KaraInfoDB) error {
	if CONFIG.Mugen.Gitlab.IsSetup() {
		if ki.Private {
			return nil
		}

		gitlabUpdateMutex.Lock()
		defer gitlabUpdateMutex.Unlock()

		// check if kara is an import
		mugen_import := &MugenImport{}
		err := tx.Where(&MugenImport{KaraID: ki.ID}).First(mugen_import).Error
		if err == nil {
			// kara was imported
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// check if kara is already exported
		mugen_export := &MugenExport{}
		err = tx.Where(&MugenExport{KaraID: ki.ID}).First(&mugen_export).Error
		if err == nil {
			if mugen_export.KaraID != 0 && mugen_export.GitlabIssue != -1 {
				// already exported, update issue
				err = updateGitlabIssue(tx.Statement.Context, tx, *ki, mugen_export)
				if err != nil {
					getLogger().Println(err)
				}
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		err = createGitlabIssue(tx.Statement.Context, tx, *ki, mugen_export)
		if err != nil {
			getLogger().Println(err)
		}
	}

	return nil
}

func (ki *KaraInfoDB) AfterUpdate(tx *gorm.DB) error {
	// update kara just in case
	err := tx.First(&ki).Error
	if err != nil {
		return err
	}

	if ki.CurrentKaraInfoID == nil {
		if CONFIG.Dakara.BaseURL != "" && ki.VideoUploaded && ki.SubtitlesUploaded {
			SyncDakaraNotify()
		}

		err = UploadHookGitlab(tx, ki)
		if err != nil {
			return err
		}

		if isNewKaraUpdate(tx, ki) {

			// ignore imported karas
			mugen_import := &MugenImport{}
			err := tx.Where(&MugenImport{KaraID: ki.ID}).First(mugen_import).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				go PostWebhooks(*ki)
			} else if err != nil {
				return err
			}

		}
	}
	return nil
}

func (ki *KaraInfoDB) BeforeUpdate(tx *gorm.DB) error {
	if isAssociationsUpdate(tx) {
		return nil
	}
	orig_kara_info := &KaraInfoDB{}
	err := tx.First(orig_kara_info, ki.ID).Error
	if err != nil {
		return err
	}

	if isNewKaraUpdate(tx, ki) {
		ki.KaraokeCreationTime = time.Now()
	}

	// create historic entry with the current value
	orig_kara_info.ID = 0
	orig_kara_info.CurrentKaraInfo = ki
	err = tx.Create(orig_kara_info).Error

	getLogger().Printf("Updating kara %d", ki.ID)
	return err
}

func (ki *KaraInfoDB) BeforeSave(tx *gorm.DB) error {
	ki.Version = trimWhitespace(ki.Version)
	ki.Comment = trimWhitespace(ki.Comment)
	ki.Title = trimWhitespace(ki.Title)
	ki.Language = trimWhitespace(ki.Language)

	// set editor for this new version
	ki.EditorUser = getCurrentUser(tx.Statement.Context)
	if ki.EditorUser == nil {
		ki.EditorUserID = nil
	}

	if ki.SubtitlesUploaded && ki.Hardsubbed {
		ki.Hardsubbed = false
	}

	return nil
}

type MugenImport struct {
	MugenKID uuid.UUID `gorm:"primarykey"`
	KaraID   uint
	Kara     KaraInfoDB `gorm:"foreignKey:KaraID;references:ID;constraint:OnDelete:CASCADE"`
}

type MugenExport struct {
	KaraID         uint       `gorm:"primarykey" json:"kid"`
	Kara           KaraInfoDB `gorm:"foreignKey:KaraID;references:ID;constraint:OnDelete:CASCADE" json:"kara"`
	GitlabIssue    int        `json:"gitlab_issue"`
	GitlabIssueIID int        `json:"gitlab_issue_iid"`
	Closed         bool       `json:"closed"`
}

func (k KaraInfoDB) getAudioTags() ([]AudioTag, error) {
	audio_tags := make([]AudioTag, len(k.AudioTags))

	for i, tag := range k.AudioTags {
		audio_tag, err := getAudioTag(tag.ID)
		if err != nil {
			return nil, err
		}
		audio_tags[i] = *audio_tag
	}

	return audio_tags, nil
}

func (k KaraInfoDB) getVideoTags() ([]VideoTag, error) {
	video_tags := make([]VideoTag, len(k.VideoTags))

	for i, tag := range k.VideoTags {
		video_tag, err := getVideoTag(tag.ID)
		if err != nil {
			return nil, err
		}
		video_tags[i] = *video_tag
	}

	return video_tags, nil
}

func (k KaraInfoDB) VideoFilename() string {
	return fmt.Sprintf("%d.mkv", k.ID)
}

func (k KaraInfoDB) AudioFilename() string {
	return fmt.Sprintf("%d.mka", k.ID)
}

func (k KaraInfoDB) SubsFilename() string {
	return fmt.Sprintf("%d.ass", k.ID)
}

type Font struct {
	gorm.Model
	Name       string
	UploadedAt time.Time
	// TODO: font properties (family name, weight, ...)
}

type OAuthToken struct {
	Server       string `gorm:"primarykey"`
	ClientID     string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func init_model(db *gorm.DB) {
	err := db.AutoMigrate(
		&User{},
		&TimingAuthor{},
		&TokenV2{},
		&Artist{},
		&MediaDB{},
		&AdditionalName{},
		&VideoTagDB{},
		&AudioTagDB{},
		&KaraInfoDB{},
		&MugenImport{},
		&MugenExport{},
		&Font{},
		&OAuthToken{},
	)
	if err != nil {
		panic(err)
	}

	// https://github.com/Japan7/karaberus/pull/73
	// drop previous indexes
	if db.Migrator().HasIndex(&Artist{}, "idx_artist_name") {
		err = db.Migrator().DropIndex(&Artist{}, "idx_artist_name")
		if err != nil {
			panic(err)
		}
	}
	if db.Migrator().HasIndex(&MediaDB{}, "idx_media_name_type") {
		err = db.Migrator().DropIndex(&MediaDB{}, "idx_media_name_type")
		if err != nil {
			panic(err)
		}
	}

	// some karas might not have the right creation date at some point...
	fixCreationTime(db)

	// set size and crc32 for files uploaded before they were introduced
	ctx := db.Statement.Context
	if isKaraberusInit(ctx) {
		err = exportRemainingKaras(ctx, db)
		if err != nil {
			panic(err)
		}
		go initSizeCRC(db)
	}
}

func fixCreationTime(db *gorm.DB) {
	var karas []KaraInfoDB
	err := db.Scopes(CurrentKaras).Find(&karas).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if err != nil {
		panic(err)
	}

	for _, kara := range karas {
		if kara.KaraokeCreationTime.Before(time.Unix(1, 0)) {
			var kara_time time.Time
			if kara.SubtitlesUploaded && kara.VideoUploaded && kara.SubtitlesModTime.After(time.Unix(1, 0)) {
				kara_time = kara.SubtitlesModTime
			} else if !kara.KaraokeCreationTime.IsZero() {
				// reset creation time to zero if time is Unix Epoch
				// could happen for karaoke that are not uploaded (yet)
				kara_time = time.Time{}
			} else {
				continue
			}

			getLogger().Printf("setting kara %d creation time\n", kara.ID)
			kara.KaraokeCreationTime = kara.SubtitlesModTime
			err = db.Model(&kara).Update(
				"karaoke_creation_time", kara_time,
			).Error
			if err != nil {
				panic(err)
			}
		}
	}
}

func initSizeCRC(db *gorm.DB) {
	var karas []KaraInfoDB
	err := db.Scopes(CurrentKaras).Find(&karas).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if err != nil {
		panic(err)
	}

	// could be done concurrently but also probably doesn’t matter that much
	for _, kara := range karas {
		changed := false
		if kara.VideoUploaded && kara.VideoSize <= 0 {
			getLogger().Printf("%d: calculating video crc/size\n", kara.ID)
			obj, err := GetKaraObject(db.Statement.Context, kara, "video")
			if err != nil {
				panic(err)
			}
			hasher := crc32.NewIEEE()
			kara.VideoSize, err = io.Copy(hasher, obj)
			if err != nil {
				panic(err)
			}
			kara.VideoCRC32 = hasher.Sum32()
			changed = true
		}
		if kara.InstrumentalUploaded && kara.InstrumentalSize <= 0 {
			getLogger().Printf("%d: calculating inst crc/size\n", kara.ID)
			obj, err := GetKaraObject(db.Statement.Context, kara, "inst")
			if err != nil {
				panic(err)
			}
			hasher := crc32.NewIEEE()
			kara.InstrumentalSize, err = io.Copy(hasher, obj)
			if err != nil {
				panic(err)
			}
			kara.InstrumentalCRC32 = hasher.Sum32()
			changed = true
		}
		if kara.SubtitlesUploaded && kara.SubtitlesSize <= 0 {
			getLogger().Printf("%d: calculating sub crc/size\n", kara.ID)
			obj, err := GetKaraObject(db.Statement.Context, kara, "sub")
			if err != nil {
				panic(err)
			}
			hasher := crc32.NewIEEE()
			kara.SubtitlesSize, err = io.Copy(hasher, obj)
			if err != nil {
				panic(err)
			}
			kara.SubtitlesCRC32 = hasher.Sum32()
			changed = true
		}

		if changed {
			err := db.Save(&kara).Error
			if err != nil {
				panic(err)
			}
		}
	}
}

func createAdditionalNames(names []string) []AdditionalName {
	additional_names := make([]AdditionalName, len(names))
	for i, name := range names {
		additional_names[i] = AdditionalName{Name: name}
	}

	return additional_names
}
