// Package analyzer validates log messages against defined rules.
package analyzer

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/HAL-X9/go-loglint/internal/config"
	"github.com/HAL-X9/go-loglint/internal/rules"
	"golang.org/x/tools/go/analysis"
)

var configPath string
var legacyConfigPath string

func init() {
	Analyzer.Flags.StringVar(&configPath, "loglint-config", "", "path to loglint YAML config file (optional, all rules enabled by default)")
	Analyzer.Flags.StringVar(&legacyConfigPath, "config", "", "deprecated: use -loglint-config")
}

// SetConfigPath sets the config file path (used by golangci-lint plugin).
func SetConfigPath(p string) {
	configPath = p
}

// Analyzer checks log messages against defined rules.
var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages against defined rules",
	Run:  run,
}

// run walks the AST and checks function calls for log message violations.
func run(pass *analysis.Pass) (interface{}, error) {
	enabledRules := defaultEnabledRules()
	path, explicitConfig := resolveConfiguredPath()
	if explicitConfig {
		resolved := resolveConfigPath(path, pass)
		cfg, err := config.LoadConfig(resolved)
		if err != nil {
			// If user explicitly provides a config, do not silently fall back
			// to defaults on load/validation errors.
			enabledRules = map[string]bool{}
		} else {
			enabledRules = make(map[string]bool)
			for _, r := range cfg.Rules {
				enabledRules[r.Name] = r.Enable
			}
		}
	}

	rulesList := rules.FromConfig(enabledRules)

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			checkLogCall(pass, call, rulesList)
			return true
		})
	}

	return nil, nil
}

func resolveConfiguredPath() (string, bool) {
	if envPath := strings.TrimSpace(os.Getenv("LOGLINT_CONFIG_PATH")); envPath != "" {
		return envPath, true
	}
	if path := strings.TrimSpace(configPath); path != "" {
		return path, true
	}
	if path := strings.TrimSpace(legacyConfigPath); path != "" {
		// golangci-lint itself uses -config for .golangci.yml.
		// Ignore that value for legacy support to avoid false config loads.
		if looksLikeGolangCIConfig(path) {
			return "", false
		}
		return path, true
	}
	return "", false
}

func looksLikeGolangCIConfig(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	switch base {
	case ".golangci.yml", ".golangci.yaml", ".golangci.toml", ".golangci.json",
		"golangci.yml", "golangci.yaml", "golangci.toml", "golangci.json":
		return true
	default:
		return false
	}
}

func defaultEnabledRules() map[string]bool {
	return map[string]bool{
		"lowercase_start":   true,
		"english_only":      true,
		"no_special_chars":  true,
		"no_sensitive_data": true,
	}
}

// checkLogCall validates a log call's first argument against message rules.
func checkLogCall(pass *analysis.Pass, call *ast.CallExpr, rulesList []rules.RuleFunc) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	if !isLoggerCall(sel) {
		return
	}

	if len(call.Args) == 0 {
		return
	}

	messages := extractStrings(call.Args[0])
	if len(messages) == 0 {
		return
	}

	report := func(pos token.Pos, format string, args ...interface{}) {
		pass.Reportf(pos, format, args...)
	}

	for _, message := range messages {
		for _, rule := range rulesList {
			rule(message, call.Pos(), report)
		}
	}
}

// extractStrings collects string literals from an expression (handles concatenation).
func extractStrings(expr ast.Expr) []string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return []string{strings.Trim(e.Value, `"`)}
		}
		return nil
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return nil
		}
		return append(extractStrings(e.X), extractStrings(e.Y)...)
	}
	return nil
}

// resolveConfigPath resolves relative path; returns absolute path when found.
func resolveConfigPath(path string, pass *analysis.Pass) string {
	if filepath.IsAbs(path) && fileExists(path) {
		return path
	}
	// Try CWD first
	if fileExists(path) {
		if abs, err := filepath.Abs(path); err == nil {
			return abs
		}
		return path
	}
	// Walk up from first analyzed file to find config
	if len(pass.Files) > 0 {
		fpath := pass.Fset.Position(pass.Files[0].Pos()).Filename
		dir := filepath.Dir(fpath)
		for dir != "" && dir != "." {
			p := filepath.Join(dir, path)
			if fileExists(p) {
				if abs, err := filepath.Abs(p); err == nil {
					return abs
				}
				return p
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return path
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// isLoggerCall reports whether sel is a logging call (Info, Warn, Error, etc.).
func isLoggerCall(sel *ast.SelectorExpr) bool {
	method := sel.Sel.Name

	switch method {
	case "Print", "Printf", "Println":
		if ident, ok := sel.X.(*ast.Ident); ok {
			if ident.Name == "fmt" {
				return false
			}
		}
		return true

	case "Info", "Warn", "Error", "Debug", "Fatal":
		return true
	}

	return false
}
