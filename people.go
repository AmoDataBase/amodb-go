package amodb

import "context"

// PeopleService groups the person endpoints. Use it as client.People.
//
// The sections that append_to_response accepts on Details are
// combined_credits, movie_credits, tv_credits, external_ids and images.
type PeopleService service

// Details returns /person/{id}.
func (s *PeopleService) Details(ctx context.Context, id int, params Params) (*PersonDetails, error) {
	return get[PersonDetails](ctx, s.client, buildPath("person", id), params)
}

// CombinedCredits returns movie and TV credits in one list.
func (s *PeopleService) CombinedCredits(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("person", id, "combined_credits"), params)
}

// MovieCredits returns the movie credits.
func (s *PeopleService) MovieCredits(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("person", id, "movie_credits"), params)
}

// TVCredits returns the TV credits.
func (s *PeopleService) TVCredits(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("person", id, "tv_credits"), params)
}

// Images returns profile images. Only the profiles list is filled.
func (s *PeopleService) Images(ctx context.Context, id int, params Params) (Object, error) {
	return getObject(ctx, s.client, buildPath("person", id, "images"), params)
}

// ExternalIDs returns the ids of the person on other services.
func (s *PeopleService) ExternalIDs(ctx context.Context, id int, params Params) (*PersonExternalIDs, error) {
	return get[PersonExternalIDs](ctx, s.client, buildPath("person", id, "external_ids"), params)
}

// Popular returns popular people.
func (s *PeopleService) Popular(ctx context.Context, params Params) (*Paged[Person], error) {
	return get[Paged[Person]](ctx, s.client, "/person/popular", params)
}
