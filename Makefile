
NAME=live-go

PKGS ?= $(shell glide novendor)
# Many Go tools take file globs or directories as arguments instead of packages.
PKG_FILES ?= *.go log utils config


GO_VERSION := $(shell go version | cut -d " " -f 3)
GO_MINOR_VERSION := $(word 2,$(subst ., ,$(GO_VERSION)))
LINTABLE_MINOR_VERSION := 9
ifneq ($(filter $(LINTABLE_MINOR_VERSION), $(GO_MINOR_VERSION)),)
	SHOULD_LINT := true
endif

# Golang Flags
GOPATH ?= $(GOPATH:):./vendor
GOFLAGS ?= $(GOFLAGS:)
GO=go

.PHONY: all
all: lint test

.PHONY: setup-ci
setup-ci:
	@echo "Installing development dependencies..."
	go get -u -f golang.org/x/tools/cmd/goimports
	go get -u -f golang.org/x/tools/cmd/cover
	@echo "Installing dependencies..."
	go get -u -f github.com/pborman/uuid
	go get -u -f github.com/golang/lint
	go get -u -f github.com/spf13/viper
	go get -u -f github.com/gorilla/websocket
	@echo "Installing test dependencies..."
	go get -u -f github.com/axw/gocov
	go get -u -f github.com/mattn/goveralls
ifdef SHOULD_LINT
	@echo "Installing golint..."
	#go install ./vendor/github.com/golang/lint/golint
	go get -u -f github.com/golang/lint/golint
else
	@echo "Not installing golint, since we don't expect to lint on" $(GO_VERSION)
endif


# Disable printf-like invocation checking due to testify.assert.Error()
VET_RULES := -printf=false

.PHONY: lint
lint:
ifdef SHOULD_LINT
	@rm -rf lint.log
	@echo "Checking formatting..."
	@gofmt -d -s $(PKG_FILES) 2>&1 | tee lint.log
	@echo "Installing test dependencies for vet..."
	@go test -i $(PKGS)
	@echo "Checking vet..."
	@$(foreach dir,$(PKG_FILES),go tool vet $(VET_RULES) $(dir) 2>&1 | tee -a lint.log;)
	@echo "Checking lint..."
	@$(foreach dir,$(PKGS),golint $(dir) 2>&1 | tee -a lint.log;)
	@echo "Checking for unresolved FIXMEs..."
	@git grep -i fixme | grep -v -e vendor -e Makefile | tee -a lint.log
	#@echo "Checking for license headers..."
	#@./check_license.sh | tee -a lint.log
	@[ ! -s lint.log ]
else
	@echo "Skipping linters on" $(GO_VERSION)
endif

.PHONY: coveralls
coveralls:
	goveralls -service=travis-ci .

.PHONY: test
test:
	go test -race $(PKGS)


GO_EXECUTABLE ?= go
DIST_DIRS := find * -type d -exec
GIT_COMMIT=`git rev-parse --short HEAD`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
# NOTE: the `git tag` command is filtered through `grep .` so it returns non-zero when empty
GIT_TAG=`git tag --list "v*" --sort "v:refname" --points-at HEAD 2>/dev/null | tail -n 1 | grep . || echo "none"`
VERSION := $$(git describe --abbrev=0 --tags)

# Build Flags
BUILD_NUMBER ?= $(BUILD_NUMBER:)
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)
# If we don't set the build number it defaults to dev
ifeq ($(BUILD_NUMBER),)
	BUILD_NUMBER := dev
endif


# Output paths
DIST_ROOT=dist
DIST_PATH=$(DIST_ROOT)/live-go

# Tests
TESTS=.

build:
	${GO_EXECUTABLE} build -o ${NAME} -ldflags "-X main.version=${VERSION}" main.go

fmt:
	gofmt -w=true -s $$(find . -type f -name '*.go')
	goimports -w=true -d $$(find . -type f -name '*.go')

dist: | test


version:
	@echo $(REPO_VERSION)

clean:
	rm -f build/bin/*
