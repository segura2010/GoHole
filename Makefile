BUILD=go build
COMMIT=`git describe --always --long`
CDATE=`date +%FT%T%z`
LDFLAGS=-ldflags="-X main.Commit=${COMMIT} -X main.CompilationDate=${CDATE}"

default: linux

install:
	@go install ${LDFLAGS}
	@mv ${GOPATH}/bin/GoHole ${GOPATH}/bin/gohole
	@echo "Binary on ${GOPATH}/bin/gohole"

clean:
	@rm -rf bin/
	@rm -f debug debug.test web/debug web/debug.test

all: windows linux macos

windows:
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/windows_amd64.exe ${LDFLAGS}
	GOOS=windows GOARCH=386 $(BUILD) -o bin/windows_x86.exe ${LDFLAGS}
linux:
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/linux_amd64 ${LDFLAGS}
	GOOS=linux GOARCH=386 $(BUILD) -o bin/linux_x86 ${LDFLAGS}
macos:
	GOOS=darwin GOARCH=amd64 $(BUILD) -o bin/macos ${LDFLAGS}
arm:
	GOOS=linux GOARCH=arm $(BUILD) -o bin/arm ${LDFLAGS}