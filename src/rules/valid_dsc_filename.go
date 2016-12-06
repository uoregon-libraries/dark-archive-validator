package rules

import (
	"fmt"
	"os"
	"strings"
)

// These cannot be part of any filename per DSC rules
var dscInvalidChars = []rune{'&', ',', '*', '%', '#', ';', '(', ')', '!', '@',
	'$', '^', '~', '\'', '{', '}', '[', ']', '\\', '?', '<', '>'}

func init() {
	RegisterValidator("valid-dsc-filename", ValidDSCFilename)
}

// ValidDSCFilename rejects files that use characters DSC disallows
func ValidDSCFilename(path string, info os.FileInfo) error {
	var badChars []rune

	var name = info.Name()
	for _, r := range dscInvalidChars {
		if strings.ContainsRune(name, r) {
			badChars = append(badChars, r)
		}
	}

	if len(badChars) > 0 {
		return fmt.Errorf("contains invalid characters: %s", joinRunes(badChars))
	}

	return nil
}
