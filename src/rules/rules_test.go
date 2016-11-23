package rules_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rules"
)

func fakeFileWalk(root string, walkfn filepath.WalkFunc) error {
	var walk = func(dir string, i os.FileInfo) {
		var fullPath = filepath.Join(root, dir, i.Name())
		walkfn(fullPath, i, nil)
	}

	// Perfectly good file
	walk("", rules.NewFakeFile("goodfile.txt", 1024))
	// Perfectly good directory
	walk("", rules.NewFakeDir("stuff"))
	// Perfectly good file within the good dir
	walk("stuff/", rules.NewFakeFile("goodfile2.txt", 1024))
	// Symlink - we'll be excluding these
	walk("", rules.NewFakeSymlink("flarb"))
	// Windows reserved filename, even though it's got an extension
	walk("", rules.NewFakeDir("lPt2.dir"))
	// Good file even though the dir is bad
	walk("lPt2.dir", rules.NewFakeFile("file.txt", 1024))
	// Good file even though the dir is *really* bad
	walk("foo: \x05bar*baz.dir", rules.NewFakeFile("file.txt", 1024))
	// Duped filename; case insensitivity will be important
	walk("lPt2.dir", rules.NewFakeFile("fILe.txt", 2048))
	// Zero-length files are bad, mmkay?
	walk("lPt2.dir", rules.NewFakeFile("zerofile.txt", 0))

	// Multiple problems: bad characters for windows, bad characters for our own
	// sanity, too long a path, device file
	walk(strings.Repeat("blah", 100) + "/dev/", rules.NewFakeDevice(":\"thi\x05ng*"))

	return nil
}

func failFunc(path string, failures []rules.Failure) {
	for _, f := range failures {
		fmt.Printf("%s says %#v %s\n", f.V.Name, path, f.E)
	}
}

func ExampleEngine() {
	var e = rules.NewEngine()
	e.TraverseFn = fakeFileWalk
	e.AddValidator("no-special-files", rules.NoSpecialFiles)

	e.ValidateTree("/this/path/shouldn't/actually/have/any/kind/of/testing/so I can do *all kinds* of bad things in here!\x1b\x1b/", failFunc)

	// Output:
	// no-special-files says "flarb" is a symbolic link
	// valid-windows-filename says "lPt2.dir" uses a reserved file name
	// valid-windows-filename says "blahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" contains invalid characters (:, ", *)
	// no-special-files says "blahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" is a device file
}
