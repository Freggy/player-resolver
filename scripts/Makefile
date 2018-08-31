GOBUILD=go build
GOTEST=go test
GIT_REVISION= $(shell git rev-parse --short=8 HEAD)
BINARY=player-resolver-$(GIT_REVISION)
SWTCH_FMT=$(shell cd .. && gofmt -w -s)

all: test build
build:
	cd ..; $(GOBUILD) -o $(BINARY) -v
test:
	cd ..; $(GOTEST) ./... -v