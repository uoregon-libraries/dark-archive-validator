package rules

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Failure keeps a validator and the error returned in one place for easy
// reporting of issues
type Failure struct {
	V Validator
	E error
}

// Engine is the rules runner
type Engine struct {
	Validators []Validator
	TraverseFn func(string, filepath.WalkFunc) error
}

// NewEngine returns an engine with the one required rule we have
func NewEngine() *Engine {
	var e = &Engine{TraverseFn: filepath.Walk}
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
func (e *Engine) ValidateTree(root string, failFunc func(string, []Failure)) {
	e.TraverseFn(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("CRITICAL - Unable to process %#v: %s", path, err)
		}

		var basepath = strings.Replace(path, root, "", -1)
		var fl = e.Validate(basepath, info)
		if len(fl) > 0 {
			failFunc(basepath, fl)
		}
		return nil
	})
}

// Validate checks the given full path against all validators, returning an
// array of errors found
func (e *Engine) Validate(path string, info os.FileInfo) []Failure {
	var flist []Failure

	for _, v := range e.Validators {
		var err = v.Validate(path, info)
		if err != nil {
			flist = append(flist, Failure{v, err})
		}
	}

	return flist
}
