package sdk

import (
	"context"
	"net/http"

	"github.com/kaaproject/httperror"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type AuthenticationClient interface {
	// WithAuth creates http.Handler middleware for authenticating the incoming request by means of either OAuth 2.0 Access Token
	// or "X-iamcore-API-Key" HTTP header. This handler should precede the application request handling in the handlers chain.
	// It populates the request context with validated requester principal's IRN for further use.
	// Returns 401 Unauthorized HTTP error in case of unauthorized access, and stops HTTP request propagation.
	WithAuth(next http.Handler) http.Handler
}

// contextKeyType is a context.Context key type
type contextKeyType int

const principalIRNKey contextKeyType = 1

func (c *client) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := range c.authenticators {
			principal, err := c.authenticators[i].Authenticate(r.Context(), r.Header)
			if err != nil {
				httperror.Write(w, err)

				return
			}

			if principal != nil {
				ctx := context.WithValue(r.Context(), principalIRNKey, principal)

				r = r.WithContext(ctx)

				// Pass control to the next handler
				next.ServeHTTP(w, r)
				return
			}
		}

		httperror.Write(w, httperror.New(http.StatusUnauthorized, "Failed to authenticate request with any of available authenticators"))
	})
}

// PrincipalIRN extracts and returns principal's IRN from the request context.
func PrincipalIRN(ctx context.Context) (*irn.IRN, error) {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return nil, ErrPrincipalIRNIsNotSet
	}

	return principal, nil
}
