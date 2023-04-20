package iamcore

import (
	"context"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type Authenticator interface {
	Authenticate(ctx context.Context, header http.Header) (principalIRN *irn.IRN, authorizationHeader http.Header, err error)
}
