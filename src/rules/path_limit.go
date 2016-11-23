package rules

import (
	"fmt"
	"os"
)

// PathLimitFn returns a validator function which will report when the path
// exceeds n characters
func PathLimitFn(n int) ValidatorFunc {
	return func (path string, info os.FileInfo) error {
		if len(path) > n {
			return fmt.Errorf("exceeds the maximum path length of %d characters", n)
		}
		return nil
	}
}
