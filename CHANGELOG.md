# Changelog

All notable changes to this project are listed here. The format follows Keep a Changelog and the project uses semantic versioning.

## Unreleased

Align extension types with the live API.

### Changed

- `Moods`, `Soundtracks`, `Ratings`, `Wikidata` and `Showtimes` return typed structs instead of `Object`, on both `Movies` and `TV`. The fields of `MovieDetails` and `TVDetails` that append_to_response fills use the same types.
- `Moods` has a `Source` (`MoodSourceMeasured`, `MoodSourceEstimated` or empty) and `Results` with `Name`, `Avg` and `Values` on a 0 to 10 scale.
- `Ratings` holds a `*Rating` with `Value` and `Votes` per provider. A provider without a score is `nil`.
- `Wikidata` money fields use `Value` and `Currency`. Links carry `TMDBID` and `MediaType` for catalog titles.
- `Showtimes` takes an optional country. An empty string sends no country and returns the countries with counts. A country code returns the venues and their showtimes. The method is now on `TV` too.
- `Movies.ExternalIDs` and `TV.ExternalIDs` return `ExternalIDs`, with the new `KinopoiskID`. `People.ExternalIDs` returns `PersonExternalIDs`.
- `Movie`, `TVShow` and `Person` have a `MediaType` field, which find results fill.
- `Client.Find` documents the full list of external sources, including `kinopoisk_id` and `wikidata_id`.
- The README lists the valid append_to_response names and the calls that are not offered.
- Add CI, release and Dependabot workflows.
- The module path is now `github.com/amodb-dev/amodb-go`.
- The minimum Go version is now 1.22, matching the versions tested in CI.

### Added

- `Movies.Latest`, `Movies.Lists`, `TV.Latest`, `TV.Lists`, `TV.Translations`, `TV.AlternativeTitles` and `TV.EpisodeGroups`.
- `MovieDetails.OriginCountry` and `MovieDetails.Lists`. On `TVDetails`, the `Lists`, `Translations`, `AlternativeTitles`, `EpisodeGroups` and `Showtimes` sections for append_to_response.

### Removed

- `TV.AggregateCredits` and `TVDetails.AggregateCredits`. The API does not serve `aggregate_credits`.
- `TVDetails.Type`. The API sends no type field for a show.

## 0.1.0 - 2026-10-18

### Added

- Initial release.
