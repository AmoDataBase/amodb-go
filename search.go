package amodb

import "context"

// SearchService groups the search endpoints. Use it as client.Search.
//
// Each method takes the search text first. Pass "page", "year",
// "first_air_date_year", "language" or any other documented parameter in
// Params. The include_adult parameter is ignored, because adult titles and
// people are never served.
type SearchService service

// Movie searches movies.
func (s *SearchService) Movie(ctx context.Context, query string, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, "/search/movie", withParam(params, "query", query))
}

// TV searches TV shows.
func (s *SearchService) TV(ctx context.Context, query string, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, "/search/tv", withParam(params, "query", query))
}

// Person searches people.
func (s *SearchService) Person(ctx context.Context, query string, params Params) (*Paged[Person], error) {
	return get[Paged[Person]](ctx, s.client, "/search/person", withParam(params, "query", query))
}

// Multi searches movies, shows and people at once. Each result has a
// media_type field.
func (s *SearchService) Multi(ctx context.Context, query string, params Params) (*Paged[Object], error) {
	return get[Paged[Object]](ctx, s.client, "/search/multi", withParam(params, "query", query))
}

// Collection searches collections.
func (s *SearchService) Collection(ctx context.Context, query string, params Params) (*Paged[Object], error) {
	return get[Paged[Object]](ctx, s.client, "/search/collection", withParam(params, "query", query))
}

// Company searches production companies.
func (s *SearchService) Company(ctx context.Context, query string, params Params) (*Paged[Object], error) {
	return get[Paged[Object]](ctx, s.client, "/search/company", withParam(params, "query", query))
}

// Keyword searches keywords.
func (s *SearchService) Keyword(ctx context.Context, query string, params Params) (*Paged[Object], error) {
	return get[Paged[Object]](ctx, s.client, "/search/keyword", withParam(params, "query", query))
}
