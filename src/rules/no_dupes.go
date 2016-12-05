package rules

import (
	"fmt"
	"os"
	"strings"
)

var nameLookup = make(map[string]string)

func init() {
	RegisterValidator("no-duped-names", NoDupedNames)
}

// ResetDupemap clears the name lookup to allow for running multiple trees
// without fear of a false dupe warning
func ResetDupemap() {
	nameLookup = make(map[string]string)
}

func NoDupedNames(path string, info os.FileInfo) error {
	var pathUpper = strings.ToUpper(path)
	if nameLookup[pathUpper] != "" {
		return fmt.Errorf("is a duplicate of %#v", nameLookup[pathUpper])
	}

	nameLookup[pathUpper] = path
	return nil
}