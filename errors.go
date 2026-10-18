package amodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ErrTimeout is wrapped by the error returned when a request times out.
// Test for it with errors.Is(err, amodb.ErrTimeout).
var ErrTimeout = errors.New("amodb: request timed out")

// ErrNetwork is wrapped by the error returned when no response arrived,
// for example after a DNS failure or a refused connection.
var ErrNetwork = errors.New("amodb: network error")

// ErrMissingAPIKey is returned by NewClient when the API key is empty.
var ErrMissingAPIKey = errors.New("amodb: API key is required")

// Error is the error type returned for every failed call.
type Error struct {
	// Status is the HTTP status code. It is 0 when no response arrived.
	Status int
	// StatusCode is the status_code field of the error body, or 0 when absent.
	StatusCode int
	// StatusMessage is the status_message field of the error body. When the
	// body has none, it holds the OMDb Error text or the HTTP status text.
	StatusMessage string
	// Body is the raw response body.
	Body []byte
	// RetryAfter is the parsed Retry-After header, or 0 when absent.
	RetryAfter time.Duration

	cause error
}

// Error implements the error interface.
func (e *Error) Error() string {
	msg := e.StatusMessage
	if msg == "" {
		msg = "request failed"
	}
	if e.Status == 0 {
		return "amodb: " + msg
	}
	if e.StatusCode != 0 {
		return fmt.Sprintf("amodb: HTTP %d (code %d): %s", e.Status, e.StatusCode, msg)
	}
	return fmt.Sprintf("amodb: HTTP %d: %s", e.Status, msg)
}

// Unwrap returns the underlying cause. For timeouts and network failures it
// matches ErrTimeout or ErrNetwork with errors.Is.
func (e *Error) Unwrap() error {
	return e.cause
}

// IsRateLimited reports whether err is an HTTP 429 answer.
func IsRateLimited(err error) bool {
	return hasStatus(err, http.StatusTooManyRequests)
}

// IsNotFound reports whether err is an HTTP 404 answer.
func IsNotFound(err error) bool {
	return hasStatus(err, http.StatusNotFound)
}

// IsUnauthorized reports whether err is an HTTP 401 answer.
func IsUnauthorized(err error) bool {
	return hasStatus(err, http.StatusUnauthorized)
}

func hasStatus(err error, status int) bool {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.Status == status
}

// newHTTPError builds an *Error from a non 2xx response. Bodies that are not
// JSON, or that are empty, still produce an error with the HTTP status.
func newHTTPError(res *response, now time.Time) *Error {
	var parsed struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
		Error         string `json:"Error"`
	}
	_ = json.Unmarshal(res.body, &parsed)

	e := &Error{Status: res.status, Body: res.body}
	e.StatusCode = parsed.StatusCode
	e.StatusMessage = parsed.StatusMessage
	if e.StatusMessage == "" {
		e.StatusMessage = parsed.Error
	}
	if e.StatusMessage == "" {
		e.StatusMessage = http.StatusText(res.status)
	}
	if wait, ok := parseRetryAfter(res.header.Get("Retry-After"), now); ok {
		e.RetryAfter = wait
	}
	return e
}
