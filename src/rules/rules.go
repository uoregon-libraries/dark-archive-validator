package rules

import (
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

// Engine is the rules runner.  By default it will run all known validators
// except those explicitly skipped.
type Engine struct {
	TraverseFn func(string, filepath.WalkFunc) error
	skip       map[string]bool
}

// NewEngine returns an engine with the one required rule we have
func NewEngine() *Engine {
	return &Engine{TraverseFn: filepath.Walk, skip: make(map[string]bool)}
}

// Skip looks up the given validator by name and, if it exists, adds it to this
// engine's validator skip list
func (e *Engine) Skip(name string) (ok bool) {
	for _, v := range validators {
		if v.Name == name {
			e.skip[name] = true
			return true
		}
	}

	return false
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

// Validators returns a sorted list of all validators which are not explicitly
// skipped.  We sort by priority and then name in order to allow
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

// Validate checks the given full path against all validators not in the skip
// list, returning an array of errors found
func (e *Engine) Validate(path string, info os.FileInfo) []Failure {
	var flist []Failure

	var v Validator
	for _, v = range e.Validators() {
		flist = v.Validate(path, info, flist)
	}

	return flist
}
