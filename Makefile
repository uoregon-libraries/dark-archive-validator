.PHONY: default clean linux64 linux32 win win32 osx test lint fmt

SOURCES := $(shell find ./src -name "*.go")
SOURCEDIRS := $(shell find ./src -type d)

default: deps bin/validate

deps:
	go mod download

# For quick building of binaries, you can run something like "make bin/server"
# and still have a little bit of the vetting without running the entire
# validation script
bin/%: src/cmd/% $(SOURCES) $(SOURCEDIRS)
	go vet ./$<
	go build -ldflags="-s -w" -o $@ github.com/uoregon-libraries/dark-archive-validator/$<

test:
	go test -v ./src/...

clean:
	rm -rf bin/

lint:
	golint -set_exit_status src/...

fmt:
	find src -name "*.go" -exec gofmt -l -w -s {} \;
