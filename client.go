package sdk

import (
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/authenticator"
	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/iamcore"
)

// Client provides:
// - OpenID Connect-based, OAuth 2.0-compliant middleware for authenticating API calls;
type Client interface {
	AuthenticationClient
}

type client struct {
	authenticators []authenticator.Authenticator
}

func NewClient(opts *Options) (Client, error) {
	if err := opts.validate(); err != nil {
		return nil, nil
	}

	iamcoreClient := iamcore.NewClient(opts, http.DefaultClient)

	return &client{
		[]authenticator.Authenticator{
			authenticator.NewBearer(iamcoreClient),
		},
	}, nil
}
