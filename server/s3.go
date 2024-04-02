// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2024 odrling

package server

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func getS3Session() (*session.Session, error) {
	cred := credentials.NewStaticCredentials(S3_KEYID, S3_SECRET, "")

	awsConfig := aws.NewConfig()
	awsConfig.WithMaxRetries(3)
	awsConfig.WithEndpoint(S3_ENDPOINT)
	awsConfig.WithCredentials(cred)

	return session.NewSession(awsConfig)
}

func getS3Uploader() (*s3manager.Uploader, error) {
	sess, err := getS3Session()

	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return nil, err
	}

	return s3manager.NewUploader(sess), nil
}

func UploadToS3(file io.Reader, bucketName string, filename string) error {
	uploader, err := getS3Uploader();
	if (err != nil) {
		return err
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
	})

	return err
}
