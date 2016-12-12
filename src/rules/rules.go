package rules

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Failure keeps a validator and the error returned in one place for easy
// reporting of issues
type Failure struct {
	V Validator
	E error
}

// badFileValidator is a hard-coded validator with no function just for
// reporting when a file can't be processed by the walk function
var badFileValidator = Validator{
	Name:        "broken-file",
	priority:    -128,
	vf:          nil,
	Criticality: CCritical,
}

func init() {
	register(badFileValidator)
}

// Engine is the rules runner.  By default it will run all known validators
// except those explicitly skipped.
type Engine struct {
	TraverseFn func(string, filepath.WalkFunc) error
	skip       map[string]bool
}

// NewEngine returns an engine with the one required rule we have
func NewEngine() *Engine {
	return &Engine{
		TraverseFn: filepath.Walk,
		skip:       make(map[string]bool),
	}
}

// Skip looks up the given validator by name and, if it exists, adds it to this
// engine's validator skip list
//
// Note that this will NEVER remove critical checks, as those rules are in
// place so the dark archive filesystem works properly
func (e *Engine) Skip(name string) (ok bool) {
	for _, v := range validators {
		if v.Name == name && v.Criticality > CCritical {
			e.skip[name] = true
			return true
		}
	}

	return false
}

// Unskip removes a "skip" assigned to a named validator, returning true if it
// existed and was "unskipped", false otherwise.
func (e *Engine) Unskip(name string) bool {
	if e.skip[name] {
		e.skip[name] = false
		return true
	}

	return false
}

// SkipAll is to be used solely when manually setting up (via Unskip)
// individual validators for precise use-cases.  As with Skip(), this will not
// remove the Windows filename restriction validator.
func (e *Engine) SkipAll() {
	for _, v := range validators {
		e.Skip(v.Name)
	}
}

// ValidateTree walks all files under root, sending everything found to all
// registered validators, yielding to failFunc whenever a validation against a
// file returns any errors
func (e *Engine) ValidateTree(root string, failFunc func(string, []Failure)) {
	if root[len(root)-1] != filepath.Separator {
		root += string(filepath.Separator)
	}
	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		log.Fatalf("Unable to read path %#v: %s", root, err)
	}

	e.TraverseFn(root, func(path string, info os.FileInfo, err error) error {
		var basepath = strings.Replace(path, root, "", 1)

		if err != nil {
			var fl = make([]Failure, 1)
			fl[0] = Failure{V: badFileValidator, E: fmt.Errorf("critical error: %s", err)}
			failFunc(basepath, fl)
			return nil
		}

		// The root filename doesn't matter, since our goal is to validate the
		// contents of root, and then move them to the *real* dark archive root
		if root == path {
			return nil
		}

		var fl = e.Validate(basepath, info)
		if len(fl) > 0 {
			failFunc(basepath, fl)
		}

		return nil
	})
}

// Validators returns a sorted list of all validators which are not explicitly
// skipped - though the Windows filename restrictions are forcibly added to the
// list no matter what.  We sort by priority and then name in order to allow
// priority/skip-on-fail to make sense, and to ensure consistent reporting.
func (e *Engine) Validators() ValidatorList {
	var vList ValidatorList
	var v Validator

	for _, v = range validators {
		if e.skip[v.Name] {
			continue
		}
		vList = append(vList, v)
	}

	sort.Sort(vList)
	return vList
}

// Validate checks the given base path against all validators not in the skip
// list, and returns an array of errors found
func (e *Engine) Validate(basepath string, info os.FileInfo) []Failure {
	var flist []Failure

	var v Validator
	for _, v = range e.Validators() {
		flist = v.Validate(basepath, info, flist)
	}

	return flist
}
