package server

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MediaType struct {
	ID       string // used in the database/API
	Name     string // user visible name
	IconName string // font-awesome icon name
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
}

type VideoTag TagType
type AudioTag TagType

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
	ID              string `gorm:"primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Admin           bool
	TimingProfileID uint
	TimingProfile   TimingAuthor `gorm:"foreignKey:TimingProfileID;references:ID"`
}

type TimingAuthor struct {
	gorm.Model
	Name string
}

type Scopes struct {
	Kara bool `json:"kara"`
	User bool `json:"user"`
}

var AllScopes = Scopes{Kara: true, User: true}

func (scopes Scopes) HasScope(scope string) bool {
	if scope == "kara" {
		return scopes.Kara
	}
	if scope == "user" {
		return scopes.User
	}

	panic("unknown scope " + scope)
}

type Token struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    string
	User      User `gorm:"foreignKey:UserID;references:ID"`
	Admin     bool
	ReadOnly  bool
	Scopes
}

// Artists

type Artist struct {
	gorm.Model
	Name            string           `gorm:"uniqueIndex:idx_artist_name"`
	AdditionalNames []AdditionalName `gorm:"many2many:artists_additional_name"`
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
	Name            string           `json:"name" example:"Shinseiki Evangelion" gorm:"uniqueIndex:idx_media_name_type"`
	Type            string           `json:"media_type" example:"ANIME" gorm:"uniqueIndex:idx_media_name_type"`
	AdditionalNames []AdditionalName `json:"additional_name" gorm:"many2many:media_additional_name"`
}

// Video tags
var VideoTags = []VideoTag{
	{ID: "FANMADE", Name: "Fanmade", Hue: 140},
	{ID: "STREAM", Name: "Stream", Hue: 160},
	{ID: "CONCERT", Name: "Concert", Hue: 260},
	{ID: "AD", Name: "Advertisement", Hue: 120},
	{ID: "TRAILER", Name: "Trailer", Hue: 100},
	{ID: "NSFW", Name: "Not Safe For Work", Hue: 0},
	{ID: "SPOILER", Name: "Spoiler", Hue: 20},
	{ID: "MV", Name: "Music Video", Hue: 120},
}

// Audio tags
var AudioTags = []AudioTag{
	{ID: "OP", Name: "Opening", Hue: 280},
	{ID: "ED", Name: "Ending", Hue: 280},
	{ID: "INS", Name: "Insert", Hue: 280},
	{ID: "IS", Name: "Image Song", Hue: 280},
	{ID: "LIVE", Name: "Live", Hue: 240},
	{ID: "REMIX", Name: "Remix/Cover", Hue: 220},
}

type AdditionalName struct {
	gorm.Model
	Name string
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
	InstrumentalUploaded bool
	InstrumentalModTime  time.Time
	SubtitlesUploaded    bool
	SubtitlesModTime     time.Time
	Hardsubbed           bool
	Duration             int32
	// date of the first upload of the sub file
	KaraokeCreationTime time.Time
}

func NewUploadInfo() UploadInfo {
	return UploadInfo{
		VideoUploaded:        false,
		InstrumentalUploaded: false,
		SubtitlesUploaded:    false,
		Hardsubbed:           false,
		VideoModTime:         time.Unix(0, 0),
		InstrumentalModTime:  time.Unix(0, 0),
		SubtitlesModTime:     time.Unix(0, 0),
		KaraokeCreationTime:  time.Unix(0, 0),
	}
}

type KaraInfoDB struct {
	gorm.Model
	Authors       []TimingAuthor `gorm:"many2many:kara_authors_tags"`
	Artists       []Artist       `gorm:"many2many:kara_artist_tags"`
	VideoTags     []VideoTagDB   `gorm:"many2many:kara_video_tags"`
	AudioTags     []AudioTagDB   `gorm:"many2many:kara_audio_tags"`
	SourceMediaID uint
	SourceMedia   MediaDB   `gorm:"foreignKey:SourceMediaID;references:ID"`
	Medias        []MediaDB `gorm:"many2many:kara_media_tags"`
	Title         string
	ExtraTitles   []AdditionalName `gorm:"many2many:kara_info_additional_name"`
	Version       string
	Comment       string
	SongOrder     uint
	Language      string
	UploadInfo
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

func init_model(db *gorm.DB) {
	db.AutoMigrate(&KaraInfoDB{}, &User{}, &Token{}, &MediaDB{}, &Artist{}, &Font{})
}

func createAdditionalNames(names []string) []AdditionalName {
	additional_names := make([]AdditionalName, len(names))
	for i, name := range names {
		additional_names[i] = AdditionalName{Name: name}
	}

	return additional_names
}
