package rules

import "unicode"

func noSpecialChars(msg string) bool {
	for _, r := range msg {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != ' ' {
			return false
		}
	}
	return true
}
