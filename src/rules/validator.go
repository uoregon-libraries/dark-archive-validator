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

// validatorLookup holds all validators mapped by name
var validatorLookup = make(map[string]Validator)

// RegisterValidator maps a string name to the given validator
func RegisterValidator(name string, vf ValidatorFunc) {
	var v = Validator{name, vf}
	validatorLookup[name] = v
}

// NukeValidatorList erases all entries from the list of known validators.
// This is primarily for custom use cases where a whitelist approach is
// preferable to the validators which auto-register themselves
func NukeValidatorList() {
	validatorLookup = make(map[string]Validator)
}
