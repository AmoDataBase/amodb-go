package amodb_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/netlistenany-bot/amodb-go"
)

func ExampleNewClient() {
	client, err := amodb.NewClient(
		os.Getenv("AMODB_API_KEY"),
		amodb.WithLanguage("ru-RU"),
		amodb.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	movie, err := client.Movies.Details(context.Background(), 550, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(movie.Title, movie.ReleaseDate)
}

func ExampleMoviesService_Details() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	params := amodb.Params{"append_to_response": []string{"credits", "images"}}
	movie, err := client.Movies.Details(context.Background(), 550, params)
	if err != nil {
		log.Fatal(err)
	}
	for _, member := range movie.Credits.Cast {
		fmt.Println(member.Name, member.Character)
	}
}

func ExampleSearchService_Movie() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	page, err := client.Search.Movie(context.Background(), "inception", amodb.Params{"year": 2010})
	if err != nil {
		log.Fatal(err)
	}
	for _, movie := range page.Results {
		fmt.Println(movie.ID, movie.Title)
	}
}

func ExampleTVService_Season() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	season, err := client.TV.Season(context.Background(), 1399, 1, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, episode := range season.Episodes {
		fmt.Println(episode.EpisodeNumber, episode.Name)
	}
}

func ExampleTrendingService_Get() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	page, err := client.Trending.Get(context.Background(), "movie", "week", nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range page.Results {
		fmt.Println(item["title"])
	}
}

func ExampleClient_Find() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	found, err := client.Find(context.Background(), "tt0137523", "imdb_id")
	if err != nil {
		log.Fatal(err)
	}
	for _, movie := range found.MovieResults {
		fmt.Println(movie.ID, movie.Title)
	}
}

func ExampleClient_ImageURL() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	movie, err := client.Movies.Details(context.Background(), 550, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client.ImageURL(movie.PosterPath, "w500"))
}

func ExampleOMDbService_ByID() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	title, err := client.OMDb.ByID(context.Background(), "tt0137523", amodb.Params{"plot": "full"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(title.Title, title.Year, title.IMDbRating)
}

func ExampleClient_Plans() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	catalog, err := client.Plans(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, plan := range catalog.Plans {
		fmt.Println(plan.ID, plan.RPM)
	}
}

func ExampleClient_Do() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	var out struct {
		TotalResults int `json:"total_results"`
	}
	err = client.Do(context.Background(), "/movie/popular", amodb.Params{"page": 2}, &out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.TotalResults)
}

func ExampleError() {
	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Movies.Details(context.Background(), 1, nil)

	var apiErr *amodb.Error
	switch {
	case err == nil:
		fmt.Println("found")
	case amodb.IsNotFound(err):
		fmt.Println("no such movie")
	case errors.As(err, &apiErr):
		fmt.Println(apiErr.Status, apiErr.StatusCode, apiErr.StatusMessage)
	case errors.Is(err, amodb.ErrTimeout):
		fmt.Println("timed out")
	default:
		fmt.Println(err)
	}
}

func ExampleIsRateLimited() {
	client, err := amodb.NewClient(
		os.Getenv("AMODB_API_KEY"),
		amodb.WithMaxRetries(1),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Movies.Popular(context.Background(), nil)
	if amodb.IsRateLimited(err) {
		fmt.Println("still over the limit after one retry")
	}
}
