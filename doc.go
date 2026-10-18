// Package amodb is a client for the amodb movie and TV metadata API.
//
// The API is compatible with TMDB v3 and OMDb. amodb is a paid API, so you
// need a key from https://amodb.dev/pricing. Create a client with NewClient
// and call the resource services on it:
//
//	client, err := amodb.NewClient(os.Getenv("AMODB_API_KEY"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	movie, err := client.Movies.Details(ctx, 550, nil)
//
// Every network method takes a context.Context first and an optional Params
// map last. Params values are sent as query parameters: booleans become true
// or false, slices are joined with commas and nil values are dropped.
//
// Endpoints with a stable, commonly used shape return typed structs such as
// MovieDetails or Paged[Movie]. The others return Object, a plain map decoded
// from the JSON body. Client.Do is the escape hatch for any other path.
//
// Failed calls return *Error. Use errors.As to read the status fields, or the
// helpers IsRateLimited, IsNotFound and IsUnauthorized. Requests that get an
// HTTP 429 answer are retried up to three times, waiting for the Retry-After
// header when it is present.
//
// A Client is safe for concurrent use by multiple goroutines.
package amodb
