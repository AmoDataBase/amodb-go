package amodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestErrorMapping(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		code    int
		message string
	}{
		{"invalid key", 401, 7, "Invalid API key: You must be granted a valid key."},
		{"suspended key", 401, 10, "Suspended API key: Access to your account has been suspended, contact amodb."},
		{"expired plan", 402, 36, "Your plan has expired. Renew at https://amodb.dev/pricing"},
		{"not found", 404, 34, "The resource you requested could not be found."},
		{"method not offered", 405, 37, "This method is not offered by amodb."},
		{"rate limited", 429, 25, "Your request count (21) is over the allowed limit of (20)."},
		{"internal error", 500, 11, "Internal error: Something went wrong, contact amodb."},
		{"backend timeout", 504, 24, "Your request to the backend server timed out. Try again."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"success":false,"status_code":%d,"status_message":%q}`, tc.code, tc.message)
			c, _ := newTestClient(t, jsonHandler(tc.status, body), WithMaxRetries(0))
			_, err := c.Movies.Details(context.Background(), 550, nil)
			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("error %v is not an *Error", err)
			}
			if apiErr.Status != tc.status {
				t.Errorf("Status = %d, want %d", apiErr.Status, tc.status)
			}
			if apiErr.StatusCode != tc.code {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tc.code)
			}
			if apiErr.StatusMessage != tc.message {
				t.Errorf("StatusMessage = %q, want %q", apiErr.StatusMessage, tc.message)
			}
			if string(apiErr.Body) != body {
				t.Errorf("Body = %q, want %q", apiErr.Body, body)
			}
			if got := IsRateLimited(err); got != (tc.status == 429) {
				t.Errorf("IsRateLimited = %v", got)
			}
			if got := IsNotFound(err); got != (tc.status == 404) {
				t.Errorf("IsNotFound = %v", got)
			}
			if got := IsUnauthorized(err); got != (tc.status == 401) {
				t.Errorf("IsUnauthorized = %v", got)
			}
			text := err.Error()
			if !strings.Contains(text, tc.message) || !strings.Contains(text, fmt.Sprint(tc.status)) {
				t.Errorf("Error() = %q", text)
			}
		})
	}
}

func TestNonJSONErrorBody(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"html page", 502, "<html>Bad gateway</html>", "Bad Gateway"},
		{"empty body", 500, "", "Internal Server Error"},
		{"plain text", 503, "upstream down", "Service Unavailable"},
		{"json object without fields", 400, `{}`, "Bad Request"},
		{"json array", 400, `[1]`, "Bad Request"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := newTestClient(t, jsonHandler(tc.status, tc.body), WithMaxRetries(0))
			_, err := c.Movies.Details(context.Background(), 550, nil)
			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("error %v is not an *Error", err)
			}
			if apiErr.Status != tc.status {
				t.Errorf("Status = %d, want %d", apiErr.Status, tc.status)
			}
			if apiErr.StatusCode != 0 {
				t.Errorf("StatusCode = %d, want 0", apiErr.StatusCode)
			}
			if apiErr.StatusMessage != tc.want {
				t.Errorf("StatusMessage = %q, want %q", apiErr.StatusMessage, tc.want)
			}
			if string(apiErr.Body) != tc.body {
				t.Errorf("Body = %q, want %q", apiErr.Body, tc.body)
			}
		})
	}
}

func TestRateLimitErrorWithTextBody(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, "slow down")
	}
	c, _ := newTestClient(t, handler, WithMaxRetries(0))
	_, err := c.Movies.Popular(context.Background(), nil)
	if !IsRateLimited(err) {
		t.Fatalf("IsRateLimited(%v) = false", err)
	}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %v is not an *Error", err)
	}
	if apiErr.RetryAfter != 30*time.Second {
		t.Errorf("RetryAfter = %v, want 30s", apiErr.RetryAfter)
	}
	if string(apiErr.Body) != "slow down" {
		t.Errorf("Body = %q", apiErr.Body)
	}
}

func TestErrorString(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{"status code and message", &Error{Status: 404, StatusCode: 34, StatusMessage: "The resource you requested could not be found."}, "amodb: HTTP 404 (code 34): The resource you requested could not be found."},
		{"status and message", &Error{Status: 502, StatusMessage: "Bad Gateway"}, "amodb: HTTP 502: Bad Gateway"},
		{"transport failure", &Error{StatusMessage: "request timed out"}, "amodb: request timed out"},
		{"status only", &Error{Status: 500}, "amodb: HTTP 500: request failed"},
		{"empty", &Error{}, "amodb: request failed"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHelpersSeeWrappedErrors(t *testing.T) {
	wrapped := func(status int) error {
		return fmt.Errorf("lookup failed: %w", &Error{Status: status})
	}
	tests := []struct {
		name         string
		err          error
		rateLimited  bool
		notFound     bool
		unauthorized bool
	}{
		{"nil", nil, false, false, false},
		{"plain error", errors.New("boom"), false, false, false},
		{"wrapped 429", wrapped(429), true, false, false},
		{"wrapped 404", wrapped(404), false, true, false},
		{"wrapped 401", wrapped(401), false, false, true},
		{"wrapped 500", wrapped(500), false, false, false},
		{"transport failure", &Error{StatusMessage: "x"}, false, false, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsRateLimited(tc.err); got != tc.rateLimited {
				t.Errorf("IsRateLimited = %v, want %v", got, tc.rateLimited)
			}
			if got := IsNotFound(tc.err); got != tc.notFound {
				t.Errorf("IsNotFound = %v, want %v", got, tc.notFound)
			}
			if got := IsUnauthorized(tc.err); got != tc.unauthorized {
				t.Errorf("IsUnauthorized = %v, want %v", got, tc.unauthorized)
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	if (&Error{Status: 404}).Unwrap() != nil {
		t.Error("an HTTP error has no cause")
	}
	cause := errors.New("root cause")
	e := &Error{Status: 200, cause: cause}
	if !errors.Is(e, cause) {
		t.Error("errors.Is should reach the cause")
	}
}
