package rules

import (
	"fmt"
	"os"
)

func init() {
	RegisterValidator("no-hidden-files", NoHiddenFiles)
}

// NoHiddenFiles verifies that the file doesn't start with "." - TODO: need to
// read attrs for Windows files, too
func NoHiddenFiles(path string, info os.FileInfo) error {
	if info.Name()[0] == '.' {
		return fmt.Errorf("is hidden (starts with a period)")
	}

	return nil
}
