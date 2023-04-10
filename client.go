package sdk

import (
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/iamcore"
)

// Client provides:
// - OpenID Connect-based, OAuth 2.0-compliant middleware for authenticating API calls;
type Client interface {
	AuthenticationClient
}

type client struct {
	authenticators []Authenticator
}

func NewClient(opts *iamcore.Options) (Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil
	}

	return &client{
		[]Authenticator{
			NewBearer(iamcore.NewClient(opts, http.DefaultClient)),
		},
	}, nil
}
