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
type VideoType KaraberusType
type AudioType KaraberusType

// Tags
var (
	KaraTagTitle   TagType = TagType{Type: "Title", Value: 1}
	KaraTagAuthor  TagType = TagType{Type: "Author", Value: 2}
	KaraTagArtist  TagType = TagType{Type: "Artist", Value: 3}
	KaraTagVersion TagType = TagType{Type: "Version", Value: 4}
)

var TagTypes []TagType = []TagType{
	KaraTagTitle, KaraTagAuthor, KaraTagArtist, KaraTagVersion,
}

type Tag struct {
	gorm.Model
	Name            string           `gorm:"unique_index:idx_name_type"`
	Type            TagType          `gorm:"unique_index:idx_name_type"`
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
	VideoTypeOpening       VideoType = VideoType{Type: "Opening", Value: 1}
	VideoTypeEnding        VideoType = VideoType{Type: "Ending", Value: 2}
	VideoTypeInsert        VideoType = VideoType{Type: "Insert", Value: 3}
	VideoTypeFanmade       VideoType = VideoType{Type: "Fanmade", Value: 4}
	VideoTypeStream        VideoType = VideoType{Type: "Stream", Value: 5}
	VideoTypeConcert       VideoType = VideoType{Type: "Concert", Value: 6}
	VideoTypeMusicVideo    VideoType = VideoType{Type: "Promotional Video", Value: 7}
	VideoTypeAdvertisement VideoType = VideoType{Type: "Advertisement", Value: 8}
	VideoTypeTrailer       VideoType = VideoType{Type: "Trailer", Value: 9}
)

var VideoTypes []VideoType = []VideoType{
	VideoTypeOpening,
	VideoTypeEnding,
	VideoTypeInsert,
	VideoTypeMusicVideo,
	VideoTypeFanmade,
	VideoTypeConcert,
	VideoTypeAdvertisement,
}

// Audio tags
var (
	AudioTypeOpening AudioType = AudioType{Type: "Opening", Value: 1}
	AudioTypeEnding  AudioType = AudioType{Type: "Ending", Value: 2}
	AudioTypeInsert  AudioType = AudioType{Type: "Insert", Value: 3}
	AudioTypeLive    AudioType = AudioType{Type: "Live", Value: 4}
	AudioTypeCover   AudioType = AudioType{Type: "Cover", Value: 5}
)

var AudioTypes []AudioType = []AudioType{
	AudioTypeOpening,
	AudioTypeEnding,
	AudioTypeInsert,
	AudioTypeLive,
}

type AdditionalName struct {
	gorm.Model
	Name string
}

type KaraInfoDB struct {
	gorm.Model
	Tags        []Tag `gorm:"many2many:kara_info_tags"`
	Title       string
	ExtraTitles []AdditionalName `gorm:"many2many:kara_info_additional_name"`
	Comment     string
	SongOrder   int
}

func init_model() {
	db := GetDB()
	db.AutoMigrate(&KaraInfoDB{})
}

// Helper functions

func getTag(name string, tag_type TagType) Tag {
	tag := Tag{}
	tx := GetDB().Where(&Tag{Name: name, Type: tag_type}).FirstOrCreate(&tag)
	if tx.Error != nil {
		panic(tx.Error.Error())
	}
	return tag
}

func getAuthor(author_name string) Tag {
	return getTag(author_name, KaraTagAuthor)
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

func getVideoTag(video_type string) VideoType {
	for _, v := range VideoTypes {
		if v.Type == video_type {
			return v
		}
	}

	panic("unknown kara type " + video_type)
}

func getAudioTag(audio_type string) AudioType {
	for _, v := range AudioTypes {
		if v.Type == audio_type {
			return v
		}
	}

	panic("unknown kara type " + audio_type)
}

// Public/API functions

func GetVideoTags(ctx context.Context, input *struct{}) (*[]VideoType, error) {
	return &VideoTypes, nil
}

func GetAudioTags(ctx context.Context, input *struct{}) (*[]AudioType, error) {
	return &AudioTypes, nil
}
