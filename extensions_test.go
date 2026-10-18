package amodb

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// The bodies in this file are trimmed copies of real replies of the API.

const moodsMeasuredBody = `{"id":155,"source":"measured","results":[{"name":"action","avg":3.9,"values":[3.6,3.4,4.2,4.2,4.8,3]},{"name":"fear","avg":5.2,"values":[3.6,6,4.8,6.4,6.4,3.8]}]}`

const moodsEstimatedTVBody = `{"id":1396,"source":"estimated","results":[{"name":"action","avg":4,"values":[3,4,5,4,3]},{"name":"fear","avg":2,"values":[1,2,3,2,1]}]}`

const moodsEmptyBody = `{"id":1798888,"source":null,"results":[]}`

const soundtracksMusicBrainzBody = `{"id":155,"results":[{"source":"musicbrainz","album":"The Dark Knight","artist":"Hans Zimmer & James Newton Howard","released":"2008-07-15","url":"https://musicbrainz.org/release/814b2f7b-c27c-451e-abac-a316d2e6cea9","track_count":14,"tracks":[{"disc":1,"number":1,"title":"Why So Serious?","artist":"Hans Zimmer","duration_ms":554000},{"disc":1,"number":2,"title":"I'm Not a Hero","artist":"Hans Zimmer","duration_ms":395000}]}]}`

const soundtracksDeezerBody = `{"id":27205,"results":[{"source":"deezer","album":"Inception (Music from the Motion Picture)","artist":"Hans Zimmer","released":"2010-07-09","url":"https://www.deezer.com/album/601778","track_count":12,"tracks":[{"disc":1,"number":1,"title":"Half Remembered Dream","artist":"Hans Zimmer","duration_ms":71000},{"disc":1,"number":2,"title":"We Built Our Own World","artist":"Hans Zimmer","duration_ms":115000}]}]}`

const soundtracksNoArtistBody = `{"id":13,"results":[{"source":"musicbrainz","album":"Forrest Gump: The Soundtrack","artist":null,"released":"1994-06-27","url":"https://musicbrainz.org/release/e4a08d83-945c-3eb8-ab87-099752e4d94c","track_count":32,"tracks":[{"disc":1,"number":1,"title":"Hound Dog","artist":"Elvis Presley","duration_ms":138493}]}]}`

const soundtracksEmptyBody = `{"id":1396,"results":[]}`

const ratingsMovieBody = `{"id":550,"imdb":{"value":8.8,"votes":2669989},"kinopoisk":{"value":8.7,"votes":null},"metacritic":{"value":67,"votes":null},"rotten_tomatoes":{"value":81,"votes":null},"tmdb":{"value":8.437,"votes":33039}}`

const ratingsTVBody = `{"id":1396,"imdb":{"value":9.5,"votes":2671962},"kinopoisk":{"value":8.9,"votes":null},"metacritic":null,"rotten_tomatoes":{"value":96,"votes":null},"tmdb":{"value":9,"votes":18863}}`

const ratingsEmptyBody = `{"id":1798888,"imdb":null,"kinopoisk":null,"metacritic":null,"rotten_tomatoes":null,"tmdb":null}`

const wikidataSequelBody = `{"id":604,"wikidata_id":"Q189600","budget":{"value":150000000,"currency":"USD"},"box_office":{"value":739412035,"currency":"USD","worldwide":true},"follows":{"wikidata_id":"Q83495","label":"The Matrix","tmdb_id":603,"media_type":"movie"},"followed_by":{"wikidata_id":"Q207536","label":"The Matrix Revolutions","tmdb_id":605,"media_type":"movie"},"sitelinks":{"en":"The Matrix Reloaded","ru":"Матрица: Перезагрузка","uk":"Матриця: Перезавантаження"},"awards":[],"nominations":[],"based_on":[],"filming_locations":[{"wikidata_id":"Q3130","label":"Sydney"}],"narrative_locations":[],"original_broadcasters":[]}`

const wikidataAwardsBody = `{"id":550,"wikidata_id":"Q190050","budget":{"value":63000000,"currency":"USD"},"box_office":{"value":101209702,"currency":"USD","worldwide":true},"follows":null,"followed_by":null,"sitelinks":{"en":"Fight Club","ru":"Бойцовский клуб (фильм)","uk":"Бійцівський клуб (фільм)"},"awards":[{"wikidata_id":"Q28044798","label":"Empire Awards, UK for best british actress"},{"wikidata_id":"Q943908","label":"Jupiter Award","year":2000}],"nominations":[{"wikidata_id":"Q488645","label":"Academy Award for Best Sound Editing","year":2000}],"based_on":[{"wikidata_id":"Q837283","label":"Fight Club","author":{"wikidata_id":"Q268181","label":"Chuck Palahniuk"}}],"filming_locations":[{"wikidata_id":"Q65","label":"Los Angeles"}],"narrative_locations":[{"wikidata_id":"Q65","label":"Los Angeles"}],"original_broadcasters":[]}`

const wikidataNoAuthorBody = `{"id":155,"wikidata_id":"Q163872","budget":{"value":185000000,"currency":"USD"},"box_office":{"value":1006234167,"currency":"USD","worldwide":true},"follows":{"wikidata_id":"Q166262","label":"Batman Begins","tmdb_id":272,"media_type":"movie"},"followed_by":{"wikidata_id":"Q189330","label":"The Dark Knight Rises","tmdb_id":49026,"media_type":"movie"},"sitelinks":{"en":"The Dark Knight","ru":"Тёмный рыцарь (фильм)","uk":"Темний лицар (фільм)"},"awards":[{"wikidata_id":"Q106291","label":"Academy Award for Best Supporting Actor","year":2009}],"nominations":[{"wikidata_id":"Q106291","label":"Academy Award for Best Supporting Actor","year":2009}],"based_on":[{"wikidata_id":"Q2695156","label":"Batman","author":null}],"filming_locations":[{"wikidata_id":"Q65","label":"Los Angeles"}],"narrative_locations":[{"wikidata_id":"Q732858","label":"Gotham City"}],"original_broadcasters":[]}`

const showtimesCountriesBody = `{"id":550,"country":null,"results":[],"countries":[{"country":"RU","name":"Russia","venues":2,"cities":1,"showtimes":7,"first_day":"2026-10-18","last_day":"2026-10-21"}]}`

const showtimesRussiaBody = `{"id":550,"country":"RU","results":[{"id":1114,"name":"Мираж Синема НА БОЛЬШОМ","chain":{"key":"mirage-ru","name":"Мираж Синема"},"city":"Санкт-Петербург","region":null,"address":"Большой пр. П.С., д. 35А","postal_code":null,"country":"RU","timezone":"Europe/Moscow","latitude":null,"longitude":null,"phone":null,"url":"https://www.mirage.ru/spb/cinema/1/","showtimes":[{"date":"2026-10-18","start":"2026-10-18T21:40","end":null,"auditorium":null,"formats":["2D"],"tags":[],"audio":null,"subtitles":null,"language":null,"price":null,"currency":null,"prices":null,"booking_url":"https://www.mirage.ru/ticket_new/aca7b114-d587-4fc5-9df6-28568e8fa9b3/","sold_out":false}]}],"countries":[{"country":"RU","name":"Russia","venues":2,"cities":1,"showtimes":7,"first_day":"2026-10-18","last_day":"2026-10-21"}]}`

const showtimesEmptyTVBody = `{"id":1396,"country":null,"results":[],"countries":[]}`

const findKinopoiskBody = `{"movie_results":[{"adult":false,"backdrop_path":"/c6OLXfKAk5BKeR6broC8pYiCquX.jpg","genre_ids":[18],"id":550,"original_language":"en","original_title":"Fight Club","overview":"A ticking-time-bomb insomniac and a slippery soap salesma...","popularity":39.0054,"poster_path":"/jSziioSwPVrOy9Yow3XhWIBDjq1.jpg","release_date":"1999-10-15","title":"Fight Club","video":false,"vote_average":8.437,"vote_count":33039,"media_type":"movie"}],"person_results":[],"tv_results":[],"tv_episode_results":[],"tv_season_results":[]}`

const findPersonBody = `{"movie_results":[],"person_results":[{"adult":false,"gender":2,"id":287,"known_for_department":"Acting","name":"Brad Pitt","original_name":"Brad Pitt","popularity":13.5158,"profile_path":"/ajNaPmXVVMJFg9GWmu6MJzTaXdV.jpg","known_for":[],"media_type":"person"}],"tv_results":[],"tv_episode_results":[],"tv_season_results":[]}`

const externalIDsMovieBody = `{"id":550,"imdb_id":"tt0137523","wikidata_id":"Q190050","facebook_id":"FightClub","instagram_id":null,"twitter_id":null,"kinopoisk_id":"361"}`

const externalIDsTVBody = `{"id":1396,"imdb_id":"tt0903747","freebase_mid":null,"freebase_id":null,"tvdb_id":81189,"tvrage_id":null,"wikidata_id":"Q1079","facebook_id":"BreakingBad","instagram_id":"breakingbad","twitter_id":"BreakingBad","kinopoisk_id":"404900"}`

const externalIDsPersonBody = `{"id":287,"imdb_id":"nm0000093","wikidata_id":"Q35332","facebook_id":null,"instagram_id":null,"tiktok_id":null,"twitter_id":null,"youtube_id":null}`

const listsMovieBody = `{"id":550,"page":1,"results":[{"description":"All of my ratings are the movies that I have seen.","favorite_count":0,"id":8302254,"item_count":2585,"iso_639_1":"en","iso_3166_1":"US","name":"Movies That I have watched all","poster_path":null}],"total_pages":147,"total_results":2924}`

const latestMovieBody = `{"adult":false,"backdrop_path":null,"belongs_to_collection":null,"budget":0,"genres":[{"id":35,"name":"Comedy"}],"homepage":"","id":1798888,"imdb_id":"tt49276703","origin_country":["US"],"original_language":"en","original_title":"In the Mood to Leave","overview":"Molly and Derek need to break up.","popularity":0.1489,"poster_path":"/b55uObskVIJm8n562djLphnljvO.jpg","production_companies":[],"production_countries":[],"release_date":"","revenue":0,"runtime":5,"spoken_languages":[],"status":"Released","tagline":"","title":"In the Mood to Leave","video":false,"vote_average":0,"vote_count":0}`

const latestTVBody = `{"adult":false,"backdrop_path":"/3ax3arGvBqwrvAXTu71hepUGJ0Q.jpg","created_by":[],"episode_run_time":[],"first_air_date":"2000-05-01","genres":[{"id":18,"name":"Drama"}],"homepage":"","id":338695,"in_production":false,"languages":["ru"],"last_air_date":null,"last_episode_to_air":null,"name":"Бандитский Петербург. Барон","next_episode_to_air":null,"networks":[{"id":1199,"logo_path":"/9wVcHld7xpIZ5e6CQ3LJyFSiFhz.png","name":"NTV ","origin_country":"RU"}],"number_of_episodes":5,"number_of_seasons":1,"origin_country":["RU"],"original_language":"ru","original_name":"Бандитский Петербург. Барон","overview":"Absolutely unique storyline, brilliant characters, unforg...","popularity":0.2277,"poster_path":"/gfelnmlZqi3zlLdyybojWwtCEjd.jpg","production_companies":[{"id":19211,"logo_path":"/pjKrG5IC2BaHdE5OUi4Kpe5DlGH.png","name":"KinoMost","origin_country":"RU"}],"production_countries":[{"iso_3166_1":"RU","name":"Russia"}],"seasons":[{"air_date":"2000-05-01","episode_count":5,"id":550759,"name":"Сезон 1","overview":"","poster_path":"/gfelnmlZqi3zlLdyybojWwtCEjd.jpg","season_number":1,"vote_average":0}],"spoken_languages":[{"english_name":"Russian","iso_639_1":"ru","name":"Русский"}],"status":"Ended","tagline":"","vote_average":10,"vote_count":1}`

const appendedMovieBody = `{"id":550,"title":"Fight Club","origin_country":["US","DE"],"lists":{"page":1,"results":[{"description":"All of my ratings are the movies that I have seen.","favorite_count":0,"id":8302254,"item_count":2585,"iso_639_1":"en","iso_3166_1":"US","name":"Movies That I have watched all","poster_path":null},{"description":"Letterboxd's most popular films this week, auto-synced.","favorite_count":0,"id":8634608,"item_count":71,"iso_639_1":"en","iso_3166_1":"US","name":"second-hand prominence","poster_path":null}],"total_pages":147,"total_results":2924},"moods":{"source":"estimated","results":[{"name":"action","avg":4,"values":[3,4,5,4,3]},{"name":"fear","avg":7,"values":[6,7,8,7,6]}]},"ratings":{"imdb":{"value":8.8,"votes":2669989},"kinopoisk":{"value":8.7,"votes":null},"metacritic":{"value":67,"votes":null},"rotten_tomatoes":{"value":81,"votes":null},"tmdb":{"value":8.437,"votes":33039}},"soundtracks":{"results":[]},"external_ids":{"imdb_id":"tt0137523","wikidata_id":"Q190050","facebook_id":"FightClub","instagram_id":null,"twitter_id":null,"kinopoisk_id":"361"},"showtimes":{"id":550,"country":null,"results":[],"countries":[{"country":"RU","name":"Russia","venues":2,"cities":1,"showtimes":7,"first_day":"2026-10-18","last_day":"2026-10-21"}]}}`

const notFoundBody = `{"success":false,"status_code":34,"status_message":"The resource you requested could not be found."}`

const notOfferedBody = `{"success":false,"status_code":37,"status_message":"This method is not offered by amodb."}`

func TestMeasuredMoodsDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, moodsMeasuredBody))
	moods, err := c.Movies.Moods(context.Background(), 155, nil)
	if err != nil {
		t.Fatalf("Movies.Moods: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/movie/155/moods" {
		t.Errorf("path = %q", got)
	}
	if moods.ID != 155 || moods.Source != MoodSourceMeasured {
		t.Errorf("moods = %+v", moods)
	}
	if len(moods.Results) != 2 {
		t.Fatalf("results = %+v", moods.Results)
	}
	action := moods.Results[0]
	if action.Name != "action" || action.Avg != 3.9 || len(action.Values) != 6 {
		t.Errorf("action = %+v", action)
	}
	if action.Values[0] != 3.6 || action.Values[5] != 3 {
		t.Errorf("values = %v", action.Values)
	}
}

func TestEstimatedMoodsOfATVShowDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, moodsEstimatedTVBody))
	moods, err := c.TV.Moods(context.Background(), 1396, nil)
	if err != nil {
		t.Fatalf("TV.Moods: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/tv/1396/moods" {
		t.Errorf("path = %q", got)
	}
	if moods.Source != MoodSourceEstimated || len(moods.Results) != 2 {
		t.Fatalf("moods = %+v", moods)
	}
	if first := moods.Results[0]; first.Avg != 4 || len(first.Values) != 5 {
		t.Errorf("first = %+v", first)
	}
}

func TestMoodsWithoutDataDecoding(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, moodsEmptyBody))
	moods, err := c.Movies.Moods(context.Background(), 1798888, nil)
	if err != nil {
		t.Fatalf("Movies.Moods: %v", err)
	}
	if moods.ID != 1798888 || moods.Source != "" || len(moods.Results) != 0 {
		t.Errorf("moods = %+v, want an empty source and no results", moods)
	}
}

func TestSoundtracksDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, soundtracksMusicBrainzBody))
	soundtracks, err := c.Movies.Soundtracks(context.Background(), 155, nil)
	if err != nil {
		t.Fatalf("Movies.Soundtracks: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/movie/155/soundtracks" {
		t.Errorf("path = %q", got)
	}
	if soundtracks.ID != 155 || len(soundtracks.Results) != 1 {
		t.Fatalf("soundtracks = %+v", soundtracks)
	}
	album := soundtracks.Results[0]
	if album.Source != "musicbrainz" || album.Album != "The Dark Knight" || album.Released != "2008-07-15" {
		t.Errorf("album = %+v", album)
	}
	if album.Artist != "Hans Zimmer & James Newton Howard" || album.TrackCount != 14 {
		t.Errorf("album = %+v", album)
	}
	if len(album.Tracks) != 2 {
		t.Fatalf("tracks = %+v", album.Tracks)
	}
	track := album.Tracks[0]
	if track.Disc != 1 || track.Number != 1 || track.Title != "Why So Serious?" || track.DurationMS != 554000 {
		t.Errorf("track = %+v", track)
	}
}

func TestDeezerSoundtrackDecoding(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, soundtracksDeezerBody))
	soundtracks, err := c.Movies.Soundtracks(context.Background(), 27205, nil)
	if err != nil {
		t.Fatalf("Movies.Soundtracks: %v", err)
	}
	if len(soundtracks.Results) != 1 {
		t.Fatalf("soundtracks = %+v", soundtracks)
	}
	album := soundtracks.Results[0]
	if album.Source != "deezer" || album.Artist != "Hans Zimmer" || album.URL != "https://www.deezer.com/album/601778" {
		t.Errorf("album = %+v", album)
	}
}

func TestSoundtrackAlbumWithoutAnArtistDecoding(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, soundtracksNoArtistBody))
	soundtracks, err := c.Movies.Soundtracks(context.Background(), 13, nil)
	if err != nil {
		t.Fatalf("Movies.Soundtracks: %v", err)
	}
	if len(soundtracks.Results) != 1 || len(soundtracks.Results[0].Tracks) != 1 {
		t.Fatalf("soundtracks = %+v", soundtracks)
	}
	album := soundtracks.Results[0]
	if album.Artist != "" || album.Tracks[0].Artist != "Elvis Presley" {
		t.Errorf("album = %+v", album)
	}
}

func TestEmptySoundtracksDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, soundtracksEmptyBody))
	soundtracks, err := c.TV.Soundtracks(context.Background(), 1396, nil)
	if err != nil {
		t.Fatalf("TV.Soundtracks: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/tv/1396/soundtracks" {
		t.Errorf("path = %q", got)
	}
	if soundtracks.ID != 1396 || len(soundtracks.Results) != 0 {
		t.Errorf("soundtracks = %+v", soundtracks)
	}
}

func TestRatingsDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, ratingsMovieBody))
	ratings, err := c.Movies.Ratings(context.Background(), 550, nil)
	if err != nil {
		t.Fatalf("Movies.Ratings: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/movie/550/ratings" {
		t.Errorf("path = %q", got)
	}
	if ratings.ID != 550 {
		t.Errorf("id = %d", ratings.ID)
	}
	if ratings.IMDb == nil || ratings.IMDb.Value != 8.8 || ratings.IMDb.Votes == nil || *ratings.IMDb.Votes != 2669989 {
		t.Errorf("imdb = %+v", ratings.IMDb)
	}
	if ratings.Kinopoisk == nil || ratings.Kinopoisk.Value != 8.7 || ratings.Kinopoisk.Votes != nil {
		t.Errorf("kinopoisk = %+v, want a value and no votes", ratings.Kinopoisk)
	}
	if ratings.Metacritic == nil || ratings.Metacritic.Value != 67 {
		t.Errorf("metacritic = %+v", ratings.Metacritic)
	}
	if ratings.RottenTomatoes == nil || ratings.RottenTomatoes.Value != 81 {
		t.Errorf("rotten tomatoes = %+v", ratings.RottenTomatoes)
	}
	if ratings.TMDB == nil || ratings.TMDB.Value != 8.437 || ratings.TMDB.Votes == nil || *ratings.TMDB.Votes != 33039 {
		t.Errorf("tmdb = %+v", ratings.TMDB)
	}
}

func TestRatingsWithAMissingProviderDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, ratingsTVBody))
	ratings, err := c.TV.Ratings(context.Background(), 1396, nil)
	if err != nil {
		t.Fatalf("TV.Ratings: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/tv/1396/ratings" {
		t.Errorf("path = %q", got)
	}
	if ratings.Metacritic != nil {
		t.Errorf("metacritic = %+v, want nil", ratings.Metacritic)
	}
	if ratings.IMDb == nil || ratings.IMDb.Value != 9.5 {
		t.Errorf("imdb = %+v", ratings.IMDb)
	}
	if ratings.TMDB == nil || ratings.TMDB.Value != 9 {
		t.Errorf("tmdb = %+v", ratings.TMDB)
	}
}

func TestRatingsWithoutScoresDecoding(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, ratingsEmptyBody))
	ratings, err := c.Movies.Ratings(context.Background(), 1798888, nil)
	if err != nil {
		t.Fatalf("Movies.Ratings: %v", err)
	}
	if ratings.IMDb != nil || ratings.Kinopoisk != nil || ratings.Metacritic != nil || ratings.RottenTomatoes != nil || ratings.TMDB != nil {
		t.Errorf("ratings = %+v, want every provider nil", ratings)
	}
}

func TestWikidataLinksAndMoneyDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, wikidataSequelBody))
	wd, err := c.Movies.Wikidata(context.Background(), 604, nil)
	if err != nil {
		t.Fatalf("Movies.Wikidata: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/movie/604/wikidata" {
		t.Errorf("path = %q", got)
	}
	if wd.ID != 604 || wd.WikidataID != "Q189600" {
		t.Errorf("wikidata = %+v", wd)
	}
	if wd.Budget == nil || wd.Budget.Value != 150000000 || wd.Budget.Currency != "USD" {
		t.Errorf("budget = %+v", wd.Budget)
	}
	if wd.BoxOffice == nil || wd.BoxOffice.Value != 739412035 || !wd.BoxOffice.Worldwide {
		t.Errorf("box office = %+v", wd.BoxOffice)
	}
	if wd.Follows == nil || wd.Follows.TMDBID != 603 || wd.Follows.MediaType != "movie" || wd.Follows.Label != "The Matrix" {
		t.Errorf("follows = %+v", wd.Follows)
	}
	if wd.FollowedBy == nil || wd.FollowedBy.TMDBID != 605 {
		t.Errorf("followed by = %+v", wd.FollowedBy)
	}
	if wd.Sitelinks["en"] != "The Matrix Reloaded" {
		t.Errorf("sitelinks = %v", wd.Sitelinks)
	}
	if len(wd.FilmingLocations) != 1 || wd.FilmingLocations[0].WikidataID != "Q3130" || wd.FilmingLocations[0].TMDBID != 0 {
		t.Errorf("filming locations = %+v", wd.FilmingLocations)
	}
	if len(wd.Awards) != 0 || len(wd.BasedOn) != 0 || len(wd.OriginalBroadcasters) != 0 {
		t.Errorf("wikidata = %+v, want empty groups", wd)
	}
}

func TestWikidataAwardsAndSourcesDecoding(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, wikidataAwardsBody))
	wd, err := c.Movies.Wikidata(context.Background(), 550, nil)
	if err != nil {
		t.Fatalf("Movies.Wikidata: %v", err)
	}
	if wd.Follows != nil || wd.FollowedBy != nil {
		t.Errorf("follows = %+v, followed by = %+v, want nil", wd.Follows, wd.FollowedBy)
	}
	if len(wd.Awards) != 2 {
		t.Fatalf("awards = %+v", wd.Awards)
	}
	if wd.Awards[0].Year != 0 || wd.Awards[1].Year != 2000 || wd.Awards[1].WikidataID != "Q943908" {
		t.Errorf("awards = %+v", wd.Awards)
	}
	if len(wd.Nominations) != 1 || wd.Nominations[0].Label != "Academy Award for Best Sound Editing" {
		t.Errorf("nominations = %+v", wd.Nominations)
	}
	if len(wd.BasedOn) != 1 || wd.BasedOn[0].Label != "Fight Club" {
		t.Fatalf("based on = %+v", wd.BasedOn)
	}
	if author := wd.BasedOn[0].Author; author == nil || author.Label != "Chuck Palahniuk" {
		t.Errorf("author = %+v", author)
	}
}

func TestWikidataSourceWithoutAnAuthorDecoding(t *testing.T) {
	c, _ := newTestClient(t, jsonHandler(200, wikidataNoAuthorBody))
	wd, err := c.Movies.Wikidata(context.Background(), 155, nil)
	if err != nil {
		t.Fatalf("Movies.Wikidata: %v", err)
	}
	if len(wd.BasedOn) != 1 || wd.BasedOn[0].Label != "Batman" || wd.BasedOn[0].Author != nil {
		t.Errorf("based on = %+v, want one source without an author", wd.BasedOn)
	}
}

func TestShowtimesWithoutACountryListCountries(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, showtimesCountriesBody))
	showtimes, err := c.Movies.Showtimes(context.Background(), 550, "", nil)
	if err != nil {
		t.Fatalf("Movies.Showtimes: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/movie/550/showtimes" {
		t.Errorf("path = %q", got.Path)
	}
	if _, sent := got.Query["country"]; sent {
		t.Errorf("query = %v, an empty country must not be sent", got.Query)
	}
	if showtimes.ID != 550 || showtimes.Country != "" || len(showtimes.Results) != 0 {
		t.Errorf("showtimes = %+v", showtimes)
	}
	if len(showtimes.Countries) != 1 {
		t.Fatalf("countries = %+v", showtimes.Countries)
	}
	ru := showtimes.Countries[0]
	if ru.Country != "RU" || ru.Name != "Russia" || ru.Venues != 2 || ru.Cities != 1 || ru.Showtimes != 7 {
		t.Errorf("country = %+v", ru)
	}
	if ru.FirstDay != "2026-10-18" || ru.LastDay != "2026-10-21" {
		t.Errorf("days = %q to %q", ru.FirstDay, ru.LastDay)
	}
}

func TestShowtimesWithACountryReturnVenues(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, showtimesRussiaBody))
	showtimes, err := c.Movies.Showtimes(context.Background(), 550, "RU", Params{"language": "ru-RU"})
	if err != nil {
		t.Fatalf("Movies.Showtimes: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/movie/550/showtimes" {
		t.Errorf("path = %q", got.Path)
	}
	if got.Query.Get("country") != "RU" || got.Query.Get("language") != "ru-RU" {
		t.Errorf("query = %v", got.Query)
	}
	if showtimes.Country != "RU" || len(showtimes.Results) != 1 || len(showtimes.Countries) != 1 {
		t.Fatalf("showtimes = %+v", showtimes)
	}
	venue := showtimes.Results[0]
	if venue.ID != 1114 || venue.Chain.Key != "mirage-ru" || venue.Country != "RU" || venue.Timezone != "Europe/Moscow" {
		t.Errorf("venue = %+v", venue)
	}
	if venue.Latitude != nil || venue.Longitude != nil || venue.Region != "" {
		t.Errorf("venue = %+v, want no coordinates and no region", venue)
	}
	if len(venue.Showtimes) != 1 {
		t.Fatalf("showtimes = %+v", venue.Showtimes)
	}
	show := venue.Showtimes[0]
	if show.Date != "2026-10-18" || show.Start != "2026-10-18T21:40" || show.End != "" {
		t.Errorf("showtime = %+v", show)
	}
	if len(show.Formats) != 1 || show.Formats[0] != "2D" || len(show.Tags) != 0 {
		t.Errorf("formats = %v, tags = %v", show.Formats, show.Tags)
	}
	if show.Price != nil || show.Prices != nil || show.SoldOut {
		t.Errorf("showtime = %+v", show)
	}
}

func TestShowtimesTakeTheRegionAliasInParams(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, showtimesRussiaBody))
	if _, err := c.Movies.Showtimes(context.Background(), 550, "", Params{"region": "RU"}); err != nil {
		t.Fatalf("Movies.Showtimes: %v", err)
	}
	got := rec.last(t)
	if got.Query.Get("region") != "RU" {
		t.Errorf("query = %v", got.Query)
	}
	if _, sent := got.Query["country"]; sent {
		t.Errorf("query = %v, no country was given", got.Query)
	}
}

func TestShowtimesOfATVShowWithoutDataDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, showtimesEmptyTVBody))
	showtimes, err := c.TV.Showtimes(context.Background(), 1396, "", nil)
	if err != nil {
		t.Fatalf("TV.Showtimes: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/tv/1396/showtimes" {
		t.Errorf("path = %q", got)
	}
	if showtimes.ID != 1396 || len(showtimes.Results) != 0 || len(showtimes.Countries) != 0 {
		t.Errorf("showtimes = %+v", showtimes)
	}
}

func TestAppendedExtensionSectionsDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, appendedMovieBody))
	params := Params{"append_to_response": []string{"lists", "moods", "ratings", "soundtracks", "external_ids", "showtimes"}}
	movie, err := c.Movies.Details(context.Background(), 550, params)
	if err != nil {
		t.Fatalf("Movies.Details: %v", err)
	}
	if got := rec.last(t).Query.Get("append_to_response"); got != "lists,moods,ratings,soundtracks,external_ids,showtimes" {
		t.Errorf("append_to_response = %q", got)
	}
	if movie.ID != 550 || len(movie.OriginCountry) != 2 || movie.OriginCountry[0] != "US" {
		t.Errorf("movie = %+v", movie)
	}
	if movie.Moods == nil || movie.Moods.ID != 0 || movie.Moods.Source != MoodSourceEstimated || len(movie.Moods.Results) != 2 {
		t.Errorf("moods = %+v, the appended form has no id", movie.Moods)
	}
	if movie.Ratings == nil || movie.Ratings.ID != 0 || movie.Ratings.IMDb == nil || movie.Ratings.IMDb.Value != 8.8 {
		t.Errorf("ratings = %+v", movie.Ratings)
	}
	if movie.Soundtracks == nil || movie.Soundtracks.ID != 0 || len(movie.Soundtracks.Results) != 0 {
		t.Errorf("soundtracks = %+v", movie.Soundtracks)
	}
	if movie.ExternalIDs == nil || movie.ExternalIDs.KinopoiskID != "361" || movie.ExternalIDs.IMDbID != "tt0137523" {
		t.Errorf("external ids = %+v", movie.ExternalIDs)
	}
	if movie.Showtimes == nil || movie.Showtimes.ID != 550 || len(movie.Showtimes.Countries) != 1 {
		t.Errorf("showtimes = %+v, the appended form keeps the id", movie.Showtimes)
	}
	if movie.Lists == nil || len(movie.Lists.Results) != 2 || movie.Lists.TotalResults != 2924 {
		t.Errorf("lists = %+v", movie.Lists)
	}
}

func TestExternalIDsDecoding(t *testing.T) {
	cm, recm := newTestClient(t, jsonHandler(200, externalIDsMovieBody))
	movie, err := cm.Movies.ExternalIDs(context.Background(), 550, nil)
	if err != nil {
		t.Fatalf("Movies.ExternalIDs: %v", err)
	}
	if got := recm.last(t).Path; got != "/3/movie/550/external_ids" {
		t.Errorf("path = %q", got)
	}
	if movie.ID != 550 || movie.IMDbID != "tt0137523" || movie.WikidataID != "Q190050" || movie.KinopoiskID != "361" {
		t.Errorf("movie ids = %+v", movie)
	}
	if movie.FacebookID != "FightClub" || movie.InstagramID != "" || movie.TwitterID != "" {
		t.Errorf("movie ids = %+v", movie)
	}

	ct, rect := newTestClient(t, jsonHandler(200, externalIDsTVBody))
	show, err := ct.TV.ExternalIDs(context.Background(), 1396, nil)
	if err != nil {
		t.Fatalf("TV.ExternalIDs: %v", err)
	}
	if got := rect.last(t).Path; got != "/3/tv/1396/external_ids" {
		t.Errorf("path = %q", got)
	}
	if show.TVDBID != 81189 || show.KinopoiskID != "404900" || show.IMDbID != "tt0903747" || show.FreebaseMID != "" {
		t.Errorf("tv ids = %+v", show)
	}

	cp, recp := newTestClient(t, jsonHandler(200, externalIDsPersonBody))
	person, err := cp.People.ExternalIDs(context.Background(), 287, nil)
	if err != nil {
		t.Fatalf("People.ExternalIDs: %v", err)
	}
	if got := recp.last(t).Path; got != "/3/person/287/external_ids" {
		t.Errorf("path = %q", got)
	}
	if person.ID != 287 || person.IMDbID != "nm0000093" || person.WikidataID != "Q35332" || person.TikTokID != "" {
		t.Errorf("person ids = %+v", person)
	}
}

func TestFindByKinopoiskIDDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, findKinopoiskBody))
	res, err := c.Find(context.Background(), "361", "kinopoisk_id")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/find/361" || got.Query.Get("external_source") != "kinopoisk_id" {
		t.Errorf("request = %s ? %v", got.Path, got.Query)
	}
	if len(res.MovieResults) != 1 || res.MovieResults[0].ID != 550 || res.MovieResults[0].MediaType != "movie" {
		t.Fatalf("movie results = %+v", res.MovieResults)
	}
	if len(res.PersonResults) != 0 || len(res.TVResults) != 0 || len(res.TVEpisodeResults) != 0 || len(res.TVSeasonResults) != 0 {
		t.Errorf("results = %+v", res)
	}
}

func TestFindPersonByWikidataIDDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, findPersonBody))
	res, err := c.Find(context.Background(), "Q35332", "wikidata_id")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/find/Q35332" || got.Query.Get("external_source") != "wikidata_id" {
		t.Errorf("request = %s ? %v", got.Path, got.Query)
	}
	if len(res.PersonResults) != 1 {
		t.Fatalf("person results = %+v", res.PersonResults)
	}
	person := res.PersonResults[0]
	if person.ID != 287 || person.Name != "Brad Pitt" || person.MediaType != "person" || len(person.KnownFor) != 0 {
		t.Errorf("person = %+v", person)
	}
}

func TestListsDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, listsMovieBody))
	res, err := c.Movies.Lists(context.Background(), 550, Params{"page": 2})
	if err != nil {
		t.Fatalf("Movies.Lists: %v", err)
	}
	got := rec.last(t)
	if got.Path != "/3/movie/550/lists" || got.Query.Get("page") != "2" {
		t.Errorf("request = %s ? %v", got.Path, got.Query)
	}
	if res.Page != 1 || res.TotalPages != 147 || res.TotalResults != 2924 || len(res.Results) != 1 {
		t.Fatalf("lists = %+v", res)
	}
	list := res.Results[0]
	if list.ID != 8302254 || list.Name != "Movies That I have watched all" || list.ItemCount != 2585 {
		t.Errorf("list = %+v", list)
	}
	if list.ISO31661 != "US" || list.ISO6391 != "en" || list.PosterPath != "" {
		t.Errorf("list = %+v", list)
	}
}

func TestLatestMovieDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, latestMovieBody))
	movie, err := c.Movies.Latest(context.Background(), nil)
	if err != nil {
		t.Fatalf("Movies.Latest: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/movie/latest" {
		t.Errorf("path = %q", got)
	}
	if movie.ID != 1798888 || movie.Title != "In the Mood to Leave" || movie.IMDbID != "tt49276703" || movie.Runtime != 5 {
		t.Errorf("movie = %+v", movie)
	}
	if len(movie.OriginCountry) != 1 || movie.OriginCountry[0] != "US" {
		t.Errorf("origin country = %v", movie.OriginCountry)
	}
}

func TestLatestTVShowDecoding(t *testing.T) {
	c, rec := newTestClient(t, jsonHandler(200, latestTVBody))
	show, err := c.TV.Latest(context.Background(), nil)
	if err != nil {
		t.Fatalf("TV.Latest: %v", err)
	}
	if got := rec.last(t).Path; got != "/3/tv/latest" {
		t.Errorf("path = %q", got)
	}
	if show.ID != 338695 || show.NumberOfSeasons != 1 || show.NumberOfEpisodes != 5 || show.LastEpisodeToAir != nil {
		t.Errorf("show = %+v", show)
	}
	if len(show.Languages) != 1 || show.Languages[0] != "ru" {
		t.Errorf("languages = %v", show.Languages)
	}
}

func TestExtensionEndpointErrors(t *testing.T) {
	t.Run("unknown title", func(t *testing.T) {
		c, _ := newTestClient(t, jsonHandler(http.StatusNotFound, notFoundBody))
		_, err := c.TV.Wikidata(context.Background(), 999999999, nil)
		if !IsNotFound(err) {
			t.Fatalf("err = %v, want a not found error", err)
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Fatalf("err = %v, want *Error", err)
		}
		if apiErr.Status != 404 || apiErr.StatusCode != 34 || apiErr.StatusMessage != "The resource you requested could not be found." {
			t.Errorf("error = %+v", apiErr)
		}
	})
	t.Run("method not offered", func(t *testing.T) {
		c, _ := newTestClient(t, jsonHandler(http.StatusMethodNotAllowed, notOfferedBody))
		err := c.Do(context.Background(), "/account", nil, nil)
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Fatalf("err = %v, want *Error", err)
		}
		if apiErr.Status != 405 || apiErr.StatusCode != 37 || apiErr.StatusMessage != "This method is not offered by amodb." {
			t.Errorf("error = %+v", apiErr)
		}
	})
}
