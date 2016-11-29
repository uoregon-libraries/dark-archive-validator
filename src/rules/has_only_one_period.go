package rules

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	RegisterValidator("has-only-one-period", HasOnlyOnePeriod)
}

// HasOnlyOnePeriod enforces the rule that we can have one extension, but no
// periods in the "base" of the filename
func HasOnlyOnePeriod(path string, info os.FileInfo) error {
	var c = strings.Count(info.Name(), ".")
	if c > 1 {
		return fmt.Errorf("has %d periods (maximum is 1)", c)
	}

	return nil
}
