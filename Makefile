.PHONY: build install test clean

build:
	CGO_ENABLED=0 go build -trimpath -o gh-runner

install: build
	gh extension install .

test:
	go test ./...

clean:
	rm -f gh-runner

.DEFAULT_GOAL := build
