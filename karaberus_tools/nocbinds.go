//go:build !cgo

package karaberus_tools

import "github.com/minio/minio-go/v7"

func DakaraCheckResults(obj *minio.Object) DakaraCheckResultsOutput {
	out := DakaraCheckResultsOutput{Passed: true}
	return out
}
