package config

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr string
	}{
		{
			name:    "nil_config",
			cfg:     nil,
			wantErr: "config is nil",
		},
		{
			name: "valid_full_config",
			cfg: &Config{
				AutoFix: true,
				Rules: []Rule{
					{Name: "lowercase_start", Enable: true},
					{Name: "english_only", Enable: true},
				},
				Loggers: Loggers{
					Include: []string{"log", "slog"},
				},
			},
			wantErr: "",
		},
		{
			name: "invalid_rules_wrapped",
			cfg: &Config{
				AutoFix: true,
				Rules: []Rule{
					{Name: "unknown_rule", Enable: true},
				},
				Loggers: Loggers{
					Include: []string{"log"},
				},
			},
			wantErr: "rule:",
		},
		{
			name: "invalid_loggers_wrapped",
			cfg: &Config{
				AutoFix: true,
				Rules: []Rule{
					{Name: "lowercase_start", Enable: true},
				},
				Loggers: Loggers{
					Include: []string{"invalid_logger"},
				},
			},
			wantErr: "loggers:",
		},
		{
			name: "empty_rules",
			cfg: &Config{
				AutoFix: true,
				Rules:   []Rule{},
				Loggers: Loggers{
					Include: []string{"log"},
				},
			},
			wantErr: "rule:",
		},
		{
			name: "nil_rules",
			cfg: &Config{
				AutoFix: true,
				Rules:   nil,
				Loggers: Loggers{
					Include: []string{"log"},
				},
			},
			wantErr: "rule:",
		},
		{
			name: "empty_loggers",
			cfg: &Config{
				AutoFix: true,
				Rules: []Rule{
					{Name: "lowercase_start", Enable: true},
				},
				Loggers: Loggers{
					Include: []string{},
				},
			},
			wantErr: "loggers:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Errorf("Validate() expected error containing %q, got nil", tt.wantErr)
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRule(t *testing.T) {
	tests := []struct {
		name    string
		rules   []Rule
		wantErr string
	}{
		{
			name:    "nil_rules",
			rules:   nil,
			wantErr: "no rules defined in the configuration",
		},
		{
			name:    "empty_rules",
			rules:   []Rule{},
			wantErr: "no rules defined in the configuration",
		},
		{
			name: "valid_lowercase_start",
			rules: []Rule{
				{Name: "lowercase_start", Enable: true},
			},
			wantErr: "",
		},
		{
			name: "valid_english_only",
			rules: []Rule{
				{Name: "english_only", Enable: false},
			},
			wantErr: "",
		},
		{
			name: "valid_no_special_chars",
			rules: []Rule{
				{Name: "no_special_chars", Enable: true},
			},
			wantErr: "",
		},
		{
			name: "valid_no_sensitive_data",
			rules: []Rule{
				{Name: "no_sensitive_data", Enable: true},
			},
			wantErr: "",
		},
		{
			name: "valid_all_rules",
			rules: []Rule{
				{Name: "lowercase_start", Enable: true},
				{Name: "english_only", Enable: true},
				{Name: "no_special_chars", Enable: true},
				{Name: "no_sensitive_data", Enable: false},
			},
			wantErr: "",
		},
		{
			name: "invalid_rule_name",
			rules: []Rule{
				{Name: "unknown_rule", Enable: true},
			},
			wantErr: "rule name must be",
		},
		{
			name: "empty_rule_name",
			rules: []Rule{
				{Name: "", Enable: true},
			},
			wantErr: "rule name must be",
		},
		{
			name: "first_invalid_stops_validation",
			rules: []Rule{
				{Name: "lowercase_start", Enable: true},
				{Name: "bad_rule", Enable: true},
				{Name: "english_only", Enable: true},
			},
			wantErr: "rule name must be",
		},
		{
			name: "typo_in_rule_name",
			rules: []Rule{
				{Name: "lowercase_starts", Enable: true},
			},
			wantErr: "rule name must be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRule(tt.rules)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("ValidateRule() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Errorf("ValidateRule() expected error containing %q, got nil", tt.wantErr)
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ValidateRule() error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLoggers(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Loggers
		wantErr string
	}{
		{
			name:    "nil_config",
			cfg:     nil,
			wantErr: "no loggers defined in the configuration",
		},
		{
			name: "empty_include",
			cfg: &Loggers{
				Include: []string{},
			},
			wantErr: "no loggers defined in the configuration",
		},
		{
			name: "valid_log",
			cfg: &Loggers{
				Include: []string{"log"},
			},
			wantErr: "",
		},
		{
			name: "valid_slog",
			cfg: &Loggers{
				Include: []string{"slog"},
			},
			wantErr: "",
		},
		{
			name: "valid_logrus_logger",
			cfg: &Loggers{
				Include: []string{"logrus.Logger"},
			},
			wantErr: "",
		},
		{
			name: "valid_zap_sugared_logger",
			cfg: &Loggers{
				Include: []string{"zap.SugaredLogger"},
			},
			wantErr: "",
		},
		{
			name: "valid_all_loggers",
			cfg: &Loggers{
				Include: []string{"log", "slog", "logrus.Logger", "zap.SugaredLogger"},
			},
			wantErr: "",
		},
		{
			name: "invalid_logger_name",
			cfg: &Loggers{
				Include: []string{"logrus"},
			},
			wantErr: "invalid logger name",
		},
		{
			name: "invalid_logger_includes_allowed",
			cfg: &Loggers{
				Include: []string{"log", "invalid", "slog"},
			},
			wantErr: "invalid logger name",
		},
		{
			name: "wrong_case_log",
			cfg: &Loggers{
				Include: []string{"Log"},
			},
			wantErr: "invalid logger name",
		},
		{
			name: "empty_string_logger",
			cfg: &Loggers{
				Include: []string{""},
			},
			wantErr: "invalid logger name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLoggers(tt.cfg)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("ValidateLoggers() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Errorf("ValidateLoggers() expected error containing %q, got nil", tt.wantErr)
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ValidateLoggers() error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_error_wrapping(t *testing.T) {
	cfg := &Config{
		Rules: []Rule{
			{Name: "bad_rule", Enable: true},
		},
		Loggers: Loggers{
			Include: []string{"log"},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Validate() expected error, got nil")
	}

	// Verify error chain: top-level contains "rule:", unwrapped contains rule validation message
	if !strings.Contains(err.Error(), "rule:") {
		t.Errorf("Validate() error should wrap rule error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "rule name must be") {
		t.Errorf("Validate() error should contain underlying rule error, got: %v", err)
	}
}

func TestValidateLoggers_error_message_format(t *testing.T) {
	cfg := &Loggers{
		Include: []string{"wrong_logger"},
	}

	err := ValidateLoggers(cfg)
	if err == nil {
		t.Fatal("ValidateLoggers() expected error, got nil")
	}

	// Error must include the invalid value in quotes
	if !strings.Contains(err.Error(), `"wrong_logger"`) {
		t.Errorf("ValidateLoggers() error should quote invalid logger name, got: %v", err)
	}
	if !strings.Contains(err.Error(), "allowed:") {
		t.Errorf("ValidateLoggers() error should list allowed values, got: %v", err)
	}
}
