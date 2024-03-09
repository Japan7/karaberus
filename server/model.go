package server

import (
	"context"

	"gorm.io/gorm"
)

type KaraberusType struct {
	Type  string
	Value uint
}

type TagType KaraberusType
type MediaType KaraberusType
type VideoTag KaraberusType
type AudioTag KaraberusType

// Tags

type TimingAuthor struct {
	gorm.Model
	Name            string
	AdditionalNames []AdditionalName `gorm:"many2many:tags_additional_name"`
}

var (
	KaraTagArtist TagType = TagType{Type: "Artist", Value: 1}
)

var TagTypes []TagType = []TagType{
	KaraTagArtist,
}

type Tag struct {
	gorm.Model
	Name            string           `gorm:"unique_index:idx_name_type"`
	Type            uint             `gorm:"unique_index:idx_name_type"`
	AdditionalNames []AdditionalName `gorm:"many2many:tags_additional_name"`
}

// Media types
var (
	ANIME   MediaType = MediaType{Type: "Anime", Value: 1}
	GAME    MediaType = MediaType{Type: "Game", Value: 2}
	LIVE    MediaType = MediaType{Type: "Live", Value: 3}
	CARTOON MediaType = MediaType{Type: "Cartoon", Value: 4}
)

type MediaDB struct {
	gorm.Model
	Name string `json:"name" example:"Shinseiki Evangelion"`
	Type uint   `json:"media_type" example:"1"`
}

var MediaTypes []MediaType = []MediaType{ANIME, GAME, LIVE, CARTOON}

// Video tags
var (
	VideoTypeOpening       VideoTag = VideoTag{Type: "Opening", Value: 1}
	VideoTypeEnding        VideoTag = VideoTag{Type: "Ending", Value: 2}
	VideoTypeInsert        VideoTag = VideoTag{Type: "Insert", Value: 3}
	VideoTypeFanmade       VideoTag = VideoTag{Type: "Fanmade", Value: 4}
	VideoTypeStream        VideoTag = VideoTag{Type: "Stream", Value: 5}
	VideoTypeConcert       VideoTag = VideoTag{Type: "Concert", Value: 6}
	VideoTypeMusicVideo    VideoTag = VideoTag{Type: "Promotional Video", Value: 7}
	VideoTypeAdvertisement VideoTag = VideoTag{Type: "Advertisement", Value: 8}
	VideoTypeTrailer       VideoTag = VideoTag{Type: "Trailer", Value: 9}
)

var VideoTags []VideoTag = []VideoTag{
	VideoTypeOpening,
	VideoTypeEnding,
	VideoTypeInsert,
	VideoTypeMusicVideo,
	VideoTypeFanmade,
	VideoTypeConcert,
	VideoTypeAdvertisement,
}

type VideoTagDB struct {
	ID uint
}

// Audio tags
var (
	AudioTypeOpening AudioTag = AudioTag{Type: "Opening", Value: 1}
	AudioTypeEnding  AudioTag = AudioTag{Type: "Ending", Value: 2}
	AudioTypeInsert  AudioTag = AudioTag{Type: "Insert", Value: 3}
	AudioTypeLive    AudioTag = AudioTag{Type: "Live", Value: 4}
	AudioTypeCover   AudioTag = AudioTag{Type: "Cover", Value: 5}
)

var AudioTags []AudioTag = []AudioTag{
	AudioTypeOpening,
	AudioTypeEnding,
	AudioTypeInsert,
	AudioTypeLive,
}

type AudioTagDB struct {
	ID uint
}

type AdditionalName struct {
	gorm.Model
	Name string
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

	for _, tag := range AudioTags {
		db.FirstOrCreate(&AudioTagDB{ID: tag.Value})
	}
	for _, tag := range VideoTags {
		db.FirstOrCreate(&VideoTagDB{ID: tag.Value})
	}
}

// Helper functions

func getTag(name string, tag_type TagType) Tag {
	tag := Tag{}
	tx := GetDB().Where(&Tag{Name: name, Type: tag_type.Value}).FirstOrCreate(&tag)
	if tx.Error != nil {
		panic(tx.Error.Error())
	}
	return tag
}

func getAuthor(author_name string) TimingAuthor {
	author := TimingAuthor{}
	tx := GetDB().Where(&TimingAuthor{Name: author_name}).FirstOrCreate(&author)
	if tx.Error != nil {
		panic(tx.Error.Error())
	}
	return author
}

func getMediaType(media_type_name string) MediaType {
	for _, v := range MediaTypes {
		if v.Type == media_type_name {
			return v
		}
	}

	// TODO: make huma check the input
	panic("unknown media type " + media_type_name)
}

func getMedia(name string, media_type_str string) MediaDB {
	media_type := getMediaType(media_type_str)
	media := MediaDB{}
	tx := GetDB().Where(&MediaDB{Name: name, Type: media_type.Value}).FirstOrCreate(&media)
	if tx.Error != nil {
		panic(tx.Error.Error())
	}

	return media
}

func getVideoTag(video_type string) VideoTagDB {
	for _, v := range VideoTags {
		if v.Type == video_type {
			return VideoTagDB{ID: v.Value}
		}
	}

	panic("unknown kara type " + video_type)
}

func getAudioTag(audio_type string) AudioTagDB {
	for _, v := range AudioTags {
		if v.Type == audio_type {
			return AudioTagDB{ID: v.Value}
		}
	}

	panic("unknown kara type " + audio_type)
}

// Public/API functions

func GetVideoTags(ctx context.Context, input *struct{}) (*[]VideoTag, error) {
	return &VideoTags, nil
}

func GetAudioTags(ctx context.Context, input *struct{}) (*[]AudioTag, error) {
	return &AudioTags, nil
}
