package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// GetPathFromGolangciLintSettings extracts config path from golangci-lint custom linter settings.
// Handles both flat {config: "..."} and nested {settings: {config: "..."}}.
func GetPathFromGolangciLintSettings(settings any) string {
	m := toMap(settings)
	if m == nil {
		return ""
	}
	if c, ok := m["config"].(string); ok && c != "" {
		return c
	}
	if s, ok := m["settings"].(map[string]any); ok {
		if c, ok := s["config"].(string); ok && c != "" {
			return c
		}
	}
	return ""
}

func toMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m map[string]any
	if json.Unmarshal(data, &m) != nil {
		return nil
	}
	return m
}

// LoadConfig reads a YAML configuration file from the given path,
// unmarshals it into a Config struct, and validates it.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var config Config

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal config %s: %w", path, err)
	}

	if err = Validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// LoadFromEnv loads the configuration from the path specified by
// the LOGLINT_CONFIG_PATH environment variable.
// Returns an error if the variable is not set or the file is invalid.
func LoadFromEnv() (*Config, error) {
	path := os.Getenv("LOGLINT_CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("LOGLINT_CONFIG_PATH env var is not set")
	}

	return LoadConfig(path)
}
