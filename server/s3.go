// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Japan7/karaberus/karaberus_tools"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

var S3_CLIENT *minio.Client = nil

func getS3Client() (*minio.Client, error) {
	var err error = nil
	if S3_CLIENT == nil {
		S3_CLIENT, err = minio.New(CONFIG.S3.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(CONFIG.S3.KeyID, CONFIG.S3.Secret, ""),
			Secure: CONFIG.S3.Secure,
		})
	}
	return S3_CLIENT, err
}

func UploadToS3(ctx context.Context, file io.Reader, filename string, filesize int64, user_metadata map[string]string) error {
	client, err := getS3Client()
	if err != nil {
		return err
	}

	info, err := client.PutObject(ctx, CONFIG.S3.BucketName, filename, file, filesize, minio.PutObjectOptions{
		UserMetadata: user_metadata,
		PartSize:     5 * 1024 * 1024,
	})
	getLogger().Printf("upload info: %+v\n", info)

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

func SaveFileToS3WithMetadata(ctx context.Context, tx *gorm.DB, fd io.Reader, kara *KaraInfoDB, type_directory string, filesize int64, user_metadata map[string]string) (*CheckKaraOutput, error) {
	kid := kara.ID

	if !CheckValidFiletype(type_directory) {
		return nil, errors.New("Unknown file type " + type_directory)
	}
	filename := fmt.Sprintf("%s/%d", type_directory, kid)
	err := UploadToS3(ctx, fd, filename, filesize, user_metadata)
	if err != nil {
		return nil, err
	}

	err = updateKaraokeAfterUpload(tx, kara, type_directory)
	if err != nil {
		return nil, err
	}

	res, err := CheckKara(ctx, *kara)
	if err != nil {
		return nil, err
	}

	if res.Video != nil {
		if res.Video.Duration != kara.Duration {
			kara.Duration = res.Video.Duration
		}

		if CONFIG.Dakara.BaseURL != "" && kara.UploadInfo.VideoUploaded && kara.UploadInfo.SubtitlesUploaded {
			go SyncDakara(context.Background())
		}
	}

	err = tx.Save(&kara).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}

	return res, nil
}

func SaveFileToS3(ctx context.Context, tx *gorm.DB, fd io.Reader, kara *KaraInfoDB, type_directory string, filesize int64) (*CheckKaraOutput, error) {
	return SaveFileToS3WithMetadata(ctx, tx, fd, kara, type_directory, filesize, nil)
}

func SaveFontToS3(ctx context.Context, fd io.Reader, id uint, filesize int64) error {
	filename := getS3FontFilename(id)
	return UploadToS3(ctx, fd, filename, filesize, nil)
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
		defer obj.Close()
		video_check_res := CheckS3File(ctx, obj)
		if !video_check_res.Passed {
			return nil, errors.New(fmt.Sprintf("Checks failed for kara %d", kara.ID))
		}
		out.Video = video_check_res
	}
	if kara.SubtitlesUploaded {
		obj, err := GetKaraObject(ctx, kara, "sub")
		if err != nil {
			return nil, err
		}
		defer obj.Close()
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
		defer obj.Close()
		inst_check_res := CheckS3File(ctx, obj)
		if !inst_check_res.Passed {
			return nil, errors.New(fmt.Sprintf("Checks failed for kara %d", kara.ID))
		}
		out.Instrumental = inst_check_res
	}

	return out, nil
}

func getS3FontFilename(id uint) string {
	return fmt.Sprintf("font/%d", id)
}

func GetFontObject(ctx context.Context, id uint) (*minio.Object, error) {
	client, err := getS3Client()
	if err != nil {
		return nil, err
	}

	filename := getS3FontFilename(id)
	obj, err := client.GetObject(ctx, CONFIG.S3.BucketName, filename, minio.GetObjectOptions{})
	return obj, err
}

func GetKaraObject(ctx context.Context, kara KaraInfoDB, filetype string) (*minio.Object, error) {
	client, err := getS3Client()
	if err != nil {
		return nil, err
	}

	if !CheckValidFiletype(filetype) {
		return nil, errors.New("Unknown file type " + filetype)
	}

	filename := fmt.Sprintf("%s/%d", filetype, kara.ID)
	obj, err := client.GetObject(ctx, CONFIG.S3.BucketName, filename, minio.GetObjectOptions{})
	return obj, err
}

func CheckS3File(ctx context.Context, obj *minio.Object) *karaberus_tools.DakaraCheckResultsOutput {
	res := karaberus_tools.DakaraCheckResults(obj)
	return &res
}

func CheckS3Ass(ctx context.Context, obj *minio.Object) (*karaberus_tools.DakaraCheckResultsOutput, error) {
	return &karaberus_tools.DakaraCheckResultsOutput{Passed: true}, nil
}
