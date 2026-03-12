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
git clone https://github.com/loglint/loglint.git
cd loglint

# Build the binary
make build
```

### Usage

```bash
# All rules enabled by default
./bin/go-loglint ./...

# With config
./bin/go-loglint -config=loglint.yaml ./...

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

## golangci-lint Plugin (Module Plugin — recommended)

Module plugin works with any golangci-lint installation. No CGO required.

### 1. Create `.custom-gcl.yml` in your project

```yaml
# Get your version: golangci-lint version
version: v2.10.0

plugins:
  - module: github.com/go-loglint
    path: /path/to/loglint   # or use version: v1.0.0 when published
    import: github.com/go-loglint/plugin/module
```

### 2. Build custom golangci-lint

```bash
golangci-lint custom
```

This creates `./custom-gcl` (or `./golangci-lint` depending on config).

### 3. Add to `.golangci.yml`

```yaml
version: "2"
linters:
  enable:
    - loglint
    # ... other linters
  settings:
    custom:
      loglint:
        type: module
        description: Log message linter
```

### 4. Run

```bash
./custom-gcl run
```

---

### Go Plugin (legacy, requires CGO)

```bash
CGO_ENABLED=1 make build-plugin
cp ./bin/loglint.so ~/.golangci-lint/plugins/
```

Add to `.golangci.yml`:

```yaml
linters:
  settings:
    custom:
      loglint:
        path: /Users/sergey/.golangci-lint/plugins/loglint.so
        description: Log message linter
```

Note: Go plugin requires golangci-lint built with CGO; pre-built binaries often don't support it.

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
make build-plugin  # Build golangci-lint plugin
make test          # Run tests
make clean         # Remove build artifacts
```
