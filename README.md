# go-loglint

A Go linter for log statements that enforces best practices and consistency across your logging code.

## Features

- **Configurable rules** — Enable or disable rules per project
- **Auto-fix support** — Automatically fix issues where possible
- **Multiple logger support** — Works with `log`, `slog`, `logrus`, `zap`, and more
- **YAML configuration** — Simple, human-readable config format

### Built-in Rules

| Rule | Description |
|------|-------------|
| `lowercase_start` | Ensures log messages start with lowercase |
| `english_only` | Enforces English-only log messages |
| `no_special_chars` | Flags inappropriate special characters |
| `no_sensitive_data` | Detects potential sensitive data in logs |

## Quick Start

### Prerequisites

- Go 1.26 or later

### Installation

```bash
# Clone the repository
git clone https://github.com/loglint/loglint.git
cd loglint

# Build the binary
make build
```

### Usage

1. **Create a config file** (or use the example in `configs/loglint.yaml`):

```yaml
version: 1.0
auto_fix: true

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
  include:
    - log
    - slog
    - logrus.Logger
    - zap.SugaredLogger
```

2. **Run loglint**:

```bash
# Using Make (recommended)
LOGLINT_CONFIG=./configs/loglint.yaml make run-loglint

# Or with environment variable
export LOGLINT_CONFIG_PATH=./configs/loglint.yaml
go run ./cmd
```

## Configuration

| Option | Type | Description |
|--------|------|-------------|
| `version` | string | Config schema version |
| `auto_fix` | bool | Enable automatic fixes when possible |
| `rules` | array | List of rules with `name` and `enable` |
| `loggers.include` | array | Logger names/packages to analyze |

## Development

```bash
make help          # Show available commands
make run-loglint   # Run loglint locally
make build         # Build binary to ./bin/
make test          # Run tests
make clean         # Remove build artifacts
```
