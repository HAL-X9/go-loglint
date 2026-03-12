package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validYAML = `
auto_fix: true
rules:
  - name: lowercase_start
    enable: true
  - name: english_only
    enable: false
loggers:
  include:
    - log
    - slog
`

func TestLoadConfig(t *testing.T) {
	t.Run("file_not_found", func(t *testing.T) {
		_, err := LoadConfig("/nonexistent/path/config.yaml")
		if err == nil {
			t.Fatal("LoadConfig() expected error for missing file, got nil")
		}
		if !strings.Contains(err.Error(), "read config") {
			t.Errorf("LoadConfig() error should mention read, got: %v", err)
		}
		if !strings.Contains(err.Error(), "/nonexistent/path/config.yaml") {
			t.Errorf("LoadConfig() error should include path, got: %v", err)
		}
	})

	t.Run("invalid_yaml", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.yaml")
		if err := os.WriteFile(path, []byte("invalid: yaml: content: ["), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		_, err := LoadConfig(path)
		if err == nil {
			t.Fatal("LoadConfig() expected error for invalid YAML, got nil")
		}
		if !strings.Contains(err.Error(), "unmarshal") {
			t.Errorf("LoadConfig() error should mention unmarshal, got: %v", err)
		}
	})

	t.Run("invalid_config_validation_fails", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "invalid.yaml")
		yaml := `
auto_fix: true
rules:
  - name: unknown_rule
    enable: true
loggers:
  include:
    - log
`
		if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		_, err := LoadConfig(path)
		if err == nil {
			t.Fatal("LoadConfig() expected error for invalid config, got nil")
		}
		if !strings.Contains(err.Error(), "config validation failed") {
			t.Errorf("LoadConfig() error should mention validation, got: %v", err)
		}
		if !strings.Contains(err.Error(), "rule name must be") {
			t.Errorf("LoadConfig() error should wrap validation error, got: %v", err)
		}
	})

	t.Run("empty_rules_validation_fails", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "empty_rules.yaml")
		yaml := `
auto_fix: false
rules: []
loggers:
  include:
    - slog
`
		if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		_, err := LoadConfig(path)
		if err == nil {
			t.Fatal("LoadConfig() expected error for empty rules, got nil")
		}
		if !strings.Contains(err.Error(), "config validation failed") {
			t.Errorf("LoadConfig() error should mention validation, got: %v", err)
		}
	})

	t.Run("empty_loggers_validation_fails", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "empty_loggers.yaml")
		yaml := `
auto_fix: true
rules:
  - name: lowercase_start
    enable: true
loggers:
  include: []
`
		if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		_, err := LoadConfig(path)
		if err == nil {
			t.Fatal("LoadConfig() expected error for empty loggers, got nil")
		}
		if !strings.Contains(err.Error(), "config validation failed") {
			t.Errorf("LoadConfig() error should mention validation, got: %v", err)
		}
	})

	t.Run("valid_config", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "valid.yaml")
		if err := os.WriteFile(path, []byte(validYAML), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		cfg, err := LoadConfig(path)
		if err != nil {
			t.Fatalf("LoadConfig() unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadConfig() returned nil config")
		}
		if !cfg.AutoFix {
			t.Error("LoadConfig() AutoFix should be true")
		}
		if len(cfg.Rules) != 2 {
			t.Errorf("LoadConfig() Rules len = %d, want 2", len(cfg.Rules))
		}
		if cfg.Rules[0].Name != "lowercase_start" || !cfg.Rules[0].Enable {
			t.Errorf("LoadConfig() first rule = %+v", cfg.Rules[0])
		}
		if cfg.Rules[1].Name != "english_only" || cfg.Rules[1].Enable {
			t.Errorf("LoadConfig() second rule = %+v", cfg.Rules[1])
		}
		if len(cfg.Loggers.Include) != 2 {
			t.Errorf("LoadConfig() Loggers.Include len = %d, want 2", len(cfg.Loggers.Include))
		}
		if cfg.Loggers.Include[0] != "log" || cfg.Loggers.Include[1] != "slog" {
			t.Errorf("LoadConfig() Loggers.Include = %v", cfg.Loggers.Include)
		}
	})

	t.Run("valid_config_all_options", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "full.yaml")
		yaml := `
auto_fix: true
rules:
  - name: lowercase_start
    enable: true
  - name: english_only
    enable: true
  - name: no_special_chars
    enable: false
  - name: no_sensitive_data
    enable: true
loggers:
  include:
    - log
    - slog
    - logrus.Logger
    - zap.SugaredLogger
`
		if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		cfg, err := LoadConfig(path)
		if err != nil {
			t.Fatalf("LoadConfig() unexpected error: %v", err)
		}
		if len(cfg.Rules) != 4 {
			t.Errorf("LoadConfig() Rules len = %d, want 4", len(cfg.Rules))
		}
		if len(cfg.Loggers.Include) != 4 {
			t.Errorf("LoadConfig() Loggers.Include len = %d, want 4", len(cfg.Loggers.Include))
		}
	})
}

func TestLoadFromEnv(t *testing.T) {
	saveEnv := func() (string, bool) {
		v, ok := os.LookupEnv("LOGLINT_CONFIG_PATH")
		return v, ok
	}
	restoreEnv := func(prev string, wasSet bool) {
		if wasSet {
			os.Setenv("LOGLINT_CONFIG_PATH", prev)
		} else {
			os.Unsetenv("LOGLINT_CONFIG_PATH")
		}
	}

	t.Run("env_not_set", func(t *testing.T) {
		prev, wasSet := saveEnv()
		defer restoreEnv(prev, wasSet)
		os.Unsetenv("LOGLINT_CONFIG_PATH")

		_, err := LoadFromEnv()
		if err == nil {
			t.Fatal("LoadFromEnv() expected error when env not set, got nil")
		}
		if !strings.Contains(err.Error(), "LOGLINT_CONFIG_PATH") {
			t.Errorf("LoadFromEnv() error should mention env var, got: %v", err)
		}
		if !strings.Contains(err.Error(), "not set") {
			t.Errorf("LoadFromEnv() error should mention not set, got: %v", err)
		}
	})

	t.Run("env_set_empty", func(t *testing.T) {
		prev, wasSet := saveEnv()
		defer restoreEnv(prev, wasSet)
		os.Setenv("LOGLINT_CONFIG_PATH", "")

		_, err := LoadFromEnv()
		if err == nil {
			t.Fatal("LoadFromEnv() expected error when env empty, got nil")
		}
		if !strings.Contains(err.Error(), "LOGLINT_CONFIG_PATH") {
			t.Errorf("LoadFromEnv() error should mention env var, got: %v", err)
		}
	})

	t.Run("env_set_valid_path", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		if err := os.WriteFile(path, []byte(validYAML), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		prev, wasSet := saveEnv()
		defer restoreEnv(prev, wasSet)
		os.Setenv("LOGLINT_CONFIG_PATH", path)

		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("LoadFromEnv() unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadFromEnv() returned nil config")
		}
		if len(cfg.Rules) != 2 {
			t.Errorf("LoadFromEnv() Rules len = %d, want 2", len(cfg.Rules))
		}
	})

	t.Run("env_set_invalid_path", func(t *testing.T) {
		prev, wasSet := saveEnv()
		defer restoreEnv(prev, wasSet)
		os.Setenv("LOGLINT_CONFIG_PATH", "/nonexistent/config.yaml")

		_, err := LoadFromEnv()
		if err == nil {
			t.Fatal("LoadFromEnv() expected error for invalid path, got nil")
		}
		if !strings.Contains(err.Error(), "read config") {
			t.Errorf("LoadFromEnv() error should propagate LoadConfig error, got: %v", err)
		}
	})

	t.Run("env_set_invalid_yaml", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.yaml")
		if err := os.WriteFile(path, []byte("not: valid: yaml: [[["), 0644); err != nil {
			t.Fatalf("setup: write temp file: %v", err)
		}

		prev, wasSet := saveEnv()
		defer restoreEnv(prev, wasSet)
		os.Setenv("LOGLINT_CONFIG_PATH", path)

		_, err := LoadFromEnv()
		if err == nil {
			t.Fatal("LoadFromEnv() expected error for invalid YAML, got nil")
		}
		if !strings.Contains(err.Error(), "unmarshal") {
			t.Errorf("LoadFromEnv() error should propagate unmarshal error, got: %v", err)
		}
	})
}
