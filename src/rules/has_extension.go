package rules

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	RegisterValidator("has-extension", HasExtension)
}

// HasExtension verifies that an extension exists for regular files
func HasExtension(path string, info os.FileInfo) error {
	// Only regular files need extensions
	if !info.Mode().IsRegular() {
		return nil
	}

	if filepath.Ext(info.Name()) == "" {
		return fmt.Errorf("doesn't have an extension")
	}

	return nil
}
