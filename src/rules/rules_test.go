package rules_test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"checksum"
	"rules"
)

// fakeFileWalk fires off the walkfn on a variety of paths to test out most of
// the validators
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
	// File with no extension
	walk("lPt2.dir", rules.NewFakeFile("fILe", 2048))
	// File with space
	walk("", rules.NewFakeFile("this isbad.txt", 1024))
	// File with space at the end
	walk("", rules.NewFakeFile("thisisbad.txt ", 1024))
	// File with wonky space
	walk("", rules.NewFakeFile("this\u202fisbad.txt", 1024))
	// File with wonky space at the end
	walk("", rules.NewFakeFile("thisisbad.txt\u202f", 1024))
	// File that doesn't start with an alpha character
	walk("", rules.NewFakeFile("0.txt", 1024))
	// File that violates DSC conventions
	walk("", rules.NewFakeFile("abc@foo.bar", 1024))
	// Too many periods
	walk("", rules.NewFakeFile("foo.bar.txt", 1024))
	walk("", rules.NewFakeDir("foo.bar.dir"))
	// Hidden
	walk("", rules.NewFakeFile(".hiddenfile", 1024))
	walk("", rules.NewFakeDir(".hiddendir"))
	// UTF-8
	walk("", rules.NewFakeFile("foo‣‡•.txt", 1024))
	// Bad UTF-8
	walk("", rules.NewFakeFile("foo\xed\x88.txt", 1024))
	// Mac OSX garbage
	walk("", rules.NewFakeFile(".DS_Store", 1024))
	walk("", rules.NewFakeFile("._foo.txt", 1024))
	// Windows garbage
	walk("", rules.NewFakeFile("Thumbs.db", 1024))
	walk("", rules.NewFakeFile("desktop.ini", 1024))

	// Multiple problems: bad characters for windows, bad characters for our own
	// sanity, too long a path, device file
	walk(strings.Repeat("blah", 10)+"/dev/", rules.NewFakeDevice(":\"thi\x05ng*"))

	return nil
}

// fakeFileWalk2 tests that our "restrictive naming" catches things missed by
// other validators.  It's separate from the above function because the most
// obvious way to test this is by manually skipping other validators that would
// otherwise trap the error.
func fakeFileWalk2(root string, walkfn filepath.WalkFunc) error {
	var walk = func(dir string, i os.FileInfo) {
		var fullPath = filepath.Join(root, dir, i.Name())
		walkfn(fullPath, i, nil)
	}

	// File that violates DSC conventions
	walk("", rules.NewFakeFile("abc@foo.bar", 1024))

	return nil
}

// fakeFileWalkChecksum walks files that will cause a fake checksum collision.
// The fake checksum function runs against the base filename.
func fakeFileWalkChecksum(root string, walkfn filepath.WalkFunc) error {
	var walk = func(dir string, i os.FileInfo) {
		var fullPath = filepath.Join(root, dir, i.Name())
		walkfn(fullPath, i, nil)
	}

	walk("a", rules.NewFakeFile("one.txt", 1024))
	walk("b", rules.NewFakeFile("one.txt", 1024))
	walk("b", rules.NewFakeFile("two.txt", 1024))

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

	// For testing, we have to register a shorter path-limit validation
	rules.RegisterValidatorHigh("path-limit", rules.PathLimitFn(50))

	e.ValidateTree("/this/path/shouldn't/actually/have/any/kind/of/testing/so I can do *all kinds* of bad things in here!\x1b\x1b/", failFunc)

	// Output:
	// no-special-files says "flarb" is a symbolic link
	// valid-windows-filename says "lPt2.dir" uses a reserved file name
	// no-duped-names says "lPt2.dir/fILe.txt" is a duplicate of "lPt2.dir/file.txt"
	// nonzero-filesize says "lPt2.dir/zerofile.txt" is an empty file
	// has-extension says "lPt2.dir/fILe" doesn't have an extension
	// no-spaces says "this isbad.txt" has a space in the filename
	// valid-windows-filename says "thisisbad.txt " has a trailing space
	// no-spaces says "thisisbad.txt " ends with a space
	// no-spaces says "this\u202fisbad.txt" has a space in the filename
	// no-utf8 says "this\u202fisbad.txt" contains unicode characters (" ")
	// no-spaces says "thisisbad.txt\u202f" ends with a space
	// no-utf8 says "thisisbad.txt\u202f" contains unicode characters (" ")
	// starts-with-alpha says "0.txt" starts with a non-alphabetic character
	// valid-dsc-filename says "abc@foo.bar" contains invalid characters: @
	// has-only-one-period says "foo.bar.txt" has 2 periods (maximum is 1)
	// has-only-one-period says "foo.bar.dir" has 2 periods (maximum is 1)
	// no-hidden-files says ".hiddenfile" is hidden (starts with a period)
	// starts-with-alpha says ".hiddenfile" starts with a non-alphabetic character
	// no-hidden-files says ".hiddendir" is hidden (starts with a period)
	// starts-with-alpha says ".hiddendir" starts with a non-alphabetic character
	// no-utf8 says "foo‣‡•.txt" contains unicode characters ("‣", "‡", "•")
	// invalid-utf8 says "foo\xed\x88.txt" contains invalid unicode
	// no-extraneous-files says ".DS_Store" may be an extraneous file; consider deletion
	// no-extraneous-files says "._foo.txt" may be an extraneous file; consider deletion
	// no-extraneous-files says "Thumbs.db" may be an extraneous file; consider deletion
	// no-extraneous-files says "desktop.ini" may be an extraneous file; consider deletion
	// valid-windows-filename says "blahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" contains invalid characters: : " *
	// no-control-chars says "blahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" contains one or more control characters
	// no-special-files says "blahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" is a device file
	// path-limit says "blahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" exceeds the maximum path length of 50 characters
	// starts-with-alpha says "blahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" starts with a non-alphabetic character
	// valid-dsc-filename says "blahblahblahblahblahblahblahblahblahblah/dev/:\"thi\x05ng*" contains invalid characters: *
}

// This example skips valid-dsc-filename in order to let restrictive-naming
// actually catch an error, since that validator only catches items that have
// no other failures.  For simplicity, we use fakeFileWalk2, which only has one
// fake file to test.
func ExampleEngine_skipDSCForRestrictiveTest() {
	rules.ResetDupemap()
	var e = rules.NewEngine()
	e.TraverseFn = fakeFileWalk2
	e.Skip("valid-dsc-filename")
	e.ValidateTree("/blah", failFunc)

	// Output:
	// restrictive-naming says "abc@foo.bar" doesn't match required filename pattern
}

// This example shows how serious we are about not skipping critical
// validations.  DEADLY SERIOUS, FOLKS.
func ExampleEngine_noSkippingCriticalValidations() {
	var e = rules.NewEngine()
	for _, v := range e.Validators() {
		e.Skip(v.Name)
	}

	for _, v := range e.Validators() {
		fmt.Println("After manually running Skip, found", v.Name)
	}

	e.SkipAll()
	for _, v := range e.Validators() {
		fmt.Println("After SkipAll, found", v.Name)
	}

	// Output:
	// After manually running Skip, found no-duped-names
	// After manually running Skip, found valid-windows-filename
	// After SkipAll, found no-duped-names
	// After SkipAll, found valid-windows-filename
}

func fakeBlockWrite(path string, w io.Writer) error {
	var basename = filepath.Base(path)
	w.Write([]byte(basename))
	return nil
}

// This example verifies checksums are working properly
func ExampleEngine_onlyTestChecksums() {
	var e = rules.NewEngine()
	e.TraverseFn = fakeFileWalkChecksum
	e.SkipAll()
	rules.RegisterChecksumValidator("/blah", &checksum.Checksum{sha256.New(), fakeBlockWrite})
	e.ValidateTree("/blah", failFunc)

	// Output:
	// no-duped-content says "b/one.txt" duplicates the content of "a/one.txt"
}
