package rules

import "testing"

func TestIsLowercase(t *testing.T) {
	tests := []struct {
		msg  string
		want bool
	}{
		{"", true},
		{"a", true},
		{"hello", true},
		{"A", false},
		{"Hello", false},
	}
	for _, tt := range tests {
		if got := isLowercase(tt.msg); got != tt.want {
			t.Errorf("isLowercase(%q) = %v, want %v", tt.msg, got, tt.want)
		}
	}
}

func TestIsEnglish(t *testing.T) {
	tests := []struct {
		msg  string
		want bool
	}{
		{"", true},
		{"hello", true},
		{"Hello world 123", true},
		{"café", false},
		{"привет", false},
		{"日本語", false},
	}
	for _, tt := range tests {
		if got := isEnglish(tt.msg); got != tt.want {
			t.Errorf("isEnglish(%q) = %v, want %v", tt.msg, got, tt.want)
		}
	}
}

func TestNoSpecialChars(t *testing.T) {
	tests := []struct {
		msg  string
		want bool
	}{
		{"", true},
		{"hello", true},
		{"hello world 123", true},
		{"hello\n", false},
		{"\thello", false},
		{"connection failed!!!", false},
		{"something went wrong...", false},
		{"hello!", false},
		{"really?", false},
		{"done.", false},
		{"error: failed", false},
		{"user@host", false},
	}
	for _, tt := range tests {
		if got := noSpecialChars(tt.msg); got != tt.want {
			t.Errorf("noSpecialChars(%q) = %v, want %v", tt.msg, got, tt.want)
		}
	}
}

func TestContainsSensitive(t *testing.T) {
	tests := []struct {
		msg  string
		want bool
	}{
		{"", false},
		{"hello", false},
		{"password", true},
		{"PASSWORD", true},
		{"api_key=xxx", true},
		{"token expired", true},
		{"secret key", true},
		{"my password is 123", true},
	}
	for _, tt := range tests {
		if got := containsSensitive(tt.msg); got != tt.want {
			t.Errorf("containsSensitive(%q) = %v, want %v", tt.msg, got, tt.want)
		}
	}
}
