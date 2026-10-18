package amodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Params holds query parameters for a call. A nil Params is allowed.
//
// Booleans are sent as true or false, slices are joined with commas, nil
// values are dropped and time.Time values are sent as YYYY-MM-DD. Keys are
// sent exactly as given, in the snake_case form of the API.
type Params map[string]any

const maxBodyBytes = 64 << 20

const maxRetryWait = 60 * time.Second

const maxParsedRetryAfter = 24 * 60 * 60

type authScope int

const (
	scopeTMDB authScope = iota
	scopeOMDb
	scopePublic
)

type response struct {
	status int
	header http.Header
	body   []byte
}

// Do sends a GET request to the base URL plus path and decodes the JSON
// answer into out. The path starts with a slash, for example "/movie/550".
// It applies the same authentication, language default, retries and error
// handling as the resource methods. A nil out discards the body.
func (c *Client) Do(ctx context.Context, path string, params Params, out any) error {
	target := c.baseURL + apiPath(path)
	res, err := c.fetch(ctx, target, params, scopeTMDB)
	if err != nil {
		return err
	}
	return decodeJSON(res, out)
}

// get calls Do and returns the decoded value as a pointer to T.
func get[T any](ctx context.Context, c *Client, path string, params Params) (*T, error) {
	var out T
	if err := c.Do(ctx, path, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// getObject calls Do and returns the decoded JSON object.
func getObject(ctx context.Context, c *Client, path string, params Params) (Object, error) {
	var out Object
	if err := c.Do(ctx, path, params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// withParam returns a copy of params with key set to value.
func withParam(params Params, key string, value any) Params {
	out := make(Params, len(params)+1)
	for k, v := range params {
		out[k] = v
	}
	out[key] = value
	return out
}

// withCountry returns params with the country set. An empty country sends no
// country parameter, and params is returned as it is.
func withCountry(params Params, country string) Params {
	if country == "" {
		return params
	}
	return withParam(params, "country", country)
}

func apiPath(path string) string {
	if path == "" || path[0] != '/' {
		return "/" + path
	}
	return path
}

// buildPath joins path segments. Strings are escaped with url.PathEscape and
// integers are written as decimal numbers.
func buildPath(segments ...any) string {
	var b strings.Builder
	for _, segment := range segments {
		b.WriteByte('/')
		switch v := segment.(type) {
		case int:
			b.WriteString(strconv.Itoa(v))
		case string:
			b.WriteString(url.PathEscape(v))
		default:
			b.WriteString(url.PathEscape(fmt.Sprint(v)))
		}
	}
	return b.String()
}

// buildURL appends the query string, including the default language and the
// key when the client is in query mode.
func (c *Client) buildURL(rawURL string, params Params, scope authScope) string {
	values := url.Values{}
	for key, value := range params {
		if s, ok := formatParam(value); ok {
			values.Set(key, s)
		}
	}
	switch scope {
	case scopeTMDB:
		if c.language != "" && values.Get("language") == "" {
			values.Set("language", c.language)
		}
		if c.auth == AuthQuery {
			values.Set("api_key", c.apiKey)
		}
	case scopeOMDb:
		values.Set("apikey", c.apiKey)
	}
	if len(values) == 0 {
		return rawURL
	}
	sep := "?"
	if strings.Contains(rawURL, "?") {
		sep = "&"
	}
	return rawURL + sep + values.Encode()
}

// formatParam turns a parameter value into its query string form. The second
// result is false when the value must be dropped.
func formatParam(v any) (string, bool) {
	if v == nil {
		return "", false
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return "", false
	}
	switch x := v.(type) {
	case string:
		return x, true
	case bool:
		return strconv.FormatBool(x), true
	case int:
		return strconv.Itoa(x), true
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64), true
	case time.Time:
		return x.Format("2006-01-02"), true
	case fmt.Stringer:
		return x.String(), true
	}
	switch rv.Kind() {
	case reflect.Ptr:
		return formatParam(rv.Elem().Interface())
	case reflect.Slice, reflect.Array:
		parts := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			if s, ok := formatParam(rv.Index(i).Interface()); ok {
				parts = append(parts, s)
			}
		}
		if len(parts) == 0 {
			return "", false
		}
		return strings.Join(parts, ","), true
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), true
	case reflect.String:
		return rv.String(), true
	}
	return fmt.Sprint(v), true
}

// fetch sends the request and retries HTTP 429 answers. It returns the
// response for 2xx statuses and an *Error for everything else.
func (c *Client) fetch(ctx context.Context, rawURL string, params Params, scope authScope) (*response, error) {
	target := c.buildURL(rawURL, params, scope)
	for attempt := 0; ; attempt++ {
		res, err := c.send(ctx, target, scope)
		if err != nil {
			return nil, err
		}
		if res.status >= 200 && res.status < 300 {
			return res, nil
		}
		now := c.now()
		failure := newHTTPError(res, now)
		if res.status != http.StatusTooManyRequests || attempt >= c.maxRetries {
			return nil, failure
		}
		wait := retryDelay(res.header, attempt, now)
		if serr := c.sleep(ctx, wait); serr != nil {
			return nil, contextError(serr)
		}
	}
}

// send performs one HTTP attempt. Each attempt gets its own timeout.
func (c *Client) send(ctx context.Context, target string, scope authScope) (*response, error) {
	attemptCtx, cancel := c.attemptContext(ctx)
	defer cancel()

	req, err := http.NewRequestWithContext(attemptCtx, http.MethodGet, target, nil)
	if err != nil {
		return nil, fmt.Errorf("amodb: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if scope == scopeTMDB && c.auth == AuthBearer {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, transportError(ctx, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return nil, transportError(ctx, err)
	}
	return &response{status: resp.StatusCode, header: resp.Header, body: body}, nil
}

// attemptContext applies the client timeout. context.WithTimeout keeps the
// caller deadline when it is earlier than the client timeout.
func (c *Client) attemptContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.timeout <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, c.timeout)
}

// transportError maps a failure from the HTTP client. A canceled caller
// context is returned unchanged so errors.Is(err, context.Canceled) works.
func transportError(ctx context.Context, err error) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}
	if errors.Is(err, context.DeadlineExceeded) || isTimeout(err) {
		return timeoutError(err)
	}
	e := &Error{StatusMessage: "network error: " + err.Error()}
	e.cause = fmt.Errorf("%w: %w", ErrNetwork, err)
	return e
}

// contextError maps an error from a wait between retries.
func contextError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return timeoutError(err)
	}
	return err
}

func timeoutError(err error) *Error {
	e := &Error{StatusMessage: "request timed out"}
	e.cause = fmt.Errorf("%w: %w", ErrTimeout, err)
	return e
}

func isTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// decodeJSON decodes a successful response into out.
func decodeJSON(res *response, out any) error {
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(res.body, out); err != nil {
		e := &Error{Status: res.status, Body: res.body, cause: err}
		e.StatusMessage = "invalid JSON in response: " + err.Error()
		return e
	}
	return nil
}

// retryDelay returns how long to wait before the next attempt. It honours
// Retry-After (seconds or an HTTP date) clamped to 0..60 seconds. Without the
// header it waits 2^attempt seconds, so 1, 2 and 4 seconds.
func retryDelay(h http.Header, attempt int, now time.Time) time.Duration {
	wait, ok := parseRetryAfter(h.Get("Retry-After"), now)
	if !ok {
		shift := attempt
		if shift < 0 {
			shift = 0
		}
		if shift > 6 {
			shift = 6
		}
		secs := 1 << uint(shift)
		wait = time.Duration(secs) * time.Second
	}
	if wait < 0 {
		wait = 0
	}
	if wait > maxRetryWait {
		wait = maxRetryWait
	}
	return wait
}

// parseRetryAfter reads a Retry-After value. The second result is false when
// the header is missing or not understood. Values in the past give 0.
func parseRetryAfter(value string, now time.Time) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if secs, err := strconv.ParseInt(value, 10, 64); err == nil {
		if secs < 0 {
			secs = 0
		}
		if secs > maxParsedRetryAfter {
			secs = maxParsedRetryAfter
		}
		return time.Duration(secs) * time.Second, true
	}
	if when, err := http.ParseTime(value); err == nil {
		wait := when.Sub(now)
		if wait < 0 {
			wait = 0
		}
		return wait, true
	}
	return 0, false
}

// sleepContext waits for d or until ctx is done.
func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
