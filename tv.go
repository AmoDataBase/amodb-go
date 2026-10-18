package amodb

import "context"

// TVService groups the TV endpoints. Use it as client.TV.
//
// Every method takes an optional Params map. Pass "page" for list endpoints,
// "language" to override the client language and "append_to_response" on
// Details to merge more sections into one answer.
type TVService service

// Details returns /tv/{id}.
func (s *TVService) Details(ctx context.Context, id int, params Params) (*TVDetails, error) {
	return get[TVDetails](ctx, s.client, buildPath("tv", id), params)
}

// AggregateCredits returns the cast and crew across all seasons.
func (s *TVService) AggregateCredits(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "aggregate_credits"), params)
}

// Credits returns the cast and crew of the latest season.
func (s *TVService) Credits(ctx context.Context, id int, params Params) (*Credits, error) {
	return get[Credits](ctx, s.client, buildPath("tv", id, "credits"), params)
}

// Images returns posters, backdrops and logos.
func (s *TVService) Images(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "images"), params)
}

// Videos returns trailers and clips.
func (s *TVService) Videos(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "videos"), params)
}

// Keywords returns the keywords of the show.
func (s *TVService) Keywords(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "keywords"), params)
}

// ExternalIDs returns the ids of the show on other services.
func (s *TVService) ExternalIDs(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "external_ids"), params)
}

// Recommendations returns recommended shows.
func (s *TVService) Recommendations(ctx context.Context, id int, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, buildPath("tv", id, "recommendations"), params)
}

// Similar returns similar shows.
func (s *TVService) Similar(ctx context.Context, id int, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, buildPath("tv", id, "similar"), params)
}

// Reviews returns user reviews.
func (s *TVService) Reviews(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "reviews"), params)
}

// ContentRatings returns the age ratings per country.
func (s *TVService) ContentRatings(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "content_ratings"), params)
}

// WatchProviders returns where the show streams, per country.
func (s *TVService) WatchProviders(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "watch", "providers"), params)
}

// Season returns one season with its episodes.
func (s *TVService) Season(ctx context.Context, tvID, seasonNumber int, params Params) (*Season, error) {
	return get[Season](ctx, s.client, buildPath("tv", tvID, "season", seasonNumber), params)
}

// Episode returns one episode.
func (s *TVService) Episode(ctx context.Context, tvID, seasonNumber, episodeNumber int, params Params) (*Episode, error) {
	return get[Episode](ctx, s.client, buildPath("tv", tvID, "season", seasonNumber, "episode", episodeNumber), params)
}

// Popular returns popular shows.
func (s *TVService) Popular(ctx context.Context, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, "/tv/popular", params)
}

// TopRated returns the top rated shows.
func (s *TVService) TopRated(ctx context.Context, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, "/tv/top_rated", params)
}

// AiringToday returns shows with an episode on the air today.
func (s *TVService) AiringToday(ctx context.Context, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, "/tv/airing_today", params)
}

// OnTheAir returns shows with an episode in the next 7 days.
func (s *TVService) OnTheAir(ctx context.Context, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, "/tv/on_the_air", params)
}

// Moods returns the mood meters of the show.
func (s *TVService) Moods(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "moods"), params)
}

// Soundtracks returns the soundtrack of the show.
func (s *TVService) Soundtracks(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "soundtracks"), params)
}

// Ratings returns the ratings from several sources. Each source holds a value
// and a vote count, or null.
func (s *TVService) Ratings(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "ratings"), params)
}

// Wikidata returns selected facts about the show.
func (s *TVService) Wikidata(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("tv", id, "wikidata"), params)
}
