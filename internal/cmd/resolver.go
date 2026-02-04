package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dedene/delijn-cli/internal/api"
	"github.com/dedene/delijn-cli/internal/config"
)

// ResolveStop resolves a stop reference to a stop number.
// Supports:
//   - @alias - lookup from favorites
//   - numeric - use directly
//   - string - search and pick (single result) or error (multiple results)
func ResolveStop(ctx context.Context, client *api.Client, ref string) (int, error) {
	// @alias - lookup from favorites
	if alias, ok := strings.CutPrefix(ref, "@"); ok {
		stopNum, err := config.GetFavorite(alias)
		if err != nil {
			return 0, fmt.Errorf("resolve favorite %q: %w", alias, err)
		}

		return stopNum, nil
	}

	// Numeric - use directly
	if stopNum, err := strconv.Atoi(ref); err == nil {
		return stopNum, nil
	}

	// String - search
	resp, err := client.SearchStops(ctx, ref)
	if err != nil {
		return 0, fmt.Errorf("search stops: %w", err)
	}

	if len(resp.Stops) == 0 {
		return 0, fmt.Errorf("no stops found matching %q", ref)
	}

	if len(resp.Stops) == 1 {
		return resp.Stops[0].Number, nil
	}

	// Multiple results - for now, return error with options
	// TODO: Interactive picker
	return 0, &AmbiguousStopError{
		Query: ref,
		Stops: resp.Stops,
	}
}

// AmbiguousStopError is returned when multiple stops match a query.
type AmbiguousStopError struct {
	Query string
	Stops []api.Stop
}

func (e *AmbiguousStopError) Error() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Multiple stops match %q:\n", e.Query)

	for i, s := range e.Stops {
		if i >= 5 {
			fmt.Fprintf(&sb, "  ... and %d more\n", len(e.Stops)-5)

			break
		}

		fmt.Fprintf(&sb, "  %d - %s, %s\n", s.Number, s.Description, s.Municipality)
	}

	fmt.Fprintf(&sb, "\nUse the stop number directly, or set a favorite with:\n")
	fmt.Fprintf(&sb, "  delijn config set-favorite <name> %d", e.Stops[0].Number)

	return sb.String()
}
