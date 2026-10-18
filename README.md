<!-- GitHub topics: amodb, tmdb, omdb, go, golang, movie-api, tv-api, api-client, sdk -->

# amodb-go

Go client for the amodb movie and TV metadata API. It is compatible with TMDB v3 and OMDb, uses only the standard library and targets Go 1.21.

amodb is a paid API. Get a key at https://amodb.dev/pricing.

Existing TMDB clients work by changing the base URL to `https://api.amodb.dev/3`.

## Install

```sh
go get github.com/netlistenany-bot/amodb-go
```

The import path ends in `amodb-go`, the package name is `amodb`.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/netlistenany-bot/amodb-go"
)

func main() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	movie, err := client.Movies.Details(context.Background(), 550, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(movie.Title, movie.ReleaseDate)
}
```

`NewClient` returns an error when the key is empty. Read the key from the environment and never commit it.

Every network method takes a `context.Context` first and an optional `amodb.Params` map last. Pass `nil` when you have no extra parameters.

## Search

```go
page, err := client.Search.Movie(ctx, "inception", amodb.Params{"year": 2010, "page": 1})
if err != nil {
	log.Fatal(err)
}
for _, m := range page.Results {
	fmt.Println(m.ID, m.Title, m.ReleaseDate)
}
```

`Search.TV`, `Search.Person`, `Search.Multi`, `Search.Collection`, `Search.Company` and `Search.Keyword` work the same way. List endpoints such as `Movies.Popular` and `Discover.Movie` take `page` in `Params`.

## Details with append_to_response

Slices are joined with commas, booleans become `true` or `false` and `nil` values are dropped.

```go
movie, err := client.Movies.Details(ctx, 550, amodb.Params{
	"append_to_response": []string{"credits", "images"},
})
if err != nil {
	log.Fatal(err)
}
for _, member := range movie.Credits.Cast {
	fmt.Println(member.Name, member.Character)
}
```

Sections without a typed field, such as `Images`, are `amodb.Object` values. TV works the same way: `client.TV.Details(ctx, 1399, nil)`, `client.TV.Season(ctx, 1399, 1, nil)` and `client.TV.Episode(ctx, 1399, 1, 2, nil)`.

Endpoints with a stable shape return structs such as `MovieDetails`, `TVDetails`, `Person` and `Paged[Movie]`. The others return `amodb.Object`, which is `map[string]any`. `client.Do` is the escape hatch for any path:

```go
var out map[string]any
err := client.Do(ctx, "/movie/550/credits", nil, &out)
```

## Language

```go
client, _ := amodb.NewClient(key, amodb.WithLanguage("ru-RU"))

// A language in Params overrides the client default for one call.
movie, _ := client.Movies.Details(ctx, 550, amodb.Params{"language": "pt-BR"})
```

The service serves 14 base languages: en ru uk be zh pt bg he cs ro fr pl it es. Any other tag falls back to `en`.

## Images

```go
movie, _ := client.Movies.Details(ctx, 550, nil)
fmt.Println(client.ImageURL(movie.PosterPath, "w500"))
// https://img.amodb.dev/t/p/w500/...
```

Sizes are `w92 w154 w185 w300 w342 w500 w780 w1280 h632 original`. An empty size means `w500`. An empty path returns an empty string. Image URLs need no key.

## OMDb

```go
title, err := client.OMDb.ByID(ctx, "tt0137523", amodb.Params{"plot": "full"})
if err != nil {
	log.Fatal(err)
}
fmt.Println(title.Title, title.Year, title.IMDbRating)

found, err := client.OMDb.Search(ctx, "fight club", amodb.Params{"type": "movie", "page": 1})
```

`OMDb.ByTitle` looks up an exact title. Params use the OMDb names `plot`, `type`, `y` and `page`, and `year` is accepted as an alias of `y`. OMDb calls always send the key as `apikey`. An answer with `"Response":"False"` becomes an `*amodb.Error` with `Status` 200 and the OMDb text in `StatusMessage`.

## Find by external id

```go
res, err := client.Find(ctx, "tt0137523", "imdb_id")
if err != nil {
	log.Fatal(err)
}
fmt.Println(res.MovieResults[0].Title)
```

The source is `imdb_id` for IMDb ids. The API reference at https://amodb.dev lists the other accepted sources.

## Extras

These amodb endpoints work for movies and TV: `Moods`, `Soundtracks`, `Ratings` and `Wikidata`. `Movies.Showtimes` takes a country code. They also work in `append_to_response`.

```go
ratings, err := client.Movies.Ratings(ctx, 550, nil)
showtimes, err := client.Movies.Showtimes(ctx, 550, "US", nil)
```

## Error handling

Failed calls return `*amodb.Error` with `Status` (HTTP status, 0 when no response arrived), `StatusCode`, `StatusMessage`, `Body` and `RetryAfter`.

```go
_, err := client.Movies.Details(ctx, 1, nil)

var apiErr *amodb.Error
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.Status, apiErr.StatusCode, apiErr.StatusMessage)
}

switch {
case amodb.IsNotFound(err):
	fmt.Println("no such movie")
case amodb.IsRateLimited(err):
	fmt.Println("over the limit")
case amodb.IsUnauthorized(err):
	fmt.Println("check the key")
case errors.Is(err, amodb.ErrTimeout):
	fmt.Println("timed out")
case errors.Is(err, amodb.ErrNetwork):
	fmt.Println("network failure")
}
```

A canceled context returns `context.Canceled` unchanged.

## Retries and timeouts

Only HTTP 429 answers are retried. The client waits for `Retry-After` (seconds or an HTTP date, clamped to 0 to 60 seconds). Without the header it waits 1, 2 and 4 seconds. The default is 3 retries, which is also the maximum.

Each HTTP attempt has a 15 second timeout. A deadline on your context applies when it is earlier.

```go
client, err := amodb.NewClient(key,
	amodb.WithTimeout(5*time.Second),
	amodb.WithMaxRetries(1),
)
```

## Custom http.Client

```go
hc := &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
client, err := amodb.NewClient(key, amodb.WithHTTPClient(hc))
```

Other options: `WithBaseURL`, `WithOMDbURL`, `WithImageBaseURL`, `WithAuth(amodb.AuthQuery)`, `WithUserAgent` and `WithSleep`. The default auth mode sends `Authorization: Bearer KEY`. `AuthQuery` sends `api_key=KEY` instead.

## Plans and rate limits

Every plan gets all data. Plans differ only in requests per minute: lite 20, pro 40, proplus 120. The live catalog is at https://api.amodb.dev/plans and `client.Plans(ctx)` returns it without sending a key.

## Development

```sh
gofmt -l .
go vet ./...
go test ./...
```

Tests use `net/http/httptest` servers and never call the live API. They need Go 1.21 or newer.

## License

MIT. See `LICENSE`.
