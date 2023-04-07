package sdk

import (
	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/authenticator"
	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/iamcore"
	"net/http"
)

// Client provides:
// - OpenID Connect-based, OAuth 2.0-compliant middleware for authenticating API calls;
type Client struct {
	httpClient *http.Client

	authenticator authenticator.AuthenticationClient
}

func NewClient(opts *Options) (*Client, error) {
	if err := opts.validate(); err != nil {
		return nil, nil
	}

	httpClient := http.DefaultClient

	iamcoreClient := iamcore.NewClient(opts, httpClient)

	return &Client{
		httpClient,
		authenticator.NewClient([]authenticator.Authenticator{
			authenticator.NewBearer(iamcoreClient),
		}),
	}, nil
}
