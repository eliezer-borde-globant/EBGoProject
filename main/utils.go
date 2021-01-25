package main

import (
	"os"

	"github.com/rs/zerolog"
)

type forkResponseParams struct {
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	GitURL string `json:"git_url"`
}

type UpdateParams struct {
	Repo    string                              `json:"repo" xml:"repo" form:"repo"`
	Owner   string                              `json:"owner" xml:"owner" form:"owner"`
	Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
}

type CreateParams struct {
	Repo    string `json:"repo" xml:"repo" form:"repo"`
	Owner   string `json:"owner" xml:"owner" form:"owner"`
	Content string `json:"content" xml:"content" form:"content"`
}

type SecretUpdateMap map[string][]map[string]interface{}

var (
	gitHubToken     = "99f85d692d80031fe3e2bbc67c8ff903608ce5aa"
	zeroLogger      = zerolog.New(os.Stdout).With().Timestamp().Logger()
	SecretsFileName = ".secrets.baseline"
)
