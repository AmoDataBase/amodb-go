package amodb

import "context"

// DiscoverService groups the discover endpoints. Use it as client.Discover.
type DiscoverService service

// Movie discovers movies. Filters go in Params, for example
// Params{"with_genres": []int{28, 12}, "sort_by": "popularity.desc"}.
func (s *DiscoverService) Movie(ctx context.Context, params Params) (*Paged[Movie], error) {
	return get[Paged[Movie]](ctx, s.client, "/discover/movie", params)
}

// TV discovers TV shows. Filters go in Params.
func (s *DiscoverService) TV(ctx context.Context, params Params) (*Paged[TVShow], error) {
	return get[Paged[TVShow]](ctx, s.client, "/discover/tv", params)
}

// TrendingService groups the trending endpoint. Use it as client.Trending.
type TrendingService service

// Get returns trending items. mediaType is all, movie, tv or person. window is
// day or week. Each result is an Object because the shape depends on the
// media type.
func (s *TrendingService) Get(ctx context.Context, mediaType, window string, params Params) (*Paged[Object], error) {
	return get[Paged[Object]](ctx, s.client, buildPath("trending", mediaType, window), params)
}

// GenresService groups the genre lists. Use it as client.Genres.
type GenresService service

type genreList struct {
	Genres []Genre `json:"genres"`
}

// Movies returns the movie genres.
func (s *GenresService) Movies(ctx context.Context, params Params) ([]Genre, error) {
	return s.list(ctx, "/genre/movie/list", params)
}

// TV returns the TV genres.
func (s *GenresService) TV(ctx context.Context, params Params) ([]Genre, error) {
	return s.list(ctx, "/genre/tv/list", params)
}

func (s *GenresService) list(ctx context.Context, path string, params Params) ([]Genre, error) {
	out, err := get[genreList](ctx, s.client, path, params)
	if err != nil {
		return nil, err
	}
	return out.Genres, nil
}

// CollectionsService groups the collection endpoints. Use it as
// client.Collections.
type CollectionsService service

// Details returns a collection with its parts.
func (s *CollectionsService) Details(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("collection", id), params)
}

// CompaniesService groups the company endpoints. Use it as client.Companies.
type CompaniesService service

// Details returns a production company.
func (s *CompaniesService) Details(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("company", id), params)
}

// NetworksService groups the network endpoints. Use it as client.Networks.
type NetworksService service

// Details returns a TV network.
func (s *NetworksService) Details(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("network", id), params)
}

// KeywordsService groups the keyword endpoints. Use it as client.Keywords.
type KeywordsService service

// Details returns a keyword.
func (s *KeywordsService) Details(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("keyword", id), params)
}

// CertificationsService groups the certification lists. Use it as
// client.Certifications.
type CertificationsService service

// Movie returns the movie certifications per country.
func (s *CertificationsService) Movie(ctx context.Context) (Object, error) {
	return getObject(ctx, s.client, "/certification/movie/list", nil)
}

// TV returns the TV certifications per country.
func (s *CertificationsService) TV(ctx context.Context) (Object, error) {
	return getObject(ctx, s.client, "/certification/tv/list", nil)
}

// WatchProvidersService groups the watch provider lists. Use it as
// client.WatchProviders.
type WatchProvidersService service

// Regions returns the countries that have watch provider data.
func (s *WatchProvidersService) Regions(ctx context.Context, params Params) (Object, error) {
	return getObject(ctx, s.client, "/watch/providers/regions", params)
}

// Movie returns the movie watch providers.
func (s *WatchProvidersService) Movie(ctx context.Context, params Params) (Object, error) {
	return getObject(ctx, s.client, "/watch/providers/movie", params)
}

// TV returns the TV watch providers.
func (s *WatchProvidersService) TV(ctx context.Context, params Params) (Object, error) {
	return getObject(ctx, s.client, "/watch/providers/tv", params)
}
