package sdk

import (
	"context"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/iamcore-sdk-go.git/iamcore"
	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

const apiKeyHeaderName = "X-iamcore-API-Key"

type APIKey struct {
	iamcore *iamcore.Client
}

func NewAPIKey(iamcore *iamcore.Client) *APIKey {
	return &APIKey{
		iamcore: iamcore,
	}
}

func (a *APIKey) Authenticate(ctx context.Context, header http.Header) (*irn.IRN, error) {
	apiKeyHeader := header.Get(apiKeyHeaderName)
	if len(apiKeyHeader) == 0 {
		return nil, nil
	}

	return a.iamcore.GetPrincipalIRN(apiKeyHeaderName, apiKeyHeader)
}
