// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Japan7/karaberus/karaberus_tools"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func getS3Client() *minio.Client {

	client, err := minio.New(CONFIG.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(CONFIG.S3.KeyID, CONFIG.S3.Secret, ""),
		Secure: CONFIG.S3.Secure,
	})
	if err != nil {
		panic(err)
	}

	return client
}

func UploadToS3(ctx context.Context, file io.Reader, filename string, filesize int64) error {
	client := getS3Client()
	info, err := client.PutObject(ctx, CONFIG.S3.BucketName, filename, file, filesize, minio.PutObjectOptions{})
	fmt.Printf("info: %v\n", info)

	return err
}

func CheckValidFiletype(type_directory string) bool {
	switch type_directory {
	case "video", "sub", "inst":
		return true
	default:
		return false
	}
}

func SaveFileToS3(ctx context.Context, fd io.Reader, kid string, type_directory string, filesize int64) error {
	if !CheckValidFiletype(type_directory) {
		return &KaraberusError{"Unknown file type " + type_directory}
	}
	filename := filepath.Join(type_directory, "/", kid)
	return UploadToS3(ctx, fd, filename, filesize)
}

type CheckS3FileOutput struct {
	Passed bool `json:"passed" example:"true" doc:"true if file passed all checks"`
}

type CheckKaraOutput struct {
	Video        *karaberus_tools.DakaraCheckResultsOutput
	Instrumental *karaberus_tools.DakaraCheckResultsOutput
	Subtitles    *karaberus_tools.DakaraCheckResultsOutput
}

func CheckKara(ctx context.Context, kara KaraInfoDB) (*CheckKaraOutput, error) {
	out := &CheckKaraOutput{}

	if kara.VideoUploaded {
		video_filename := fmt.Sprintf("video/%d", kara.ID)
		video_check_res, err := CheckS3File(ctx, video_filename)
		if err != nil {
			return nil, err
		}
		out.Video = video_check_res
	}
	if kara.SubtitlesUploaded {
		sub_filename := fmt.Sprintf("sub/%d", kara.ID)
		sub_check_res, err := CheckS3File(ctx, sub_filename)
		if err != nil {
			return nil, err
		}
		out.Subtitles = sub_check_res
	}
	if kara.InstrumentalUploaded {
		inst_filename := fmt.Sprintf("inst/%d", kara.ID)
		inst_check_res, err := CheckS3File(ctx, inst_filename)
		if err != nil {
			return nil, err
		}
		out.Instrumental = inst_check_res
	}

	return out, nil
}

func CheckS3File(ctx context.Context, video_filename string) (*karaberus_tools.DakaraCheckResultsOutput, error) {
	client := getS3Client()

	obj, err := client.GetObject(ctx, CONFIG.S3.BucketName, video_filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	res := karaberus_tools.DakaraCheckResults(obj)

	return &res, nil
}
