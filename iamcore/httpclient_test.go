package iamcore

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// stubRoundTripper returns pre-programmed results per call and records the
// request body seen on every attempt, so tests can assert both the retry count
// and that a replayed request carries the full body.
type stubRoundTripper struct {
	results []func() (*http.Response, error)
	calls   int
	bodies  []string
}

func (s *stubRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		s.bodies = append(s.bodies, string(b))
	} else {
		s.bodies = append(s.bodies, "")
	}

	idx := s.calls
	s.calls++

	if idx >= len(s.results) {
		return nil, errors.New("stub: unexpected extra call")
	}

	return s.results[idx]()
}

func okResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
		Header:     make(http.Header),
	}
}

func TestRetryRoundTripperRetriesOnStaleConn(t *testing.T) {
	stub := &stubRoundTripper{results: []func() (*http.Response, error){
		func() (*http.Response, error) { return nil, io.EOF },
		func() (*http.Response, error) { return okResponse(), nil },
	}}
	rt := &retryRoundTripper{base: stub, maxRetries: 1}

	req, _ := http.NewRequest(http.MethodPost, "http://iamcore/api/v1/evaluate",
		strings.NewReader(`{"a":1}`))

	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if stub.calls != 2 {
		t.Fatalf("expected exactly 2 attempts (1 retry), got %d", stub.calls)
	}
}

func TestRetryRoundTripperRewindsBodyOnRetry(t *testing.T) {
	payload := `{"resources":["a","b"],"action":"read"}`
	stub := &stubRoundTripper{results: []func() (*http.Response, error){
		func() (*http.Response, error) { return nil, io.ErrUnexpectedEOF },
		func() (*http.Response, error) { return okResponse(), nil },
	}}
	rt := &retryRoundTripper{base: stub, maxRetries: 1}

	req, _ := http.NewRequest(http.MethodPost, "http://iamcore/x", strings.NewReader(payload))

	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stub.bodies) != 2 {
		t.Fatalf("expected 2 recorded bodies, got %d", len(stub.bodies))
	}
	if stub.bodies[0] != payload || stub.bodies[1] != payload {
		t.Fatalf("body not fully replayed on retry: first=%q retry=%q want=%q",
			stub.bodies[0], stub.bodies[1], payload)
	}
}

func TestRetryRoundTripperNoRetryOnOtherError(t *testing.T) {
	sentinel := errors.New("bad request")
	stub := &stubRoundTripper{results: []func() (*http.Response, error){
		func() (*http.Response, error) { return nil, sentinel },
	}}
	rt := &retryRoundTripper{base: stub, maxRetries: 1}

	req, _ := http.NewRequest(http.MethodGet, "http://iamcore/x", nil)

	_, err := rt.RoundTrip(req)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected original error, got %v", err)
	}
	if stub.calls != 1 {
		t.Fatalf("expected no retry on a non-connection error, got %d attempts", stub.calls)
	}
}

func TestRetryRoundTripperNoRetryOnSuccess(t *testing.T) {
	stub := &stubRoundTripper{results: []func() (*http.Response, error){
		func() (*http.Response, error) { return okResponse(), nil },
	}}
	rt := &retryRoundTripper{base: stub, maxRetries: 1}

	req, _ := http.NewRequest(http.MethodGet, "http://iamcore/x", nil)

	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.calls != 1 {
		t.Fatalf("expected a single attempt on success, got %d", stub.calls)
	}
}

func TestIsStaleConnError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"eof", io.EOF, true},
		{"unexpected eof", io.ErrUnexpectedEOF, true},
		{"server closed idle", errors.New(`Post "http://x": http: server closed idle connection`), true},
		{"connection reset", errors.New("read tcp: connection reset by peer"), true},
		{"generic", errors.New("some business error"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isStaleConnError(c.err); got != c.want {
				t.Fatalf("isStaleConnError(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

// TestNewDefaultHTTPClientKeepsPoolingBoundsIdle documents the intent of the
// fix: connection pooling stays ON, but idle connections are reaped before a
// typical ingress (~60s keep-alive) would silently close them, and every
// request is bounded by a total timeout.
func TestNewDefaultHTTPClientKeepsPoolingBoundsIdle(t *testing.T) {
	c := newDefaultHTTPClient()

	if c.Timeout <= 0 {
		t.Fatal("client must set a total request Timeout (http.DefaultClient has none)")
	}

	rt, ok := c.Transport.(*retryRoundTripper)
	if !ok {
		t.Fatalf("transport must be *retryRoundTripper, got %T", c.Transport)
	}

	tr, ok := rt.base.(*http.Transport)
	if !ok {
		t.Fatalf("base transport must be *http.Transport, got %T", rt.base)
	}
	if tr.IdleConnTimeout <= 0 || tr.IdleConnTimeout >= 60*time.Second {
		t.Fatalf("IdleConnTimeout must be set below the typical 60s ingress keep-alive, got %v",
			tr.IdleConnTimeout)
	}
	if tr.MaxIdleConnsPerHost <= 0 {
		t.Fatalf("connection pooling must stay enabled (MaxIdleConnsPerHost>0), got %d",
			tr.MaxIdleConnsPerHost)
	}
}

// TestClientRetriesWhenServerDropsConnection exercises the whole client stack
// against a server that drops the first connection before answering (the
// production idle-keep-alive race) and answers the retry.
func TestClientRetriesWhenServerDropsConnection(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&hits, 1) == 1 {
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Error("test server does not support hijacking")
				return
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				t.Errorf("hijack failed: %v", err)
				return
			}
			_ = conn.Close() // drop the connection without responding
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	}))
	defer srv.Close()

	client := newDefaultHTTPClient()

	req, _ := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(`{"a":1}`))
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 after retry, got %d", resp.StatusCode)
	}
	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("expected server to be hit twice (drop + retry), got %d", got)
	}
}
