package rules

import (
	"fmt"
	"os"
	"unicode"
)

func init() {
	RegisterValidator("no-spaces", NoSpaces)
}

// NoSpaces verifies that the file described by info has no spaces.  Trailing
// spaces are reported with a special message since those can be a pain to
// notice.
func NoSpaces(path string, info os.FileInfo) error {
	var name = info.Name()
	var hasSpace bool
	var spaceAtEnd bool
	for _, r := range name {
		spaceAtEnd = false
		if unicode.IsSpace(r) {
			hasSpace = true
			spaceAtEnd = true
		}
	}

	if spaceAtEnd {
		return fmt.Errorf("ends with a space")
	}

	if hasSpace {
		return fmt.Errorf("has a space in the filename")
	}

	return nil
}
