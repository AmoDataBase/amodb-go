package amodb

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestQueryEncoding(t *testing.T) {
	var nilPtr *int
	five := 5
	when := time.Date(2020, time.March, 5, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		params Params
		want   map[string]string
		absent []string
	}{
		{"booleans", Params{"include_adult": true, "include_video": false}, map[string]string{"include_adult": "true", "include_video": "false"}, nil},
		{"string slice is comma joined", Params{"append_to_response": []string{"credits", "images"}}, map[string]string{"append_to_response": "credits,images"}, nil},
		{"int slice is comma joined", Params{"with_genres": []int{28, 12}}, map[string]string{"with_genres": "28,12"}, nil},
		{"mixed slice skips nil", Params{"with_keywords": []any{1, nil, "x"}}, map[string]string{"with_keywords": "1,x"}, nil},
		{"nil values are dropped", Params{"a": nil, "b": nilPtr, "page": 2}, map[string]string{"page": "2"}, []string{"a", "b"}},
		{"pointer is dereferenced", Params{"page": &five}, map[string]string{"page": "5"}, nil},
		{"floats", Params{"vote_average.gte": 7.5, "weight": float32(0.5)}, map[string]string{"vote_average.gte": "7.5", "weight": "0.5"}, nil},
		{"sized integers", Params{"a": int64(7), "b": uint(8), "c": int8(-3)}, map[string]string{"a": "7", "b": "8", "c": "-3"}, nil},
		{"time is a date", Params{"primary_release_date.gte": when}, map[string]string{"primary_release_date.gte": "2020-03-05"}, nil},
		{"special characters survive", Params{"query": "a b&c=d+e"}, map[string]string{"query": "a b&c=d+e"}, nil},
		{"empty slice is dropped", Params{"x": []string{}}, nil, []string{"x"}},
		{"nil params", nil, nil, []string{"language", "api_key"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := mustClient(t)
			base := c.baseURL + "/x"
			raw := c.buildURL(base, tc.params, scopeTMDB)
			u, err := url.Parse(raw)
			if err != nil {
				t.Fatalf("parse %q: %v", raw, err)
			}
			q := u.Query()
			for key, want := range tc.want {
				if got := q.Get(key); got != want {
					t.Errorf("%s = %q, want %q", key, got, want)
				}
			}
			for _, key := range tc.absent {
				if _, ok := q[key]; ok {
					t.Errorf("%s should be dropped, query is %q", key, u.RawQuery)
				}
			}
		})
	}
}

func TestBuildPath(t *testing.T) {
	tests := []struct {
		name     string
		segments []any
		want     string
	}{
		{"ints", []any{"movie", 550, "credits"}, "/movie/550/credits"},
		{"nested", []any{"tv", 1399, "season", 1, "episode", 2}, "/tv/1399/season/1/episode/2"},
		{"slash is escaped", []any{"find", "tt 01/2"}, "/find/tt%2001%2F2"},
		{"query characters are escaped", []any{"search", "a?b#c"}, "/search/a%3Fb%23c"},
		{"unicode is escaped", []any{"x", "caf\u00e9"}, "/x/caf%C3%A9"},
		{"other integer types", []any{"x", int64(5)}, "/x/5"},
		{"no segments", nil, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildPath(tc.segments...); got != tc.want {
				t.Errorf("buildPath = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAppendToResponseJoined(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, `{"id":550}`))
	params := Params{"append_to_response": []string{"credits", "images", "moods"}}
	if _, err := c.Movies.Details(context.Background(), 550, params); err != nil {
		t.Fatalf("Details: %v", err)
	}
	got := rec.last(t)
	if v := got.Query.Get("append_to_response"); v != "credits,images,moods" {
		t.Errorf("append_to_response = %q", v)
	}
	if got.Path != "/3/movie/550" {
		t.Errorf("path = %q", got.Path)
	}
}

func TestAppendToResponseIsPercentEncoded(t *testing.T) {
	c := mustClient(t)
	raw := c.buildURL("http://x.test/3/movie/1", Params{"append_to_response": []string{"credits", "images"}}, scopeTMDB)
	if !strings.Contains(raw, "append_to_response=credits%2Cimages") {
		t.Errorf("url = %q, want a percent encoded comma", raw)
	}
}

func TestPathEscapingOnTheWire(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, `{}`))
	ctx := context.Background()

	if _, err := c.Find(ctx, "tt 01/2", "imdb_id"); err != nil {
		t.Fatalf("Find: %v", err)
	}
	got := rec.last(t)
	if got.Escaped != "/3/find/tt%2001%2F2" {
		t.Errorf("escaped path = %q", got.Escaped)
	}
	if got.Path != "/3/find/tt 01/2" {
		t.Errorf("decoded path = %q", got.Path)
	}
	if v := got.Query.Get("external_source"); v != "imdb_id" {
		t.Errorf("external_source = %q", v)
	}

	if _, err := c.Trending.Get(ctx, "mo vie", "we/ek", nil); err != nil {
		t.Fatalf("Trending.Get: %v", err)
	}
	if got := rec.last(t).Escaped; got != "/3/trending/mo%20vie/we%2Fek" {
		t.Errorf("escaped path = %q", got)
	}
}

func TestKeyStaysOutOfTheURLInBearerMode(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, `{}`))
	if _, err := c.Configuration(context.Background()); err != nil {
		t.Fatalf("Configuration: %v", err)
	}
	got := rec.last(t)
	if _, ok := got.Query["api_key"]; ok {
		t.Error("api_key must not be in the URL in bearer mode")
	}
	if _, ok := got.Query["apikey"]; ok {
		t.Error("apikey must not be in the URL in bearer mode")
	}
}

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, time.October, 18, 12, 0, 0, 0, time.UTC)
	future := now.Add(30 * time.Second).Format(http.TimeFormat)
	past := now.Add(-time.Minute).Format(http.TimeFormat)
	tests := []struct {
		name   string
		value  string
		want   time.Duration
		wantOK bool
	}{
		{"empty", "", 0, false},
		{"seconds", "7", 7 * time.Second, true},
		{"seconds with spaces", " 7 ", 7 * time.Second, true},
		{"zero", "0", 0, true},
		{"negative", "-3", 0, true},
		{"http date in the future", future, 30 * time.Second, true},
		{"http date in the past", past, 0, true},
		{"not understood", "tomorrow", 0, false},
		{"huge value is capped at one day", "99999999999", 24 * time.Hour, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := parseRetryAfter(tc.value, now)
			if ok != tc.wantOK || got != tc.want {
				t.Errorf("parseRetryAfter(%q) = %v, %v; want %v, %v", tc.value, got, ok, tc.want, tc.wantOK)
			}
		})
	}
}

func TestRetryDelay(t *testing.T) {
	now := time.Date(2026, time.October, 18, 12, 0, 0, 0, time.UTC)
	future := now.Add(30 * time.Second).Format(http.TimeFormat)
	tests := []struct {
		name    string
		header  string
		attempt int
		want    time.Duration
	}{
		{"header seconds", "5", 0, 5 * time.Second},
		{"header wins over the attempt number", "5", 2, 5 * time.Second},
		{"long header is clamped to 60 seconds", "300", 0, 60 * time.Second},
		{"http date", future, 0, 30 * time.Second},
		{"missing header, attempt 0", "", 0, 1 * time.Second},
		{"missing header, attempt 1", "", 1, 2 * time.Second},
		{"missing header, attempt 2", "", 2, 4 * time.Second},
		{"missing header, attempt 3", "", 3, 8 * time.Second},
		{"missing header, huge attempt is clamped", "", 20, 60 * time.Second},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := http.Header{}
			if tc.header != "" {
				h.Set("Retry-After", tc.header)
			}
			if got := retryDelay(h, tc.attempt, now); got != tc.want {
				t.Errorf("retryDelay = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRetryOnRateLimit(t *testing.T) {
	past := time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat)
	tests := []struct {
		name       string
		retryAfter string
		failures   int
		wantWaits  []time.Duration
	}{
		{"one second", "1", 1, []time.Duration{time.Second}},
		{"two failures", "2", 2, []time.Duration{2 * time.Second, 2 * time.Second}},
		{"missing header waits 2^attempt", "", 3, []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}},
		{"zero", "0", 1, []time.Duration{0}},
		{"clamped to 60 seconds", "120", 1, []time.Duration{60 * time.Second}},
		{"negative becomes zero", "-4", 1, []time.Duration{0}},
		{"junk waits 2^attempt", "soon", 2, []time.Duration{1 * time.Second, 2 * time.Second}},
		{"http date in the past", past, 1, []time.Duration{0}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sleeps := &sleepRecorder{}
			c, rec := newTestClient(t, rateLimitedThenOK(tc.failures, tc.retryAfter), WithSleep(sleeps.sleep))
			movie, err := c.Movies.Details(context.Background(), 550, nil)
			if err != nil {
				t.Fatalf("Details: %v", err)
			}
			if movie.Title != "Fight Club" {
				t.Errorf("title = %q", movie.Title)
			}
			if got := rec.count(); got != tc.failures+1 {
				t.Errorf("requests = %d, want %d", got, tc.failures+1)
			}
			if got := sleeps.recorded(); !equalDurations(got, tc.wantWaits) {
				t.Errorf("waits = %v, want %v", got, tc.wantWaits)
			}
		})
	}
}

func TestRetryStopsAfterMaxRetries(t *testing.T) {
	tests := []struct {
		name         string
		opts         []Option
		wantRequests int
		wantWaits    int
	}{
		{"default is 3 retries", nil, 4, 3},
		{"two retries", []Option{WithMaxRetries(2)}, 3, 2},
		{"one retry", []Option{WithMaxRetries(1)}, 2, 1},
		{"no retries", []Option{WithMaxRetries(0)}, 1, 0},
		{"the cap is 3", []Option{WithMaxRetries(50)}, 4, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sleeps := &sleepRecorder{}
			opts := append([]Option{WithSleep(sleeps.sleep)}, tc.opts...)
			c, rec := newTestClient(t, rateLimitedThenOK(1000, "1"), opts...)
			_, err := c.Movies.Details(context.Background(), 550, nil)
			if err == nil {
				t.Fatal("expected a rate limit error")
			}
			if !IsRateLimited(err) {
				t.Errorf("IsRateLimited(%v) = false", err)
			}
			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("error %v is not an *Error", err)
			}
			if apiErr.Status != 429 || apiErr.StatusCode != 25 {
				t.Errorf("status = %d, status_code = %d", apiErr.Status, apiErr.StatusCode)
			}
			if apiErr.RetryAfter != time.Second {
				t.Errorf("RetryAfter = %v, want 1s", apiErr.RetryAfter)
			}
			if got := rec.count(); got != tc.wantRequests {
				t.Errorf("requests = %d, want %d", got, tc.wantRequests)
			}
			if got := len(sleeps.recorded()); got != tc.wantWaits {
				t.Errorf("waits = %d, want %d", got, tc.wantWaits)
			}
		})
	}
}

func TestNoRetryForOtherStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{"unauthorized", 401, `{"success":false,"status_code":7,"status_message":"Invalid API key: You must be granted a valid key."}`},
		{"expired plan", 402, `{"success":false,"status_code":36,"status_message":"Your plan has expired. Renew at https://amodb.dev/pricing"}`},
		{"not found", 404, `{"success":false,"status_code":34,"status_message":"The resource you requested could not be found."}`},
		{"server error", 500, `{"success":false,"status_code":11,"status_message":"Internal error: Something went wrong, contact amodb."}`},
		{"backend timeout", 504, `{"success":false,"status_code":24,"status_message":"Your request to the backend server timed out. Try again."}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sleeps := &sleepRecorder{}
			c, rec := newTestClient(t, jsonHandler(tc.status, tc.body), WithSleep(sleeps.sleep))
			_, err := c.Movies.Details(context.Background(), 550, nil)
			if err == nil {
				t.Fatal("expected an error")
			}
			if rec.count() != 1 {
				t.Errorf("requests = %d, want 1", rec.count())
			}
			if n := len(sleeps.recorded()); n != 0 {
				t.Errorf("waits = %d, want 0", n)
			}
		})
	}
}

func TestRetryWaitInterrupted(t *testing.T) {
	stop := func(context.Context, time.Duration) error {
		return context.Canceled
	}
	c, rec := newTestClient(t, rateLimitedThenOK(5, "1"), WithSleep(stop))
	_, err := c.Movies.Details(context.Background(), 550, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if rec.count() != 1 {
		t.Errorf("requests = %d, want 1", rec.count())
	}
}

func TestTimeoutWithBlockingTransport(t *testing.T) {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return nil, req.Context().Err()
	})
	c := mustClient(t,
		WithHTTPClient(&http.Client{Transport: transport}),
		WithTimeout(20*time.Millisecond),
	)
	start := time.Now()
	_, err := c.Movies.Details(context.Background(), 550, nil)
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("error = %v, want ErrTimeout", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want context.DeadlineExceeded", err)
	}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %v is not an *Error", err)
	}
	if apiErr.Status != 0 {
		t.Errorf("Status = %d, want 0", apiErr.Status)
	}
	if time.Since(start) > 5*time.Second {
		t.Error("the timeout did not fire")
	}
}

func TestTimeoutAgainstSlowServer(t *testing.T) {
	release := make(chan struct{})
	slow := func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}
	c, _ := newTestClient(t, slow, WithTimeout(50*time.Millisecond))
	defer close(release)
	_, err := c.Movies.Details(context.Background(), 550, nil)
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("error = %v, want ErrTimeout", err)
	}
}

func TestCallerDeadlineWins(t *testing.T) {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return nil, req.Context().Err()
	})
	c := mustClient(t, WithHTTPClient(&http.Client{Transport: transport}))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := c.Movies.Details(ctx, 550, nil)
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("error = %v, want ErrTimeout", err)
	}
}

func TestContextCanceled(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, `{}`))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.Movies.Details(ctx, 550, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if errors.Is(err, ErrTimeout) {
		t.Error("a canceled context must not look like a timeout")
	}
	if rec.count() != 0 {
		t.Errorf("requests = %d, want 0", rec.count())
	}
}

func TestNetworkError(t *testing.T) {
	transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})
	c := mustClient(t, WithHTTPClient(&http.Client{Transport: transport}))
	_, err := c.Movies.Details(context.Background(), 550, nil)
	if !errors.Is(err, ErrNetwork) {
		t.Fatalf("error = %v, want ErrNetwork", err)
	}
	if errors.Is(err, ErrTimeout) {
		t.Error("a refused connection must not look like a timeout")
	}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %v is not an *Error", err)
	}
	if apiErr.Status != 0 {
		t.Errorf("Status = %d, want 0", apiErr.Status)
	}
	if !strings.Contains(apiErr.StatusMessage, "connection refused") {
		t.Errorf("StatusMessage = %q", apiErr.StatusMessage)
	}
}

func TestDo(t *testing.T) {
	type item struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	var out struct {
		Page    int    `json:"page"`
		Results []item `json:"results"`
	}
	c, rec := newTestClient(t,
		jsonHandler(200, `{"page":2,"results":[{"id":1,"title":"A"}]}`),
		WithLanguage("ru-RU"),
	)
	err := c.Do(context.Background(), "/movie/popular", Params{"page": 2}, &out)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	got := rec.last(t)
	if got.Method != http.MethodGet {
		t.Errorf("method = %q", got.Method)
	}
	if got.Path != "/3/movie/popular" {
		t.Errorf("path = %q", got.Path)
	}
	if got.Query.Get("page") != "2" || got.Query.Get("language") != "ru-RU" {
		t.Errorf("query = %v", got.Query)
	}
	if got.Header.Get("Authorization") != "Bearer "+testKey {
		t.Errorf("Authorization = %q", got.Header.Get("Authorization"))
	}
	if out.Page != 2 || len(out.Results) != 1 || out.Results[0].Title != "A" {
		t.Errorf("decoded = %+v", out)
	}
}

func TestDoAddsLeadingSlashAndAcceptsNilOut(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, `{"ignored":true}`))
	if err := c.Do(context.Background(), "movie/popular", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/movie/popular" {
		t.Errorf("path = %q", got)
	}
}

func TestDoRejectsInvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"not json", "not json"},
		{"empty body", ""},
		{"wrong shape", `[1,2,3]`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := newTestClient(t, jsonHandler(200, tc.body))
			var out map[string]any
			err := c.Do(context.Background(), "/configuration", nil, &out)
			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("error %v is not an *Error", err)
			}
			if apiErr.Status != 200 {
				t.Errorf("Status = %d, want 200", apiErr.Status)
			}
			if !strings.Contains(apiErr.StatusMessage, "invalid JSON") {
				t.Errorf("StatusMessage = %q", apiErr.StatusMessage)
			}
			if apiErr.Unwrap() == nil {
				t.Error("the decode error should be wrapped")
			}
		})
	}
}
