package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/uoregon-libraries/dark-archive-validator/src/rules"
)

// FileValidationFailure combines path and failure list
type FileValidationFailure struct {
	Filepath string
	Failures []rules.Failure
}

var engine *rules.Engine
var rootPath string
var allValidatorNames []string
var validatorNameIndices = make(map[string]int)
var fileValidationFailures = make([]FileValidationFailure, 0)
var checksums = make(map[string][]string)

func main() {
	engine = rules.NewEngine()
	processCLI()
	getAllValidators()
	engine.ValidateTree(rootPath, failfunc)
	exportValidationFailures()

	if opts.SHAOutput != "" {
		writeSha()
	}

	if len(fileValidationFailures) != 0 {
		os.Exit(1)
	}

	os.Exit(0)
}

// getAllValidators puts together the complete list of validator names from a
// fresh rules engine.  This is done to ensure a consistent report of which
// validators exist, not just which failed for a given tree
func getAllValidators() {
	var vList = rules.NewEngine().Validators()
	allValidatorNames = make([]string, len(vList))
	for i, v := range vList {
		allValidatorNames[i] = v.Name
		validatorNameIndices[v.Name] = i
	}
}

func failfunc(path string, fList []rules.Failure) {
	fileValidationFailures = append(fileValidationFailures, FileValidationFailure{path, fList})
}

// exportValidationFailures prints out a CSV of failure data
func exportValidationFailures() {
	var header = make([]string, len(allValidatorNames)+1)
	header[0] = "Filename"
	for i, vName := range allValidatorNames {
		header[i+1] = vName
	}
	printTSV(header)

	for _, fvf := range fileValidationFailures {
		// Prep the columns
		var columns = make([]string, len(allValidatorNames)+1)
		columns[0] = fmt.Sprintf("%#v", fvf.Filepath)

		// Build the failure message for the appropriate column
		for _, f := range fvf.Failures {
			columns[validatorNameIndices[f.V.Name]+1] = f.E.Error()
		}

		printTSV(columns)
	}
}

// printTSV just prints the strings in cols as-is, tab-separated
func printTSV(cols []string) {
	fmt.Println(strings.Join(cols, "\t"))
}

func writeSha() {
	var lines = make([]string, 0)
	for sha, filenames := range checksums {
		for _, filename := range filenames {
			lines = append(lines, fmt.Sprintf("%s  %s", sha, filename))
		}
	}
	sort.Strings(lines)

	var f, err = os.Create(opts.SHAOutput)
	if err == nil {
		_, err = f.WriteString(strings.Join(lines, "\n"))
	}

	if err != nil {
		log.Fatalf("Unable to write SHA checksums: %s", err)
	}

	f.Close()
}
