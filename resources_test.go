package amodb

import (
	"context"
	"testing"
)

const movieDetailsBody = `{"id":550,"title":"Fight Club","original_title":"Fight Club","overview":"An office worker starts a club.","tagline":"Mischief.","status":"Released","release_date":"1999-10-15","runtime":139,"budget":63000000,"revenue":100853753,"imdb_id":"tt0137523","poster_path":"/p.jpg","backdrop_path":null,"vote_average":8.4,"vote_count":28000,"popularity":61.4,"adult":false,"genres":[{"id":18,"name":"Drama"}],"production_companies":[{"id":508,"name":"Example Pictures","logo_path":null,"origin_country":"US"}],"production_countries":[{"iso_3166_1":"US","name":"United States of America"}],"spoken_languages":[{"iso_639_1":"en","english_name":"English","name":"English"}],"belongs_to_collection":null,"credits":{"id":550,"cast":[{"id":819,"name":"Jane Doe","character":"Narrator","order":0,"gender":2}],"crew":[{"id":7467,"name":"John Roe","job":"Director","department":"Directing"}]},"images":{"posters":[{"file_path":"/p.jpg"}]},"ratings":{"imdb":{"value":8.8,"votes":2300000}}}`

func TestMovieDetailsDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, movieDetailsBody))
	params := Params{"append_to_response": []string{"credits", "images", "ratings"}}
	movie, err := c.Movies.Details(context.Background(), 550, params)
	if err != nil {
		t.Fatalf("Details: %v", err)
	}
	if got := rec.last(t).Query.Get("append_to_response"); got != "credits,images,ratings" {
		t.Errorf("append_to_response = %q", got)
	}
	if movie.ID != 550 || movie.Title != "Fight Club" || movie.Runtime != 139 {
		t.Errorf("movie = %+v", movie)
	}
	if movie.Budget != 63000000 || movie.Revenue != 100853753 {
		t.Errorf("budget = %d, revenue = %d", movie.Budget, movie.Revenue)
	}
	if movie.IMDbID != "tt0137523" || movie.ReleaseDate != "1999-10-15" {
		t.Errorf("imdb id = %q, release date = %q", movie.IMDbID, movie.ReleaseDate)
	}
	if movie.PosterPath != "/p.jpg" || movie.BackdropPath != "" {
		t.Errorf("poster = %q, backdrop = %q", movie.PosterPath, movie.BackdropPath)
	}
	if movie.VoteAverage != 8.4 || movie.VoteCount != 28000 {
		t.Errorf("vote average = %v, count = %d", movie.VoteAverage, movie.VoteCount)
	}
	if len(movie.Genres) != 1 || movie.Genres[0].Name != "Drama" {
		t.Errorf("genres = %+v", movie.Genres)
	}
	if len(movie.ProductionCompanies) != 1 || movie.ProductionCompanies[0].Name != "Example Pictures" {
		t.Errorf("companies = %+v", movie.ProductionCompanies)
	}
	if len(movie.ProductionCountries) != 1 || movie.ProductionCountries[0].ISO31661 != "US" {
		t.Errorf("countries = %+v", movie.ProductionCountries)
	}
	if len(movie.SpokenLanguages) != 1 || movie.SpokenLanguages[0].ISO6391 != "en" {
		t.Errorf("languages = %+v", movie.SpokenLanguages)
	}
	if movie.BelongsToCollection != nil {
		t.Errorf("collection = %+v, want nil", movie.BelongsToCollection)
	}
	if movie.Credits == nil || len(movie.Credits.Cast) != 1 || len(movie.Credits.Crew) != 1 {
		t.Fatalf("credits = %+v", movie.Credits)
	}
	if movie.Credits.Cast[0].Character != "Narrator" || movie.Credits.Crew[0].Job != "Director" {
		t.Errorf("credits = %+v", movie.Credits)
	}
	if movie.Images == nil || movie.Images["posters"] == nil {
		t.Errorf("images = %v", movie.Images)
	}
	if movie.Ratings == nil || movie.Ratings.IMDb == nil {
		t.Fatalf("ratings = %+v", movie.Ratings)
	}
	imdb := movie.Ratings.IMDb
	if imdb.Value != 8.8 || imdb.Votes == nil || *imdb.Votes != 2300000 {
		t.Errorf("imdb = %+v", imdb)
	}
	if movie.Ratings.Kinopoisk != nil || movie.Ratings.TMDB != nil {
		t.Errorf("ratings = %+v, the missing providers must be nil", movie.Ratings)
	}
}

func TestPagedSearchDecoding(t *testing.T) {
	body := `{"page":1,"results":[{"id":27205,"title":"Inception","release_date":"2010-07-15","vote_average":8.4,"genre_ids":[28,878]}],"total_pages":3,"total_results":55}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	res, err := c.Search.Movie(context.Background(), "inception", Params{"query": "ignored", "year": 2010})
	if err != nil {
		t.Fatalf("Search.Movie: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/search/movie" {
		t.Errorf("path = %q", got.Path)
	}
	if got.Query.Get("query") != "inception" {
		t.Errorf("query = %q, the argument must win over Params", got.Query.Get("query"))
	}
	if got.Query.Get("year") != "2010" {
		t.Errorf("year = %q", got.Query.Get("year"))
	}
	if res.Page != 1 || res.TotalPages != 3 || res.TotalResults != 55 {
		t.Errorf("paging = %d / %d / %d", res.Page, res.TotalPages, res.TotalResults)
	}
	if len(res.Results) != 1 {
		t.Fatalf("results = %+v", res.Results)
	}
	first := res.Results[0]
	if first.ID != 27205 || first.Title != "Inception" || len(first.GenreIDs) != 2 {
		t.Errorf("first = %+v", first)
	}
}

func TestTVDetailsDecoding(t *testing.T) {
	body := `{"id":1399,"name":"Example Show","original_name":"Example Show","first_air_date":"2011-04-17","in_production":false,"number_of_seasons":8,"number_of_episodes":73,"episode_run_time":[60],"origin_country":["US"],"genres":[{"id":18,"name":"Drama"}],"networks":[{"id":49,"name":"Example Network","logo_path":"/n.png","origin_country":"US"}],"created_by":[{"id":9813,"name":"Jane Doe"}],"seasons":[{"id":3624,"name":"Season 1","season_number":1,"episode_count":10,"air_date":"2011-04-17"}],"last_episode_to_air":{"id":1,"name":"Finale","season_number":8,"episode_number":6},"next_episode_to_air":null}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	show, err := c.TV.Details(context.Background(), 1399, nil)
	if err != nil {
		t.Fatalf("TV.Details: %v", err)
	}
	if rec.last(t).Path != "/3/tv/1399" {
		t.Errorf("path = %q", rec.last(t).Path)
	}
	if show.Name != "Example Show" || show.NumberOfSeasons != 8 || show.NumberOfEpisodes != 73 {
		t.Errorf("show = %+v", show)
	}
	if len(show.EpisodeRunTime) != 1 || show.EpisodeRunTime[0] != 60 {
		t.Errorf("episode run time = %v", show.EpisodeRunTime)
	}
	if len(show.Networks) != 1 || show.Networks[0].Name != "Example Network" {
		t.Errorf("networks = %+v", show.Networks)
	}
	if len(show.CreatedBy) != 1 || show.CreatedBy[0].Name != "Jane Doe" {
		t.Errorf("created by = %+v", show.CreatedBy)
	}
	if len(show.Seasons) != 1 || show.Seasons[0].EpisodeCount != 10 {
		t.Errorf("seasons = %+v", show.Seasons)
	}
	if show.LastEpisodeToAir == nil || show.LastEpisodeToAir.Name != "Finale" {
		t.Errorf("last episode = %+v", show.LastEpisodeToAir)
	}
	if show.NextEpisodeToAir != nil {
		t.Errorf("next episode = %+v, want nil", show.NextEpisodeToAir)
	}
}

func TestSeasonAndEpisodeDecoding(t *testing.T) {
	seasonBody := `{"id":3624,"name":"Season 1","season_number":1,"air_date":"2011-04-17","episodes":[{"id":63056,"name":"Pilot","episode_number":1,"season_number":1,"runtime":62,"vote_average":8.2}]}`
	c, rec := newTestClient(t, jsonHandler(200, seasonBody))
	season, err := c.TV.Season(context.Background(), 1399, 1, nil)
	if err != nil {
		t.Fatalf("TV.Season: %v", err)
	}
	if rec.last(t).Path != "/3/tv/1399/season/1" {
		t.Errorf("path = %q", rec.last(t).Path)
	}
	if season.SeasonNumber != 1 || len(season.Episodes) != 1 {
		t.Fatalf("season = %+v", season)
	}
	if ep := season.Episodes[0]; ep.Name != "Pilot" || ep.Runtime != 62 || ep.EpisodeNumber != 1 {
		t.Errorf("episode = %+v", ep)
	}

	episodeBody := `{"id":63056,"name":"Pilot","episode_number":2,"season_number":1,"show_id":1399}`
	c2, rec2 := newTestClient(t, jsonHandler(200, episodeBody))
	episode, err := c2.TV.Episode(context.Background(), 1399, 1, 2, nil)
	if err != nil {
		t.Fatalf("TV.Episode: %v", err)
	}
	if rec2.last(t).Path != "/3/tv/1399/season/1/episode/2" {
		t.Errorf("path = %q", rec2.last(t).Path)
	}
	if episode.ShowID != 1399 || episode.EpisodeNumber != 2 {
		t.Errorf("episode = %+v", episode)
	}
}

func TestPersonDetailsDecoding(t *testing.T) {
	body := `{"id":287,"name":"Jane Doe","also_known_as":["J. Doe"],"biography":"Short biography.","birthday":"1963-12-18","deathday":null,"gender":2,"imdb_id":"nm0000093","known_for_department":"Acting","place_of_birth":"Example City","popularity":12.5}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	person, err := c.People.Details(context.Background(), 287, nil)
	if err != nil {
		t.Fatalf("People.Details: %v", err)
	}
	if rec.last(t).Path != "/3/person/287" {
		t.Errorf("path = %q", rec.last(t).Path)
	}
	if person.Name != "Jane Doe" || person.Birthday != "1963-12-18" || person.Deathday != "" {
		t.Errorf("person = %+v", person)
	}
	if person.IMDbID != "nm0000093" || person.Popularity != 12.5 || len(person.AlsoKnownAs) != 1 {
		t.Errorf("person = %+v", person)
	}
}

func TestCreditsDecoding(t *testing.T) {
	body := `{"id":550,"cast":[{"id":1,"name":"A","character":"X","order":0},{"id":2,"name":"B","character":"Y","order":1}],"crew":[]}`
	c, _ := newTestClient(t, jsonHandler(200, body))
	credits, err := c.Movies.Credits(context.Background(), 550, nil)
	if err != nil {
		t.Fatalf("Credits: %v", err)
	}
	if credits.ID != 550 || len(credits.Cast) != 2 || len(credits.Crew) != 0 {
		t.Errorf("credits = %+v", credits)
	}
	if credits.Cast[1].Order != 1 || credits.Cast[1].Character != "Y" {
		t.Errorf("cast = %+v", credits.Cast)
	}
}

func TestFindDecoding(t *testing.T) {
	body := `{"movie_results":[{"id":550,"title":"Fight Club"}],"person_results":[],"tv_results":[],"tv_episode_results":[],"tv_season_results":[]}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	res, err := c.Find(context.Background(), "tt0137523", "imdb_id")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/find/tt0137523" || got.Query.Get("external_source") != "imdb_id" {
		t.Errorf("request = %s ? %v", got.Path, got.Query)
	}
	if len(res.MovieResults) != 1 || res.MovieResults[0].ID != 550 {
		t.Errorf("movie results = %+v", res.MovieResults)
	}
	if len(res.PersonResults) != 0 || len(res.TVResults) != 0 {
		t.Errorf("results = %+v", res)
	}
}

func TestConfigurationDecoding(t *testing.T) {
	body := `{"images":{"base_url":"http://img.amodb.dev/t/p/","secure_base_url":"https://img.amodb.dev/t/p/","backdrop_sizes":["w300","w780","w1280","original"],"logo_sizes":["w45"],"poster_sizes":["w92","w500","original"],"profile_sizes":["w185"],"still_sizes":["w92"]},"change_keys":["adult"]}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	cfg, err := c.Configuration(context.Background())
	if err != nil {
		t.Fatalf("Configuration: %v", err)
	}
	if rec.last(t).Path != "/3/configuration" {
		t.Errorf("path = %q", rec.last(t).Path)
	}
	if cfg.Images.SecureBaseURL != "https://img.amodb.dev/t/p/" || cfg.Images.BaseURL != "http://img.amodb.dev/t/p/" {
		t.Errorf("images = %+v", cfg.Images)
	}
	if len(cfg.Images.BackdropSizes) != 4 || len(cfg.Images.PosterSizes) != 3 {
		t.Errorf("sizes = %+v", cfg.Images)
	}
	if len(cfg.ChangeKeys) != 1 {
		t.Errorf("change keys = %v", cfg.ChangeKeys)
	}
}

func TestGenresReturnsTheList(t *testing.T) {
	body := `{"genres":[{"id":28,"name":"Action"},{"id":12,"name":"Adventure"}]}`
	c, rec := newTestClient(t, jsonHandler(200, body), WithLanguage("ru-RU"))
	genres, err := c.Genres.Movies(context.Background(), nil)
	if err != nil {
		t.Fatalf("Genres.Movies: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/genre/movie/list" || got.Query.Get("language") != "ru-RU" {
		t.Errorf("request = %s ? %v", got.Path, got.Query)
	}
	if len(genres) != 2 || genres[0].ID != 28 || genres[1].Name != "Adventure" {
		t.Errorf("genres = %+v", genres)
	}
}

func TestTrendingReturnsObjects(t *testing.T) {
	body := `{"page":1,"results":[{"id":550,"media_type":"movie","title":"Fight Club"},{"id":1399,"media_type":"tv","name":"Example Show"}],"total_pages":1,"total_results":2}`
	c, rec := newTestClient(t, jsonHandler(200, body))
	res, err := c.Trending.Get(context.Background(), "all", "week", Params{"page": 1})
	if err != nil {
		t.Fatalf("Trending.Get: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/trending/all/week" || got.Query.Get("page") != "1" {
		t.Errorf("request = %s ? %v", got.Path, got.Query)
	}
	if len(res.Results) != 2 || res.Results[1]["media_type"] != "tv" {
		t.Errorf("results = %+v", res.Results)
	}
}

func TestParamsAreNotMutated(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, `{"results":[]}`))
	params := Params{"year": 2010}
	if _, err := c.Search.Movie(context.Background(), "inception", params); err != nil {
		t.Fatalf("Search.Movie: %v", err)
	}
	if len(params) != 1 {
		t.Errorf("params = %v, the caller map must stay unchanged", params)
	}
	omdbParamsMap := Params{"year": 2010}
	if _, err := c.OMDb.ByID(context.Background(), "tt0137523", omdbParamsMap); err != nil {
		t.Fatalf("OMDb.ByID: %v", err)
	}
	if _, ok := omdbParamsMap["year"]; !ok || len(omdbParamsMap) != 1 {
		t.Errorf("params = %v, the caller map must stay unchanged", omdbParamsMap)
	}
}
