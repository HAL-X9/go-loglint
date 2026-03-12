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

type loglintSettings struct {
	Config string `json:"config"`
}

type loglintSettingsNested struct {
	Settings struct {
		Config string `json:"config"`
	} `json:"settings"`
}

func New(settings any) (register.LinterPlugin, error) {
	// Try flat structure first: {config: "..."}
	if s, err := register.DecodeSettings[loglintSettings](settings); err == nil && s.Config != "" {
		analyzer.SetConfigPath(s.Config)
		return &loglintPlugin{}, nil
	}
	// Try nested: {settings: {config: "..."}}
	if s, err := register.DecodeSettings[loglintSettingsNested](settings); err == nil && s.Settings.Config != "" {
		analyzer.SetConfigPath(s.Settings.Config)
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
