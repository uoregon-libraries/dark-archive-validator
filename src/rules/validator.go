package rules

import (
	"os"
)

// ValidatorFunc is the function called by a validator to determine if a path
// is invalid in any way
type ValidatorFunc func(path string, info os.FileInfo) error

// A Validator is basically a named function which takes a full path to a file,
// and returns an error if any was found.  priority is used to order validators
// in an Engine's ValidationNames list.  Setting it below zero pushing a
// validator to run before the default, whereas setting it greater than zero
// makes validators run after the default validators run.  skipOnFail should be
// set to true if the validator shouldn't be run if other failures have already
// happened, allowing for validators which are really there just to catch
// unexpected problems.
type Validator struct {
	Name       string
	vf         ValidatorFunc
	priority   int8
	skipOnFail bool
}

// Validate checks for errors in the validator function and returns the
// (potentially updated) failure list
func (v Validator) Validate(path string, info os.FileInfo, fList []Failure) []Failure {
	if v.skipOnFail && len(fList) > 0 {
		return fList
	}

	var err = v.vf(path, info)
	if err != nil {
		return append(fList, Failure{v, err})
	}
	return fList
}

// ValidatorList encapsulators a slice of validators primarily for sorting
type ValidatorList []Validator
func (vl ValidatorList) Len() int { return len(vl) }
func (vl ValidatorList) Swap(i, j int) { vl[i], vl[j] = vl[j], vl[i] }

// Less is defined as having a numerically lower priority - or equal priority
// but alphabetically earlier name
func (vl ValidatorList) Less(i, j int) bool {
	if vl[i].priority == vl[j].priority {
		return vl[i].Name < vl[j].Name
	}

	return vl[i].priority < vl[j].priority
}

// validators holds all known validators
var validators ValidatorList

func register(v Validator) {
	validators = append(validators, v)
}

// RegisterValidator creates a simple validator with default priority and no
// skipping on failure, then puts it in the validator list
func RegisterValidator(name string, validate ValidatorFunc) {
	register(Validator{Name: name, vf: validate})
}

// RegisterCustomValidator creates a validator with explicitly set values for
// priority and skipOnFail, and puts that in the validator list
func RegisterCustomValidator(name string, validate ValidatorFunc, priority int8, skipOnFail bool) {
	register(Validator{name, validate, priority, skipOnFail})
}

// NukeValidatorList erases all entries from the list of known validators.
// This is primarily for custom use cases where a whitelist approach is
// preferable to the validators which auto-register themselves
func NukeValidatorList() {
	validators = nil
}
