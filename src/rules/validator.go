package rules

import (
	"os"
)

// Criticality defines how important a validator is
type Criticality int

// Criticalities, sorted such that low-importance items have a higher number to
// ensure they're sorted after high-importance items.  Normal is explicitly set
// to zero since that's the int default.
const (
	CCritical Criticality = -2
	CHigh     Criticality = -1
	CNormal   Criticality = 0
	CLow      Criticality = 1
)

func (c Criticality) String() string {
	switch c {
	case CCritical:
		return "Critical"
	case CHigh:
		return "High"
	case CNormal:
		return "Normal"
	case CLow:
		return "Low"
	default:
		return "UNKNOWN"
	}
}

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
	Name        string
	vf          ValidatorFunc
	priority    int8
	Criticality Criticality

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

// IsImportant reports whether this validator should be considered necessary
// even in quick runs.  We define this as being > low criticality or any
// criticality with stopOnFailure behavior, as those tests can prevent other
// less meaningful tests from running.
func (v Validator) IsImportant() bool {
	return v.Criticality < CLow || v.stopOnFailure
}

// ValidatorList encapsulators a slice of validators primarily for sorting
type ValidatorList []Validator

func (vl ValidatorList) Len() int      { return len(vl) }
func (vl ValidatorList) Swap(i, j int) { vl[i], vl[j] = vl[j], vl[i] }

// i is less than j if i is prioritized earlier (numerically lower priority).
// If both have the same priority, criticality is then considered, which VLHigh
// being sorted first, etc.  If priority and criticality are the same, sorting
// is alphabetic by name.
func (vl ValidatorList) Less(i, j int) bool {
	if vl[i].priority != vl[j].priority {
		return vl[i].priority < vl[j].priority
	}

	if vl[i].Criticality != vl[j].Criticality {
		return vl[i].Criticality < vl[j].Criticality
	}

	return vl[i].Name < vl[j].Name
}

// validators holds all known validators
var validators ValidatorList

func register(v Validator) {
	validators = append(validators, v)
}

// RegisterValidator creates a simple validator with default criticality and
// priority, and no skipping on failure, then puts it in the validator list
func RegisterValidator(name string, validate ValidatorFunc) {
	register(Validator{Name: name, vf: validate})
}

// RegisterValidatorCritical registers a critical validator
func RegisterValidatorCritical(name string, validate ValidatorFunc) {
	register(Validator{Name: name, vf: validate, Criticality: CCritical})
}

// RegisterValidatorHigh registers a high-criticality validator
func RegisterValidatorHigh(name string, validate ValidatorFunc) {
	register(Validator{Name: name, vf: validate, Criticality: CHigh})
}

// RegisterValidatorLow registers a low-criticality validator
func RegisterValidatorLow(name string, validate ValidatorFunc) {
	register(Validator{Name: name, vf: validate, Criticality: CLow})
}

// RegisterCustomValidator creates a validator with explicitly set values for
// priority and failure modes, and puts that in the validator list
func RegisterCustomValidator(name string, validate ValidatorFunc, c Criticality, priority int8, skipOnPreviousFailures, stopOnFailure bool) {
	register(Validator{name, validate, priority, c, skipOnPreviousFailures, stopOnFailure})
}

// NukeValidatorList erases all entries from the list of known validators.
// This is primarily for custom use cases where a whitelist approach is
// preferable to the validators which auto-register themselves
func NukeValidatorList() {
	validators = nil
}
