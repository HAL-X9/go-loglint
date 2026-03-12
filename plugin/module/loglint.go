// Package module provides the golangci-lint module plugin for loglint.
package module

import (
	"github.com/HAL-X9/go-loglint/internal/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

func New(settings any) (register.LinterPlugin, error) {
	// settings.config is not passed reliably to module plugins and may contain
	// unrelated config values. Prefer LOGLINT_CONFIG_PATH for plugin mode.
	_ = settings
	return &loglintPlugin{}, nil
}

type loglintPlugin struct{}

func (p *loglintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *loglintPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
