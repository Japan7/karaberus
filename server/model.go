package server

import (
	"time"

	"gorm.io/gorm"
)

type KaraberusType struct {
	ID   string // used in the database/API
	Name string // user visible name
}

type MediaType KaraberusType
type VideoTag KaraberusType
type AudioTag KaraberusType

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
	{ID: "ANIME", Name: "Anime"},
	{ID: "GAME", Name: "Game"},
	{ID: "LIVE", Name: "Live"},
	{ID: "CARTOON", Name: "Cartoon"},
}

type MediaDB struct {
	gorm.Model
	Name            string           `json:"name" example:"Shinseiki Evangelion" gorm:"uniqueIndex:idx_media_name_type"`
	Type            string           `json:"media_type" example:"ANIME" gorm:"uniqueIndex:idx_media_name_type"`
	AdditionalNames []AdditionalName `json:"additional_name" gorm:"many2many:media_additional_name"`
}

// Video tags
var VideoTags []VideoTag = []VideoTag{
	{ID: "FANMADE", Name: "Fanmade"},
	{ID: "STREAM", Name: "Stream"},
	{ID: "CONCERT", Name: "Concert"},
	{ID: "AD", Name: "Advertisement"},
	{ID: "TRAILER", Name: "Trailer"},
	{ID: "NSFW", Name: "Not Safe For Work"},
	{ID: "SPOILER", Name: "Spoiler"},
	{ID: "MV", Name: "Music Video"},
}

// Audio tags
var AudioTags []AudioTag = []AudioTag{
	{ID: "OP", Name: "Opening"},
	{ID: "ED", Name: "Ending"},
	{ID: "INS", Name: "Insert"},
	{ID: "IS", Name: "Image Song"},
	{ID: "LIVE", Name: "Live"},
	{ID: "REMIX", Name: "Remix/Cover"},
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
