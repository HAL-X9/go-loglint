package config

import (
	"fmt"
)

func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if err := ValidateLogger(&cfg.Loggers); err != nil {
		return fmt.Errorf("loggers: %w", err)
	}

	return nil
}

func ValidateLogger(cfg *Loggers) error {

	return nil
}
