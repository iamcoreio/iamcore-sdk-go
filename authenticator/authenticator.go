package authenticator

import (
	"context"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type Authenticator interface {
	Authenticate(ctx context.Context, header http.Header) (*irn.IRN, error)
}
