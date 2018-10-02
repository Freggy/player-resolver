GOBUILD=go build
GOTEST=go test
GIT_REVISION=$(shell git rev-parse --short=8 HEAD)
PLATFORMS=darwin linux windows
BINARY=player-resolver-$(GIT_REVISION)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${GIT_REVISION}"

all: deps dep test build

deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint

dep:
	dep ensure -vendor-only

build:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); gofmt -w -s .; CGO_ENABLED=0 $(GOBUILD) -o $(BINARY) -v

test:
	$(GOTEST) ./... -v