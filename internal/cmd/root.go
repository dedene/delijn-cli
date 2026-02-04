package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type RootFlags struct {
	JSON    bool `help:"Output JSON to stdout"`
	Plain   bool `help:"Output plain TSV (for scripting)"`
	NoColor bool `help:"Disable colors" env:"NO_COLOR"`
}

type CLI struct {
	RootFlags `embed:""`

	Version    kong.VersionFlag `help:"Print version and exit"`
	VersionCmd VersionCmd       `cmd:"" name:"version" help:"Print version"`
	Auth       AuthCmd          `cmd:"" help:"Manage API key"`
	Config     ConfigCmd        `cmd:"" help:"Manage configuration"`
	Stops      StopsCmd         `cmd:"" help:"Search and view stops"`
	Lines      LinesCmd         `cmd:"" help:"Search and view lines"`
	Departures DeparturesCmd    `cmd:"" help:"Show realtime departures"`
	Info       InfoCmd          `cmd:"" help:"Show CLI and API info"`
	Completion CompletionCmd    `cmd:"" help:"Generate shell completions"`
}

type exitPanic struct{ code int }

func Execute(args []string) (err error) {
	parser, err := newParser()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if ep, ok := r.(exitPanic); ok {
				if ep.code == 0 {
					err = nil

					return
				}

				err = &ExitError{Code: ep.code, Err: errors.New("exited")}

				return
			}

			panic(r)
		}
	}()

	if len(args) == 0 {
		args = []string{"--help"}
	}

	kctx, err := parser.Parse(args)
	if err != nil {
		parsedErr := wrapParseError(err)
		_, _ = fmt.Fprintln(os.Stderr, parsedErr)

		return parsedErr
	}

	err = kctx.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		return err
	}

	return nil
}

func wrapParseError(err error) error {
	if err == nil {
		return nil
	}

	var parseErr *kong.ParseError
	if errors.As(err, &parseErr) {
		return &ExitError{Code: 2, Err: parseErr}
	}

	return err
}

func newParser() (*kong.Kong, error) {
	vars := kong.Vars{
		"version": VersionString(),
	}

	cli := &CLI{}
	parser, err := kong.New(
		cli,
		kong.Name("delijn"),
		kong.Description("De Lijn CLI - Flemish public transport from the command line"),
		kong.Vars(vars),
		kong.Writers(os.Stdout, os.Stderr),
		kong.Exit(func(code int) { panic(exitPanic{code: code}) }),
		kong.Bind(&cli.RootFlags),
		kong.Help(helpPrinter),
		kong.ConfigureHelp(helpOptions()),
	)
	if err != nil {
		return nil, err
	}

	return parser, nil
}
