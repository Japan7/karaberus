// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
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
	getLogger().Printf("upload info: %v\n", info)

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
		return errors.New("Unknown file type " + type_directory)
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
		obj, err := GetKaraObject(ctx, kara, "video")
		if err != nil {
			return nil, err
		}
		video_check_res, err := CheckS3File(ctx, obj)
		if err != nil {
			return nil, err
		}
		out.Video = video_check_res
	}
	if kara.SubtitlesUploaded {
		obj, err := GetKaraObject(ctx, kara, "sub")
		if err != nil {
			return nil, err
		}
		sub_check_res, err := CheckS3Ass(ctx, obj)
		if err != nil {
			return nil, err
		}
		out.Subtitles = sub_check_res
	}
	if kara.InstrumentalUploaded {
		obj, err := GetKaraObject(ctx, kara, "inst")
		if err != nil {
			return nil, err
		}
		inst_check_res, err := CheckS3File(ctx, obj)
		if err != nil {
			return nil, err
		}
		out.Instrumental = inst_check_res
	}

	return out, nil
}

func GetKaraObject(ctx context.Context, kara KaraInfoDB, filetype string) (*minio.Object, error) {
	client := getS3Client()

	if !CheckValidFiletype(filetype) {
		return nil, errors.New("Unknown file type " + filetype)
	}

	filename := fmt.Sprintf("%s/%d", filetype, kara.ID)
	obj, err := client.GetObject(ctx, CONFIG.S3.BucketName, filename, minio.GetObjectOptions{})
	return obj, err
}

func CheckS3File(ctx context.Context, obj *minio.Object) (*karaberus_tools.DakaraCheckResultsOutput, error) {
	res := karaberus_tools.DakaraCheckResults(obj)
	return &res, nil
}

func CheckS3Ass(ctx context.Context, obj *minio.Object) (*karaberus_tools.DakaraCheckResultsOutput, error) {
	return &karaberus_tools.DakaraCheckResultsOutput{Passed: true}, nil
}
