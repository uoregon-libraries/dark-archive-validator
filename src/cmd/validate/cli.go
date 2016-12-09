package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"checksum"
	"rules"

	"github.com/jessevdk/go-flags"
)

var parser *flags.Parser
var opts struct {
	SkipList       []string `short:"s" long:"skip" description:"Skip a particular validator.  Cannot be used to skip critical validations.  Can be repeated to skip multiple validations."`
	Quick          bool     `long:"quick" description:"Skip checksum and lowest-criticality validators"`
	ListValidators bool     `short:"l" long:"list-validators" description:"List all validators this command would have run"`
	SHAOutput      string   `short:"o" long:"sha-output" description:"Filename for writing all files' SHA256 hashes"`
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
	fmt.Printf("Validators to run:\n")
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

func skipUnimportantValidators() {
	for _, v := range engine.Validators() {
		if !v.IsImportant() {
			engine.Skip(v.Name)
		}
	}
}

func processCLI() {
	parser = flags.NewParser(&opts, flags.HelpFlag)
	parser.Usage = "[OPTIONS] <path to validate>"
	var more, err = parser.Parse()
	if err != nil {
		usage(err)
	}

	if len(more) > 0 {
		rootPath = more[0]
	}
	rules.RegisterChecksumValidator(rootPath, checksum.New(sha256.New()), checksums)

	if opts.SHAOutput != "" {
		// Make sure the given file can be created and written
		var _, err = os.OpenFile(opts.SHAOutput, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			usage(fmt.Errorf("Unable to write to %s", opts.SHAOutput))
		}
	}

	// --quick skips non-critical validators and the very slow checksumming
	// validator
	if opts.Quick {
		if opts.SHAOutput != "" {
			usage(fmt.Errorf("Cannot combine --quick and --sha-output"))
		}
		skipUnimportantValidators()
		opts.SkipList = append(opts.SkipList, "no-duped-content")
	}

	// Check for skips so we can verify those quickly
	var invalids = processSkipList()
	if len(invalids) != 0 {
		usage(fmt.Errorf("Invalid --skip value(s): %s", strings.Join(invalids, ", ")))
	}

	// If listing validators, nothing else matters
	if opts.ListValidators {
		listValidatorsAndExit()
	}

	// Not listing validators; need to check for valid path
	if err == nil && len(more) != 1 {
		usage(fmt.Errorf("must specify exactly one path to validate"))
	}

	var info os.FileInfo
	info, err = os.Stat(rootPath)
	if err != nil {
		usage(err)
	}

	if !info.Mode().IsDir() {
		usage(fmt.Errorf("%s is not a valid path to validate", rootPath))
	}
}
