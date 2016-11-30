package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var reExample = regexp.MustCompile(`\A--- FAIL: (Example.*) \(.*\)\z`)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <test output filename>")
		os.Exit(1)
	}

	var f, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Unable to open %#v: %s", os.Args[1], err)
		os.Exit(1)
	}
	defer f.Close()

	var exampleName string
	var seenGot, seenWant bool
	var gotFile, wantFile *os.File
	var s = bufio.NewScanner(f)
	for s.Scan() {
		var line = s.Text()
		if line == "FAIL" {
			break
		}

		var m = reExample.FindStringSubmatch(line)
		if m != nil {
			exampleName = m[1]
			fmt.Printf("Found failed example: %#v\n", exampleName)
			seenGot = false
			seenWant = false

			gotFile, err = os.Create("test." + exampleName + ".got")
			if err != nil {
				panic(err)
			}

			wantFile, err = os.Create("test." + exampleName + ".want")
			if err != nil {
				panic(err)
			}
		}

		if exampleName != "" {
			if !seenGot && line == "got:" {
				seenGot = true
				continue
			}
			if !seenWant && line == "want:" {
				seenWant = true
				continue
			}

			if seenGot && !seenWant {
				gotFile.Write([]byte(line+"\n"))
			}
			if seenGot && seenWant {
				wantFile.Write([]byte(line+"\n"))
			}
		}
	}
}
