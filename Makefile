GO := go
GLIDE := glide
COMPOSE := docker-compose
TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

test:
	@$(GO) test ./...
