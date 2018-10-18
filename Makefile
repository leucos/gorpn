VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w
SRC=$(shell find . -name '*.go')

.PHONY: build install sign release clean

build:
	go build -o gorpn -ldflags="$(FLAGS)" .

install:
	go install -ldflags="$(FLAGS)" .

gorpn-linux-amd64: $(SRC)
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

gorpn-darwin-amd64: $(SRC)
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

gorpn-windows-386.exe: $(SRC)
	GOOS=windows GOARCH=386 go build -o $@ -ldflags="$(FLAGS)" .

release: gorpn-linux-amd64 gorpn-darwin-amd64 gorpn-windows-386.exe

clean:
	rm -f gorpn gorpn-linux-amd64 gorpn-darwin-amd64 gorpn-windows-386.exe
