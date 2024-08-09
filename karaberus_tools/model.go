package karaberus_tools

type DakaraCheckResultsOutput struct {
	Passed   bool  `json:"passed" example:"true" doc:"true if file passed all checks"`
	Duration int32 `json:"duration" example:"90" doc:"file duration"`
}

type DakaraCheckSubResultsOutput struct {
	Lyrics string
	Passed bool
}
