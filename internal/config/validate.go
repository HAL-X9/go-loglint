package config

import (
	"fmt"
)

func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if err := ValidateRule(cfg.Rules); err != nil {
		return fmt.Errorf("rule: %w", err)
	}

	if err := ValidateLoggers(&cfg.Loggers); err != nil {
		return fmt.Errorf("loggers: %w", err)
	}

	return nil
}

func ValidateRule(rules []Rule) error {
	if len(rules) == 0 {
		return fmt.Errorf("no rules defined in the configuration")
	}

	for _, r := range rules {
		switch r.Name {
		case "lowercase_start", "english_only", "no_special_chars", "no_sensitive_data":
		default:
			return fmt.Errorf("rule name must be: lowercase_start, english_only, no_special_chars, no_sensitive_data")
		}
	}
	return nil
}

func ValidateLoggers(cfg *Loggers) error {
	if cfg == nil || len(cfg.Include) == 0 {
		return fmt.Errorf("no loggers defined in the configuration")
	}

	for _, l := range cfg.Include {
		switch l {
		case "log", "slog", "logrus.Logger", "zap.SugaredLogger":
		default:
			return fmt.Errorf("invalid logger name: %q, allowed: log, slog, logrus.Logger, zap.SugaredLogger", l)
		}
	}

	return nil
}
