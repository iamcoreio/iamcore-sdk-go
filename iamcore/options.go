package iamcore

import (
	"errors"
	"os"
)

var ErrEmptyAPIKey = errors.New("empty API key")

type Options struct {
	// iamcoreHost to access the iamcore; "cloud.iamcore.io" by default.
	iamcoreHost string
	// API key for outbound HTTP requests to secured by iamcore applications or iamcore itself
	apiKey string
}

const (
	apiKeyEnvKey = "API_KEY"

	iamcoreHostEnvKey  = "IAMCORE_HOST"
	iamcoreDefaultHost = "cloud.iamcore.io"
)

func newOptions(apiKey, iamcoreHost string) (*Options, error) {
	if apiKey == "" {
		apiKey = os.Getenv(apiKeyEnvKey)
	}

	if apiKey == "" {
		return nil, ErrEmptyAPIKey
	}

	if iamcoreHost == "" {
		iamcoreHost = os.Getenv(iamcoreHostEnvKey)
	}

	if iamcoreHost == "" {
		iamcoreHost = iamcoreDefaultHost
	}

	return &Options{
		iamcoreHost: iamcoreHost,
		apiKey:      apiKey,
	}, nil
}
