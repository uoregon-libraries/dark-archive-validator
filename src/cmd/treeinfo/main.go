package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	var root = os.Args[1]
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("CRITICAL - Unable to process %#v: %s", path, err)
		}

		var basepath = strings.Replace(path, root+"/", "", -1)
		log.Printf("%#v: %#v", basepath, info)
		return nil
	})
}
