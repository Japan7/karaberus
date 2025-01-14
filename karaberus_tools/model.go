package karaberus_tools

import (
	"github.com/danielgtaylor/huma/v2"
)

type DakaraCheckResultsOutput struct {
	Passed   bool     `json:"passed" example:"true" doc:"true if file passed all checks"`
	Duration int32    `json:"duration" example:"90" doc:"file duration"`
	Messages []string `json:"messages" doc:"error messages"`
}

func (res DakaraCheckResultsOutput) Error() error {
	msg := ""
	for _, message := range res.Messages {
		msg += message + "\n"
	}
	return huma.Error422UnprocessableEntity(msg)
}

type DakaraCheckSubResultsOutput struct {
	Lyrics string `json:"lyrics" doc:"lyrics extracted from the subtitles"`
	Passed bool   `json:"passed" example:"true" doc:"true if file passed all checks"`
}
