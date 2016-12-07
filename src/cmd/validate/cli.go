package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

var parser *flags.Parser
var opts struct {
	SkipList       []string `short:"s" long:"skip" description:"Skip a particular validator.  Cannot be used to skip Windows filename validations.  Can be repeated to skip multiple validations."`
	SHA256         bool     `long:"sha256" description:"Use SHA-256 checksum to look for duplicate content"`
	ListValidators bool     `short:"l" long:"list-validators" description:"List all validators' names"`
}

func usage(err error) {
	var status int
	if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
		err = nil
	}
	if err != nil {
		os.Stderr.Write([]byte("ERROR: " + err.Error() + "\n\n"))
		status = 1
	}

	parser.WriteHelp(os.Stderr)
	os.Exit(status)
}

func listValidatorsAndExit() {
	fmt.Printf("Known validators:\n")
	for _, v := range engine.Validators() {
		fmt.Printf("  %s (%s)\n", v.Name, v.Criticality)
	}
	os.Exit(0)
}

func processSkipList() []string {
	var invalids []string
	for _, skip := range opts.SkipList {
		if !engine.Skip(skip) {
			invalids = append(invalids, skip)
		}
	}

	return invalids
}

func processCLI() {
	parser = flags.NewParser(&opts, flags.HelpFlag)
	parser.Usage = "[OPTIONS] <path to validate>"
	var more, err = parser.Parse()
	if err != nil {
		usage(err)
	}

	// If listing validators, nothing else matters
	if opts.ListValidators {
		listValidatorsAndExit()
	}

	// Check for skips so we can verify those quickly
	var invalids = processSkipList()
	if len(invalids) != 0 {
		usage(fmt.Errorf("Invalid --skip value(s): %s", strings.Join(invalids, ", ")))
	}

	// Not listing validators; need to check for valid path
	if err == nil && len(more) != 1 {
		usage(fmt.Errorf("must specify exactly one path to validate"))
	}

	var info os.FileInfo
	var vPath = more[0]
	info, err = os.Stat(vPath)
	if err != nil {
		usage(err)
	}

	if !info.Mode().IsDir() {
		usage(fmt.Errorf("%s is not a valid path to validate", vPath))
	}

	// All seems well; set the rootPath and we're good to go
	rootPath = more[0]
}
