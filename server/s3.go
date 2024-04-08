// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

/*
#cgo pkg-config: dakara_check
#include <stdlib.h>
#include <unistd.h>
#include <dakara_check.h>

const size_t BUFSIZE = 1024*4;

struct fdpipe {
  int fdr;
  int fdw;
};

struct fdpipe *create_pipe(void) {
  int pipefd[2];
  if (pipe(pipefd) < 0) {
    perror("failed to create pipe");
    return NULL;
  }
  struct fdpipe *fdpipe = malloc(sizeof(struct fdpipe));
  if (fdpipe == NULL)
    return NULL;

  fdpipe->fdr = pipefd[0];
  fdpipe->fdw = pipefd[1];
  return fdpipe;
}

int read_piped(void *opaque, uint8_t *buf, int n) {
  int *fd = (int*) opaque;
  return read(*fd, buf, n);
}

struct dakara_check_results *karaberus_dakara_check(int fdr) {
	return dakara_check_avio(BUFSIZE, &fdr, read_piped, NULL);
}
*/
import "C"
import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func getS3Client() (*s3.Client, error) {
	s3_creds := credentials.NewStaticCredentialsProvider(S3_KEYID, S3_SECRET, "")
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	config.Credentials = s3_creds
	config.RetryMaxAttempts = 3
	config.BaseEndpoint = &S3_ENDPOINT

	client := s3.NewFromConfig(config)
	return client, nil
}

func UploadToS3(ctx context.Context, file io.Reader, filename string) error {
	client, err := getS3Client()
	if err != nil {
		return err
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
		Body:   file,
	})

	return err
}

func SaveFileToS3(ctx context.Context, fd io.Reader, kid uuid.UUID, type_directory string) error {
	filename := filepath.Join(type_directory, "/", kid.String())
	return UploadToS3(ctx, fd, filename)
}

type CheckS3FileOutput struct {
	Passed bool `json:"passed" example:"true" doc:"true if file passed all checks"`
}

func CheckKara(ctx context.Context, kid uuid.UUID) (*CheckS3FileOutput, error) {
	// TODO: find all related files to check
	video_filename := filepath.Join("video/", kid.String())
	return CheckS3File(ctx, video_filename)
}

func CheckS3File(ctx context.Context, video_filename string) (*CheckS3FileOutput, error) {

	client, err := getS3Client()
	if err != nil {
		return nil, err
	}

	obj, err := client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(BUCKET_NAME),
			Key:    aws.String(video_filename),
		},
	)

	fdpipe := C.create_pipe()
	if fdpipe == nil {
		return nil, errors.New("failed to create pipe")
	}
	defer C.close(fdpipe.fdr)
	defer C.close(fdpipe.fdw)
	defer C.free(unsafe.Pointer(fdpipe))

	go func(fdw C.int, objreader io.Reader) {
		for {
			buf := make([]byte, C.BUFSIZE)
			n, err := objreader.Read(buf)
			if err != nil {
				panic(err)
			}
			if n == 0 {
				break
			}
			written := C.write(fdw, C.CBytes(buf), C.size_t(n))
			if int(written) != n {
				panic(fmt.Sprintf("wrote less bytes than expected: %d != %d", written, n))
			}
		}
	}(fdpipe.fdw, obj.Body)

	res := C.karaberus_dakara_check(fdpipe.fdr)
	defer C.dakara_check_results_free(res)

	out := &CheckS3FileOutput{
		Passed: bool(res.passed),
	}

	return out, nil
}
