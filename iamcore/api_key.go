package iamcore

import (
	"context"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

const apiKeyHeaderName = "X-iamcore-API-Key" //#nosec

type APIKey struct {
	iamcore *ServerClient
}

func NewAPIKey(iamcore *ServerClient) *APIKey {
	return &APIKey{
		iamcore: iamcore,
	}
}

func (a *APIKey) Authenticate(ctx context.Context, header http.Header) (*irn.IRN, http.Header, error) {
	apiKeyHeader := header.Get(apiKeyHeaderName)
	if len(apiKeyHeader) == 0 {
		return nil, nil, nil
	}

	authorizationHeader := http.Header{
		apiKeyHeaderName: {apiKeyHeader},
	}

	principalIRN, err := a.iamcore.GetPrincipalIRN(ctx, authorizationHeader)
	if err != nil {
		return nil, nil, err
	}

	return principalIRN, authorizationHeader, nil
}
