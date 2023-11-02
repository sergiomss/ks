NAME = ks
VERSION = $(shell git describe --tags --always)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT = $(shell git rev-parse HEAD)
LD_FLAGS = "-X main.version=$(VERSION) -X main.date=$(DATE) -X main.commit=$(GIT_COMMIT)"

.PHONY: build

build:
	@go get github.com/mitchellh/gox
	@gox -output "dist/$(NAME)_{{.OS}}_{{.Arch}}" -ldflags $(LD_FLAGS) -arch "amd64 arm64" -os "linux windows darwin"

release: build
	@for f in $(shell ls dist/); do shasum -a 256 dist/$${f}; done

