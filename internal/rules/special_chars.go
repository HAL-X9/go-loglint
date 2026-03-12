package rules

func noSpecialChars(msg string) bool {
	for _, r := range msg {
		if !(r >= 32 && r <= 126) {
			return false
		}
	}
	return true
}
