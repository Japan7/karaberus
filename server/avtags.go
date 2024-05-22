package server

import "context"

func getVideoTag(video_type string) VideoTag {
	for _, v := range VideoTags {
		if v.ID == video_type {
			return v
		}
	}

	panic("unknown video type " + video_type)
}

func getAudioTag(audio_type string) AudioTag {
	for _, v := range AudioTags {
		if v.ID == audio_type {
			return v
		}
	}

	panic("unknown audio type " + audio_type)
}

// Public/API functions

type VideoTagsOutput struct {
	Body []VideoTag `json:"video_tags"`
}

func GetVideoTags(ctx context.Context, input *struct{}) (*VideoTagsOutput, error) {
	return &VideoTagsOutput{VideoTags}, nil
}

type AudioTagsOutput struct {
	Body []AudioTag `json:"audio_tags"`
}

func GetAudioTags(ctx context.Context, input *struct{}) (*AudioTagsOutput, error) {
	return &AudioTagsOutput{AudioTags}, nil
}
