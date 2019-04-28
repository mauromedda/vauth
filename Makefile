.PHONY: fmt build
# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GO_VERSION_MIN=1.11
LD_FLAGS="-s -w"
fmt:
	gofmt -w $(GOFMT_FILES)

build:
	go build -ldflags=$(LD_FLAGS)

