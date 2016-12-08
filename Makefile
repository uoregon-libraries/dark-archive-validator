.PHONY: bin test clean lint fmt

bin:
	gb build

test:
	gb test -v

clean:
	rm -rf pkg/ bin/

lint:
	GOPATH=$(PWD) gometalinter --disable gotype --deadline 10s src/...

fmt:
	find src -name "*.go" -exec gofmt -l -w -s {} \;
