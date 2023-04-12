package sdk

import (
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/iamcore"
)

// Client provides:
// - OpenID Connect-based, OAuth 2.0-compliant middleware for authenticating API calls;
type Client interface {
	AuthenticationClient
	AuthorizationClient
}

type client struct {
	authenticators []Authenticator
	iamcoreClient  *iamcore.Client
}

func NewClient(apiKey, iamcoreHost string) (Client, error) {
	options, err := newOptions(apiKey, iamcoreHost)
	if err != nil {
		return nil, err
	}

	iamcoreClient := iamcore.NewClient(options.iamcoreHost, http.DefaultClient)

	return &client{
		[]Authenticator{
			NewBearer(iamcoreClient),
			NewAPIKey(iamcoreClient),
		},
		iamcoreClient,
	}, nil
}
