GO := go
GLIDE := glide
COMPOSE := docker-compose
TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

test:
	@$(GO) test ./...

build-local:
	@$(GO) build -o gowatch main.go
	@mv gowatch ~/go/bin
