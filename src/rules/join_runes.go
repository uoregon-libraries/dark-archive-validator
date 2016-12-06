package rules

// joinRunes creates a human-readable string listing the given runes
func joinRunes(rList []rune) string {
	var out string
	for i, r := range rList {
		out += string(r)
		if i != len(rList) - 1 {
			out += " "
		}
	}

	return out
}
