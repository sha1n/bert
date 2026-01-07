ifndef VERSION
# Set VERSION to the latest version tag name. Assuming version tags are formatted 'v*'
VERSION := $(shell git describe --always --abbrev=0 --tags --match "v*" $(git rev-list --tags --max-count=1))
BUILD := $(shell git rev-parse $(VERSION))
endif
PROJECTNAME := "bert"
# We pass that to the main module to generate the correct help text
PROGRAMNAME := $(PROJECTNAME)

# Go related variables.
GOHOSTOS := $(shell go env GOHOSTOS)
GOHOSTARCH := $(shell go env GOHOSTARCH)
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GODIST := $(GOBASE)/dist
GOBUILD := $(GOBASE)/build
GOFILES := $(shell find . -type f -name '*.go' -not -path './vendor/*')
GOOS_DARWIN := "darwin"
GOOS_LINUX := "linux"
GOOS_WINDOWS := "windows"
GOARCH_AMD64 := "amd64"
GOARCH_ARM64 := "arm64"
GOARCH_ARM := "arm"

MODFLAGS=-mod=readonly

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-w -s -X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.ProgramName=$(PROGRAMNAME)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := $(GOBUILD)/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID := $(GOBUILD)/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: install lint format test build

## install: Checks for missing dependencies and installs them
.PHONY: install
install: go-get

## format: Formats Go source files
.PHONY: format
format: go-format

## lint: Runs all linters including go vet and golangci-lint
.PHONY: lint
lint: go-lint golangci-lint

## build: Builds binaries for all supported platforms
.PHONY: build
build:
	@[ -d $(GOBUILD) ] || mkdir -p $(GOBUILD)
	@-mkdir -p $(GOBUILD)/completions
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-build #2> $(STDERR)
	# generate completions
	bin/bert-$(GOHOSTOS)-$(GOHOSTARCH) completion zsh > $(GOBUILD)/completions/_bert
	bin/bert-$(GOHOSTOS)-$(GOHOSTARCH) completion bash > $(GOBUILD)/completions/bert.bash
	bin/bert-$(GOHOSTOS)-$(GOHOSTARCH) completion fish > $(GOBUILD)/completions/bert.fish

	#@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*\]/ /' | sed "/^/s/^/     /" 1>&2


## test: Runs all Go tests
.PHONY: test
test: install go-test

## clean: Removes build artifacts
.PHONY: clean
clean:
	@-rm $(GOBIN)/$(PROGRAMNAME)* 2> /dev/null
	@-$(MAKE) go-clean

.PHONY: go-lint
go-lint:
	@echo "  >  Linting source files..."
	go vet $(MODFLAGS) -c=10 `go list $(MODFLAGS) ./...`

## golangci-lint: Runs golangci-lint
.PHONY: golangci-lint
golangci-lint:
	@echo "  >  Running golangci-lint..."
	golangci-lint run

.PHONY: go-format
go-format:
	@echo "  >  Formating source files..."
	gofmt -s -w $(GOFILES)

.PHONY: go-build-current
go-build-current:
	@echo "  >  Building $(GOHOSTOS)/$(GOHOSTARCH) binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME) $(GOBASE)/cmd

.PHONY: go-build
go-build: go-get go-build-linux-amd64 go-build-linux-arm64 go-build-linux-arm go-build-darwin-amd64 go-build-darwin-arm64 go-build-windows-amd64 go-build-windows-arm

.PHONY: go-test
go-test:
	@echo "  >  Running Go tests..."
	go test $(MODFLAGS) -covermode=count `go list $(MODFLAGS) ./...`

.PHONY: go-build-linux-amd64
go-build-linux-amd64:
	@echo "  >  Building linux amd64 binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_AMD64) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_LINUX)-$(GOARCH_AMD64) $(GOBASE)/cmd

.PHONY: go-build-linux-arm64
go-build-linux-arm64:
	@echo "  >  Building linux arm64 binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_ARM64) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_LINUX)-$(GOARCH_ARM64) $(GOBASE)/cmd

.PHONY: go-build-linux-arm
go-build-linux-arm:
	@echo "  >  Building linux arm binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_ARM) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_LINUX)-$(GOARCH_ARM) $(GOBASE)/cmd

.PHONY: go-build-darwin-amd64
go-build-darwin-amd64:
	@echo "  >  Building darwin amd64 binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_DARWIN) GOARCH=$(GOARCH_AMD64) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_DARWIN)-$(GOARCH_AMD64) $(GOBASE)/cmd

.PHONY: go-build-darwin-arm64
go-build-darwin-arm64:
	@echo "  >  Building darwin arm64 binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_DARWIN) GOARCH=$(GOARCH_ARM64) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_DARWIN)-$(GOARCH_ARM64) $(GOBASE)/cmd

.PHONY: go-build-windows-amd64
go-build-windows-amd64:
	@echo "  >  Building windows amd64 binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_WINDOWS) GOARCH=$(GOARCH_AMD64) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_WINDOWS)-$(GOARCH_AMD64).exe $(GOBASE)/cmd

.PHONY: go-build-windows-arm
go-build-windows-arm:
	@echo "  >  Building windows arm binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_WINDOWS) GOARCH=$(GOARCH_ARM) GOBIN=$(GOBIN) go build $(MODFLAGS) $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_WINDOWS)-$(GOARCH_ARM).exe $(GOBASE)/cmd

.PHONY: go-generate
go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

.PHONY: go-get
go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy

.PHONY: go-install
go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

.PHONY: go-clean
go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean $(MODFLAGS) $(GOBASE)/cmd
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean -modcache

.PHONY: run-sanity-tests
run-sanity-tests: build-docker run-linux-dockerized-tests

.PHONY: build-docker
build-docker:
	@echo "  >  Building docker image..."
	docker build -t sha1n/$(PROJECTNAME):latest .
	docker tag sha1n/$(PROJECTNAME):latest sha1n/$(PROJECTNAME):$(VERSION:v%=%)

.PHONY: run-linux-dockerized-tests
run-linux-dockerized-tests:
	@echo "  >  Running with experimental UI..."
	docker run --rm -ti sha1n/bert /bert/bin/bert-linux-amd64 -c /bert/test/data/spec_test_load.yaml
	@echo "  >  Running with experimental UI + debug..."
	docker run --rm -ti sha1n/bert /bert/bin/bert-linux-amd64 -c /bert/test/data/spec_test_load.yaml -d
	@echo "  >  Running with experimental UI + silent..."
	docker run --rm -ti sha1n/bert /bert/bin/bert-linux-amd64 -c /bert/test/data/spec_test_load.yaml -s
	@echo "  >  Running with experimental UI + piped stdout..."
	docker run --rm -ti sha1n/bert /bert/bin/bert-linux-amd64 -c /bert/test/data/spec_test_load.yaml --pipe-stdout
	@echo "  >  Running with experimental UI + piped stderr..."
	docker run --rm -ti sha1n/bert /bert/bin/bert-linux-amd64 -c /bert/test/data/spec_test_load.yaml --pipe-stderr
	@echo "  >  Running with ad-hoc commands..."
	docker run --rm -ti sha1n/bert /bert/bin/bert-linux-amd64 'ls' 'ls -laH' --executions 10

.PHONY: release
release:
ifdef GITHUB_TOKEN
	@echo "  >  Releasing..."
	goreleaser release --clean
else
	$(error GITHUB_TOKEN is not set)
endif

.PHONY: all
all: help

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
