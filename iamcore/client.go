package iamcore

import (
	"log"
	"net/http"
)

type Client interface {
	AuthenticationClient
	AuthorizationClient
	ResourceManager
}

type сlient struct {
	authenticators []Authenticator
	iamcoreClient  *ServerClient
	disabled       bool

	apiKey string
}

func NewClient(apiKey, serverURL string, disabled bool) (Client, error) {
	if disabled {
		log.Println("iamcore SDK is DISABLED")

		return &сlient{
			disabled: true,
		}, nil
	}

	options, err := newOptions(apiKey, serverURL)
	if err != nil {
		return nil, err
	}

	iamcoreClient := NewServerClient(options.serverURL, http.DefaultClient)

	return &сlient{
		authenticators: []Authenticator{
			NewBearer(iamcoreClient),
			NewAPIKey(iamcoreClient),
		},
		iamcoreClient: iamcoreClient,
		disabled:      false,

		apiKey: options.apiKey,
	}, nil
}
