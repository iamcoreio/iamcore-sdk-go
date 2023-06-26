package iamcore

import (
	"errors"
	"os"
)

var ErrEmptyAPIKey = errors.New("empty API key")

type Options struct {
	// serverURL to access the iamcore; "https://cloud.iamcore.io" by default.
	serverURL string
	// API key for outbound HTTP requests to secured by iamcore applications or iamcore itself
	apiKey string
}

const (
	apiKeyEnvKey = "IAMCORE_API_KEY" //#nosec

	iamcoreURLEnvKey  = "IAMCORE_URL"
	iamcoreDefaultURL = "https://cloud.iamcore.io"
)

func newOptions(apiKey, serverURL string) (*Options, error) {
	if apiKey == "" {
		apiKey = os.Getenv(apiKeyEnvKey)
	}

	if apiKey == "" {
		return nil, ErrEmptyAPIKey
	}

	if serverURL == "" {
		serverURL = os.Getenv(iamcoreURLEnvKey)
	}

	if serverURL == "" {
		serverURL = iamcoreDefaultURL
	}

	return &Options{
		serverURL: serverURL,
		apiKey:    apiKey,
	}, nil
}
