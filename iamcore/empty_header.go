package iamcore

import (
	"context"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type EmptyHeader struct {
	iamcore *ServerClient
}

func NewEmptyHeader(iamcore *ServerClient) *EmptyHeader {
	return &EmptyHeader{
		iamcore: iamcore,
	}
}

func (a *EmptyHeader) Authenticate(ctx context.Context, _ http.Header) (*irn.IRN, http.Header, error) {
	principalIRN, err := a.iamcore.GetPrincipalIRN(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	return principalIRN, nil, nil
}
