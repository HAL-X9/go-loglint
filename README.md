# go-loglint

A Go linter for log statements that enforces best practices and consistency across your logging code. Works as a golangci-lint plugin and as a standalone CLI.

## Features

- **Configurable rules** — Enable or disable rules per project
- **Multiple logger support** — Works with `log`, `slog`, `logrus`, `zap`, and more
- **YAML configuration** — Simple, human-readable config format

### Built-in Rules

| Rule | Description |
|------|-------------|
| `lowercase_start` | Ensures log messages start with lowercase |
| `english_only` | Enforces English-only log messages |
| `no_special_chars` | Only letters, digits, space allowed |
| `no_sensitive_data` | Detects potential sensitive data in logs |

## Quick Start

### Prerequisites

- Go 1.22 or later

### Installation

```bash
# Clone the repository
git clone https://github.com/HAL-X9/go-loglint.git
cd loglint

# Build the binary
make build
```

### Usage

```bash
# All rules enabled by default
./bin/go-loglint ./...

# With config
./bin/go-loglint -loglint-config=loglint.yaml ./...

# Run on testdata
./bin/go-loglint ./testdata/example
```

## Examples

### Rule 1: Lowercase start

```go
// Wrong
log.Print("Starting server on port 8080")
slog.Error("Failed to connect to database")

// Correct
log.Print("starting server on port 8080")
slog.Error("failed to connect to database")
```

### Rule 2: English only

```go
// Wrong
slog.Info("запуск сервера")
slog.Error("ошибка подключения к базе данных")

// Correct
slog.Info("starting server")
slog.Error("failed to connect to database")
```

### Rule 3: No special characters

```go
// Wrong
slog.Info("server started!")
slog.Error("connection failed!!!")
slog.Warn("warning: something went wrong...")

// Correct
slog.Info("server started")
slog.Error("connection failed")
slog.Warn("something went wrong")
```

### Rule 4: No sensitive data

```go
// Wrong
slog.Info("user password: " + password)
slog.Debug("api_key=" + apiKey)
slog.Info("token: " + token)

// Correct
slog.Info("user authenticated successfully")
slog.Debug("api request completed")
slog.Info("authentication completed")
```

## golangci-lint Plugin

Uses **Module Plugin** — the recommended golangci-lint integration. Works with any installation (brew, go install), no CGO required.

### 1. Create `.custom-gcl.yml` in your project

```yaml
# Get version: golangci-lint version
version: v2.11.3

plugins:
  # Remote (from GitHub)
  - module: github.com/HAL-X9/go-loglint
    import: github.com/HAL-X9/go-loglint/plugin/module
    version: v1.1.0

  # Or local
  # - module: github.com/HAL-X9/go-loglint
  #   path: /path/to/loglint
  #   import: github.com/HAL-X9/go-loglint/plugin/module
```

### 2. Build custom golangci-lint

```bash
golangci-lint custom
```

Creates `./custom-gcl` binary.

### 3. Add to `.golangci.yml`

```yaml
version: "2"
linters:
  enable:
    - loglint
  settings:
    custom:
      loglint:
        type: module
        description: Log message linter
```

### 4. Run

```bash
# Without config — all rules enabled by default
./custom-gcl cache clean
./custom-gcl run

# With config — use LOGLINT_CONFIG_PATH
./custom-gcl cache clean
LOGLINT_CONFIG_PATH=./loglint.yaml ./custom-gcl run
```

**Tip:** Add `custom-gcl` to `.gitignore` — each developer builds locally. In CI: add `golangci-lint custom` step before running.

**Note:** golangci-lint does not pass `settings.config` to module plugins. Use `LOGLINT_CONFIG_PATH` for config.
If you change `loglint.yaml` and get inconsistent results, run `./custom-gcl cache clean` before the next run.

## Configuration

Config file example (`loglint.yaml`):

```yaml
rules:
  - name: lowercase_start
    enable: true
  - name: english_only
    enable: true
  - name: no_special_chars
    enable: true
  - name: no_sensitive_data
    enable: true
loggers:
  include: [log, slog, logrus.Logger, zap.SugaredLogger]
```

| Option | Type | Description |
|--------|------|-------------|
| `rules` | array | List of rules with `name` and `enable` |
| `loggers.include` | array | Logger names/packages to analyze |

## Testdata

Example files in `testdata/example/` demonstrate all rules:

```bash
./bin/go-loglint ./testdata/example
```

## Development

```bash
make help          # Show available commands
make run-loglint   # Run loglint locally
make build         # Build binary to ./bin/
make build-plugin  # Build Go plugin
make test          # Run tests
make clean         # Remove build artifacts
```
