package utils

import (
	"github.com/rs/zerolog"
	"os"
)


var (
	GitHubToken = os.Getenv("GITHUB_TOKEN")
	ZeroLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
)

const (
	SecretsFileName = ".secrets.baseline"
)

// map defined to read the parse json, used while trying to edit the secrets file
type SecretUpdateMap map[string][]map[string]interface{}