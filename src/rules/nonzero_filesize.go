package rules

import (
	"fmt"
	"os"
)

func init() {
	RegisterValidator("nonzero-filesize", NonzeroFilesize)
}

// NonzeroFilesize enforces that all regular files are at least 1 byte
func NonzeroFilesize(path string, info os.FileInfo) error {
	if info.Size() == 0 && info.Mode().IsRegular() {
		return fmt.Errorf("is an empty file")
	}

	return nil
}
