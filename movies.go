package amodb

import "context"

// MoviesService groups the movie endpoints. Use it as client.Movies.
//
// Every method takes an optional Params map. Pass "page" for list endpoints,
// "language" to override the client language and "append_to_response" on
// Details to merge more sections into one answer.
type MoviesService service

// Details returns /movie/{id}.
func (s *MoviesService) Details(ctx context.Context, id int, params Params) (*MovieDetails, error) {
	return get[MovieDetails](ctx, s.client, buildPath("movie", id), params)
}

// Credits returns the cast and crew.
func (s *MoviesService) Credits(ctx context.Context, id int, params Params) (*Credits, error) {
	return get[Credits](ctx, s.client, buildPath("movie", id, "credits"), params)
}

// Images returns posters, backdrops and logos.
func (s *MoviesService) Images(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "images"), params)
}

// Videos returns trailers and clips.
func (s *MoviesService) Videos(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "videos"), params)
}

// Keywords returns the keywords of the movie.
func (s *MoviesService) Keywords(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "keywords"), params)
}

// ExternalIDs returns the ids of the movie on other services.
func (s *MoviesService) ExternalIDs(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "external_ids"), params)
}

// Recommendations returns recommended movies.
func (s *MoviesService) Recommendations(ctx context.Context, id int, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, buildPath("movie", id, "recommendations"), params)
}

// Similar returns similar movies.
func (s *MoviesService) Similar(ctx context.Context, id int, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, buildPath("movie", id, "similar"), params)
}

// Reviews returns user reviews.
func (s *MoviesService) Reviews(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "reviews"), params)
}

// ReleaseDates returns release dates and certifications per country.
func (s *MoviesService) ReleaseDates(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "release_dates"), params)
}

// WatchProviders returns where the movie streams, per country.
func (s *MoviesService) WatchProviders(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "watch", "providers"), params)
}

// Translations returns the translated titles and overviews.
func (s *MoviesService) Translations(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "translations"), params)
}

// AlternativeTitles returns the titles used in other countries.
func (s *MoviesService) AlternativeTitles(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "alternative_titles"), params)
}

// Popular returns popular movies.
func (s *MoviesService) Popular(ctx context.Context, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, "/movie/popular", params)
}

// TopRated returns the top rated movies.
func (s *MoviesService) TopRated(ctx context.Context, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, "/movie/top_rated", params)
}

// Upcoming returns movies that open soon.
func (s *MoviesService) Upcoming(ctx context.Context, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, "/movie/upcoming", params)
}

// NowPlaying returns movies that are in cinemas now.
func (s *MoviesService) NowPlaying(ctx context.Context, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, "/movie/now_playing", params)
}

// Moods returns the mood meters of the movie.
func (s *MoviesService) Moods(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "moods"), params)
}

// Soundtracks returns the soundtrack of the movie.
func (s *MoviesService) Soundtracks(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "soundtracks"), params)
}

// Ratings returns the ratings from several sources. Each source holds a value
// and a vote count, or null.
func (s *MoviesService) Ratings(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "ratings"), params)
}

// Wikidata returns selected facts about the movie.
func (s *MoviesService) Wikidata(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "wikidata"), params)
}

// Showtimes returns cinema showtimes for a two letter country code, when known.
func (s *MoviesService) Showtimes(ctx context.Context, id int, country string, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("movie", id, "showtimes"), withParam(params, "country", country))
}
