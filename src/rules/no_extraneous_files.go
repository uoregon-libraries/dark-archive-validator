package rules

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	RegisterValidator("no-extraneous-files", NoExtraneousFiles)
}

// NoExtraneousFiles validates a variety of file patterns to ensure various
// unnecessary file types aren't included, such as Thumbs.db, .DS_Store, etc.
func NoExtraneousFiles(path string, info os.FileInfo) error {
	var n = strings.ToUpper(info.Name())
	var genericError = fmt.Errorf("is an extraneous file and should be deleted")

	if n == ".DS_STORE" || n == "THUMBS.DB" || n == "DESKTOP.INI" {
		return genericError
	}

	if n[0] == '.' && n[1] == '_' {
		return genericError
	}

	return nil
}
