//go:build !cgo

package karaberus_tools

import "github.com/minio/minio-go/v7"

func DakaraCheckResults(obj *minio.Object, ftype string) DakaraCheckResultsOutput {
	out := DakaraCheckResultsOutput{Passed: true}
	return out
}

func DakaraCheckSub(obj *minio.Object) (DakaraCheckSubResultsOutput, error) {
	out := DakaraCheckSubResultsOutput{
		Lyrics: "",
		Passed: true,
	}
	return out, nil
}
