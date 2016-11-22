package rules

import (
	"testing"
)

var filenameTests = []struct {
	name string
	good bool
}{
	{"hello.txt", true},
	{"CON", false},
	{"CON.txt", false},
	{"con.txt", false},
	{"LpT7.foobar", false},
	{"t:est.txt", false},
	{"hello.", false},
	{"hello.txt ", false},
}

func TestInvalidWindowsFilename(t *testing.T) {
	for _, ft := range filenameTests {
		var info = FakeInfo{name: ft.name}
		var err = ValidWindowsFilename("", info)
		if err != nil && ft.good {
			t.Errorf("Expected %#v to be good, but got error: %s", ft.name, err)
			continue
		}
		if err == nil && !ft.good {
			t.Errorf("Expected %#v to be bad, but no error", ft.name)
			continue
		}

		if ft.good {
			t.Logf("OK - %#v (valid filename)", ft.name)
		} else {
			t.Logf("OK - %#v %s (invalid filename)", ft.name, err)
		}
	}
}
