//go:build !cgo

package karaberus_tools

import "io"

func DakaraCheckResultsVideo(obj io.ReadSeeker, size int64) DakaraCheckResultsOutput {
	out := DakaraCheckResultsOutput{Passed: true}
	return out
}

func DakaraCheckResultsNoVideo(obj io.ReadSeeker, size int64) DakaraCheckResultsOutput {
	out := DakaraCheckResultsOutput{Passed: true}
	return out
}

func DakaraCheckResultsInst(obj io.ReadSeeker, size int64) DakaraCheckResultsOutput {
	out := DakaraCheckResultsOutput{Passed: true}
	return out
}

func DakaraCheckSub(obj io.ReadSeeker, size int64) (DakaraCheckSubResultsOutput, error) {
	out := DakaraCheckSubResultsOutput{
		Lyrics: "",
		Passed: true,
	}
	return out, nil
}
