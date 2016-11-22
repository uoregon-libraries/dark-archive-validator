package rules

import (
	"os"
	"time"
)

// FakeInfo will be used by multiple tests to supply fake file info data
type FakeInfo struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
	isDir bool
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
