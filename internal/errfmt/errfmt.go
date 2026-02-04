package errfmt

import (
	"errors"
	"fmt"

	"github.com/dedene/delijn-cli/internal/api"
)

// Format returns a user-friendly error message with actionable hints.
func Format(err error) string {
	if err == nil {
		return ""
	}

	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		return formatAPIError(apiErr)
	}

	var authErr *api.AuthError
	if errors.As(err, &authErr) {
		return formatAuthError(authErr)
	}

	var rateLimitErr *api.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return formatRateLimitError(rateLimitErr)
	}

	var cbErr *api.CircuitBreakerError
	if errors.As(err, &cbErr) {
		return formatCircuitBreakerError()
	}

	return err.Error()
}

func formatAPIError(err *api.APIError) string {
	switch err.StatusCode {
	case 401, 403:
		return fmt.Sprintf("Authentication failed: %s\n\nRun 'delijn auth set-key' to configure your API key.\nGet your key from https://data.delijn.be/", err.Message)
	case 404:
		return "Resource not found. Check that the stop/line number is correct."
	case 429:
		return "Rate limit exceeded. Please wait before trying again."
	default:
		if err.Details != "" {
			return fmt.Sprintf("API error (%d): %s - %s", err.StatusCode, err.Message, err.Details)
		}

		return fmt.Sprintf("API error (%d): %s", err.StatusCode, err.Message)
	}
}

func formatAuthError(err *api.AuthError) string {
	if errors.Is(err.Err, api.ErrNotAuthenticated) {
		return "No API key configured.\n\nRun 'delijn auth set-key' to configure your API key.\nGet your key from https://data.delijn.be/"
	}

	return fmt.Sprintf("Authentication error: %v\n\nRun 'delijn auth set-key' to reconfigure your API key.", err.Err)
}

func formatRateLimitError(err *api.RateLimitError) string {
	if err.RetryAfter > 0 {
		return fmt.Sprintf("Rate limit exceeded. Please wait %d seconds before trying again.", err.RetryAfter)
	}

	return "Rate limit exceeded. Please wait before trying again."
}

func formatCircuitBreakerError() string {
	return "Too many consecutive failures. The De Lijn API may be experiencing issues.\nPlease try again later."
}

// ExitCode returns the appropriate exit code for an error.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ExitCode()
	}

	var authErr *api.AuthError
	if errors.As(err, &authErr) {
		return api.ExitAuth
	}

	var rateLimitErr *api.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return api.ExitRateLimit
	}

	return api.ExitError
}
