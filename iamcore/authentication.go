package iamcore

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

var ErrNoAuthContext = errors.New("no auth context")

type AuthenticationClient interface {
	// WithAuth creates http.Handler middleware for authenticating the incoming request by means of either OAuth 2.0 Access Token
	// or "X-iamcore-API-Key" HTTP header. This handler should precede the application request handling in the handlers chain.
	// It populates the request context with validated requester principal's IRN for further use.
	// Returns 401 Unauthorized HTTP error in case of unauthorized access, and stops HTTP request propagation.
	WithAuth(next http.Handler) http.Handler
}

// contextKeyType is a context.Context key type.
type contextKeyType int

const (
	authorizationHeaderKey contextKeyType = 0
	principalIRNKey        contextKeyType = 1
)

func (c *сlient) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.disabled {
			next.ServeHTTP(w, r)

			return
		}

		for i := range c.authenticators {
			principal, authorizationHeader, err := c.authenticators[i].Authenticate(r.Context(), r.Header)
			switch {
			case err != nil && errors.Is(err, ErrUnauthenticated):
				writeResponseMessage(w, http.StatusUnauthorized, err.Error())

				return
			case err != nil:
				writeResponseMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

				return
			}

			if principal != nil {
				ctx := r.Context()

				ctx = context.WithValue(ctx, principalIRNKey, principal)
				ctx = context.WithValue(ctx, authorizationHeaderKey, authorizationHeader)

				r = r.WithContext(ctx)

				// Pass control to the next handler
				next.ServeHTTP(w, r)

				return
			}
		}

		writeResponseMessage(w, http.StatusUnauthorized, "Failed to authenticate request with any of available authenticators")
	})
}

// PrincipalIRN extracts and returns principal's IRN from the request context.
func PrincipalIRN(ctx context.Context) (*irn.IRN, error) {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return nil, ErrNoAuthContext
	}

	return principal, nil
}

// AccountID extracts and returns principal's account ID from the request context.
func AccountID(ctx context.Context) (string, error) {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return "", ErrNoAuthContext
	}

	return principal.GetAccountID(), nil
}

// TenantID extracts and returns principal's tenant ID from the request context.
func TenantID(ctx context.Context) (string, error) {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return "", ErrNoAuthContext
	}

	return principal.GetTenantID(), nil
}

// Path extracts and returns principal's path from the request context.
func Path(ctx context.Context) (string, error) {
	principal, ok := ctx.Value(principalIRNKey).(*irn.IRN)
	if !ok {
		return "", ErrNoAuthContext
	}

	return principal.GetPath(), nil
}

// AuthorizationHeader extracts and returns authorization header from the request context.
func AuthorizationHeader(ctx context.Context) (http.Header, error) {
	authorizationHeader, ok := ctx.Value(authorizationHeaderKey).(http.Header)
	if !ok {
		return nil, ErrNoAuthContext
	}

	return authorizationHeader, nil
}

func writeResponseMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	responseDTO := &ErrorResponseDTO{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(responseDTO); err != nil {
		log.Printf("Error writing response message: %v", err)
	}
}
