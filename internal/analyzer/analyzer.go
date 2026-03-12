// Package analyzer validates log messages against defined rules.
package analyzer

import (
	"go/ast"
	"go/token"
	"os"
	"strings"

	"github.com/loglint/loglint/internal/config"
	"github.com/loglint/loglint/internal/rules"
	"golang.org/x/tools/go/analysis"
)

var configPath string

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", "", "path to YAML config file (optional, all rules enabled by default)")
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
	path := configPath
	if path == "" {
		path = os.Getenv("LOGLINT_CONFIG_PATH")
	}
	if path != "" {
		cfg, err := config.LoadConfig(path)
		if err != nil {
			return nil, err
		}
		enabledRules = make(map[string]bool)
		for _, r := range cfg.Rules {
			enabledRules[r.Name] = r.Enable
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
