package analyzer

import (
	"os"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	prevConfigPath := configPath
	prevEnv, envWasSet := os.LookupEnv("LOGLINT_CONFIG_PATH")
	configPath = ""
	os.Unsetenv("LOGLINT_CONFIG_PATH")
	t.Cleanup(func() {
		configPath = prevConfigPath
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
	prevEnv, envWasSet := os.LookupEnv("LOGLINT_CONFIG_PATH")
	configPath = ""
	_ = os.Setenv("LOGLINT_CONFIG_PATH", "/definitely/missing/loglint.yaml")
	t.Cleanup(func() {
		configPath = prevConfigPath
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
