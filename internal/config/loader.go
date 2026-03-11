package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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

func LoadFromEnv() (*Config, error) {
	path := os.Getenv("LOGLINT_CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("LOGLINT_CONFIG_PATH env var is not set")
	}

	return LoadConfig(path)
}
