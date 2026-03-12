package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	prevConfigPath := configPath
	prevLegacyConfigPath := legacyConfigPath
	prevEnv, envWasSet := os.LookupEnv("LOGLINT_CONFIG_PATH")
	configPath = ""
	legacyConfigPath = ""
	os.Unsetenv("LOGLINT_CONFIG_PATH")
	t.Cleanup(func() {
		configPath = prevConfigPath
		legacyConfigPath = prevLegacyConfigPath
		if envWasSet {
			_ = os.Setenv("LOGLINT_CONFIG_PATH", prevEnv)
			return
		}
		os.Unsetenv("LOGLINT_CONFIG_PATH")
	})

	filemap := map[string]string{
		"pkg/pkg.go": `package pkg

import (
	"fmt"
	"log"
	"log/slog"
)

func F(password, apiKey string) {
	log.Print("Uppercase")   // want "log message must start with lowercase"
	log.Print("lowercase")
	log.Print("")
	slog.Info("Also bad")    // want "log message must start with lowercase"
	slog.Info("also ok")
	fmt.Print("Uppercase")
	slog.Info("user password " + password)  // want "log message must not contain sensitive data"
	slog.Debug("token " + apiKey)           // want "log message must not contain sensitive data"
}
`,
	}
	dir, cleanup, err := analysistest.WriteFiles(filemap)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	analysistest.Run(t, dir, Analyzer, "pkg")
}

func TestAnalyzer_ConfigPathInvalid_DisablesRules(t *testing.T) {
	prevConfigPath := configPath
	prevLegacyConfigPath := legacyConfigPath
	prevEnv, envWasSet := os.LookupEnv("LOGLINT_CONFIG_PATH")
	configPath = ""
	legacyConfigPath = ""
	_ = os.Setenv("LOGLINT_CONFIG_PATH", "/definitely/missing/loglint.yaml")
	t.Cleanup(func() {
		configPath = prevConfigPath
		legacyConfigPath = prevLegacyConfigPath
		if envWasSet {
			_ = os.Setenv("LOGLINT_CONFIG_PATH", prevEnv)
			return
		}
		os.Unsetenv("LOGLINT_CONFIG_PATH")
	})

	filemap := map[string]string{
		"pkg/pkg.go": `package pkg

import (
	"log"
)

func F() {
	log.Print("Uppercase")
	log.Print("message with !!!")
}
`,
	}
	dir, cleanup, err := analysistest.WriteFiles(filemap)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// No diagnostics expected: explicit config path is invalid.
	analysistest.Run(t, dir, Analyzer, "pkg")
}

func TestResolveConfiguredPath(t *testing.T) {
	prevConfigPath := configPath
	prevLegacyConfigPath := legacyConfigPath
	prevEnv, envWasSet := os.LookupEnv("LOGLINT_CONFIG_PATH")
	t.Cleanup(func() {
		configPath = prevConfigPath
		legacyConfigPath = prevLegacyConfigPath
		if envWasSet {
			_ = os.Setenv("LOGLINT_CONFIG_PATH", prevEnv)
			return
		}
		os.Unsetenv("LOGLINT_CONFIG_PATH")
	})

	t.Run("env_has_highest_priority", func(t *testing.T) {
		configPath = "flag.yaml"
		legacyConfigPath = "legacy.yaml"
		_ = os.Setenv("LOGLINT_CONFIG_PATH", "env.yaml")

		got, explicit := resolveConfiguredPath()
		if !explicit {
			t.Fatal("resolveConfiguredPath() explicit = false, want true")
		}
		if got != "env.yaml" {
			t.Fatalf("resolveConfiguredPath() path = %q, want env.yaml", got)
		}
	})

	t.Run("uses_new_flag_when_env_empty", func(t *testing.T) {
		_ = os.Unsetenv("LOGLINT_CONFIG_PATH")
		configPath = "flag.yaml"
		legacyConfigPath = "legacy.yaml"

		got, explicit := resolveConfiguredPath()
		if !explicit {
			t.Fatal("resolveConfiguredPath() explicit = false, want true")
		}
		if got != "flag.yaml" {
			t.Fatalf("resolveConfiguredPath() path = %q, want flag.yaml", got)
		}
	})

	t.Run("ignores_legacy_golangci_config_path", func(t *testing.T) {
		_ = os.Unsetenv("LOGLINT_CONFIG_PATH")
		configPath = ""
		legacyConfigPath = ".golangci.yml"

		got, explicit := resolveConfiguredPath()
		if explicit {
			t.Fatal("resolveConfiguredPath() explicit = true, want false")
		}
		if got != "" {
			t.Fatalf("resolveConfiguredPath() path = %q, want empty", got)
		}
	})
}

func TestResolveConfigPath_UsesPWDForRelativePath(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "loglint.yaml")
	cfgYAML := `
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
`
	if err := os.WriteFile(cfgPath, []byte(cfgYAML), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	prevPWD, hadPWD := os.LookupEnv("PWD")
	_ = os.Setenv("PWD", dir)
	t.Cleanup(func() {
		if hadPWD {
			_ = os.Setenv("PWD", prevPWD)
			return
		}
		os.Unsetenv("PWD")
	})

	got := resolveConfigPath("./loglint.yaml", &analysis.Pass{})
	want, err := filepath.Abs(cfgPath)
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if got != want {
		t.Fatalf("resolveConfigPath() = %q, want %q", got, want)
	}
}
