package rules

import (
	"os"
	"path/filepath"
)

// Failure keeps a validator and the error returned in one place for easy
// reporting of issues
type Failure struct {
	V Validator
	E error
}

// FailureList acts as a proxy into a slice of failures to allow for "magic"
// append capabilities
type FailureList struct {
	list []Failure
}

// AppendIfError adds a Failure to the list if v returns an error when run
// against info
func (fl FailureList) AppendIfError(v Validator, path string, info os.FileInfo) {
	var err = v.Validate(path, info)
	if err != nil {
		fl.list = append(fl.list, Failure{v, err})
	}
}

// Any returns true if the failure list has any elements, which only happens
// when AppendIfError finds an error from a validator
func (fl FailureList) Any() bool {
	return len(fl.list) > 0
}

// Engine is the rules runner
type Engine struct {
	Validators []Validator
}

// NewEngine returns an engine with the one required rule we have
func NewEngine() *Engine {
	var e = &Engine{}
	e.AddValidator("valid-windows-filename", ValidWindowsFilename)
	return e
}

// AddValidator maps a name to a validator
func (e *Engine) AddValidator(name string, v ValidatorFunc) {
	e.Validators = append(e.Validators, Validator{name, v})
}

// ValidateTree walks all files under root, sending everything found to all
// registered validators, yielding to failFunc whenever a validation against a
// file returns any errors
func (e *Engine) ValidateTree(root string, failFunc func(string, FailureList)) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		var fl = e.Validate(path, info)
		if fl.Any() {
			failFunc(path, fl)
		}
		return nil
	})
}

// Validate checks the given full path against all validators, returning an
// array of errors found
func (e *Engine) Validate(path string, info os.FileInfo) FailureList {
	var flist FailureList

	for _, v := range e.Validators {
		flist.AppendIfError(v, path, info)
	}

	return flist
}
