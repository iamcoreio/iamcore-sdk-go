package iamcore

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

const (
	idleConnTimeout     = 30 * time.Second
	requestTimeout      = 30 * time.Second
	dialTimeout         = 5 * time.Second
	keepAliveProbe      = 30 * time.Second
	tlsHandshakeTimeout = 5 * time.Second
	maxIdleConns        = 100
	maxIdleConnsPerHost = 10
	maxConnRetries      = 1
)

var errNoRewind = errors.New("iamcore: request body cannot be rewound for retry")

func newDefaultHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: keepAliveProbe,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          maxIdleConns,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		IdleConnTimeout:       idleConnTimeout,
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: &retryRoundTripper{base: transport, maxRetries: maxConnRetries},
		Timeout:   requestTimeout,
	}
}

type retryRoundTripper struct {
	base       http.RoundTripper
	maxRetries int
}

func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.base.RoundTrip(req)
	if err == nil || !isStaleConnError(err) {
		return resp, err
	}

	for attempt := 0; attempt < rt.maxRetries; attempt++ {
		retryReq, rewindErr := rewindRequest(req)
		if rewindErr != nil {
			return resp, err
		}

		resp, err = rt.base.RoundTrip(retryReq)
		if err == nil || !isStaleConnError(err) {
			break
		}
	}

	return resp, err
}

func rewindRequest(req *http.Request) (*http.Request, error) {
	clone := req.Clone(req.Context())

	if req.Body == nil || req.Body == http.NoBody {
		return clone, nil
	}

	if req.GetBody == nil {
		return nil, errNoRewind
	}

	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	clone.Body = body

	return clone, nil
}

func isStaleConnError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.EPIPE) {
		return true
	}

	// Go's transport reports some of these conditions only as error strings.
	msg := err.Error()

	return strings.Contains(msg, "server closed idle connection") ||
		strings.Contains(msg, "connection reset by peer") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "unexpected EOF")
}
