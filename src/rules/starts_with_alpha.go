package rules

import (
	"fmt"
	"os"
)

func init() {
	RegisterValidator("starts-with-alpha", StartsWithAlpha)
}

// StartsWithAlpha enforces that a filename starts with ASCII A-Z or a-z
func StartsWithAlpha(path string, info os.FileInfo) error {
	var r = info.Name()[0]
	if r >= 'A' && r <= 'Z' {
		return nil
	}
	if r >= 'a' && r <= 'z' {
		return nil
	}

	return fmt.Errorf("starts with a non-alphabetic character")
}
