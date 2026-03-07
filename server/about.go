package server

import (
	"context"
	"fmt"
)

type KaraberusVersion struct {
	Revision    string `json:"revision"`
	RepoURL     string `json:"repository_url"`
	RevisionURL string `json:"revision_url"`
}

var karaberus_revision string
var karaberus_repository_url string

type KaraberusVersionOutput struct {
	Body struct {
		Info KaraberusVersion `json:"info"`
	}
}

func RevisionURL(repo_url string, revision string) string {
	return fmt.Sprintf("%s/commit/%s", repo_url, revision)
}

func GetVersion(ctx context.Context, input *struct{}) (*KaraberusVersionOutput, error) {
	version := KaraberusVersion{
		Revision:    karaberus_revision,
		RepoURL:     karaberus_repository_url,
		RevisionURL: RevisionURL(karaberus_repository_url, karaberus_revision),
	}

	out := &KaraberusVersionOutput{}
	out.Body.Info = version

	return out, nil
}
