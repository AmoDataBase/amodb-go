package amodb

import (
	"context"
	"errors"
	"testing"
	"time"
)

const omdbTitleBody = `{"Title":"Fight Club","Year":"1999","Rated":"R","Released":"15 Oct 1999","Runtime":"139 min","Genre":"Drama","Director":"Jane Doe","Writer":"John Roe","Actors":"A, B, C","Plot":"An office worker starts a club.","Language":"English","Country":"United States","Awards":"Won 1 Oscar.","Poster":"https://img.amodb.dev/t/p/w300/abc.jpg","Ratings":[{"Source":"Example source","Value":"8.8/10"}],"Metascore":"67","imdbRating":"8.8","imdbVotes":"2,300,000","imdbID":"tt0137523","Type":"movie","DVD":"N/A","BoxOffice":"$37,030,102","Production":"N/A","Website":"N/A","Response":"True"}`

func TestOMDbByID(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, omdbTitleBody), WithLanguage("ru-RU"))
	title, err := c.OMDb.ByID(context.Background(), "tt0137523", Params{"plot": "full", "year": 1999})
	if err != nil {
		t.Fatalf("ByID: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/omdb/" {
		t.Errorf("path = %q, want /omdb/", got.Path)
	}
	wantQuery := map[string]string{
		"apikey": testKey,
		"i":      "tt0137523",
		"plot":   "full",
		"y":      "1999",
	}
	for key, want := range wantQuery {
		if v := got.Query.Get(key); v != want {
			t.Errorf("%s = %q, want %q", key, v, want)
		}
	}
	if _, ok := got.Query["year"]; ok {
		t.Error("the year alias must be mapped to y")
	}
	if _, ok := got.Query["language"]; ok {
		t.Error("language must not be sent to the OMDb endpoint")
	}
	if got.Header.Get("Authorization") != "" {
		t.Error("OMDb calls must not send a bearer header")
	}
	if title.Title != "Fight Club" || title.Year != "1999" || title.IMDbID != "tt0137523" {
		t.Errorf("title = %+v", title)
	}
	if title.IMDbRating != "8.8" || title.IMDbVotes != "2,300,000" || title.Metascore != "67" {
		t.Errorf("ratings = %+v", title)
	}
	if len(title.Ratings) != 1 || title.Ratings[0].Value != "8.8/10" {
		t.Errorf("Ratings = %+v", title.Ratings)
	}
	if title.Response != "True" {
		t.Errorf("Response = %q", title.Response)
	}
}

func TestOMDbByTitle(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, omdbTitleBody))
	title, err := c.OMDb.ByTitle(context.Background(), "Fight Club", Params{"type": "movie", "y": 1999, "plot": "short"})
	if err != nil {
		t.Fatalf("ByTitle: %v", err)
	}
	got := rec.last(t)
	wantQuery := map[string]string{
		"t":    "Fight Club",
		"type": "movie",
		"y":    "1999",
		"plot": "short",
	}
	for key, want := range wantQuery {
		if v := got.Query.Get(key); v != want {
			t.Errorf("%s = %q, want %q", key, v, want)
		}
	}
	if title.Title != "Fight Club" {
		t.Errorf("Title = %q", title.Title)
	}
}

func TestOMDbSearch(t *testing.T) {
	body := `{"Search":[{"Title":"Fight Club","Year":"1999","imdbID":"tt0137523","Type":"movie","Poster":"https://img.amodb.dev/t/p/w300/abc.jpg"}],"totalResults":"1","Response":"True"}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	res, err := c.OMDb.Search(context.Background(), "fight", Params{"type": "movie", "page": 2})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	got := rec.last(t)
	wantQuery := map[string]string{
		"s":    "fight",
		"type": "movie",
		"page": "2",
	}
	for key, want := range wantQuery {
		if v := got.Query.Get(key); v != want {
			t.Errorf("%s = %q, want %q", key, v, want)
		}
	}
	if len(res.Search) != 1 || res.Search[0].IMDbID != "tt0137523" {
		t.Errorf("Search = %+v", res.Search)
	}
	if res.TotalResults != "1" {
		t.Errorf("TotalResults = %q", res.TotalResults)
	}
}

func TestOMDbAlwaysUsesTheQueryKey(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{"bearer mode", nil},
		{"query mode", []Option{WithAuth(AuthQuery)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := newTestClient(t, jsonHandler(200, omdbTitleBody), tc.opts...)
			if _, err := c.OMDb.ByID(context.Background(), "tt0137523", nil); err != nil {
				t.Fatalf("ByID: %v", err)
			}
			got := rec.last(t)
			if got.Query.Get("apikey") != testKey {
				t.Errorf("apikey = %q", got.Query.Get("apikey"))
			}
			if _, ok := got.Query["api_key"]; ok {
				t.Error("api_key must not be sent to the OMDb endpoint")
			}
			if got.Header.Get("Authorization") != "" {
				t.Error("OMDb calls must not send a bearer header")
			}
		})
	}
}

func TestOMDbYearAlias(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		wantY  string
	}{
		{"year becomes y", Params{"year": 2001}, "2001"},
		{"y wins over year", Params{"year": 2001, "y": 1999}, "1999"},
		{"only y", Params{"y": 1999}, "1999"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := newTestClient(t, jsonHandler(200, omdbTitleBody))
			if _, err := c.OMDb.ByID(context.Background(), "tt0137523", tc.params); err != nil {
				t.Fatalf("ByID: %v", err)
			}
			got := rec.last(t)
			if v := got.Query.Get("y"); v != tc.wantY {
				t.Errorf("y = %q, want %q", v, tc.wantY)
			}
			if _, ok := got.Query["year"]; ok {
				t.Error("year must not be sent")
			}
		})
	}
}

func TestOMDbResponseFalse(t *testing.T) {
	body := `{"Response":"False","Error":"Movie not found!"}`
	c, _ := newTestClient(t, jsonHandler(200, body))
	ctx := context.Background()
	tests := []struct {
		name string
		call func() error
	}{
		{"by id", func() error {
			return errOf(c.OMDb.ByID(ctx, "tt0000000", nil))
		}},
		{"by title", func() error {
			return errOf(c.OMDb.ByTitle(ctx, "No such title", nil))
		}},
		{"search", func() error {
			return errOf(c.OMDb.Search(ctx, "zzzz", nil))
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("error %v is not an *Error", err)
			}
			if apiErr.Status != 200 {
				t.Errorf("Status = %d, want 200", apiErr.Status)
			}
			if apiErr.StatusCode != 0 {
				t.Errorf("StatusCode = %d, want 0", apiErr.StatusCode)
			}
			if apiErr.StatusMessage != "Movie not found!" {
				t.Errorf("StatusMessage = %q", apiErr.StatusMessage)
			}
			if string(apiErr.Body) != body {
				t.Errorf("Body = %q", apiErr.Body)
			}
			if IsNotFound(err) || IsRateLimited(err) || IsUnauthorized(err) {
				t.Error("a Response False answer is not an HTTP error class")
			}
		})
	}
}

func TestOMDbResponseFalseWithoutText(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, `{"Response":"False"}`))
	_, err := c.OMDb.ByID(context.Background(), "tt0137523", nil)
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %v is not an *Error", err)
	}
	if apiErr.StatusMessage == "" {
		t.Error("StatusMessage should not be empty")
	}
}

func TestOMDbHTTPError(t *testing.T) {
	body := `{"Response":"False","Error":"Invalid API key!"}`
	c, _ := newTestClient(t, jsonHandler(401, body))
	_, err := c.OMDb.ByID(context.Background(), "tt0137523", nil)
	if !IsUnauthorized(err) {
		t.Fatalf("IsUnauthorized(%v) = false", err)
	}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %v is not an *Error", err)
	}
	if apiErr.StatusMessage != "Invalid API key!" {
		t.Errorf("StatusMessage = %q", apiErr.StatusMessage)
	}
	if apiErr.StatusCode != 0 {
		t.Errorf("StatusCode = %d, want 0", apiErr.StatusCode)
	}
}

func TestOMDbRetriesOnRateLimit(t *testing.T) {
	sleeps := &sleepRecorder{}
	c, rec := newTestClient(t, rateLimitedThenOK(1, "2"), WithSleep(sleeps.sleep))
	if _, err := c.OMDb.ByID(context.Background(), "tt0137523", nil); err != nil {
		t.Fatalf("ByID: %v", err)
	}
	if rec.count() != 2 {
		t.Errorf("requests = %d, want 2", rec.count())
	}
	if got := sleeps.recorded(); !equalDurations(got, []time.Duration{2 * time.Second}) {
		t.Errorf("waits = %v", got)
	}
}
