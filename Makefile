BUILD=go build

default: linux

clean:
	@rm -rf bin/
	@rm -f debug debug.test web/debug web/debug.test

all: windows linux macos

windows:
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/windows_amd64.exe
	GOOS=windows GOARCH=386 $(BUILD) -o bin/windows_x86.exe
linux:
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/linux_amd64
	GOOS=linux GOARCH=386 $(BUILD) -o bin/linux_x86
macos:
	GOOS=darwin GOARCH=amd64 $(BUILD) -o bin/macos
arm:
	GOOS=linux GOARCH=arm $(BUILD) -o bin/arm