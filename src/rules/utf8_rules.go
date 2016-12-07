// The two UTF8 checks are separated even though it would be far easier to have
// a single function - we may have cases where we want to skip the no-utf8
// check, but we probably never want to skip the invalid-utf8 check

package rules

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func init() {
	RegisterValidator("no-utf8", NoUTF8)
	RegisterValidatorHigh("invalid-utf8", InvalidUTF8)
}

// runeListErrorString converts a list of runes into a useful string for
// formatted error output
func runeListErrorString(rList []rune) string {
	var out []string
	for _, r := range rList {
		out = append(out, fmt.Sprintf(`"%c"`, r))
	}

	return strings.Join(out, ", ")
}

// runeValid tells us if a rune is a valid UTF-8 rune, and not utf8.RuneError
func runeValid(r rune) bool {
	return utf8.ValidRune(r) && r != utf8.RuneError
}

// NoUTF8 validates that the filename doesn't have "long" characters (e.g.,
// runes must be one byte UTF8 characters, which only encompasses low ASCII)
func NoUTF8(path string, info os.FileInfo) error {
	var utfRunes []rune
	for _, r := range info.Name() {
		if runeValid(r) && utf8.RuneLen(r) > 1 {
			utfRunes = append(utfRunes, r)
		}
	}

	if len(utfRunes) > 0 {
		return fmt.Errorf("contains unicode characters (%s)", runeListErrorString(utfRunes))
	}

	return nil
}

// InvalidUTF8 validates that all runes are valid UTF8 characters.  The message
// isn't terribly helpful, but it'll be hard to display which characters are
// invalid, given that they won't be displayable.
func InvalidUTF8(path string, info os.FileInfo) error {
	for _, r := range info.Name() {
		if !runeValid(r) {
			return fmt.Errorf("contains invalid unicode")
		}
	}

	return nil
}
