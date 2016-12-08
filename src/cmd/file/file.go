// file.go is a simple example to verify that the magicmime package is doing
// what I want.  It is not meant as a replacement of the file command.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rakyll/magicmime"
)

func usage() {
	fmt.Printf("usage: %s <filename> [--mime]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	var l = len(os.Args)
	if l < 2 || l > 3 {
		usage()
	}

	var flags = magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR
	var fName = os.Args[1]

	if l == 3 {
		if os.Args[2] == "--mime" {
			flags |= magicmime.MAGIC_MIME_TYPE
		} else {
			usage()
		}
	}

	var err = magicmime.Open(flags)
	if err != nil {
		log.Fatal(err)
	}
	defer magicmime.Close()

	var mt string
	mt, err = magicmime.TypeByFile(fName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(mt)
}
