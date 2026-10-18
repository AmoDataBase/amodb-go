package amodb

import (
	"context"
	"encoding/json"
	"strings"
)

// OMDbService calls the OMDb compatible endpoint. Use it as client.OMDb.
//
// The key is always sent as the apikey query parameter. Params use the OMDb
// names: plot (short or full), type (movie, series or episode), y (the year)
// and page. The name year is accepted as an alias of y.
//
// An answer with "Response":"False" is returned as an *Error with Status 200
// and the OMDb text in StatusMessage.
type OMDbService service

// ByID looks up a title by IMDb id, for example "tt0137523".
func (s *OMDbService) ByID(ctx context.Context, imdbID string, params Params) (*OMDbTitle, error) {
	return omdbGet[OMDbTitle](ctx, s.client, omdbParams(params, "i", imdbID))
}

// ByTitle looks up a title by its exact name.
func (s *OMDbService) ByTitle(ctx context.Context, title string, params Params) (*OMDbTitle, error) {
	return omdbGet[OMDbTitle](ctx, s.client, omdbParams(params, "t", title))
}

// Search searches titles by words. Use the page parameter to walk through
// the results.
func (s *OMDbService) Search(ctx context.Context, query string, params Params) (*OMDbSearch, error) {
	return omdbGet[OMDbSearch](ctx, s.client, omdbParams(params, "s", query))
}

// omdbParams sets key and maps the alias year to y.
func omdbParams(params Params, key string, value any) Params {
	out := withParam(params, key, value)
	if year, ok := out["year"]; ok {
		if _, has := out["y"]; !has {
			out["y"] = year
		}
		delete(out, "year")
	}
	return out
}

// omdbGet sends an OMDb request and decodes the answer into T.
func omdbGet[T any](ctx context.Context, c *Client, params Params) (*T, error) {
	res, err := c.fetch(ctx, c.omdbURL, params, scopeOMDb)
	if err != nil {
		return nil, err
	}
	if oerr := checkOMDb(res); oerr != nil {
		return nil, oerr
	}
	var out T
	if derr := decodeJSON(res, &out); derr != nil {
		return nil, derr
	}
	return &out, nil
}

// checkOMDb turns a Response False body into an *Error.
func checkOMDb(res *response) error {
	var probe struct {
		Response string `json:"Response"`
		Error    string `json:"Error"`
	}
	if json.Unmarshal(res.body, &probe) != nil {
		return nil
	}
	if !strings.EqualFold(probe.Response, "false") {
		return nil
	}
	e := &Error{Status: res.status, Body: res.body, StatusMessage: probe.Error}
	if e.StatusMessage == "" {
		e.StatusMessage = "OMDb reported an error"
	}
	return e
}
