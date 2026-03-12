package rules

func isEnglish(msg string) bool {
	for _, r := range msg {
		if r > 127 {
			return false
		}
	}

	return true
}
