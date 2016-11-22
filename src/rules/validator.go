package rules

import (
	"os"
)

// ValidatorFunc is the function called by a validator to determine if a path
// is invalid in any way
type ValidatorFunc func(path string, info os.FileInfo) error

// A Validator is basically a named function which takes a full path to a file,
// and returns an error if any was found
type Validator struct {
	Name string
	vf   ValidatorFunc
}

// Validate simply delegates the filename into the validator function
func (v Validator) Validate(path string, info os.FileInfo) error {
	return v.vf(path, info)
}
