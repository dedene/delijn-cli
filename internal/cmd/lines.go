package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/dedene/delijn-cli/internal/api"
	"github.com/dedene/delijn-cli/internal/output"
)

type LinesCmd struct {
	Search LinesSearchCmd `cmd:"" help:"Search lines by number or name"`
	Get    LinesGetCmd    `cmd:"" help:"Get line details"`
}

type LinesSearchCmd struct {
	Query string `arg:"" required:"" help:"Line number or name"`
}

func (c *LinesSearchCmd) Run(root *RootFlags) error {
	output.SetNoColor(root.NoColor)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	resp, err := client.SearchLines(ctx, c.Query)
	if err != nil {
		return fmt.Errorf("search lines: %w", err)
	}

	if root.JSON {
		return outputJSON(resp.Lines)
	}

	if root.Plain {
		outputLinesPlain(resp.Lines)

		return nil
	}

	outputLinesTable(resp.Lines)

	return nil
}

type LinesGetCmd struct {
	Entity int `arg:"" required:"" help:"Entity number (1-5)"`
	Line   int `arg:"" required:"" help:"Line number"`
}

func (c *LinesGetCmd) Run(root *RootFlags) error {
	output.SetNoColor(root.NoColor)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	line, err := client.GetLine(ctx, c.Entity, c.Line)
	if err != nil {
		return fmt.Errorf("get line: %w", err)
	}

	if root.JSON {
		return outputJSON(line)
	}

	if root.Plain {
		outputLinePlain(line)

		return nil
	}

	outputLineDetails(line)

	return nil
}

func outputLinesTable(lines []api.Line) {
	if len(lines) == 0 {
		fmt.Fprintln(os.Stdout, "No lines found.")

		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "ENTITY\tNUMBER\tPUBLIC\tTYPE\tDESCRIPTION")

	for _, l := range lines {
		fmt.Fprintf(w, "%d\t%d\t%s\t%s\t%s\n",
			l.EntityNumber,
			l.LineNumber,
			l.PublicNumber,
			l.TransportType,
			l.Description,
		)
	}
}

func outputLinesPlain(lines []api.Line) {
	for _, l := range lines {
		fmt.Fprintf(os.Stdout, "%d\t%d\t%s\t%s\t%s\n",
			l.EntityNumber,
			l.LineNumber,
			l.PublicNumber,
			l.TransportType,
			l.Description,
		)
	}
}

func outputLinePlain(line *api.Line) {
	fmt.Fprintf(os.Stdout, "%d\t%d\t%s\t%s\t%s\n",
		line.EntityNumber,
		line.LineNumber,
		line.PublicNumber,
		line.TransportType,
		line.Description,
	)
}

func outputLineDetails(line *api.Line) {
	fmt.Fprintf(os.Stdout, "Line:        %s\n", line.PublicNumber)
	fmt.Fprintf(os.Stdout, "Description: %s\n", line.Description)
	fmt.Fprintf(os.Stdout, "Type:        %s\n", line.TransportType)
	fmt.Fprintf(os.Stdout, "Entity:      %d\n", line.EntityNumber)
	fmt.Fprintf(os.Stdout, "Number:      %d\n", line.LineNumber)
}
