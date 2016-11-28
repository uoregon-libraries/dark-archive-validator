package rules

import (
	"fmt"
	"os"
)

// NoControlChars rejects files that use anything below ASCII space, or the
// delete code (0-31, 127)
func NoControlChars(path string, info os.FileInfo) error {
	var name = info.Name()
	for _, r := range name {
		if r < 32 || r == 127 {
			return fmt.Errorf("contains one or more control characters")
		}
	}

	return nil
}
