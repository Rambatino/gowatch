GO := go
GLIDE := glide
COMPOSE := docker-compose
TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG
VERSION = snapshot
GHRFLAGS =

.PHONY: build release
default: build

test:
	@$(GO) test ./...

build-local:
	@$(GO) build -o gowatch main.go
	@mv gowatch ~/go/bin

build:
	goxc -d=pkg -pv=$(VERSION) -os="linux darwin windows"

release:
	ghr -u Rambatino $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)
	cd homebrew && git commit -am "Updated gowatch" && git push
