// Package analyzer validates log messages against defined rules.
package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer checks log messages against defined rules.
var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages against defined rules",
	Run:  run,
}

// run walks the AST and checks function calls for log message violations.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {

			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			checkLogCall(pass, call)

			return true
		})
	}

	return nil, nil
}

// checkLogCall validates a log call's first argument against message rules.
func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {

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

	if len(message) > 0 && message[0] >= 'A' && message[0] <= 'Z' {
		pass.Reportf(call.Pos(), "log message must start with lowercase")
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
