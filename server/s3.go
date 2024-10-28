// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Japan7/karaberus/karaberus_tools"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

var S3_CLIENTS map[string]*minio.Client = map[string]*minio.Client{}
var S3_BEST_CLIENT *minio.Client

var clientsMutex = sync.Mutex{}

type TestedClient struct {
	Client      *minio.Client
	ListLatency int64
}

func initS3Clients(ctx context.Context) {
	if len(CONFIG.S3.Endpoints) == 0 {
		panic("No S3 endpoints configured")
	}

	var err error
	for _, endpoint := range CONFIG.S3.Endpoints {
		S3_CLIENTS[endpoint], err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(CONFIG.S3.KeyID, CONFIG.S3.Secret, ""),
			Secure: CONFIG.S3.Secure,
		})
		if err != nil {
			panic(err)
		}
	}

	pickBestClient(ctx)

	go func() {
		for {
			pickBestClient(ctx)
			time.Sleep(60 * time.Second)
		}
	}()
}

// We’re assuming that a garage node among one of the addresses is on the
// local host, which would have the lowest latency and probably offer the best
// bandwidth.
func pickBestClient(ctx context.Context) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	var best_client *TestedClient = nil

	var err error
	for _, client := range S3_CLIENTS {
		begin_time := time.Now()
		_, err = client.ListBuckets(ctx)
		if err != nil {
			continue
		}
		tested := TestedClient{client, time.Since(begin_time).Nanoseconds()}

		if best_client == nil || tested.ListLatency < best_client.ListLatency {
			best_client = &tested
		}
	}

	if best_client == nil {
		panic(err)
	}

	S3_BEST_CLIENT = best_client.Client
}

func getS3Client() *minio.Client {
	return S3_BEST_CLIENT
}

func UploadToS3(ctx context.Context, file io.Reader, filename string, filesize int64, user_metadata map[string]string) error {
	client := getS3Client()

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

func SaveFileToS3WithMetadata(ctx context.Context, tx *gorm.DB, fd io.Reader, kara *KaraInfoDB, type_directory string, filesize int64, crc32 uint32, user_metadata map[string]string) (*CheckKaraOutput, error) {
	if kara.ID == 0 {
		return nil, errors.New("trying to upload to a karaoke that doesn't exist")
	}

	if !CheckValidFiletype(type_directory) {
		return nil, errors.New("Unknown file type " + type_directory)
	}
	filename := fmt.Sprintf("%s/%d", type_directory, kara.ID)
	err := UploadToS3(ctx, fd, filename, filesize, user_metadata)
	if err != nil {
		return nil, err
	}

	err = updateKaraokeAfterUpload(tx, kara, type_directory, filesize, crc32)
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
	}

	err = tx.Model(&kara).Updates(&kara).Error
	if err != nil {
		return nil, DBErrToHumaErr(err)
	}

	return res, nil
}

func SaveTempFileToS3WithMetadata(ctx context.Context, tx *gorm.DB, tempfile UploadTempFile, kara *KaraInfoDB, type_directory string, user_metadata map[string]string) (*CheckKaraOutput, error) {
	return SaveFileToS3WithMetadata(ctx, tx, tempfile.Fd, kara, type_directory, tempfile.Size, tempfile.CRC32, user_metadata)
}

func SaveTempFileToS3(ctx context.Context, tx *gorm.DB, tempfile UploadTempFile, kara *KaraInfoDB, type_directory string) (*CheckKaraOutput, error) {
	return SaveTempFileToS3WithMetadata(ctx, tx, tempfile, kara, type_directory, nil)
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
	Subtitles    *karaberus_tools.DakaraCheckSubResultsOutput
}

func CheckKara(ctx context.Context, kara KaraInfoDB) (*CheckKaraOutput, error) {
	out := &CheckKaraOutput{}

	if kara.VideoUploaded {
		obj, err := GetKaraObject(ctx, kara, "video")
		if err != nil {
			return nil, err
		}
		defer Closer(obj)
		stat, err := obj.Stat()
		if err != nil {
			return nil, err
		}
		video_check_res := CheckS3Video(ctx, obj, stat.Size)
		if !video_check_res.Passed {
			return nil, fmt.Errorf("checks failed for kara %d", kara.ID)
		}
		out.Video = &video_check_res
	}
	if kara.SubtitlesUploaded {
		obj, err := GetKaraObject(ctx, kara, "sub")
		if err != nil {
			return nil, err
		}
		defer Closer(obj)
		stat, err := obj.Stat()
		if err != nil {
			return nil, err
		}
		sub_check_res, err := CheckS3Ass(ctx, obj, stat.Size)
		if err != nil {
			return nil, err
		}
		out.Subtitles = &sub_check_res
	}
	if kara.InstrumentalUploaded {
		obj, err := GetKaraObject(ctx, kara, "inst")
		if err != nil {
			return nil, err
		}
		defer Closer(obj)
		stat, err := obj.Stat()
		if err != nil {
			return nil, err
		}
		inst_check_res := CheckS3Inst(ctx, obj, stat.Size)
		if !inst_check_res.Passed {
			return nil, fmt.Errorf("checks failed for kara %d", kara.ID)
		}
		out.Instrumental = &inst_check_res
	}

	return out, nil
}

func getS3FontFilename(id uint) string {
	return fmt.Sprintf("font/%d", id)
}

func GetFontObject(ctx context.Context, id uint) (*minio.Object, error) {
	client := getS3Client()
	filename := getS3FontFilename(id)
	obj, err := client.GetObject(ctx, CONFIG.S3.BucketName, filename, minio.GetObjectOptions{})
	return obj, err
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

func GetKaraLyrics(ctx context.Context, kara KaraInfoDB) (string, error) {
	if !kara.SubtitlesUploaded {
		return "", nil
	}

	obj, err := GetKaraObject(ctx, kara, "sub")
	if err != nil {
		return "", err
	}
	stat, err := obj.Stat()
	if err != nil {
		return "", err
	}
	res, err := CheckS3Ass(ctx, obj, stat.Size)
	if err != nil {
		return "", err
	}
	return res.Lyrics, nil
}

func CheckS3Video(ctx context.Context, obj io.ReadSeeker, size int64) karaberus_tools.DakaraCheckResultsOutput {
	res := karaberus_tools.DakaraCheckResults(obj, "video", size)
	return res
}

func CheckS3Inst(ctx context.Context, obj io.ReadSeeker, size int64) karaberus_tools.DakaraCheckResultsOutput {
	res := karaberus_tools.DakaraCheckResults(obj, "inst", size)
	return res
}

func CheckS3Ass(ctx context.Context, obj io.ReadSeeker, size int64) (karaberus_tools.DakaraCheckSubResultsOutput, error) {
	out, err := karaberus_tools.DakaraCheckSub(obj, size)
	return out, err
}
