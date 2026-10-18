package amodb

// Object is a decoded JSON object. Endpoints without a typed result return it.
type Object = map[string]any

// Paged is the list envelope used by search, discover, trending and lists.
type Paged[T any] struct {
	Page         int `json:"page"`
	Results      []T `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// Genre is a genre id and name.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Company is a production company.
type Company struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path,omitempty"`
	OriginCountry string `json:"origin_country,omitempty"`
}

// Network is a TV network.
type Network struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path,omitempty"`
	OriginCountry string `json:"origin_country,omitempty"`
}

// Country is a production country.
type Country struct {
	ISO31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

// SpokenLanguage is a language spoken in a title.
type SpokenLanguage struct {
	ISO6391     string `json:"iso_639_1"`
	EnglishName string `json:"english_name,omitempty"`
	Name        string `json:"name"`
}

// CollectionRef is the collection a movie belongs to.
type CollectionRef struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path,omitempty"`
	BackdropPath string `json:"backdrop_path,omitempty"`
}

// Movie is the short movie form used in lists and search results.
type Movie struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title,omitempty"`
	OriginalLanguage string  `json:"original_language,omitempty"`
	Overview         string  `json:"overview,omitempty"`
	ReleaseDate      string  `json:"release_date,omitempty"`
	PosterPath       string  `json:"poster_path,omitempty"`
	BackdropPath     string  `json:"backdrop_path,omitempty"`
	GenreIDs         []int   `json:"genre_ids,omitempty"`
	Popularity       float64 `json:"popularity,omitempty"`
	VoteAverage      float64 `json:"vote_average,omitempty"`
	VoteCount        int     `json:"vote_count,omitempty"`
	Adult            bool    `json:"adult,omitempty"`
	Video            bool    `json:"video,omitempty"`
}

// MovieDetails is the answer of Movies.Details. The sections that
// append_to_response adds are decoded into the fields at the end. A section
// without a field here can be read with Client.Do into your own struct.
type MovieDetails struct {
	ID                  int              `json:"id"`
	Title               string           `json:"title"`
	OriginalTitle       string           `json:"original_title,omitempty"`
	OriginalLanguage    string           `json:"original_language,omitempty"`
	Overview            string           `json:"overview,omitempty"`
	Tagline             string           `json:"tagline,omitempty"`
	Status              string           `json:"status,omitempty"`
	ReleaseDate         string           `json:"release_date,omitempty"`
	Runtime             int              `json:"runtime,omitempty"`
	Budget              int64            `json:"budget,omitempty"`
	Revenue             int64            `json:"revenue,omitempty"`
	Homepage            string           `json:"homepage,omitempty"`
	IMDbID              string           `json:"imdb_id,omitempty"`
	PosterPath          string           `json:"poster_path,omitempty"`
	BackdropPath        string           `json:"backdrop_path,omitempty"`
	Popularity          float64          `json:"popularity,omitempty"`
	VoteAverage         float64          `json:"vote_average,omitempty"`
	VoteCount           int              `json:"vote_count,omitempty"`
	Adult               bool             `json:"adult,omitempty"`
	Video               bool             `json:"video,omitempty"`
	Genres              []Genre          `json:"genres,omitempty"`
	ProductionCompanies []Company        `json:"production_companies,omitempty"`
	ProductionCountries []Country        `json:"production_countries,omitempty"`
	SpokenLanguages     []SpokenLanguage `json:"spoken_languages,omitempty"`
	BelongsToCollection *CollectionRef   `json:"belongs_to_collection,omitempty"`
	Credits             *Credits         `json:"credits,omitempty"`
	Recommendations     *Paged[Movie]    `json:"recommendations,omitempty"`
	Similar             *Paged[Movie]    `json:"similar,omitempty"`
	Images              Object           `json:"images,omitempty"`
	Videos              Object           `json:"videos,omitempty"`
	Keywords            Object           `json:"keywords,omitempty"`
	ExternalIDs         Object           `json:"external_ids,omitempty"`
	ReleaseDates        Object           `json:"release_dates,omitempty"`
	Reviews             Object           `json:"reviews,omitempty"`
	Translations        Object           `json:"translations,omitempty"`
	AlternativeTitles   Object           `json:"alternative_titles,omitempty"`
	WatchProviders      Object           `json:"watch/providers,omitempty"`
	Moods               Object           `json:"moods,omitempty"`
	Soundtracks         Object           `json:"soundtracks,omitempty"`
	Ratings             Object           `json:"ratings,omitempty"`
	Wikidata            Object           `json:"wikidata,omitempty"`
	Showtimes           Object           `json:"showtimes,omitempty"`
}

// TVShow is the short TV form used in lists and search results.
type TVShow struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name,omitempty"`
	OriginalLanguage string   `json:"original_language,omitempty"`
	Overview         string   `json:"overview,omitempty"`
	FirstAirDate     string   `json:"first_air_date,omitempty"`
	PosterPath       string   `json:"poster_path,omitempty"`
	BackdropPath     string   `json:"backdrop_path,omitempty"`
	GenreIDs         []int    `json:"genre_ids,omitempty"`
	OriginCountry    []string `json:"origin_country,omitempty"`
	Popularity       float64  `json:"popularity,omitempty"`
	VoteAverage      float64  `json:"vote_average,omitempty"`
	VoteCount        int      `json:"vote_count,omitempty"`
	Adult            bool     `json:"adult,omitempty"`
}

// CreatedBy is a creator of a TV show.
type CreatedBy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CreditID    string `json:"credit_id,omitempty"`
	Gender      int    `json:"gender,omitempty"`
	ProfilePath string `json:"profile_path,omitempty"`
}

// SeasonSummary is a season entry inside TVDetails.
type SeasonSummary struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview,omitempty"`
	AirDate      string  `json:"air_date,omitempty"`
	SeasonNumber int     `json:"season_number"`
	EpisodeCount int     `json:"episode_count,omitempty"`
	PosterPath   string  `json:"poster_path,omitempty"`
	VoteAverage  float64 `json:"vote_average,omitempty"`
}

// TVDetails is the answer of TV.Details. The sections that append_to_response
// adds are decoded into the fields at the end.
type TVDetails struct {
	ID                  int              `json:"id"`
	Name                string           `json:"name"`
	OriginalName        string           `json:"original_name,omitempty"`
	OriginalLanguage    string           `json:"original_language,omitempty"`
	Overview            string           `json:"overview,omitempty"`
	Tagline             string           `json:"tagline,omitempty"`
	Status              string           `json:"status,omitempty"`
	Type                string           `json:"type,omitempty"`
	Homepage            string           `json:"homepage,omitempty"`
	FirstAirDate        string           `json:"first_air_date,omitempty"`
	LastAirDate         string           `json:"last_air_date,omitempty"`
	InProduction        bool             `json:"in_production,omitempty"`
	NumberOfSeasons     int              `json:"number_of_seasons,omitempty"`
	NumberOfEpisodes    int              `json:"number_of_episodes,omitempty"`
	EpisodeRunTime      []int            `json:"episode_run_time,omitempty"`
	Languages           []string         `json:"languages,omitempty"`
	OriginCountry       []string         `json:"origin_country,omitempty"`
	PosterPath          string           `json:"poster_path,omitempty"`
	BackdropPath        string           `json:"backdrop_path,omitempty"`
	Popularity          float64          `json:"popularity,omitempty"`
	VoteAverage         float64          `json:"vote_average,omitempty"`
	VoteCount           int              `json:"vote_count,omitempty"`
	Adult               bool             `json:"adult,omitempty"`
	Genres              []Genre          `json:"genres,omitempty"`
	Networks            []Network        `json:"networks,omitempty"`
	ProductionCompanies []Company        `json:"production_companies,omitempty"`
	ProductionCountries []Country        `json:"production_countries,omitempty"`
	SpokenLanguages     []SpokenLanguage `json:"spoken_languages,omitempty"`
	CreatedBy           []CreatedBy      `json:"created_by,omitempty"`
	Seasons             []SeasonSummary  `json:"seasons,omitempty"`
	LastEpisodeToAir    *Episode         `json:"last_episode_to_air,omitempty"`
	NextEpisodeToAir    *Episode         `json:"next_episode_to_air,omitempty"`
	Credits             *Credits         `json:"credits,omitempty"`
	Recommendations     *Paged[TVShow]   `json:"recommendations,omitempty"`
	Similar             *Paged[TVShow]   `json:"similar,omitempty"`
	AggregateCredits    Object           `json:"aggregate_credits,omitempty"`
	Images              Object           `json:"images,omitempty"`
	Videos              Object           `json:"videos,omitempty"`
	Keywords            Object           `json:"keywords,omitempty"`
	ExternalIDs         Object           `json:"external_ids,omitempty"`
	ContentRatings      Object           `json:"content_ratings,omitempty"`
	Reviews             Object           `json:"reviews,omitempty"`
	WatchProviders      Object           `json:"watch/providers,omitempty"`
	Moods               Object           `json:"moods,omitempty"`
	Soundtracks         Object           `json:"soundtracks,omitempty"`
	Ratings             Object           `json:"ratings,omitempty"`
	Wikidata            Object           `json:"wikidata,omitempty"`
}

// Season is the answer of TV.Season.
type Season struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Overview     string    `json:"overview,omitempty"`
	AirDate      string    `json:"air_date,omitempty"`
	SeasonNumber int       `json:"season_number"`
	PosterPath   string    `json:"poster_path,omitempty"`
	VoteAverage  float64   `json:"vote_average,omitempty"`
	Episodes     []Episode `json:"episodes,omitempty"`
}

// Episode is a TV episode.
type Episode struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview,omitempty"`
	AirDate        string  `json:"air_date,omitempty"`
	EpisodeNumber  int     `json:"episode_number"`
	SeasonNumber   int     `json:"season_number"`
	EpisodeType    string  `json:"episode_type,omitempty"`
	Runtime        int     `json:"runtime,omitempty"`
	ShowID         int     `json:"show_id,omitempty"`
	ProductionCode string  `json:"production_code,omitempty"`
	StillPath      string  `json:"still_path,omitempty"`
	VoteAverage    float64 `json:"vote_average,omitempty"`
	VoteCount      int     `json:"vote_count,omitempty"`
}

// Person is the short person form used in lists and search results.
type Person struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	OriginalName       string   `json:"original_name,omitempty"`
	Gender             int      `json:"gender,omitempty"`
	KnownForDepartment string   `json:"known_for_department,omitempty"`
	ProfilePath        string   `json:"profile_path,omitempty"`
	Popularity         float64  `json:"popularity,omitempty"`
	Adult              bool     `json:"adult,omitempty"`
	KnownFor           []Object `json:"known_for,omitempty"`
}

// PersonDetails is the answer of People.Details.
type PersonDetails struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	AlsoKnownAs        []string `json:"also_known_as,omitempty"`
	Biography          string   `json:"biography,omitempty"`
	Birthday           string   `json:"birthday,omitempty"`
	Deathday           string   `json:"deathday,omitempty"`
	Gender             int      `json:"gender,omitempty"`
	Homepage           string   `json:"homepage,omitempty"`
	IMDbID             string   `json:"imdb_id,omitempty"`
	KnownForDepartment string   `json:"known_for_department,omitempty"`
	PlaceOfBirth       string   `json:"place_of_birth,omitempty"`
	ProfilePath        string   `json:"profile_path,omitempty"`
	Popularity         float64  `json:"popularity,omitempty"`
	Adult              bool     `json:"adult,omitempty"`
}

// Cast is a cast member of a movie or TV show.
type Cast struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name,omitempty"`
	Character          string  `json:"character,omitempty"`
	CreditID           string  `json:"credit_id,omitempty"`
	Order              int     `json:"order"`
	Gender             int     `json:"gender,omitempty"`
	KnownForDepartment string  `json:"known_for_department,omitempty"`
	ProfilePath        string  `json:"profile_path,omitempty"`
	Popularity         float64 `json:"popularity,omitempty"`
	Adult              bool    `json:"adult,omitempty"`
}

// Crew is a crew member of a movie or TV show.
type Crew struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name,omitempty"`
	Job                string  `json:"job,omitempty"`
	Department         string  `json:"department,omitempty"`
	CreditID           string  `json:"credit_id,omitempty"`
	Gender             int     `json:"gender,omitempty"`
	KnownForDepartment string  `json:"known_for_department,omitempty"`
	ProfilePath        string  `json:"profile_path,omitempty"`
	Popularity         float64 `json:"popularity,omitempty"`
	Adult              bool    `json:"adult,omitempty"`
}

// Credits is the cast and crew of a title.
type Credits struct {
	ID   int    `json:"id,omitempty"`
	Cast []Cast `json:"cast"`
	Crew []Crew `json:"crew"`
}

// ImagesConfiguration holds the image base URLs and size lists.
type ImagesConfiguration struct {
	BaseURL       string   `json:"base_url"`
	SecureBaseURL string   `json:"secure_base_url"`
	BackdropSizes []string `json:"backdrop_sizes"`
	LogoSizes     []string `json:"logo_sizes"`
	PosterSizes   []string `json:"poster_sizes"`
	ProfileSizes  []string `json:"profile_sizes"`
	StillSizes    []string `json:"still_sizes"`
}

// Configuration is the answer of Client.Configuration.
type Configuration struct {
	Images     ImagesConfiguration `json:"images"`
	ChangeKeys []string            `json:"change_keys,omitempty"`
}

// FindResults is the answer of Client.Find.
type FindResults struct {
	MovieResults     []Movie   `json:"movie_results"`
	PersonResults    []Person  `json:"person_results"`
	TVResults        []TVShow  `json:"tv_results"`
	TVEpisodeResults []Episode `json:"tv_episode_results"`
	TVSeasonResults  []Object  `json:"tv_season_results"`
}

// OMDbRating is one entry of the Ratings list of an OMDb title.
type OMDbRating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

// OMDbTitle is the answer of OMDb.ByID and OMDb.ByTitle. All values are
// strings, as in the OMDb format.
type OMDbTitle struct {
	Title        string       `json:"Title"`
	Year         string       `json:"Year"`
	Rated        string       `json:"Rated,omitempty"`
	Released     string       `json:"Released,omitempty"`
	Runtime      string       `json:"Runtime,omitempty"`
	Genre        string       `json:"Genre,omitempty"`
	Director     string       `json:"Director,omitempty"`
	Writer       string       `json:"Writer,omitempty"`
	Actors       string       `json:"Actors,omitempty"`
	Plot         string       `json:"Plot,omitempty"`
	Language     string       `json:"Language,omitempty"`
	Country      string       `json:"Country,omitempty"`
	Awards       string       `json:"Awards,omitempty"`
	Poster       string       `json:"Poster,omitempty"`
	Ratings      []OMDbRating `json:"Ratings,omitempty"`
	Metascore    string       `json:"Metascore,omitempty"`
	IMDbRating   string       `json:"imdbRating,omitempty"`
	IMDbVotes    string       `json:"imdbVotes,omitempty"`
	IMDbID       string       `json:"imdbID"`
	Type         string       `json:"Type,omitempty"`
	DVD          string       `json:"DVD,omitempty"`
	BoxOffice    string       `json:"BoxOffice,omitempty"`
	Production   string       `json:"Production,omitempty"`
	Website      string       `json:"Website,omitempty"`
	TotalSeasons string       `json:"totalSeasons,omitempty"`
	Response     string       `json:"Response,omitempty"`
}

// OMDbSearchItem is one hit of an OMDb search.
type OMDbSearchItem struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	IMDbID string `json:"imdbID"`
	Type   string `json:"Type,omitempty"`
	Poster string `json:"Poster,omitempty"`
}

// OMDbSearch is the answer of OMDb.Search.
type OMDbSearch struct {
	Search       []OMDbSearchItem `json:"Search"`
	TotalResults string           `json:"totalResults"`
	Response     string           `json:"Response,omitempty"`
}

// PlanPeriod is a billing period of the plan catalog.
type PlanPeriod struct {
	ID     string  `json:"id"`
	Months int     `json:"months"`
	Save   float64 `json:"save"`
}

// Plan is one plan of the catalog. Plans differ only in requests per minute.
type Plan struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	RPM          int            `json:"rpm"`
	MonthlyCents int            `json:"monthlyCents"`
	Totals       map[string]int `json:"totals,omitempty"`
}

// PlanCatalog is the answer of Client.Plans. Names and prices are live, so do
// not hard code them.
type PlanCatalog struct {
	Currency  string       `json:"currency"`
	Periods   []PlanPeriod `json:"periods"`
	Plans     []Plan       `json:"plans"`
	FetchedAt any          `json:"fetchedAt,omitempty"`
}
