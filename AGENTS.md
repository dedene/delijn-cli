# Repository Guidelines

## Project Structure

- `cmd/delijn/`: CLI entrypoint
- `internal/`: implementation
  - `cmd/`: command routing (kong CLI framework)
  - `api/`: De Lijn API client + circuit breaker + rate limiting
  - `auth/`: keyring-based credential storage
  - `config/`: YAML config + favorites management
  - `output/`: terminal rendering (termenv)
  - `errfmt/`: error formatting
- `bin/`: build outputs

## Build, Test, and Development Commands

- `make build`: compile to `bin/delijn`
- `make delijn -- <args>`: build + run in one step
- `make fmt` / `make lint` / `make test` / `make ci`: format, lint, test, full local gate
- `make tools`: install pinned dev tools into `.tools/`
- `make clean`: remove bin/ and .tools/

## Coding Style & Naming Conventions

- Formatting: `make fmt` (goimports local prefix `github.com/dedene/delijn-cli` + gofumpt)
- Output: keep stdout parseable (`--json`); send human hints/progress to stderr
- Linting: golangci-lint with 47 linters enabled (.golangci.yml)

## Testing Guidelines

- Unit tests: stdlib `testing` (files: `*_test.go` next to code)
- Coverage areas: colors, circuit breaker, rate limiter
- CI uploads coverage to Codecov

## Config & Secrets

- **Keyring**: 99designs/keyring for API credentials (macOS Keychain, Linux SecretService, Windows Credential Manager)
- **Config file**: YAML-based via `internal/config`
- **Favorites**: local store for frequently used stops
- **Env overrides**:
  - `DELIJN_KEYRING_BACKEND`: keyring backend (auto/keychain/file)
  - `DELIJN_KEYRING_PASSWORD`: password for file backend

## Commit & Pull Request Guidelines

- Conventional Commits: `feat|fix|refactor|build|ci|chore|docs|style|perf|test`
- Group related changes; avoid bundling unrelated refactors
- PR review: use `gh pr view` / `gh pr diff`; don't switch branches

## Security Tips

- Never commit API credentials or keyring passwords
- Use OS keychain; file backend only for headless environments
- Circuit breaker prevents credential leakage during API failures
