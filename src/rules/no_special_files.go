package rules

import (
	"fmt"
	"os"
)

func init() {
	RegisterValidator("no-special-files", NoSpecialFiles)
}

// NoSpecialFiles verifies that the file described by info is a regular file or
// a directory, reporting problems otherwise
func NoSpecialFiles(path string, info os.FileInfo) error {
	var m = info.Mode()

	if m.IsRegular() {
		return nil
	}

	if m.IsDir() {
		return nil
	}

	if m&os.ModeSymlink != 0 {
		return fmt.Errorf("is a symbolic link")
	}

	if m&os.ModeDevice != 0 {
		return fmt.Errorf("is a device file")
	}

	if m&os.ModeNamedPipe != 0 {
		return fmt.Errorf("is a named pipe")
	}

	if m&os.ModeSocket != 0 {
		return fmt.Errorf("is a socket")
	}

	return fmt.Errorf("is not a regular file or folder")
}
