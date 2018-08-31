SHELL := /bin/sh
GOBUILD=go build
GOTEST=go test
GIT_REVISION= $(shell git rev-parse --short=8 HEAD)
BINARY=player-resolver-$(GIT_REVISION)

all: test build
build:
	cd $(PWD)/cmd; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY) -v
test:
	$(GOTEST) ./... -v