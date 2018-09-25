################################################################################
# Codeworx Makefile for vGo
################################################################################
SHELL := /bin/bash -o pipefail

# To override the variables in this section, define them in the project's
# Makefile before including this file.

ifeq ($(.DEFAULT_GOAL),)
.DEFAULT_GOAL := test
endif

# Project name is a short name for the project, used to build the names of
# release archives, for example.
PROJECT_NAME ?= $(notdir $(shell pwd))

# Disable CGO by default.
# See https://golang.org/cmd/cgo
CGO_ENABLED ?= 0

# Build matrix configuration.
#
# MATRIX_OS is a whitespace separated set of operating systems.
# MATRIX_ARCH is a whitespace separated set of CPU architectures.
#
# The build-matrix is constructed from all possible permutations of MATRIX_OS and
# MATRIX_ARCH. The default is to build only for the current OS and architecture.
-include artifacts/make/vgo/env.in
MATRIX_OS   ?= $(GOHOSTOS)
MATRIX_ARCH ?= $(GOHOSTARCH)

# Arguments passed to "go build" for debug / release builds.
DEBUG_ARGS   ?= -v
RELEASE_ARGS ?= -v -ldflags "-s -w"

# Arguments passed to "go test".
TEST_ARGS ?=

# The REQ and USE variables are used as pre-requisites for all targets that
# execute the project's Go source, including tests.
#
# They can be used to specify additional build dependencies other than the .go
# source files, such as HTML assets, etc.
#
# REQ is used a "normal" pre-requisite, whereas USE is an "order-only"
# pre-requisite.
#
# See https://www.gnu.org/software/make/manual/html_node/Prerequisite-Types.html
REQ ?=
USE ?=

# If running on a CI server, we need to do some work before installing
# dependencies
ifdef CI
_DEP_REQ += artifacts/logs/go-ci-install
endif

################################################################################
# Internal variables
################################################################################

# _MATRIX contains all permutations of MATRIX_ARCH and MATRIX_OS
_MATRIX ?= $(foreach OS,$(MATRIX_OS),$(foreach ARCH,$(MATRIX_ARCH),$(OS)/$(ARCH)))
_MATRIX_NIX := $(filter-out windows/%,$(_MATRIX))
_MATRIX_WIN := $(filter windows/%,$(_MATRIX))

# _SRC contains the paths to all Go source files.
_SRC ?= $(shell find ./src -name *.go 2> /dev/null)

# _PKGS contains the paths to all Go packages under ./src
_PKGS ?= $(sort $(dir $(_SRC)))

# _BINS contains the names of binaries (directories under ./src/cmd)
_BINS ?= $(notdir $(shell find src/cmd -mindepth 1 -maxdepth 1 -type d 2>/dev/null))
_BINS_WIN := $(addsuffix .exe,$(_BINS))

# _STEMS contains the binary names for each entry in the MATRIX (e.g. darwin/amd64/<bin>)
_STEMS_NIX := $(foreach B,$(_MATRIX_NIX),$(foreach BIN,$(_BINS),$(B)/$(BIN)))
_STEMS_WIN := $(foreach B,$(_MATRIX_WIN),$(foreach BIN,$(_BINS_WIN),$(B)/$(BIN)))
_STEMS ?= $(_STEMS_NIX) $(_STEMS_WIN)

# _COV contains the paths to a "cover.out" file for each package.
_COV ?= $(foreach P,$(_PKGS),artifacts/tests/coverage/$(P)cover.out)

# Add the linux release build to the docker requirements, for use with the
# Docker Makefile.
DOCKER_REQ ?= $(addprefix artifacts/build/release/linux/amd64/,$(_BINS))

################################################################################
# Commands (Phony Targets)
################################################################################

# Run all tests.
.PHONY: test
test: vendor $(REQ) | $(USE)
	go test $(TEST_ARGS) ./src/...

# Build debug executables for the current OS and architecture.
.PHONY: build
build: $(addprefix artifacts/build/debug/$(GOOS)/$(GOARCH)/,$(_BINS))

# Build debug executables for all OS and architecture combinations.
.PHONY: debug
.SECONDARY: $(addprefix artifacts/build/debug/,$(_STEMS))
debug: $(addprefix artifacts/build/debug/,$(_STEMS))

# Build release executables for all OS and architecture combinations.
.PHONY: release
.SECONDARY: $(addprefix artifacts/build/release/,$(_STEMS))
release: $(addprefix artifacts/build/release/,$(_STEMS))

.PHONY: archives
archives: $(addprefix artifacts/archives/$(PROJECT_NAME)-,$(addsuffix .tar.gz,$(subst /,-,$(_MATRIX_NIX))))
archives: $(addprefix artifacts/archives/$(PROJECT_NAME)-,$(addsuffix .zip,$(subst /,-,$(_MATRIX_WIN))))

# Remove all files that match the patterns .gitignore.
.PHONY: clean-all
clean-all:: clean

# Remove files that match the patterns .gitignore, excluding the vendor folder.
.PHONY: clean
clean::
	@git check-ignore ./* | grep -v ^./vendor | xargs -t -n1 rm -rf

.PHONY: clean-coverage
clean-coverage:
	rm -rf artifacts/tests/coverage

# Generate an HTML code coverage report.
.PHONY: coverage
coverage: artifacts/tests/coverage/index.html

# Generate an HTML code coverage report and open it in the browser.
# TODO: This command only works on OSX.
.PHONY: coverage-open
coverage-open: artifacts/tests/coverage/index.html
	open "$<"

# Perform code linting, syntax formatting, etc.
.PHONY: lint
lint:: artifacts/logs/lint

# Perform pre-commit checks.
.PHONY: prepare
prepare: lint artifacts/logs/errcheck test

# Run the CI build.
.PHONY: ci
ci:: ci-lint $(_COV)
	go test -race ./src/...

# A hook for adding linting to CI builds
.PHONY: ci-lint
ci-lint:: lint

# Configure the CI environment before deploying
.PHONY: ci-before-deploy
ci-before-deploy::

.PHONY:
usage::
	@cat artifacts/make/vgo/usage.txt
	@echo

################################################################################
# File Targets
################################################################################

.DELETE_ON_ERROR:
.SECONDEXPANSION:

GOCOVMERGE := $(GOPATH)/bin/gocovmerge
$(GOCOVMERGE):
	export GO111MODULE=off
	go install github.com/wadey/gocovmerge

MISSPELL := $(GOPATH)/bin/misspell
$(MISSPELL):
	export GO111MODULE=off
	go install github.com/client9/misspell/cmd/misspell

ERRCHECK := $(GOPATH)/bin/errcheck
$(ERRCHECK):
	export GO111MODULE=off
	go install github.com/kisielk/errcheck

GOMETALINTER := $(GOPATH)/bin/gometalinter
.PRECIOUS: $(GOMETALINTER)
$(GOMETALINTER):
	export GO111MODULE=off
	# export BINDIR=$(GOPATH)/bin
	# curl https://git.io/vp6lP | sh
	go install github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install 2>/dev/null

GOLANGCILINT := $(GOPATH)/bin/golangci-lint
.PRECIOUS: $(GOLANGCILINT)
$(GOLANGCILINT):
	export GO111MODULE=off
	go install -i github.com/golangci/golangci-lint/cmd/golangci-lint


go.mod: $(_DEP_REQ)
ifneq (,$(wildcard ./go.mod))
	go mod init
endif
	go mod tidy

go.sum: go.mod
	go mod verify


.PHONY: vendor
vendor: go.mod

artifacts/build/%: vendor $(REQ) $(_SRC) | $(USE)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(word 3,$(PARTS)))
	$(eval BIN   := $(word 4,$(PARTS)))
	$(eval PKG   := $(basename $(BIN)))
	$(eval ARGS  := $(if $(findstring debug,$(BUILD)),$(DEBUG_ARGS),$(RELEASE_ARGS)))

	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(OS)" GOARCH="$(ARCH)" go build $(ARGS) -o "$@" "./src/cmd/$(PKG)"

artifacts/archives/$(PROJECT_NAME)-%.tar.gz: $$(addprefix artifacts/build/release/$$(subst -,/,$$*)/,$(_BINS))
	@mkdir -p "$(@D)"
	cd "$(dir $<)" && tar -cvzf archive.tmp $(_BINS)
	mv "$(dir $<)archive.tmp" "$@"

artifacts/archives/$(PROJECT_NAME)-%.zip: $$(addprefix artifacts/build/release/$$(subst -,/,$$*)/,$(_BINS_WIN))
	@mkdir -p "$(@D)"
	cd "$(dir $<)" && zip archive.tmp $(_BINS_WIN)
	mv "$(dir $<)archive.tmp" "$@"

artifacts/tests/coverage/index.html: artifacts/tests/coverage/merged.cover.out
	go tool cover -html="$<" -o "$@"

artifacts/tests/coverage/merged.cover.out: $(_COV) | $(GOCOVMERGE)
	@mkdir -p $(@D)
	$(GOCOVMERGE) $^ > "$@"

%/cover.out: vendor $$(subst artifacts/tests/coverage/,,$$(@D))/*.go $(REQ) | $(USE)
	$(eval PKG := $(subst artifacts/tests/coverage/,,$*))
	@mkdir -p "$(@D)"
	@touch "$@" # no file is written if there are no tests
	-go test $(TEST_ARGS) -covermode=count -coverprofile="$@" "$(PKG)"
	-go tool cover -func="$@"

artifacts/logs/lint: vendor $(_SRC) $(REQ) | $(MISSPELL) $(GOLANGCILINT) $(USE)
	@mkdir -p "$(@D)"

	go vet ./src/... | tee "$@"
	! go fmt ./src/... | tee -a "$@" | grep ^

	$(MISSPELL) -w -error -locale US ./src | tee -a "$@"

	@# $(GOMETALINTER) --disable-all --deadline=60s \
	@# 	--enable=vet \
	@# 	--enable=vetshadow \
	@# 	--enable=ineffassign \
	@# 	--enable=deadcode \
	@# 	--enable=gosimple \
	@# 	--enable=gofmt \
	@# 	./src/... | tee -a "$@"

	$(GOLANGCILINT) run ./src/... | tee -a "$@"

	-$(GOMETALINTER) --disable-all --deadline=60s --cyclo-over=15 \
		--enable=golint \
		--enable=goconst \
		--enable=gocyclo \
		./src/... | tee -a "$@"

artifacts/logs/errcheck: vendor $(wildcard .errignore) $(_SRC) $(REQ) | $(ERRCHECK) $(USE)
	@mkdir -p "$(@D)"
ifeq (,$(wildcard .errignore))
	-$(ERRCHECK) -ignoretests ./src/... | tee "$@"
else
	-$(ERRCHECK) -ignoretests -exclude .errignore ./src/... | tee "$@"
endif

# extract the go environment values, even if they're not explicitly set
# as environment variables, and convert them to Makefile format.
artifacts/make/vgo/env.in:
	env -i $(shell go env) | sed -n 's/^\([^=]*\)=\(..*\)/\1 ?= \2/gp' > "$@"

artifacts/logs/go-ci-install:
	@mkdir -p "$(@D)"
	git config --global url."https://$(GITHUB_TOKEN)@github.com/".insteadOf "git@github.com:" | tee "$@"
