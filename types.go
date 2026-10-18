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
	MediaType        string  `json:"media_type,omitempty"`
}

// MovieDetails is the answer of Movies.Details and Movies.Latest. The
// sections that append_to_response adds are decoded into the fields at the
// end. A section without a field here can be read with Client.Do into your
// own struct.
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
	OriginCountry       []string         `json:"origin_country,omitempty"`
	Genres              []Genre          `json:"genres,omitempty"`
	ProductionCompanies []Company        `json:"production_companies,omitempty"`
	ProductionCountries []Country        `json:"production_countries,omitempty"`
	SpokenLanguages     []SpokenLanguage `json:"spoken_languages,omitempty"`
	BelongsToCollection *CollectionRef   `json:"belongs_to_collection,omitempty"`
	Credits             *Credits         `json:"credits,omitempty"`
	Recommendations     *Paged[Movie]    `json:"recommendations,omitempty"`
	Similar             *Paged[Movie]    `json:"similar,omitempty"`
	Lists               *Paged[UserList] `json:"lists,omitempty"`
	Images              Object           `json:"images,omitempty"`
	Videos              Object           `json:"videos,omitempty"`
	Keywords            Object           `json:"keywords,omitempty"`
	ExternalIDs         *ExternalIDs     `json:"external_ids,omitempty"`
	ReleaseDates        Object           `json:"release_dates,omitempty"`
	Reviews             Object           `json:"reviews,omitempty"`
	Translations        Object           `json:"translations,omitempty"`
	AlternativeTitles   Object           `json:"alternative_titles,omitempty"`
	WatchProviders      Object           `json:"watch/providers,omitempty"`
	Moods               *Moods           `json:"moods,omitempty"`
	Soundtracks         *Soundtracks     `json:"soundtracks,omitempty"`
	Ratings             *Ratings         `json:"ratings,omitempty"`
	Wikidata            *Wikidata        `json:"wikidata,omitempty"`
	Showtimes           *Showtimes       `json:"showtimes,omitempty"`
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
	MediaType        string   `json:"media_type,omitempty"`
}

// CreatedBy is a creator of a TV show. The created_by list is filled for only
// a few shows, so TVDetails.CreatedBy is empty for most.
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

// TVDetails is the answer of TV.Details and TV.Latest. The sections that
// append_to_response adds are decoded into the fields at the end. The API
// sends no type field for a show.
type TVDetails struct {
	ID                  int              `json:"id"`
	Name                string           `json:"name"`
	OriginalName        string           `json:"original_name,omitempty"`
	OriginalLanguage    string           `json:"original_language,omitempty"`
	Overview            string           `json:"overview,omitempty"`
	Tagline             string           `json:"tagline,omitempty"`
	Status              string           `json:"status,omitempty"`
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
	Lists               *Paged[UserList] `json:"lists,omitempty"`
	Images              Object           `json:"images,omitempty"`
	Videos              Object           `json:"videos,omitempty"`
	Keywords            Object           `json:"keywords,omitempty"`
	ExternalIDs         *ExternalIDs     `json:"external_ids,omitempty"`
	ContentRatings      Object           `json:"content_ratings,omitempty"`
	Reviews             Object           `json:"reviews,omitempty"`
	Translations        Object           `json:"translations,omitempty"`
	AlternativeTitles   Object           `json:"alternative_titles,omitempty"`
	EpisodeGroups       Object           `json:"episode_groups,omitempty"`
	WatchProviders      Object           `json:"watch/providers,omitempty"`
	Moods               *Moods           `json:"moods,omitempty"`
	Soundtracks         *Soundtracks     `json:"soundtracks,omitempty"`
	Ratings             *Ratings         `json:"ratings,omitempty"`
	Wikidata            *Wikidata        `json:"wikidata,omitempty"`
	Showtimes           *Showtimes       `json:"showtimes,omitempty"`
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
	MediaType          string   `json:"media_type,omitempty"`
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

// FindResults is the answer of Client.Find. Every item has its MediaType set.
// The API always sends empty TVEpisodeResults and TVSeasonResults lists.
type FindResults struct {
	MovieResults     []Movie   `json:"movie_results"`
	PersonResults    []Person  `json:"person_results"`
	TVResults        []TVShow  `json:"tv_results"`
	TVEpisodeResults []Episode `json:"tv_episode_results"`
	TVSeasonResults  []Object  `json:"tv_season_results"`
}

// UserList is a user made list that contains a title.
type UserList struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	FavoriteCount int    `json:"favorite_count,omitempty"`
	ItemCount     int    `json:"item_count,omitempty"`
	ISO6391       string `json:"iso_639_1,omitempty"`
	ISO31661      string `json:"iso_3166_1,omitempty"`
	PosterPath    string `json:"poster_path,omitempty"`
}

// ExternalIDs holds the ids of a movie or TV show on other services. A movie
// has no FreebaseMID, FreebaseID, TVDBID or TVRageID. KinopoiskID is a string
// such as "361". Empty values mean the id is unknown. In append_to_response
// the id field is missing.
type ExternalIDs struct {
	ID          int    `json:"id,omitempty"`
	IMDbID      string `json:"imdb_id,omitempty"`
	WikidataID  string `json:"wikidata_id,omitempty"`
	FacebookID  string `json:"facebook_id,omitempty"`
	InstagramID string `json:"instagram_id,omitempty"`
	TwitterID   string `json:"twitter_id,omitempty"`
	KinopoiskID string `json:"kinopoisk_id,omitempty"`
	FreebaseMID string `json:"freebase_mid,omitempty"`
	FreebaseID  string `json:"freebase_id,omitempty"`
	TVDBID      int    `json:"tvdb_id,omitempty"`
	TVRageID    int    `json:"tvrage_id,omitempty"`
}

// PersonExternalIDs holds the ids of a person on other services. Empty values
// mean the id is unknown.
type PersonExternalIDs struct {
	ID          int    `json:"id"`
	IMDbID      string `json:"imdb_id,omitempty"`
	WikidataID  string `json:"wikidata_id,omitempty"`
	FacebookID  string `json:"facebook_id,omitempty"`
	InstagramID string `json:"instagram_id,omitempty"`
	TikTokID    string `json:"tiktok_id,omitempty"`
	TwitterID   string `json:"twitter_id,omitempty"`
	YouTubeID   string `json:"youtube_id,omitempty"`
}

// MoodSource tells where the meters of a title come from.
type MoodSource string

// The mood sources. MoodSourceMeasured means the meters come from viewer data.
// MoodSourceEstimated means they are derived from the genres. A title without
// mood data has an empty source.
const (
	MoodSourceMeasured  MoodSource = "measured"
	MoodSourceEstimated MoodSource = "estimated"
)

// MoodMeter is one mood of a title, for example "action" or "fear". Avg and
// Values are on a 0 to 10 scale. Values holds up to 24 readings. Estimated
// meters have whole numbers.
type MoodMeter struct {
	Name   string    `json:"name"`
	Avg    float64   `json:"avg"`
	Values []float64 `json:"values"`
}

// Moods is the answer of Movies.Moods and TV.Moods. A title without mood data
// has an empty Source and no Results. In append_to_response the id field is
// missing.
type Moods struct {
	ID      int         `json:"id,omitempty"`
	Source  MoodSource  `json:"source,omitempty"`
	Results []MoodMeter `json:"results"`
}

// Track is one track of a soundtrack album.
type Track struct {
	Disc       int    `json:"disc"`
	Number     int    `json:"number"`
	Title      string `json:"title"`
	Artist     string `json:"artist,omitempty"`
	DurationMS int    `json:"duration_ms"`
}

// SoundtrackAlbum is one album of Soundtracks. Source is "musicbrainz" or
// "deezer". Artist, Released and URL are empty when unknown.
type SoundtrackAlbum struct {
	Source     string  `json:"source"`
	Album      string  `json:"album"`
	Artist     string  `json:"artist,omitempty"`
	Released   string  `json:"released,omitempty"`
	URL        string  `json:"url,omitempty"`
	TrackCount int     `json:"track_count"`
	Tracks     []Track `json:"tracks"`
}

// Soundtracks is the answer of Movies.Soundtracks and TV.Soundtracks. Results
// is empty when no album is known. In append_to_response the id field is
// missing.
type Soundtracks struct {
	ID      int               `json:"id,omitempty"`
	Results []SoundtrackAlbum `json:"results"`
}

// Rating is the score of one provider. Votes is nil when the provider gives
// no vote count.
type Rating struct {
	Value float64 `json:"value"`
	Votes *int    `json:"votes"`
}

// Ratings is the answer of Movies.Ratings and TV.Ratings. A provider without
// a score is nil. IMDb, Kinopoisk and TMDB are out of 10. Metacritic and
// RottenTomatoes are out of 100. In append_to_response the id field is
// missing.
type Ratings struct {
	ID             int     `json:"id,omitempty"`
	IMDb           *Rating `json:"imdb"`
	Kinopoisk      *Rating `json:"kinopoisk"`
	Metacritic     *Rating `json:"metacritic"`
	RottenTomatoes *Rating `json:"rotten_tomatoes"`
	TMDB           *Rating `json:"tmdb"`
}

// WikidataMoney is an amount rounded to whole units. Currency is an ISO code
// and is empty when unknown.
type WikidataMoney struct {
	Value    int64  `json:"value"`
	Currency string `json:"currency,omitempty"`
}

// WikidataBoxOffice is the box office of a title. Worldwide is true when the
// amount is the worldwide total.
type WikidataBoxOffice struct {
	Value     int64  `json:"value"`
	Currency  string `json:"currency,omitempty"`
	Worldwide bool   `json:"worldwide"`
}

// WikidataLink points to a Wikidata item. TMDBID and MediaType are set only
// when the item is a title in the catalog. MediaType is "movie" or "tv".
type WikidataLink struct {
	WikidataID string `json:"wikidata_id"`
	Label      string `json:"label,omitempty"`
	TMDBID     int    `json:"tmdb_id,omitempty"`
	MediaType  string `json:"media_type,omitempty"`
}

// WikidataEntry is a WikidataLink with the year of an award or nomination.
type WikidataEntry struct {
	WikidataLink
	Year int `json:"year,omitempty"`
}

// WikidataRef is a Wikidata item with its label.
type WikidataRef struct {
	WikidataID string `json:"wikidata_id"`
	Label      string `json:"label,omitempty"`
}

// WikidataSource is a work that a title is based on. Author is nil when it is
// unknown.
type WikidataSource struct {
	WikidataID string       `json:"wikidata_id"`
	Label      string       `json:"label,omitempty"`
	Author     *WikidataRef `json:"author"`
}

// Wikidata is the answer of Movies.Wikidata and TV.Wikidata. Budget,
// BoxOffice, Follows and FollowedBy are nil when unknown. Sitelinks maps a
// language code to a Wikipedia page title. Labels use the request language and
// fall back to English.
type Wikidata struct {
	ID                   int                `json:"id"`
	WikidataID           string             `json:"wikidata_id,omitempty"`
	Budget               *WikidataMoney     `json:"budget"`
	BoxOffice            *WikidataBoxOffice `json:"box_office"`
	Follows              *WikidataLink      `json:"follows"`
	FollowedBy           *WikidataLink      `json:"followed_by"`
	Sitelinks            map[string]string  `json:"sitelinks"`
	Awards               []WikidataEntry    `json:"awards"`
	Nominations          []WikidataEntry    `json:"nominations"`
	BasedOn              []WikidataSource   `json:"based_on"`
	FilmingLocations     []WikidataLink     `json:"filming_locations"`
	NarrativeLocations   []WikidataLink     `json:"narrative_locations"`
	OriginalBroadcasters []WikidataLink     `json:"original_broadcasters"`
}

// ShowtimeCountry counts the venues and showtimes of one country. It is part
// of Showtimes.Countries.
type ShowtimeCountry struct {
	Country   string `json:"country"`
	Name      string `json:"name"`
	Venues    int    `json:"venues"`
	Cities    int    `json:"cities"`
	Showtimes int    `json:"showtimes"`
	FirstDay  string `json:"first_day"`
	LastDay   string `json:"last_day"`
}

// Showtime is one screening. Start and End are local times at the venue, for
// example "2026-10-18T21:40". Price is nil when unknown.
type Showtime struct {
	Date       string   `json:"date"`
	Start      string   `json:"start"`
	End        string   `json:"end,omitempty"`
	Auditorium string   `json:"auditorium,omitempty"`
	Formats    []string `json:"formats"`
	Tags       []string `json:"tags"`
	Audio      string   `json:"audio,omitempty"`
	Subtitles  string   `json:"subtitles,omitempty"`
	Language   string   `json:"language,omitempty"`
	Price      *float64 `json:"price"`
	Currency   string   `json:"currency,omitempty"`
	Prices     any      `json:"prices"`
	BookingURL string   `json:"booking_url,omitempty"`
	SoldOut    bool     `json:"sold_out"`
}

// VenueChain is the cinema chain of a venue.
type VenueChain struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
}

// Venue is a cinema with its showtimes for one title. Latitude and Longitude
// are nil when unknown.
type Venue struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Chain      VenueChain `json:"chain"`
	City       string     `json:"city,omitempty"`
	Region     string     `json:"region,omitempty"`
	Address    string     `json:"address,omitempty"`
	PostalCode string     `json:"postal_code,omitempty"`
	Country    string     `json:"country"`
	Timezone   string     `json:"timezone,omitempty"`
	Latitude   *float64   `json:"latitude"`
	Longitude  *float64   `json:"longitude"`
	Phone      string     `json:"phone,omitempty"`
	URL        string     `json:"url,omitempty"`
	Showtimes  []Showtime `json:"showtimes"`
}

// Showtimes is the answer of Movies.Showtimes and TV.Showtimes. Without a
// country, Results is empty and Countries holds the counts per country. With a
// country, Country is set and Results lists the venues. Only days from today
// on are listed.
type Showtimes struct {
	ID        int               `json:"id"`
	Country   string            `json:"country,omitempty"`
	Results   []Venue           `json:"results"`
	Countries []ShowtimeCountry `json:"countries"`
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
