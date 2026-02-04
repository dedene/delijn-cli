package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RetryTransport wraps an http.RoundTripper with retry logic.
type RetryTransport struct {
	base       http.RoundTripper
	maxRetries int
	backoff    time.Duration
}

// NewRetryTransport creates a new RetryTransport wrapping the given transport.
func NewRetryTransport(base http.RoundTripper) *RetryTransport {
	return &RetryTransport{
		base:       base,
		maxRetries: 3,
		backoff:    500 * time.Millisecond,
	}
}

// RoundTrip implements http.RoundTripper with retry logic.
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Make body replayable for retries
	if err := ensureReplayableBody(req); err != nil {
		return nil, err
	}

	var lastErr error

	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(t.backoff * time.Duration(attempt))

			// Reset body for retry
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, fmt.Errorf("get request body: %w", err)
				}

				req.Body = body
			}
		}

		resp, err := t.base.RoundTrip(req)
		if err != nil {
			lastErr = err

			continue
		}

		// Don't retry on success or client errors
		if resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// Retry on 5xx or 429
		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			drainAndClose(resp.Body)

			if attempt < t.maxRetries {
				continue
			}

			return resp, nil
		}

		return resp, nil
	}

	return nil, lastErr
}

// ensureReplayableBody ensures the request body can be read multiple times.
func ensureReplayableBody(req *http.Request) error {
	if req.Body == nil || req.GetBody != nil {
		return nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}

	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}

	return nil
}

// drainAndClose reads and closes a response body.
func drainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}
