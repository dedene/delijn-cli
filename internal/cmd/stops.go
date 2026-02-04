package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/dedene/delijn-cli/internal/api"
)

type StopsCmd struct {
	Search StopsSearchCmd `cmd:"" help:"Search stops by name"`
	Get    StopsGetCmd    `cmd:"" help:"Get stop details by number"`
}

type StopsSearchCmd struct {
	Query string `arg:"" required:"" help:"Search query (stop name)"`
}

func (c *StopsSearchCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	resp, err := client.SearchStops(ctx, c.Query)
	if err != nil {
		return fmt.Errorf("search stops: %w", err)
	}

	if root.JSON {
		return outputJSON(resp.Stops)
	}

	if root.Plain {
		outputStopsPlain(resp.Stops)

		return nil
	}

	outputStopsTable(resp.Stops)

	return nil
}

type StopsGetCmd struct {
	Number string `arg:"" required:"" help:"6-digit stop number (e.g., 200552)"`
}

func (c *StopsGetCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stopNumber, err := strconv.Atoi(c.Number)
	if err != nil {
		return fmt.Errorf("invalid stop number %q: must be a 6-digit number", c.Number)
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	stop, err := client.GetStopByNumber(ctx, stopNumber)
	if err != nil {
		return fmt.Errorf("get stop: %w", err)
	}

	if root.JSON {
		return outputJSON(stop)
	}

	if root.Plain {
		outputStopPlain(stop)

		return nil
	}

	outputStopDetails(stop)

	return nil
}

func outputJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(v)
}

func outputStopsTable(stops []api.Stop) {
	if len(stops) == 0 {
		fmt.Fprintln(os.Stdout, "No stops found.")

		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "NUMBER\tNAME\tMUNICIPALITY")

	for _, s := range stops {
		fmt.Fprintf(w, "%d\t%s\t%s\n", s.Number, s.Description, s.Municipality)
	}
}

func outputStopsPlain(stops []api.Stop) {
	for _, s := range stops {
		fmt.Fprintf(os.Stdout, "%d\t%s\t%s\n", s.Number, s.Description, s.Municipality)
	}
}

func outputStopPlain(stop *api.Stop) {
	fmt.Fprintf(os.Stdout, "%d\t%s\t%s\n", stop.Number, stop.Description, stop.Municipality)
}

func outputStopDetails(stop *api.Stop) {
	fmt.Fprintf(os.Stdout, "Stop:         %s\n", stop.Description)
	fmt.Fprintf(os.Stdout, "Number:       %d\n", stop.Number)
	fmt.Fprintf(os.Stdout, "Municipality: %s\n", stop.Municipality)
	fmt.Fprintf(os.Stdout, "Entity:       %d\n", stop.EntityNumber)

	if stop.GeoCoordinate != nil {
		fmt.Fprintf(os.Stdout, "Location:     %.6f, %.6f\n",
			stop.GeoCoordinate.Latitude, stop.GeoCoordinate.Longitude)
	}
}
