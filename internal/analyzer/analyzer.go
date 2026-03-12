// Package analyzer validates log messages against defined rules.
package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/go-loglint/internal/config"
	"github.com/go-loglint/internal/rules"
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
	if configPath != "" {
		cfg, err := config.LoadConfig(configPath)
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

	arg, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}

	message := strings.Trim(arg.Value, `"`)
	report := func(pos token.Pos, format string, args ...interface{}) {
		pass.Reportf(pos, format, args...)
	}

	for _, rule := range rulesList {
		rule(message, call.Pos(), report)
	}
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
