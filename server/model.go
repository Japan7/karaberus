package server

import "gorm.io/gorm"

type KaraberusType struct {
	ID   string // used in the database/API
	Name string // user visible name
}

type TagType KaraberusType
type MediaType KaraberusType
type VideoTag KaraberusType
type AudioTag KaraberusType

// Users
type User struct {
	gorm.Model
	Admin         bool
	TimingProfile TimingAuthor
}

type TimingAuthor struct {
	gorm.Model
	Name            string
	AdditionalNames []AdditionalName `gorm:"many2many:tags_additional_name"`
}

// Tags

var (
	KaraTagArtist TagType = TagType{ID: "ARTIST", Name: "artist"}
	KaraTagAuthor TagType = TagType{ID: "AUTHOR", Name: "author"}
)

var TagTypes []TagType = []TagType{
	KaraTagArtist,
	KaraTagAuthor,
}

type Tag struct {
	gorm.Model
	Name            string           `gorm:"unique_index:idx_name_type"`
	Type            string           `gorm:"unique_index:idx_name_type"`
	AdditionalNames []AdditionalName `gorm:"many2many:tags_additional_name"`
}

// Media types
var (
	ANIME   MediaType = MediaType{ID: "ANIME", Name: "Anime"}
	GAME    MediaType = MediaType{ID: "GAME", Name: "Game"}
	LIVE    MediaType = MediaType{ID: "LIVE", Name: "Live"}
	CARTOON MediaType = MediaType{ID: "CARTOON", Name: "Cartoon"}
)

type MediaDB struct {
	gorm.Model
	Name           string `json:"name" example:"Shinseiki Evangelion"`
	Type           string `json:"media_type" example:"ANIME"`
	AdditionalName `json:"additional_name"`
}

var MediaTypes []MediaType = []MediaType{ANIME, GAME, LIVE, CARTOON}

// Video tags
var (
	VideoTypeOpening       VideoTag = VideoTag{ID: "OP", Name: "Opening"}
	VideoTypeEnding        VideoTag = VideoTag{ID: "ED", Name: "Ending"}
	VideoTypeInsert        VideoTag = VideoTag{ID: "INSERT", Name: "Insert"}
	VideoTypeFanmade       VideoTag = VideoTag{ID: "FANMADE", Name: "Fanmade"}
	VideoTypeStream        VideoTag = VideoTag{ID: "STREAM", Name: "Stream"}
	VideoTypeConcert       VideoTag = VideoTag{ID: "CONCERT", Name: "Concert"}
	VideoTypeAdvertisement VideoTag = VideoTag{ID: "AD", Name: "Advertisement"}
	VideoTypeTrailer       VideoTag = VideoTag{ID: "TRAILER", Name: "Trailer"}
)

var VideoTags []VideoTag = []VideoTag{
	VideoTypeOpening,
	VideoTypeEnding,
	VideoTypeInsert,
	VideoTypeFanmade,
	VideoTypeStream,
	VideoTypeConcert,
	VideoTypeAdvertisement,
	VideoTypeTrailer,
}

// Audio tags
var (
	AudioTypeOpening AudioTag = AudioTag{ID: "OP", Name: "Opening"}
	AudioTypeEnding  AudioTag = AudioTag{ID: "ED", Name: "Ending"}
	AudioTypeInsert  AudioTag = AudioTag{ID: "INS", Name: "Insert"}
	AudioTypeLive    AudioTag = AudioTag{ID: "LIVE", Name: "Live"}
	AudioTypeCover   AudioTag = AudioTag{ID: "COVER", Name: "Cover"}
)

var AudioTags []AudioTag = []AudioTag{
	AudioTypeOpening,
	AudioTypeEnding,
	AudioTypeInsert,
	AudioTypeLive,
	AudioTypeCover,
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

type KaraInfoDB struct {
	gorm.Model
	Tags        []Tag          `gorm:"many2many:kara_info_tags"`
	VideoTags   []VideoTagDB   `gorm:"many2many:kara_video_tags"`
	AudioTags   []AudioTagDB   `gorm:"many2many:kara_audio_tags"`
	Authors     []TimingAuthor `gorm:"many2many:kara_authors_tags"`
	Medias      []MediaDB      `gorm:"many2many:kara_media_tags"`
	Title       string
	ExtraTitles []AdditionalName `gorm:"many2many:kara_info_additional_name"`
	Version     string
	Comment     string
	SongOrder   int
}

func init_model() {
	db := GetDB()
	db.AutoMigrate(&KaraInfoDB{})
}
