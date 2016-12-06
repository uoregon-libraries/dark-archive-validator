package rules

import (
	"fmt"
	"os"
	"regexp"
)

var allowedFilename = regexp.MustCompile(`\A[A-Za-z][A-Za-z0-9_-]*\.[A-Za-z0-9_-]+\z`)
var allowedDirname = regexp.MustCompile(`\A[A-Za-z][A-Za-z0-9_-]*(\.[A-Za-z0-9_-]+)?\z`)

func init() {
	RegisterCustomValidator("restrictive-naming", RestrictiveNaming, 100, true, false)
}

// RestrictiveNaming enforces that only VERY specific whitelisted characters
// are allowed in a filename.  This one doesn't give terribly useful reporting,
// and may need to be skipped to accommodate project-specific trees over which
// we have limited control.
func RestrictiveNaming(path string, info os.FileInfo) error {
	var re *regexp.Regexp
	var errType string

	// Dirs have one pattern, files have another, and other objects are ignored
	// for this test
	if info.Mode().IsDir() {
		re = allowedDirname
		errType = "directory"
	} else if info.Mode().IsRegular() {
		re = allowedFilename
		errType = "filename"
	} else {
		return nil
	}

	if re.MatchString(info.Name()) {
		return nil
	}

	return fmt.Errorf("doesn't match required %s pattern", errType)
}
