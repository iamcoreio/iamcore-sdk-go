package iamcore

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

const authorizationHeaderName = "Authorization"

type Bearer struct {
	iamcore *ServerClient
}

func NewBearer(iamcore *ServerClient) *Bearer {
	return &Bearer{
		iamcore: iamcore,
	}
}

func (b *Bearer) Authenticate(ctx context.Context, header http.Header) (*irn.IRN, http.Header, error) {
	bearerTokenHeader := header.Get(authorizationHeaderName)
	if len(bearerTokenHeader) == 0 {
		return nil, nil, nil
	}

	rawAccessTokenParts := strings.Split(bearerTokenHeader, " ")
	if len(rawAccessTokenParts) != 2 || strings.ToLower(rawAccessTokenParts[0]) != "bearer" {
		return nil, nil, fmt.Errorf("unexpected Authorization header format, must be 'Bearer <access-token>': %w", ErrUnauthenticated)
	}

	authorizationHeader := http.Header{
		authorizationHeaderName: {bearerTokenHeader},
	}

	principalIRN, err := b.iamcore.GetPrincipalIRN(ctx, authorizationHeader)
	if err != nil {
		return nil, nil, err
	}

	return principalIRN, authorizationHeader, nil
}
