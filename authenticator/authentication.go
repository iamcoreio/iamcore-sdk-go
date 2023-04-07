package authenticator

import (
	"context"
	"net/http"

	"github.com/kaaproject/httperror"
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

type Client struct {
	authenticators []Authenticator
}

func NewClient(authenticators []Authenticator) *Client {
	return &Client{
		authenticators: authenticators,
	}
}

const (
	principalIRNKey contextKeyType = 5
)

func (c *Client) WithAuth(next http.Handler) http.Handler {
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
