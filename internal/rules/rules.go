package rules

import "go/token"

type RuleFunc func(msg string, pos token.Pos, report func(pos token.Pos, format string, args ...interface{}))

func FromConfig(enabledRules map[string]bool) []RuleFunc {
	var rules []RuleFunc

	if enabledRules["lowercase_start"] {
		rules = append(rules, func(msg string, pos token.Pos, report func(pos token.Pos, format string, args ...interface{})) {
			if !isLowercase(msg) {
				report(pos, "log message must start with lowercase")
			}
		})
	}

	if enabledRules["english_only"] {
		rules = append(rules, func(msg string, pos token.Pos, report func(pos token.Pos, format string, args ...interface{})) {
			if !isEnglish(msg) {
				report(pos, "log message must be in English")
			}
		})
	}

	if enabledRules["no_special_chars"] {
		rules = append(rules, func(msg string, pos token.Pos, report func(pos token.Pos, format string, args ...interface{})) {
			if !noSpecialChars(msg) {
				report(pos, "log message must not contain special characters or emojis")
			}
		})
	}

	if enabledRules["no_sensitive_data"] {
		rules = append(rules, func(msg string, pos token.Pos, report func(pos token.Pos, format string, args ...interface{})) {
			if containsSensitive(msg) {
				report(pos, "log message must not contain sensitive data")
			}
		})
	}

	return rules
}
