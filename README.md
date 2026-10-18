<p align="center">
  <a href="https://amodb.dev">
    <img src="https://raw.githubusercontent.com/amodb-dev/.github/main/assets/lockup.png" alt="amodb" width="260">
  </a>
</p>

<h1 align="center">amodb-go</h1>

<p align="center">Go client for the amodb movie and TV metadata API. Standard library only.</p>

<p align="center">
  <a href="https://github.com/amodb-dev/amodb-go/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/amodb-dev/amodb-go/actions/workflows/ci.yml/badge.svg"></a>
  <a href="https://pkg.go.dev/github.com/amodb-dev/amodb-go"><img alt="Latest tag" src="https://img.shields.io/github/v/tag/amodb-dev/amodb-go?sort=semver&label=version"></a>
  <a href="https://github.com/amodb-dev/amodb-go/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/github/license/amodb-dev/amodb-go"></a>
</p>

Go client for the amodb movie and TV metadata API. It is compatible with TMDB v3 and OMDb, uses only the standard library and needs Go 1.22 or newer.

amodb is a paid API. Get a key at https://amodb.dev/pricing.

Existing TMDB clients work by changing the base URL to `https://api.amodb.dev/3`.

## Install

```sh
go get github.com/amodb-dev/amodb-go
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

	"github.com/amodb-dev/amodb-go"
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

The API reads up to 20 names and ignores unknown names without an error.

- Movies: `credits`, `images`, `videos`, `keywords`, `translations`, `external_ids`, `similar`, `recommendations`, `reviews`, `release_dates`, `watch/providers`, `alternative_titles`, `lists`, `moods`, `soundtracks`, `ratings`, `wikidata`, `showtimes`.
- TV shows: the same names, with `content_ratings` and `episode_groups` in place of `release_dates`.
- People: `combined_credits`, `movie_credits`, `tv_credits`, `external_ids`, `images`.

Sections without a typed field, such as `Images`, are `amodb.Object` values. The appended `moods`, `ratings` and `soundtracks` sections have no `id` key, so their `ID` field is 0. TV works the same way: `client.TV.Details(ctx, 1399, nil)`, `client.TV.Season(ctx, 1399, 1, nil)` and `client.TV.Episode(ctx, 1399, 1, 2, nil)`.

Other calls in the same style: `client.Movies.Latest(ctx, nil)`, `client.TV.Latest(ctx, nil)`, `client.Movies.Lists(ctx, 550, nil)` and `client.TV.EpisodeGroups(ctx, 1399, nil)`. List endpoints take `page` from 1 to 500.

Endpoints with a stable shape return structs such as `MovieDetails`, `TVDetails`, `Person`, `Moods` and `Paged[Movie]`. The others return `amodb.Object`, which is `map[string]any`. `client.Do` is the escape hatch for any path:

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
fmt.Println(res.MovieResults[0].Title) // Fight Club

res, err = client.Find(ctx, "361", "kinopoisk_id")
res, err = client.Find(ctx, "Q35332", "wikidata_id")
fmt.Println(res.PersonResults[0].Name) // Brad Pitt
```

These sources return results: `imdb_id` (`tt...` for titles, `nm...` for people), `tvdb_id`, `facebook_id`, `twitter_id`, `instagram_id`, `wikidata_id` and `kinopoisk_id`. `freebase_mid`, `freebase_id`, `tvrage_id`, `tiktok_id` and `youtube_id` are accepted and return empty lists. A missing or unknown source returns an HTTP 400 error with `StatusCode` 22.

`FindResults` has `MovieResults`, `PersonResults`, `TVResults`, `TVEpisodeResults` and `TVSeasonResults`. The last two are always empty. Every item has a `MediaType`.

## Extras

`Moods`, `Soundtracks`, `Ratings`, `Wikidata` and `Showtimes` exist on both `client.Movies` and `client.TV`. Each returns a typed struct.

### Moods

```go
moods, err := client.Movies.Moods(ctx, 155, nil)
if err != nil {
	log.Fatal(err)
}
fmt.Println(moods.Source) // measured
for _, m := range moods.Results {
	fmt.Println(m.Name, m.Avg) // action 3.9
}
```

`Source` is `amodb.MoodSourceMeasured` for viewer data, `amodb.MoodSourceEstimated` when the meters come from genres (whole numbers), or empty with no `Results` when the title has no data. `Avg` and `Values` (up to 24 readings) are on a 0 to 10 scale. The names are action, fear, humor, pace, profanity, romance, sadness, sex and violence. A title may have fewer than nine.

### Soundtracks

```go
sd, err := client.Movies.Soundtracks(ctx, 155, nil)
if err != nil {
	log.Fatal(err)
}
for _, album := range sd.Results {
	fmt.Println(album.Album, album.TrackCount, "tracks") // The Dark Knight 14 tracks
	for _, track := range album.Tracks {
		fmt.Println(track.Number, track.Title, track.DurationMS)
	}
}
```

`Results` is empty when no album is known. An album has `Source` (`musicbrainz` or `deezer`), `Album`, `Artist`, `Released`, `URL`, `TrackCount` and `Tracks`. `Artist`, `Released` and `URL` are empty when unknown.

### Ratings

```go
ratings, err := client.Movies.Ratings(ctx, 550, nil)
if err != nil {
	log.Fatal(err)
}
if r := ratings.IMDb; r != nil && r.Votes != nil {
	fmt.Println(r.Value, "from", *r.Votes, "votes") // 8.8 from 2669989 votes
}
```

The providers are `IMDb`, `Kinopoisk`, `Metacritic`, `RottenTomatoes` and `TMDB`. Each is a `*amodb.Rating` with `Value` and `Votes`. A provider without a score is `nil`, and `Votes` is `nil` when the count is unknown. IMDb, Kinopoisk and TMDB are out of 10, Metacritic and RottenTomatoes out of 100.

### Wikidata

```go
wd, err := client.Movies.Wikidata(ctx, 604, nil)
if err != nil {
	log.Fatal(err)
}
if wd.Budget != nil {
	fmt.Println(wd.Budget.Value, wd.Budget.Currency) // 150000000 USD
}
if wd.Follows != nil {
	fmt.Println(wd.Follows.Label, wd.Follows.TMDBID) // The Matrix 603
}
```

`Budget` and `BoxOffice` hold `Value` and `Currency` (`BoxOffice` adds `Worldwide`) and are `nil` when unknown. `Follows` and `FollowedBy` are `*amodb.WikidataLink` or `nil`. `Sitelinks` maps a language code to a Wikipedia page title. `Awards`, `Nominations`, `BasedOn`, `FilmingLocations`, `NarrativeLocations` and `OriginalBroadcasters` are slices. A link has `WikidataID` and `Label`, plus `TMDBID` and `MediaType` when the item is a title in the catalog.

### Showtimes

```go
// Which countries have showtimes
all, err := client.Movies.Showtimes(ctx, 550, "", nil)
if err != nil {
	log.Fatal(err)
}
for _, country := range all.Countries {
	fmt.Println(country.Country, country.Venues, country.Showtimes)
}

// Venues and showtimes for one country
ru, err := client.Movies.Showtimes(ctx, 550, "RU", nil)
if err != nil {
	log.Fatal(err)
}
for _, venue := range ru.Results {
	for _, show := range venue.Showtimes {
		fmt.Println(venue.Name, show.Start)
	}
}
```

The country is optional. An empty string sends no country: `Results` is empty and `Countries` holds the counts per country (`Venues`, `Cities`, `Showtimes`, `FirstDay`, `LastDay`). With a country code, `Results` lists the venues and each venue has a `Showtimes` slice. `Start` is the local time at the venue. Only days from today on are listed. To get the venues inside `append_to_response`, pass `"country"` in `Params`.

## Calls that are not offered

These TMDB paths answer 404 with `status_code` 34, so the client has no method for them: `aggregate_credits`, `changes`, `person/{id}/translations`, `person/{id}/tagged_images`, `credit/{id}`, `review/{id}`, `tv/episode_group/{id}`, images, translations and alternative names of collections, companies and networks, and `configuration/timezones`. Everything under `account`, `authentication`, `guest_session` and `list` answers 405 with `status_code` 37, and so does any `rating` or `account_states` path.

Some served calls return little: `Reviews` is empty for most titles, and season and episode videos and translations are empty. The `include_adult` parameter is ignored, because adult titles and people are never served.

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

Tests use `net/http/httptest` servers and never call the live API. They run on Go 1.22 and 1.23 in CI.

## Links

- Website: https://amodb.dev
- Docs: https://amodb.dev/docs
- Pricing: https://amodb.dev/pricing
- Status: https://amodb.dev/status
- Issues: https://github.com/amodb-dev/amodb-go/issues
- Changelog: https://github.com/amodb-dev/amodb-go/blob/main/CHANGELOG.md
- Security: https://github.com/amodb-dev/amodb-go/blob/main/SECURITY.md
- Organization: https://github.com/amodb-dev

## License

MIT. See `LICENSE`.
