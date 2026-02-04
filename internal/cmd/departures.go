package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/dedene/delijn-cli/internal/api"
	"github.com/dedene/delijn-cli/internal/output"
)

type DeparturesCmd struct {
	Stop  string `arg:"" required:"" help:"Stop (number, name, or @favorite)"`
	Watch bool   `help:"Auto-refresh every 30 seconds" short:"w"`
	Count int    `help:"Maximum number of departures" default:"10" short:"n"`
	Line  string `help:"Filter by line number" short:"l"`
}

func (c *DeparturesCmd) Run(root *RootFlags) error {
	output.SetNoColor(root.NoColor)

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stopNumber, err := ResolveStop(ctx, client, c.Stop)
	if err != nil {
		return err
	}

	if c.Watch {
		return c.runWatch(client, stopNumber, root)
	}

	return c.runOnce(client, stopNumber, root)
}

func (c *DeparturesCmd) runOnce(client *api.Client, stopNumber int, root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	departures, err := c.fetchDepartures(ctx, client, stopNumber)
	if err != nil {
		return err
	}

	return c.output(departures, root)
}

func (c *DeparturesCmd) runWatch(client *api.Client, stopNumber int, root *RootFlags) error {
	// Handle Ctrl+C gracefully
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// First fetch
	if err := c.fetchAndPrint(ctx, client, stopNumber, root); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	for {
		select {
		case <-sigCh:
			fmt.Fprintln(os.Stdout)

			return nil
		case <-ticker.C:
			// Clear screen and move cursor to top
			fmt.Fprint(os.Stdout, "\033[2J\033[H")

			if err := c.fetchAndPrint(ctx, client, stopNumber, root); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
	}
}

func (c *DeparturesCmd) fetchAndPrint(ctx context.Context, client *api.Client, stopNumber int, root *RootFlags) error {
	fetchCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	departures, err := c.fetchDepartures(fetchCtx, client, stopNumber)
	if err != nil {
		return err
	}

	return c.output(departures, root)
}

func (c *DeparturesCmd) fetchDepartures(ctx context.Context, client *api.Client, stopNumber int) ([]api.Departure, error) {
	resp, err := client.GetRealtimeByNumber(ctx, stopNumber)
	if err != nil {
		return nil, fmt.Errorf("get departures: %w", err)
	}

	var departures []api.Departure

	for _, passage := range resp.StopPassages {
		for _, dep := range passage.Departures {
			// Parse times
			scheduled, err := output.ParseAPITime(dep.ScheduledTimeRaw)
			if err == nil {
				dep.ScheduledTime = scheduled
			}

			if dep.RealTimeRaw != "" {
				realtime, err := output.ParseAPITime(dep.RealTimeRaw)
				if err == nil {
					dep.RealTime = &realtime
				}
			}

			// Filter by line if specified
			if c.Line != "" {
				lineNum, _ := strconv.Atoi(c.Line)
				if dep.LineNumber != lineNum && dep.LinePublicNumber != c.Line {
					continue
				}
			}

			departures = append(departures, dep)

			if len(departures) >= c.Count {
				break
			}
		}

		if len(departures) >= c.Count {
			break
		}
	}

	return departures, nil
}

func (c *DeparturesCmd) output(departures []api.Departure, root *RootFlags) error {
	if root.JSON {
		return outputJSON(departures)
	}

	if root.Plain {
		outputDeparturesPlain(departures)

		return nil
	}

	outputDeparturesTable(departures)

	return nil
}

func outputDeparturesTable(departures []api.Departure) {
	if len(departures) == 0 {
		fmt.Fprintln(os.Stdout, "No departures found.")

		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "TIME\tIN\tLINE\tDESTINATION\tDELAY")

	for _, d := range departures {
		displayTime := d.ScheduledTime
		if d.RealTime != nil {
			displayTime = *d.RealTime
		}

		timeStr := output.FormatTime(displayTime)
		relStr := output.FormatRelative(displayTime)
		lineStr := formatLineNumber(d)
		delayStr := formatDelayStr(d)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			timeStr,
			relStr,
			lineStr,
			d.Destination,
			delayStr,
		)
	}
}

func outputDeparturesPlain(departures []api.Departure) {
	for _, d := range departures {
		displayTime := d.ScheduledTime
		if d.RealTime != nil {
			displayTime = *d.RealTime
		}

		fmt.Fprintf(os.Stdout, "%s\t%s\t%s\t%d\n",
			output.FormatTime(displayTime),
			formatLineNumber(d),
			d.Destination,
			d.DelaySeconds(),
		)
	}
}

func formatLineNumber(d api.Departure) string {
	if d.LinePublicNumber != "" {
		return d.LinePublicNumber
	}

	return strconv.Itoa(d.LineNumber)
}

func formatDelayStr(d api.Departure) string {
	if !d.IsRealTime() {
		return output.Dim("scheduled")
	}

	return output.FormatDelay(d.DelaySeconds())
}
