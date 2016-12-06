package rules

import (
	"fmt"
	"os"
	"strings"
)

// VWFValidatorName is defined as a constant for use both here and in the
// Engine's skip check (this rule is hard-coded to not allow skipping)
const VWFValidatorName = "valid-windows-filename"

// These cannot be part of any filename in Windows
var winReservedChars = []rune{'<', '>', ':', '"', '/', '\\', '|', '?', '*'}

// These cannot be a file's name, nor a file's prefix (e.g., "CON.txt").  These
// are case-insensitive.
var winReservedFilenames = []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2",
	"COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3",
	"LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

func init() {
	RegisterValidator(VWFValidatorName, ValidWindowsFilename)
}

// ValidWindowsFilename is a ValidatorFunc which validates the that file's name
// matches Windows naming conventions.  The path itself is not validated here.
func ValidWindowsFilename(path string, info os.FileInfo) error {
	var badChars []rune
	var badName bool

	var name = strings.ToUpper(info.Name())
	for _, r := range winReservedChars {
		if strings.ContainsRune(name, r) {
			badChars = append(badChars, r)
		}
	}

	for _, badfname := range winReservedFilenames {
		if name == badfname {
			badName = true
			break
		}

		if strings.HasPrefix(name, badfname+".") {
			badName = true
		}
	}

	if len(badChars) > 0 {
		return fmt.Errorf("contains invalid characters: %s", joinRunes(badChars))
	}
	if badName {
		return fmt.Errorf("uses a reserved file name")
	}
	if strings.HasSuffix(name, " ") {
		return fmt.Errorf("has a trailing space")
	}
	if strings.HasSuffix(name, ".") {
		return fmt.Errorf("has a trailing period")
	}

	return nil
}
