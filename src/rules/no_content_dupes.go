package rules

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/uoregon-libraries/dark-archive-validator/src/checksum"
)

// RegisterChecksumValidator is trickier than most validators - this validator
// needs a full path as context since validators only get the path relative to
// the root.  This thing also has to have a checksummer and a map to store
// checksums in the first place.  And since checksumming is optional, we don't
// want to auto-register any particular checksum validator.  So this function
// gets all that context, builds a validator function closure and registers it.
func RegisterChecksumValidator(root string, c *checksum.Checksum, checksums map[string][]string) {
	var validateChecksum = func(path string, info os.FileInfo) error {
		// Don't try to checksum non-files
		if !info.Mode().IsRegular() {
			return nil
		}

		var fullPath = filepath.Join(root, path)
		var sum, err = c.Sum(fullPath)
		if err != nil && err != io.EOF {
			return fmt.Errorf("isn't able to be checksummed (%s)", err)
		}

		var chksum = fmt.Sprintf("%x", sum)
		var chksumExist = checksums[chksum]
		if len(chksumExist) != 0 {
			err = fmt.Errorf("duplicates the content of %#v", checksums[chksum][0])
		}

		checksums[chksum] = append(checksums[chksum], fullPath)
		return err
	}
	RegisterValidatorHigh("no-duped-content", validateChecksum)
}
