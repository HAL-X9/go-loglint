package rules

import "strings"

func containsSensitive(msg string) bool {
	keywords := []string{"password", "api_key", "token", "secret"}
	for _, kw := range keywords {
		if strings.Contains(strings.ToLower(msg), kw) {
			return true
		}
	}
	return false
}
