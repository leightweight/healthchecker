VERSION := $(shell git describe --tags --dirty)
DATE    := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

IMPORT_PATH := github.com/leightweight/healthchecker
CLI_PATH    := cmd/healthchecker

LDFLAGS_VERSION := -X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"
LDFLAGS         := -ldflags='$(LDFLAGS_VERSION)'

CONTAINER_TAGS := -t leightweight/healthchecker:latest
CONTAINER_TAGS := $(CONTAINER_TAGS) -t leightweight/healthchecker:$(VERSION)
CONTAINER_ARGS := --build-arg VERSION="$(VERSION)" --build-arg DATE="$(DATE)"

.PHONY: default
default: cli

.PHONY: clean
clean:
	rm -f healthchecker
	go clean

.PHONY: cli
cli:
	go build -v $(LDFLAGS) $(IMPORT_PATH)/$(CLI_PATH)

.PHONY: container
container:
	docker build $(CONTAINER_TAGS) $(CONTAINER_ARGS) .
