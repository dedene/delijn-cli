package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/dedene/delijn-cli/internal/auth"
	"github.com/dedene/delijn-cli/internal/config"
)

type AuthCmd struct {
	SetKey AuthSetKeyCmd `cmd:"" name:"set-key" help:"Store API key in keyring"`
	Status AuthStatusCmd `cmd:"" help:"Show authentication status"`
	Remove AuthRemoveCmd `cmd:"" help:"Remove stored API key"`
}

type AuthSetKeyCmd struct {
	Stdin bool `help:"Read API key from stdin (for scripts)"`
}

func (c *AuthSetKeyCmd) Run() error {
	var key string

	if c.Stdin {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			key = strings.TrimSpace(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read from stdin: %w", err)
		}
	} else {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return fmt.Errorf("not a terminal; use --stdin flag to read from pipe")
		}

		fmt.Print("Enter your De Lijn API key: ")

		bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err != nil {
			return fmt.Errorf("read API key: %w", err)
		}

		key = strings.TrimSpace(string(bytes))
	}

	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	store, err := auth.OpenDefault()
	if err != nil {
		return fmt.Errorf("open keyring: %w", err)
	}

	if err := store.SetAPIKey(key); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	fmt.Fprintln(os.Stdout, "API key stored successfully.")
	fmt.Fprintln(os.Stdout, "Get your API key from https://data.delijn.be/")

	return nil
}

type AuthStatusCmd struct{}

func (c *AuthStatusCmd) Run() error {
	backendInfo, err := auth.ResolveKeyringBackendInfo()
	if err != nil {
		return err
	}

	configPath, _ := config.ConfigPath()
	keyringDir, _ := config.KeyringDir()

	fmt.Fprintf(os.Stdout, "Config path:     %s\n", configPath)
	fmt.Fprintf(os.Stdout, "Keyring dir:     %s\n", keyringDir)
	fmt.Fprintf(os.Stdout, "Keyring backend: %s (source: %s)\n", backendInfo.Value, backendInfo.Source)

	store, err := auth.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stdout, "API key:         error opening keyring: %v\n", err)

		return nil
	}

	hasKey, err := store.HasAPIKey()
	if err != nil {
		fmt.Fprintf(os.Stdout, "API key:         error checking: %v\n", err)

		return nil
	}

	if hasKey {
		fmt.Fprintln(os.Stdout, "API key:         configured")
	} else {
		fmt.Fprintln(os.Stdout, "API key:         not configured")
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "Run 'delijn auth set-key' to configure your API key.")
		fmt.Fprintln(os.Stdout, "Get your API key from https://data.delijn.be/")
	}

	return nil
}

type AuthRemoveCmd struct{}

func (c *AuthRemoveCmd) Run() error {
	store, err := auth.OpenDefault()
	if err != nil {
		return fmt.Errorf("open keyring: %w", err)
	}

	if err := store.DeleteAPIKey(); err != nil {
		return fmt.Errorf("remove API key: %w", err)
	}

	fmt.Fprintln(os.Stdout, "API key removed.")

	return nil
}
