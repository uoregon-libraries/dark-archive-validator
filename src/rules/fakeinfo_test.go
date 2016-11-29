package rules

import (
	"os"
	"time"
)

// FakeInfo will be used by multiple tests to supply fake file info data
type FakeInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func NewFakeFile(n string, s int64) FakeInfo {
	return FakeInfo{name: n, size: s}
}

func NewFakeDir(n string) FakeInfo {
	return FakeInfo{name: n, isDir: true, mode: os.ModeDir}
}

func NewFakeSymlink(n string) FakeInfo {
	return FakeInfo{name: n, mode: os.ModeSymlink}
}

func NewFakeDevice(n string) FakeInfo {
	return FakeInfo{name: n, mode: os.ModeDevice}
}

func (fi FakeInfo) Name() string {
	return fi.name
}

func (fi FakeInfo) Size() int64 {
	return fi.size
}

func (fi FakeInfo) Mode() os.FileMode {
	return fi.mode
}

func (fi FakeInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi FakeInfo) IsDir() bool {
	return fi.isDir
}

func (fi FakeInfo) Sys() interface{} {
	return nil
}
