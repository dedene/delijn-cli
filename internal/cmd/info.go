package cmd

import (
	"fmt"
	"os"

	"github.com/dedene/delijn-cli/internal/auth"
	"github.com/dedene/delijn-cli/internal/config"
)

type InfoCmd struct{}

func (c *InfoCmd) Run() error {
	configPath, _ := config.ConfigPath()
	keyringDir, _ := config.KeyringDir()

	fmt.Fprintf(os.Stdout, "De Lijn CLI - %s\n", VersionString())
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Config path:     %s\n", configPath)
	fmt.Fprintf(os.Stdout, "Keyring dir:     %s\n", keyringDir)

	backendInfo, err := auth.ResolveKeyringBackendInfo()
	if err == nil {
		fmt.Fprintf(os.Stdout, "Keyring backend: %s (source: %s)\n", backendInfo.Value, backendInfo.Source)
	}

	// Check API key status
	store, err := auth.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stdout, "API key:         error opening keyring: %v\n", err)
	} else {
		hasKey, checkErr := store.HasAPIKey()

		switch {
		case checkErr != nil:
			fmt.Fprintf(os.Stdout, "API key:         error checking: %v\n", checkErr)
		case hasKey:
			fmt.Fprintln(os.Stdout, "API key:         configured")
		default:
			fmt.Fprintln(os.Stdout, "API key:         not configured")
		}
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "API endpoints:")
	fmt.Fprintln(os.Stdout, "  Core:   https://api.delijn.be/DLKernOpenData/v1/beta (240 req/min)")
	fmt.Fprintln(os.Stdout, "  Search: https://api.delijn.be/DLZoekOpenData/v1/beta (6000 req/min)")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Get your API key from https://data.delijn.be/")

	return nil
}
