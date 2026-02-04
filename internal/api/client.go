package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dedene/delijn-cli/internal/auth"
)

const (
	// BaseURLKern is the base URL for core operations (240 calls/min).
	BaseURLKern = "https://api.delijn.be/DLKernOpenData/v1/beta"

	// BaseURLSearch is the base URL for search operations (6000 calls/min).
	BaseURLSearch = "https://api.delijn.be/DLZoekOpenData/v1/beta"

	// BaseURLGTFS is the base URL for GTFS-RT operations.
	BaseURLGTFS = "https://api.delijn.be/gtfs-realtime/v3"

	// AuthHeader is the header name for the API key.
	AuthHeader = "Ocp-Apim-Subscription-Key"

	// UserAgent identifies this client.
	UserAgent = "delijn-cli/1.0"

	// ContentType for JSON requests.
	ContentType = "application/json"

	// APITimeFormat is the time format used by the API.
	APITimeFormat = "2006-01-02T15:04:05"
)

// Client is the De Lijn API client.
type Client struct {
	httpClient     *http.Client
	kernLimiter    *RateLimiter
	searchLimiter  *RateLimiter
	circuitBreaker *CircuitBreaker
	apiKey         string
}

// NewClient creates a new API client.
func NewClient() (*Client, error) {
	apiKey, err := auth.GetAPIKey()
	if err != nil {
		return nil, &AuthError{Err: err}
	}

	return NewClientWithKey(apiKey), nil
}

// NewClientWithKey creates a new API client with the given API key.
func NewClientWithKey(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: NewRetryTransport(http.DefaultTransport),
			Timeout:   30 * time.Second,
		},
		kernLimiter:    NewRateLimiter(240, time.Minute),  // 240/min
		searchLimiter:  NewRateLimiter(6000, time.Minute), // 6000/min
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
		apiKey:         apiKey,
	}
}

func (c *Client) do(ctx context.Context, baseURL, method, path string, limiter *RateLimiter, body []byte, out interface{}) error {
	if c.circuitBreaker.IsOpen() {
		return &CircuitBreakerError{}
	}

	if err := limiter.Wait(ctx); err != nil {
		return err
	}

	reqURL := baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set(AuthHeader, c.apiKey)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", ContentType)

	if body != nil {
		req.Header.Set("Content-Type", ContentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.circuitBreaker.RecordFailure()

		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		c.circuitBreaker.RecordFailure()

		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "authentication failed",
			Details:    "invalid or expired API key; run 'delijn auth set-key' to configure",
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "not found",
		}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		c.circuitBreaker.RecordFailure()
		retryAfter := 0

		if ra := resp.Header.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.Atoi(ra)
		}

		return &RateLimitError{RetryAfter: retryAfter}
	}

	if resp.StatusCode >= 400 {
		c.circuitBreaker.RecordFailure()
		bodyBytes, _ := io.ReadAll(resp.Body)

		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
			Details:    string(bodyBytes),
		}
	}

	c.circuitBreaker.RecordSuccess()

	if out != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// GetKern performs a GET request to the core API.
func (c *Client) GetKern(ctx context.Context, path string, out interface{}) error {
	return c.do(ctx, BaseURLKern, http.MethodGet, path, c.kernLimiter, nil, out)
}

// GetSearch performs a GET request to the search API.
func (c *Client) GetSearch(ctx context.Context, path string, out interface{}) error {
	return c.do(ctx, BaseURLSearch, http.MethodGet, path, c.searchLimiter, nil, out)
}

// GetGTFS performs a GET request to the GTFS API.
func (c *Client) GetGTFS(ctx context.Context, path string, out interface{}) error {
	return c.do(ctx, BaseURLGTFS, http.MethodGet, path, c.searchLimiter, nil, out)
}

// GetStop retrieves a stop by entity number and stop number.
func (c *Client) GetStop(ctx context.Context, entityNumber, stopNumber int) (*Stop, error) {
	path := fmt.Sprintf("/haltes/%d/%d", entityNumber, stopNumber)

	var stop Stop
	if err := c.GetKern(ctx, path, &stop); err != nil {
		return nil, err
	}

	return &stop, nil
}

// GetStopByNumber retrieves a stop by its 6-digit stop number.
// The first digit is the entity number.
func (c *Client) GetStopByNumber(ctx context.Context, stopNumber int) (*Stop, error) {
	entityNumber := stopNumber / 100000

	return c.GetStop(ctx, entityNumber, stopNumber)
}

// SearchStops searches for stops by name.
func (c *Client) SearchStops(ctx context.Context, query string) (*StopsResponse, error) {
	path := fmt.Sprintf("/haltes/zoek/%s", url.PathEscape(query))

	var resp StopsResponse
	if err := c.GetSearch(ctx, path, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetRealtime retrieves realtime departures for a stop.
func (c *Client) GetRealtime(ctx context.Context, entityNumber, stopNumber int) (*RealtimeResponse, error) {
	path := fmt.Sprintf("/haltes/%d/%d/real-time", entityNumber, stopNumber)

	var resp RealtimeResponse
	if err := c.GetKern(ctx, path, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetRealtimeByNumber retrieves realtime departures for a stop by its 6-digit number.
func (c *Client) GetRealtimeByNumber(ctx context.Context, stopNumber int) (*RealtimeResponse, error) {
	entityNumber := stopNumber / 100000

	return c.GetRealtime(ctx, entityNumber, stopNumber)
}

// GetLine retrieves a line by entity and line number.
func (c *Client) GetLine(ctx context.Context, entityNumber, lineNumber int) (*Line, error) {
	path := fmt.Sprintf("/lijnen/%d/%d", entityNumber, lineNumber)

	var line Line
	if err := c.GetKern(ctx, path, &line); err != nil {
		return nil, err
	}

	return &line, nil
}

// SearchLines searches for lines by number or description.
func (c *Client) SearchLines(ctx context.Context, query string) (*LinesResponse, error) {
	path := fmt.Sprintf("/lijnen/zoek/%s", url.PathEscape(query))

	var resp LinesResponse
	if err := c.GetSearch(ctx, path, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetLineColours retrieves the colours for a line.
func (c *Client) GetLineColours(ctx context.Context, entityNumber, lineNumber int) (*LineColours, error) {
	path := fmt.Sprintf("/lijnen/%d/%d/lijnkleuren", entityNumber, lineNumber)

	var colours LineColours
	if err := c.GetKern(ctx, path, &colours); err != nil {
		return nil, err
	}

	return &colours, nil
}

// ParseAPITime parses a time string from the API.
func ParseAPITime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	brussels, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		brussels = time.Local
	}

	t, err := time.ParseInLocation(APITimeFormat, s, brussels)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time %q: %w", s, err)
	}

	return t, nil
}
