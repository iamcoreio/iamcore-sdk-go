package iamcore

import (
	"log"
	"net/http"
)

type Client struct {
	authenticators []Authenticator
	iamcoreClient  *ServerClient
	disabled       bool
}

func NewClient(apiKey, iamcoreHost string, disabled bool) (*Client, error) {
	if disabled {
		log.Println("iamcore SDK is DISABLED")

		return &Client{
			disabled: true,
		}, nil
	}

	options, err := newOptions(apiKey, iamcoreHost)
	if err != nil {
		return nil, err
	}

	iamcoreClient := NewServerClient(options.iamcoreHost, http.DefaultClient)

	return &Client{
		authenticators: []Authenticator{
			NewBearer(iamcoreClient),
			NewAPIKey(iamcoreClient),
		},
		iamcoreClient: iamcoreClient,
		disabled:      false,
	}, nil
}
