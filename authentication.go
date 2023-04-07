package sdk

import (
	"context"
	"errors"
	"net/http"

	"github.com/kaaproject/httperror"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type AuthenticationClient interface {
	// WithAuth creates http.Handler middleware for validating OAuth 2.0 Access Token in the incoming request. This handler
	// should precede the application request handling in the handlers chain. It populates the request context with validated
	// Access Token, JWT "sub" claim, account ID, tenant ID, username and path for further use.
	// Returns 401 Unauthorized HTTP error in case of unauthorized access (no, invalid, fake, etc. Access Token), and stops
	// HTTP request propagation.
	WithAuth(next http.Handler) http.Handler
}

// contextKeyType is a context.Context key type
type contextKeyType int

const (
	principalIRNKey contextKeyType = 5
)

var ErrPrincipalIRNIsNotSet = errors.New("principal IRN is not set")

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

		httperror.Write(w, httperror.New(http.StatusUnauthorized, "Authentication header is missed"))
	})
}

// Principal extracts and returns principal's IRN from the request context.
func Principal(ctx context.Context) (*irn.IRN, error) {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return nil, ErrPrincipalIRNIsNotSet
	}

	return principal, nil
}

// AccountID extracts and returns account ID from the request context.
func AccountID(ctx context.Context) string {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return ""
	}

	return principal.GetAccountID()
}

// TenantID extracts and returns tenant ID from the request context.
func TenantID(ctx context.Context) string {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return ""
	}

	return principal.GetTenantID()
}

// Path extracts and returns principal's path from the request context.
func Path(ctx context.Context) string {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return ""
	}

	return principal.GetPath()
}
