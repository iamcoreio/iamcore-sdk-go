package sdk

import (
	"context"
	"net/http"
	"strings"

	"github.com/kaaproject/httperror"
	"github.com/rs/zerolog/log"

	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/iamcore"
	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

const authorizationHeaderName = "Authorization"

type Bearer struct {
	iamcore *iamcore.Client
}

func NewBearer(iamcore *iamcore.Client) *Bearer {
	return &Bearer{
		iamcore: iamcore,
	}
}

func (b *Bearer) Authenticate(ctx context.Context, header http.Header) (*irn.IRN, error) {
	// Make sure that the raw access token exists, it is well-formed, and the auth scheme is "Bearer" for OAuth 2.0
	authorizationHeader := header.Get(authorizationHeaderName)
	if len(authorizationHeader) == 0 {
		log.Debug().Msgf("'%s' header is missed. Skipping authentication by bearer token", authorizationHeaderName)

		return nil, nil
	}

	rawAccessTokenParts := strings.Split(authorizationHeader, " ")
	if len(rawAccessTokenParts) != 2 || strings.ToLower(rawAccessTokenParts[0]) != "bearer" {
		return nil, httperror.New(http.StatusUnauthorized, `Unexpected Authorization header format, must be "Bearer <access-token>"`)
	}

	return b.iamcore.GetUserIRN(authorizationHeader)
}
