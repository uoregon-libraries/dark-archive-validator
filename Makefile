.PHONY: clean linux64 linux32 win win32 osx test

linux64:
	env GOOS=linux GOARCH=amd64 gb build
linux32:
	env GOOS=linux GOARCH=386 gb build
win:
	env GOOS=windows GOARCH=amd64 gb build
win32:
	env GOOS=windows GOARCH=386 gb build
osx:
	env GOOS=darwin GOARCH=amd64 gb build

test:
	gb test -v

clean:
	rm -rf pkg/
