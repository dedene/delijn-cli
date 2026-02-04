# 🚌 delijn-cli - De Lijn in your terminal

A command-line interface for [De Lijn](https://www.delijn.be/) - Flemish public transport.

[![CI](https://github.com/dedene/delijn-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/dedene/delijn-cli/actions/workflows/ci.yml)
[![Go 1.23+](https://img.shields.io/badge/go-1.23+-00ADD8.svg)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/dedene/delijn-cli)](https://goreportcard.com/report/github.com/dedene/delijn-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

## Features

- **Realtime departures** - See when the next bus/tram arrives
- **Stop search** - Find stops by name
- **Line search** - Look up bus and tram lines
- **Watch mode** - Auto-refresh departures every 30 seconds
- **Favorites** - Save frequently used stops as aliases
- **Multiple output formats** - Human-readable, JSON, or plain TSV

## Installation

### Homebrew (macOS/Linux)

```bash
brew install dedene/tap/delijn
```

### Go install

```bash
go install github.com/dedene/delijn-cli/cmd/delijn@latest
```

### From source

```bash
git clone https://github.com/dedene/delijn-cli.git
cd delijn-cli
make build
```

## Setup

You need a De Lijn API key. Get one free at [data.delijn.be](https://data.delijn.be/).

```bash
# Store your API key securely in the system keyring
delijn auth set-key

# Verify it's configured
delijn auth status
```

## Usage

### Departures

```bash
# By stop number (6-digit)
delijn departures 200552

# By stop name (searches and picks)
delijn departures "Gent Sint-Pieters"

# Watch mode - refreshes every 30 seconds
delijn departures 200552 --watch

# Filter by line
delijn departures 200552 --line 1

# Limit results
delijn departures 200552 --limit 5
```

### Stops

```bash
# Search stops by name
delijn stops search "Gent Sint-Pieters"

# Get stop details by number
delijn stops get 200552
```

### Lines

```bash
# Search lines
delijn lines search "1"

# Get line details (entity number + line number)
delijn lines get 1 1
```

### Favorites

```bash
# Save a favorite stop
delijn config set-favorite home 200552

# Use the favorite
delijn departures @home

# List all favorites
delijn config list-favorites
```

### Output formats

```bash
# JSON output (for scripting)
delijn departures 200552 --json

# Plain TSV output
delijn departures 200552 --plain

# Disable colors
delijn departures 200552 --no-color
# Or set NO_COLOR=1 environment variable
```

## Shell completions

```bash
# Bash
delijn completion bash > /etc/bash_completion.d/delijn

# Zsh
delijn completion zsh > "${fpath[1]}/_delijn"

# Fish
delijn completion fish > ~/.config/fish/completions/delijn.fish
```

## Environment variables

| Variable                 | Description                                 |
| ------------------------ | ------------------------------------------- |
| `DELIJN_API_KEY`         | API key (overrides keyring)                 |
| `DELIJN_KEYRING_BACKEND` | Keyring backend: `keychain`, `file`, `pass` |
| `NO_COLOR`               | Disable colored output                      |

## API Rate Limits

De Lijn API has two rate limits:

- **Core API**: 240 requests/minute (stops, lines, realtime)
- **Search API**: 6000 requests/minute

The CLI handles rate limiting automatically.

## License

MIT - see [LICENSE](LICENSE)

## Credits

Data provided by [De Lijn Open Data](https://data.delijn.be/).
