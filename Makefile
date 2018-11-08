NAME = ks
VERSION = $(shell git describe --tags --always)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_FLAGS = "-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)"

.PHONY: build

build:
	go build .

release:
	go get github.com/mitchellh/gox
	CGO_ENABLED=1 gox -output "dist/$(NAME)_{{.OS}}_{{.Arch}}" -ldflags $(LD_FLAGS) -arch "amd64" -os "linux windows darwin" $(shell go list ./... | grep -v '/vendor/')
	shasum -a 256 dist/ks_darwin_amd64
