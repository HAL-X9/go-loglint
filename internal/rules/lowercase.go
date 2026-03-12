package rules

import "unicode"

func isLowercase(msg string) bool {
	if msg == "" {
		return true
	}

	r := []rune(msg)[0]
	return unicode.IsLower(r)
}
