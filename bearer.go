package sdk

import (
	"context"
	"net/http"
	"strings"

	"github.com/kaaproject/httperror"

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
	bearerTokenHeader := header.Get(authorizationHeaderName)
	if len(bearerTokenHeader) == 0 {
		return nil, nil
	}

	rawAccessTokenParts := strings.Split(bearerTokenHeader, " ")
	if len(rawAccessTokenParts) != 2 || strings.ToLower(rawAccessTokenParts[0]) != "bearer" {
		return nil, httperror.New(http.StatusUnauthorized, `Unexpected Authorization header format, must be "Bearer <access-token>"`)
	}

	return b.iamcore.GetPrincipalIRN(authorizationHeaderName, bearerTokenHeader)
}
