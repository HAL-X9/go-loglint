// Package module provides the golangci-lint module plugin for loglint.
package module

import (
	"github.com/HAL-X9/go-loglint/internal/analyzer"
	"github.com/HAL-X9/go-loglint/internal/config"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

func New(settings any) (register.LinterPlugin, error) {
	if path := config.GetPathFromGolangciLintSettings(settings); path != "" {
		analyzer.SetConfigPath(path)
	}
	return &loglintPlugin{}, nil
}

type loglintPlugin struct{}

func (p *loglintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *loglintPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
