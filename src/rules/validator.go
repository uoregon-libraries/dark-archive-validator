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
// makes validators run after the default validators run.
// skipOnPreviousFailures should be set to true if the validator shouldn't be
// run if other failures have already happened, allowing for validators which
// are really there just to catch unexpected problems.  stopOnFailure should be
// set to true if this validator is expected to tell enough information that
// further validations are going to just confuse the report.
type Validator struct {
	Name     string
	vf       ValidatorFunc
	priority int8

	// skipOnPreviousFailures flags a validator not to run if there were previous failures
	skipOnPreviousFailures bool

	// stopOnFailure flags future validators not to run if this validator failed
	stopOnFailure bool
}

// Validate checks for errors in the validator function and returns the
// (potentially updated) failure list
func (v Validator) Validate(path string, info os.FileInfo, fList []Failure) []Failure {
	var l = len(fList)
	// If this validator isn't supposed to report already-failed items, break out
	// now if there are existing failures
	if v.skipOnPreviousFailures && l > 0 {
		return fList
	}

	// If the previous validator should stop all validations, break out now
	if l > 0 && fList[l-1].V.stopOnFailure {
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

func (vl ValidatorList) Len() int      { return len(vl) }
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
// priority and failure modes, and puts that in the validator list
func RegisterCustomValidator(name string, validate ValidatorFunc, priority int8, skipOnPreviousFailures, stopOnFailure bool) {
	register(Validator{name, validate, priority, skipOnPreviousFailures, stopOnFailure})
}

// NukeValidatorList erases all entries from the list of known validators.
// This is primarily for custom use cases where a whitelist approach is
// preferable to the validators which auto-register themselves
func NukeValidatorList() {
	validators = nil
}
