package amodb

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Version is the library version. It is sent in the User-Agent header.
const Version = "0.1.0"

// Default endpoints and timeout.
const (
	DefaultBaseURL      = "https://api.amodb.dev/3"
	DefaultOMDbURL      = "https://api.amodb.dev/omdb/"
	DefaultImageBaseURL = "https://img.amodb.dev/t/p/"
	DefaultTimeout      = 15 * time.Second
)

const (
	defaultMaxRetries = 3
	maxRetriesCap     = 3
)

// AuthMode selects how the API key is sent on TMDB style calls.
type AuthMode int

const (
	// AuthBearer sends the key in an Authorization header. This is the default.
	AuthBearer AuthMode = iota
	// AuthQuery sends the key as the api_key query parameter.
	AuthQuery
)

// service is the shared base of the resource services.
type service struct {
	client *Client
}

// Client talks to the amodb API. Create it with NewClient. It is safe for
// concurrent use.
type Client struct {
	Movies         *MoviesService
	TV             *TVService
	People         *PeopleService
	Search         *SearchService
	Discover       *DiscoverService
	Trending       *TrendingService
	Genres         *GenresService
	Collections    *CollectionsService
	Companies      *CompaniesService
	Networks       *NetworksService
	Keywords       *KeywordsService
	Certifications *CertificationsService
	WatchProviders *WatchProvidersService
	OMDb           *OMDbService

	apiKey       string
	baseURL      string
	omdbURL      string
	imageBaseURL string
	language     string
	timeout      time.Duration
	maxRetries   int
	auth         AuthMode
	userAgent    string
	httpClient   *http.Client
	sleep        func(context.Context, time.Duration) error
	now          func() time.Time
}

// Option changes a setting of the client. Pass options to NewClient.
type Option func(*Client)

// WithBaseURL sets the TMDB style base URL. The default is
// https://api.amodb.dev/3.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.baseURL = baseURL
		}
	}
}

// WithOMDbURL sets the OMDb style endpoint. The default is
// https://api.amodb.dev/omdb/.
func WithOMDbURL(omdbURL string) Option {
	return func(c *Client) {
		if omdbURL != "" {
			c.omdbURL = omdbURL
		}
	}
}

// WithImageBaseURL sets the base URL used by ImageURL. The default is
// https://img.amodb.dev/t/p/.
func WithImageBaseURL(imageBaseURL string) Option {
	return func(c *Client) {
		if imageBaseURL != "" {
			c.imageBaseURL = imageBaseURL
		}
	}
}

// WithLanguage sets the default language, for example "ru-RU". A language in
// the Params of a call overrides it. It is not sent on OMDb calls.
func WithLanguage(language string) Option {
	return func(c *Client) {
		c.language = language
	}
}

// WithTimeout sets the timeout of each HTTP attempt. The default is 15
// seconds. A value of 0 or less turns the client timeout off, and only the
// deadline of the context applies.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithMaxRetries sets how many times an HTTP 429 answer is retried. The value
// is clamped to 0..3. The default is 3.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		retries := n
		if retries < 0 {
			retries = 0
		}
		if retries > maxRetriesCap {
			retries = maxRetriesCap
		}
		c.maxRetries = retries
	}
}

// WithAuth selects AuthBearer (the default) or AuthQuery.
func WithAuth(mode AuthMode) Option {
	return func(c *Client) {
		c.auth = mode
	}
}

// WithHTTPClient sets the HTTP client, for example to add a proxy or a custom
// transport. A nil value is ignored.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithUserAgent sets the User-Agent header. The default is amodb-go/0.1.0.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		if userAgent != "" {
			c.userAgent = userAgent
		}
	}
}

// WithSleep replaces the function that waits between retries. Tests use it to
// record the waits instead of sleeping. A nil value is ignored.
func WithSleep(sleep func(context.Context, time.Duration) error) Option {
	return func(c *Client) {
		if sleep != nil {
			c.sleep = sleep
		}
	}
}

// NewClient creates a client. It returns ErrMissingAPIKey when apiKey is
// empty. amodb is a paid API, get a key at https://amodb.dev/pricing.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	key := strings.TrimSpace(apiKey)
	if key == "" {
		return nil, ErrMissingAPIKey
	}
	c := &Client{apiKey: key}
	c.baseURL = DefaultBaseURL
	c.omdbURL = DefaultOMDbURL
	c.imageBaseURL = DefaultImageBaseURL
	c.timeout = DefaultTimeout
	c.maxRetries = defaultMaxRetries
	c.auth = AuthBearer
	c.userAgent = "amodb-go/" + Version
	c.httpClient = &http.Client{}
	c.sleep = sleepContext
	c.now = time.Now
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	c.baseURL = strings.TrimRight(c.baseURL, "/")

	c.Movies = &MoviesService{client: c}
	c.TV = &TVService{client: c}
	c.People = &PeopleService{client: c}
	c.Search = &SearchService{client: c}
	c.Discover = &DiscoverService{client: c}
	c.Trending = &TrendingService{client: c}
	c.Genres = &GenresService{client: c}
	c.Collections = &CollectionsService{client: c}
	c.Companies = &CompaniesService{client: c}
	c.Networks = &NetworksService{client: c}
	c.Keywords = &KeywordsService{client: c}
	c.Certifications = &CertificationsService{client: c}
	c.WatchProviders = &WatchProvidersService{client: c}
	c.OMDb = &OMDbService{client: c}
	return c, nil
}

// Find looks up a title or person by an external id. externalSource is the
// external_source value of the API, for example imdb_id. The API reference
// lists the other accepted values.
func (c *Client) Find(ctx context.Context, externalID, externalSource string) (*FindResults, error) {
	params := Params{"external_source": externalSource}
	return get[FindResults](ctx, c, buildPath("find", externalID), params)
}

// Configuration returns the image base URLs and the size lists.
func (c *Client) Configuration(ctx context.Context) (*Configuration, error) {
	return get[Configuration](ctx, c, "/configuration", nil)
}

// Plans returns the public plan catalog. The URL is built from the origin of
// the base URL, and no key is sent.
func (c *Client) Plans(ctx context.Context) (*PlanCatalog, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("amodb: cannot derive the plans URL from base URL %q", c.baseURL)
	}
	target := u.Scheme + "://" + u.Host + "/plans"
	res, err := c.fetch(ctx, target, nil, scopePublic)
	if err != nil {
		return nil, err
	}
	var out PlanCatalog
	if derr := decodeJSON(res, &out); derr != nil {
		return nil, derr
	}
	return &out, nil
}

// ImageURL builds an image URL from a path such as "/abc123.jpg". An empty
// size means "w500". Valid sizes are w92, w154, w185, w300, w342, w500, w780,
// w1280, h632 and original. It returns an empty string when path is empty.
func (c *Client) ImageURL(path, size string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if size == "" {
		size = "w500"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return strings.TrimRight(c.imageBaseURL, "/") + "/" + size + path
}
