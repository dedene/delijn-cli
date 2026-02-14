---
name: delijn-cli
description: >
  Query Belgian Flemish public transport (De Lijn) via the delijn CLI. Use when the user wants
  real-time bus/tram departures, stop lookups, line searches, or watch mode for live updates.
  Triggered by mentions of De Lijn, Flanders buses, trams, Belgian public transport, or stop schedules.
license: MIT
homepage: https://github.com/dedene/delijn-cli
metadata:
  author: dedene
  version: "1.1.0"
  openclaw:
    primaryEnv: DELIJN_API_KEY
    requires:
      env:
        - DELIJN_API_KEY
        - DELIJN_KEYRING_BACKEND
      bins:
        - delijn
    install:
      - kind: brew
        tap: dedene/tap
        formula: delijn
        bins: [delijn]
      - kind: go
        package: github.com/dedene/delijn-cli/cmd/delijn
        bins: [delijn]
---

# delijn-cli

CLI for [De Lijn](https://www.delijn.be/) - Flemish public transport (buses and trams).

## Quick Start

```bash
# Verify auth
delijn auth status

# Get departures by stop number
delijn departures 200552

# Search stops
delijn stops search "Gent Sint-Pieters"
```

## Authentication

Requires De Lijn API key from [data.delijn.be](https://data.delijn.be/) (free registration).

```bash
# Store API key in system keyring
delijn auth set-key

# Verify setup
delijn auth status
```

If not authenticated, user must run `delijn auth set-key` interactively. Do not attempt auth setup on behalf of the user.

## Core Rules

1. **Always use `--json`** when parsing output programmatically
2. **Stop numbers are 6 digits** - e.g., `200552`
3. **Use favorites with @prefix** - e.g., `delijn departures @home`
4. **Watch mode blocks** - `--watch` refreshes every 30s; avoid in scripts

## Output Formats

| Flag | Format | Use case |
|------|--------|----------|
| (default) | Table | User-facing display |
| `--json` | JSON | Agent parsing, scripting |
| `--plain` | TSV | Pipe to awk/cut |

## Workflows

### Get Departures

```bash
# By stop number
delijn departures 200552

# By stop name (searches and picks)
delijn departures "Gent Sint-Pieters"

# Filter by line
delijn departures 200552 --line 1

# Limit results
delijn departures 200552 --limit 5

# JSON output for parsing
delijn departures 200552 --json
```

### Watch Mode (Live Updates)

```bash
# Auto-refresh every 30 seconds
delijn departures 200552 --watch
```

Note: Watch mode is interactive and blocks - use for user display only, not scripting.

### Search Stops

```bash
# Search by name
delijn stops search "Gent Sint-Pieters"

# Get stop details by number
delijn stops get 200552

# JSON for scripting
delijn stops search "Brugge" --json
```

### Search Lines

```bash
# Search lines
delijn lines search "1"

# Get line details (entity + line number)
delijn lines get 1 1
```

### Favorites

```bash
# Save a favorite stop
delijn config set-favorite home 200552
delijn config set-favorite work 301234

# Use favorites with @ prefix
delijn departures @home
delijn departures @work

# List all favorites
delijn config list-favorites
```

## Scripting Examples

```bash
# Get next departure time
delijn departures 200552 --limit 1 --json | jq -r '.[0].time'

# Find stop ID from name
delijn stops search "Gent" --json | jq -r '.[0].number'

# List all stops matching pattern
delijn stops search "station" --json | jq -r '.[] | "\(.number) \(.name)"'
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DELIJN_API_KEY` | API key (overrides keyring) |
| `DELIJN_KEYRING_BACKEND` | Keyring backend: keychain, file, pass |
| `NO_COLOR` | Disable colored output |

## Rate Limits

De Lijn API has two rate limits:

- **Core API**: 240 requests/minute (stops, lines, realtime)
- **Search API**: 6000 requests/minute

The CLI handles rate limiting automatically with backoff.

## Common Issues

### "API key not configured"
Run `delijn auth set-key` to store the API key interactively.

### Stop not found
Use 6-digit stop numbers. Search with `delijn stops search "name"` to find the correct number.

## Guidelines

- Never expose or log API keys
- Favorites require user setup - do not create favorites without user consent
- Watch mode is blocking - inform user it runs continuously


## Installation

```bash
brew install dedene/tap/delijn
```
