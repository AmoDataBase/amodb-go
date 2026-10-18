package amodb

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// testKey is a made up key in the 32 hex character format.
const testKey = "0123456789abcdef0123456789abcdef"

const rateLimitBody = `{"success":false,"status_code":25,"status_message":"Your request count (21) is over the allowed limit of (20)."}`

type recordedRequest struct {
	Method  string
	Path    string
	Escaped string
	Query   url.Values
	Header  http.Header
}

// recorder keeps every request that reached the test server.
type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

func (r *recorder) add(req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reqs = append(r.reqs, recordedRequest{
		req.Method,
		req.URL.Path,
		req.URL.EscapedPath(),
		req.URL.Query(),
		req.Header.Clone(),
	})
}

func (r *recorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.reqs)
}

func (r *recorder) last(t *testing.T) recordedRequest {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.reqs) == 0 {
		t.Fatal("the server received no request")
	}
	return r.reqs[len(r.reqs)-1]
}

// sleepRecorder replaces the real sleep and keeps the requested waits.
type sleepRecorder struct {
	mu    sync.Mutex
	waits []time.Duration
}

func (s *sleepRecorder) sleep(ctx context.Context, d time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.waits = append(s.waits, d)
	return ctx.Err()
}

func (s *sleepRecorder) recorded() []time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]time.Duration, len(s.waits))
	copy(out, s.waits)
	return out
}

func equalDurations(a, b []time.Duration) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// roundTripFunc lets a test supply the transport as a function.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// errOf drops the first result of a call and keeps the error.
func errOf(_ any, err error) error {
	return err
}

func jsonHandler(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}
}

// rateLimitedThenOK answers HTTP 429 for the first failures requests and then
// answers 200. An empty retryAfter leaves the header out.
func rateLimitedThenOK(failures int, retryAfter string) http.HandlerFunc {
	var calls atomic.Int32
	return func(w http.ResponseWriter, _ *http.Request) {
		n := int(calls.Add(1))
		w.Header().Set("Content-Type", "application/json")
		if n <= failures {
			if retryAfter != "" {
				w.Header().Set("Retry-After", retryAfter)
			}
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, rateLimitBody)
			return
		}
		_, _ = io.WriteString(w, `{"id":550,"title":"Fight Club"}`)
	}
}

// newTestClient starts a test server and returns a client that points at it
// for both the TMDB style and the OMDb style endpoints.
func newTestClient(t *testing.T, h http.HandlerFunc, opts ...Option) (*Client, *recorder) {
	t.Helper()
	rec := &recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.add(r)
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	base := []Option{WithBaseURL(srv.URL + "/3"), WithOMDbURL(srv.URL + "/omdb/")}
	c, err := NewClient(testKey, append(base, opts...)...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c, rec
}

// mustClient builds a client for tests that never touch the network.
func mustClient(t *testing.T, opts ...Option) *Client {
	t.Helper()
	c, err := NewClient(testKey, opts...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}
