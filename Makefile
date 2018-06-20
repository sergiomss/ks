VERSION = $(shell git describe --tags --always)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
NAME = ks
.PHONY: build

build:
	go build .

release:
	go get github.com/mitchellh/gox
	CGO_ENABLED=0 gox -output "dist/$(NAME)_{{.OS}}_{{.Arch}}" -arch "amd64" -os "linux windows darwin" $(shell go list ./... | grep -v '/vendor/')
