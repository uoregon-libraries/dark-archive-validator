package rules

import (
	"fmt"
	"os"
)

func init() {
	RegisterCustomValidator("no-extraneous-files", NoExtraneousFiles, CNormal, -100, false, true)
}

// NoExtraneousFiles validates a variety of file patterns to ensure various
// unnecessary file types aren't included, such as Thumbs.db, .DS_Store, etc.
func NoExtraneousFiles(path string, info os.FileInfo) error {
	var n = info.Name()
	var genericError = fmt.Errorf("may be an extraneous file; consider deletion")

	if n == ".DS_Store" || n == "Thumbs.db" || n == "desktop.ini" {
		return genericError
	}

	if n[0] == '.' && n[1] == '_' {
		return genericError
	}

	return nil
}
