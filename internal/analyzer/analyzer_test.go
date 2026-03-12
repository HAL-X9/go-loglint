package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
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
