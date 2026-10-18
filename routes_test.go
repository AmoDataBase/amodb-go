package amodb

import (
	"context"
	"net/http"
	"testing"
)

// TestRoutes calls every resource method once and checks the request path.
func TestRoutes(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, `{}`))
	tests := []struct {
		name  string
		path  string
		key   string
		value string
		call  func(ctx context.Context, c *Client) error
	}{
		{"Movies.Details", "/3/movie/550", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Details(ctx, 550, nil))
		}},
		{"Movies.Credits", "/3/movie/550/credits", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Credits(ctx, 550, nil))
		}},
		{"Movies.Images", "/3/movie/550/images", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Images(ctx, 550, nil))
		}},
		{"Movies.Videos", "/3/movie/550/videos", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Videos(ctx, 550, nil))
		}},
		{"Movies.Keywords", "/3/movie/550/keywords", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Keywords(ctx, 550, nil))
		}},
		{"Movies.ExternalIDs", "/3/movie/550/external_ids", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.ExternalIDs(ctx, 550, nil))
		}},
		{"Movies.Recommendations", "/3/movie/550/recommendations", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Recommendations(ctx, 550, nil))
		}},
		{"Movies.Similar", "/3/movie/550/similar", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Similar(ctx, 550, nil))
		}},
		{"Movies.Reviews", "/3/movie/550/reviews", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Reviews(ctx, 550, nil))
		}},
		{"Movies.ReleaseDates", "/3/movie/550/release_dates", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.ReleaseDates(ctx, 550, nil))
		}},
		{"Movies.WatchProviders", "/3/movie/550/watch/providers", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.WatchProviders(ctx, 550, nil))
		}},
		{"Movies.Translations", "/3/movie/550/translations", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Translations(ctx, 550, nil))
		}},
		{"Movies.AlternativeTitles", "/3/movie/550/alternative_titles", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.AlternativeTitles(ctx, 550, nil))
		}},
		{"Movies.Moods", "/3/movie/550/moods", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Moods(ctx, 550, nil))
		}},
		{"Movies.Soundtracks", "/3/movie/550/soundtracks", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Soundtracks(ctx, 550, nil))
		}},
		{"Movies.Ratings", "/3/movie/550/ratings", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Ratings(ctx, 550, nil))
		}},
		{"Movies.Wikidata", "/3/movie/550/wikidata", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Wikidata(ctx, 550, nil))
		}},
		{"Movies.Showtimes", "/3/movie/550/showtimes", "country", "US", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Showtimes(ctx, 550, "US", nil))
		}},
		{"Movies.Popular", "/3/movie/popular", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Popular(ctx, nil))
		}},
		{"Movies.TopRated", "/3/movie/top_rated", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.TopRated(ctx, nil))
		}},
		{"Movies.Upcoming", "/3/movie/upcoming", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.Upcoming(ctx, nil))
		}},
		{"Movies.NowPlaying", "/3/movie/now_playing", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Movies.NowPlaying(ctx, nil))
		}},
		{"TV.Details", "/3/tv/1399", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Details(ctx, 1399, nil))
		}},
		{"TV.AggregateCredits", "/3/tv/1399/aggregate_credits", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.AggregateCredits(ctx, 1399, nil))
		}},
		{"TV.Credits", "/3/tv/1399/credits", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Credits(ctx, 1399, nil))
		}},
		{"TV.Images", "/3/tv/1399/images", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Images(ctx, 1399, nil))
		}},
		{"TV.Videos", "/3/tv/1399/videos", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Videos(ctx, 1399, nil))
		}},
		{"TV.Keywords", "/3/tv/1399/keywords", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Keywords(ctx, 1399, nil))
		}},
		{"TV.ExternalIDs", "/3/tv/1399/external_ids", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.ExternalIDs(ctx, 1399, nil))
		}},
		{"TV.Recommendations", "/3/tv/1399/recommendations", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Recommendations(ctx, 1399, nil))
		}},
		{"TV.Similar", "/3/tv/1399/similar", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Similar(ctx, 1399, nil))
		}},
		{"TV.Reviews", "/3/tv/1399/reviews", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Reviews(ctx, 1399, nil))
		}},
		{"TV.ContentRatings", "/3/tv/1399/content_ratings", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.ContentRatings(ctx, 1399, nil))
		}},
		{"TV.WatchProviders", "/3/tv/1399/watch/providers", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.WatchProviders(ctx, 1399, nil))
		}},
		{"TV.Moods", "/3/tv/1399/moods", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Moods(ctx, 1399, nil))
		}},
		{"TV.Soundtracks", "/3/tv/1399/soundtracks", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Soundtracks(ctx, 1399, nil))
		}},
		{"TV.Ratings", "/3/tv/1399/ratings", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Ratings(ctx, 1399, nil))
		}},
		{"TV.Wikidata", "/3/tv/1399/wikidata", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Wikidata(ctx, 1399, nil))
		}},
		{"TV.Season", "/3/tv/1399/season/1", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Season(ctx, 1399, 1, nil))
		}},
		{"TV.Episode", "/3/tv/1399/season/1/episode/2", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Episode(ctx, 1399, 1, 2, nil))
		}},
		{"TV.Popular", "/3/tv/popular", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.Popular(ctx, nil))
		}},
		{"TV.TopRated", "/3/tv/top_rated", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.TopRated(ctx, nil))
		}},
		{"TV.AiringToday", "/3/tv/airing_today", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.AiringToday(ctx, nil))
		}},
		{"TV.OnTheAir", "/3/tv/on_the_air", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.TV.OnTheAir(ctx, nil))
		}},
		{"People.Details", "/3/person/287", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.People.Details(ctx, 287, nil))
		}},
		{"People.CombinedCredits", "/3/person/287/combined_credits", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.People.CombinedCredits(ctx, 287, nil))
		}},
		{"People.MovieCredits", "/3/person/287/movie_credits", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.People.MovieCredits(ctx, 287, nil))
		}},
		{"People.TVCredits", "/3/person/287/tv_credits", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.People.TVCredits(ctx, 287, nil))
		}},
		{"People.Images", "/3/person/287/images", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.People.Images(ctx, 287, nil))
		}},
		{"People.ExternalIDs", "/3/person/287/external_ids", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.People.ExternalIDs(ctx, 287, nil))
		}},
		{"People.Popular", "/3/person/popular", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.People.Popular(ctx, nil))
		}},
		{"Search.Movie", "/3/search/movie", "query", "inception", func(ctx context.Context, c *Client) error {
			return errOf(c.Search.Movie(ctx, "inception", nil))
		}},
		{"Search.TV", "/3/search/tv", "query", "inception", func(ctx context.Context, c *Client) error {
			return errOf(c.Search.TV(ctx, "inception", nil))
		}},
		{"Search.Person", "/3/search/person", "query", "inception", func(ctx context.Context, c *Client) error {
			return errOf(c.Search.Person(ctx, "inception", nil))
		}},
		{"Search.Multi", "/3/search/multi", "query", "inception", func(ctx context.Context, c *Client) error {
			return errOf(c.Search.Multi(ctx, "inception", nil))
		}},
		{"Search.Collection", "/3/search/collection", "query", "inception", func(ctx context.Context, c *Client) error {
			return errOf(c.Search.Collection(ctx, "inception", nil))
		}},
		{"Search.Company", "/3/search/company", "query", "inception", func(ctx context.Context, c *Client) error {
			return errOf(c.Search.Company(ctx, "inception", nil))
		}},
		{"Search.Keyword", "/3/search/keyword", "query", "inception", func(ctx context.Context, c *Client) error {
			return errOf(c.Search.Keyword(ctx, "inception", nil))
		}},
		{"Discover.Movie", "/3/discover/movie", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Discover.Movie(ctx, nil))
		}},
		{"Discover.TV", "/3/discover/tv", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Discover.TV(ctx, nil))
		}},
		{"Trending.Get", "/3/trending/movie/week", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Trending.Get(ctx, "movie", "week", nil))
		}},
		{"Genres.Movies", "/3/genre/movie/list", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Genres.Movies(ctx, nil))
		}},
		{"Genres.TV", "/3/genre/tv/list", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Genres.TV(ctx, nil))
		}},
		{"Collections.Details", "/3/collection/10", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Collections.Details(ctx, 10, nil))
		}},
		{"Companies.Details", "/3/company/508", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Companies.Details(ctx, 508, nil))
		}},
		{"Networks.Details", "/3/network/213", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Networks.Details(ctx, 213, nil))
		}},
		{"Keywords.Details", "/3/keyword/818", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Keywords.Details(ctx, 818, nil))
		}},
		{"Certifications.Movie", "/3/certification/movie/list", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Certifications.Movie(ctx))
		}},
		{"Certifications.TV", "/3/certification/tv/list", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Certifications.TV(ctx))
		}},
		{"WatchProviders.Regions", "/3/watch/providers/regions", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.WatchProviders.Regions(ctx, nil))
		}},
		{"WatchProviders.Movie", "/3/watch/providers/movie", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.WatchProviders.Movie(ctx, nil))
		}},
		{"WatchProviders.TV", "/3/watch/providers/tv", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.WatchProviders.TV(ctx, nil))
		}},
		{"Find", "/3/find/tt0137523", "external_source", "imdb_id", func(ctx context.Context, c *Client) error {
			return errOf(c.Find(ctx, "tt0137523", "imdb_id"))
		}},
		{"Configuration", "/3/configuration", "", "", func(ctx context.Context, c *Client) error {
			return errOf(c.Configuration(ctx))
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(context.Background(), c); err != nil {
				t.Fatalf("call failed: %v", err)
			}
			got := rec.last(t)
			if got.Method != http.MethodGet {
				t.Errorf("method = %q, want GET", got.Method)
			}
			if got.Path != tc.path {
				t.Errorf("path = %q, want %q", got.Path, tc.path)
			}
			if tc.key != "" && got.Query.Get(tc.key) != tc.value {
				t.Errorf("%s = %q, want %q", tc.key, got.Query.Get(tc.key), tc.value)
			}
		})
	}
	if got := rec.count(); got != len(tests) {
		t.Errorf("requests = %d, want one per case (%d)", got, len(tests))
	}
}
