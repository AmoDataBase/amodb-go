package amodb

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClientEmptyKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"empty", ""},
		{"spaces", "   "},
		{"newline", "\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewClient(tc.key)
			if err == nil {
				t.Fatal("expected an error for an empty key")
			}
			if !errors.Is(err, ErrMissingAPIKey) {
				t.Errorf("error = %v, want ErrMissingAPIKey", err)
			}
			if c != nil {
				t.Errorf("client = %v, want nil", c)
			}
		})
	}
}

func TestNewClientDefaults(t *testing.T) {
	c := mustClient(t)
	if c.baseURL != "https://api.amodb.dev/3" {
		t.Errorf("baseURL = %q", c.baseURL)
	}
	if c.omdbURL != "https://api.amodb.dev/omdb/" {
		t.Errorf("omdbURL = %q", c.omdbURL)
	}
	if c.imageBaseURL != "https://img.amodb.dev/t/p/" {
		t.Errorf("imageBaseURL = %q", c.imageBaseURL)
	}
	if c.timeout != 15*time.Second {
		t.Errorf("timeout = %v, want 15s", c.timeout)
	}
	if c.maxRetries != 3 {
		t.Errorf("maxRetries = %d, want 3", c.maxRetries)
	}
	if c.auth != AuthBearer {
		t.Errorf("auth = %v, want AuthBearer", c.auth)
	}
	if c.userAgent != "amodb-go/0.1.0" {
		t.Errorf("userAgent = %q", c.userAgent)
	}
	if c.language != "" {
		t.Errorf("language = %q, want empty", c.language)
	}
	services := []struct {
		name string
		ok   bool
	}{
		{"Movies", c.Movies != nil},
		{"TV", c.TV != nil},
		{"People", c.People != nil},
		{"Search", c.Search != nil},
		{"Discover", c.Discover != nil},
		{"Trending", c.Trending != nil},
		{"Genres", c.Genres != nil},
		{"Collections", c.Collections != nil},
		{"Companies", c.Companies != nil},
		{"Networks", c.Networks != nil},
		{"Keywords", c.Keywords != nil},
		{"Certifications", c.Certifications != nil},
		{"WatchProviders", c.WatchProviders != nil},
		{"OMDb", c.OMDb != nil},
	}
	for _, svc := range services {
		if !svc.ok {
			t.Errorf("service %s is nil", svc.name)
		}
	}
}

func TestOptionsApplied(t *testing.T) {
	hc := &http.Client{}
	c := mustClient(t,
		WithBaseURL("http://localhost:9/3/"),
		WithOMDbURL("http://localhost:9/omdb"),
		WithImageBaseURL("http://localhost:9/img/"),
		WithLanguage("ru-RU"),
		WithTimeout(2*time.Second),
		WithAuth(AuthQuery),
		WithHTTPClient(hc),
		WithUserAgent("my-app/1.0"),
	)
	if c.baseURL != "http://localhost:9/3" {
		t.Errorf("baseURL = %q, want trailing slash removed", c.baseURL)
	}
	if c.omdbURL != "http://localhost:9/omdb" {
		t.Errorf("omdbURL = %q", c.omdbURL)
	}
	if c.imageBaseURL != "http://localhost:9/img/" {
		t.Errorf("imageBaseURL = %q", c.imageBaseURL)
	}
	if c.language != "ru-RU" {
		t.Errorf("language = %q", c.language)
	}
	if c.timeout != 2*time.Second {
		t.Errorf("timeout = %v", c.timeout)
	}
	if c.auth != AuthQuery {
		t.Errorf("auth = %v", c.auth)
	}
	if c.httpClient != hc {
		t.Error("custom http client was not used")
	}
	if c.userAgent != "my-app/1.0" {
		t.Errorf("userAgent = %q", c.userAgent)
	}
}

func TestOptionsIgnoreEmptyValues(t *testing.T) {
	c := mustClient(t,
		WithBaseURL(""),
		WithOMDbURL(""),
		WithImageBaseURL(""),
		WithUserAgent(""),
		WithHTTPClient(nil),
		WithSleep(nil),
	)
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q", c.baseURL)
	}
	if c.omdbURL != DefaultOMDbURL {
		t.Errorf("omdbURL = %q", c.omdbURL)
	}
	if c.imageBaseURL != DefaultImageBaseURL {
		t.Errorf("imageBaseURL = %q", c.imageBaseURL)
	}
	if c.httpClient == nil {
		t.Error("httpClient is nil")
	}
	if c.sleep == nil {
		t.Error("sleep is nil")
	}
}

func TestMaxRetriesClamp(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"negative", -5, 0},
		{"zero", 0, 0},
		{"one", 1, 1},
		{"three", 3, 3},
		{"above the cap", 4, 3},
		{"far above the cap", 100, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := mustClient(t, WithMaxRetries(tc.in))
			if c.maxRetries != tc.want {
				t.Errorf("maxRetries = %d, want %d", c.maxRetries, tc.want)
			}
		})
	}
}

func TestImageURL(t *testing.T) {
	c := mustClient(t)
	custom := mustClient(t, WithImageBaseURL("http://cdn.test/t/p"))
	tests := []struct {
		name   string
		client *Client
		path   string
		size   string
		want   string
	}{
		{"default size", c, "/abc.jpg", "", "https://img.amodb.dev/t/p/w500/abc.jpg"},
		{"explicit size", c, "/abc.jpg", "w92", "https://img.amodb.dev/t/p/w92/abc.jpg"},
		{"original", c, "/abc.jpg", "original", "https://img.amodb.dev/t/p/original/abc.jpg"},
		{"adds a leading slash", c, "abc.jpg", "w185", "https://img.amodb.dev/t/p/w185/abc.jpg"},
		{"empty path", c, "", "w500", ""},
		{"blank path", c, "  ", "w500", ""},
		{"custom base without slash", custom, "/x.png", "w300", "http://cdn.test/t/p/w300/x.png"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.client.ImageURL(tc.path, tc.size)
			if got != tc.want {
				t.Errorf("ImageURL(%q, %q) = %q, want %q", tc.path, tc.size, got, tc.want)
			}
		})
	}
}

func TestAuthModes(t *testing.T) {
	bearer := "Bearer " + testKey
	tests := []struct {
		name       string
		opts       []Option
		wantHeader string
		wantQuery  string
	}{
		{"bearer by default", nil, bearer, ""},
		{"explicit bearer", []Option{WithAuth(AuthBearer)}, bearer, ""},
		{"query", []Option{WithAuth(AuthQuery)}, "", testKey},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := newTestClient(t, jsonHandler(200, `{"id":550}`), tc.opts...)
			if _, err := c.Movies.Details(context.Background(), 550, nil); err != nil {
				t.Fatalf("Details: %v", err)
			}
			got := rec.last(t)
			if h := got.Header.Get("Authorization"); h != tc.wantHeader {
				t.Errorf("Authorization = %q, want %q", h, tc.wantHeader)
			}
			if q := got.Query.Get("api_key"); q != tc.wantQuery {
				t.Errorf("api_key = %q, want %q", q, tc.wantQuery)
			}
		})
	}
}

func TestLanguage(t *testing.T) {
	tests := []struct {
		name   string
		opts   []Option
		params Params
		want   string
	}{
		{"no language", nil, nil, ""},
		{"client default", []Option{WithLanguage("ru-RU")}, nil, "ru-RU"},
		{"call overrides default", []Option{WithLanguage("ru-RU")}, Params{"language": "pt-BR"}, "pt-BR"},
		{"call without default", nil, Params{"language": "fr-FR"}, "fr-FR"},
		{"nil language keeps default", []Option{WithLanguage("uk-UA")}, Params{"language": nil}, "uk-UA"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := newTestClient(t, jsonHandler(200, `{"id":550}`), tc.opts...)
			if _, err := c.Movies.Details(context.Background(), 550, tc.params); err != nil {
				t.Fatalf("Details: %v", err)
			}
			got := rec.last(t).Query.Get("language")
			if got != tc.want {
				t.Errorf("language = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestUserAgentHeader(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want string
	}{
		{"default", nil, "amodb-go/0.1.0"},
		{"custom", []Option{WithUserAgent("shop-backend/2")}, "shop-backend/2"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := newTestClient(t, jsonHandler(200, `{}`), tc.opts...)
			if _, err := c.Configuration(context.Background()); err != nil {
				t.Fatalf("Configuration: %v", err)
			}
			if got := rec.last(t).Header.Get("User-Agent"); got != tc.want {
				t.Errorf("User-Agent = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCustomHTTPClient(t *testing.T) {
	var calls atomic.Int32
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, errors.New("blocked by test transport")
	})
	c := mustClient(t, WithHTTPClient(&http.Client{Transport: transport}))
	_, err := c.Movies.Details(context.Background(), 550, nil)
	if err == nil {
		t.Fatal("expected an error from the test transport")
	}
	if calls.Load() != 1 {
		t.Errorf("transport calls = %d, want 1", calls.Load())
	}
	if !errors.Is(err, ErrNetwork) {
		t.Errorf("error = %v, want ErrNetwork", err)
	}
}

func TestPlans(t *testing.T) {
	body := `{"currency":"USD","periods":[{"id":"m1","months":1,"save":0}],` +
		`"plans":[{"id":"lite","name":"Lite","rpm":20,"monthlyCents":900,"totals":{"m1":900}}],` +
		`"fetchedAt":"2026-10-18T10:00:00Z"}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	plans, err := c.Plans(context.Background())
	if err != nil {
		t.Fatalf("Plans: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/plans" {
		t.Errorf("path = %q, want /plans", got.Path)
	}
	if got.Header.Get("Authorization") != "" {
		t.Error("the plans call must not send a key")
	}
	if got.Query.Get("api_key") != "" || got.Query.Get("apikey") != "" {
		t.Error("the plans call must not put a key in the URL")
	}
	if plans.Currency != "USD" {
		t.Errorf("currency = %q", plans.Currency)
	}
	if len(plans.Plans) != 1 || plans.Plans[0].ID != "lite" || plans.Plans[0].RPM != 20 {
		t.Errorf("plans = %+v", plans.Plans)
	}
	if plans.Plans[0].Totals["m1"] != 900 {
		t.Errorf("totals = %v", plans.Plans[0].Totals)
	}
	if len(plans.Periods) != 1 || plans.Periods[0].Months != 1 {
		t.Errorf("periods = %+v", plans.Periods)
	}
}

func TestPlansBadBaseURL(t *testing.T) {
	c := mustClient(t, WithBaseURL("not a url"))
	if _, err := c.Plans(context.Background()); err == nil {
		t.Fatal("expected an error when the origin cannot be derived")
	}
}
