.PHONY: all

all: fmt test build
# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

TEST?=$$(go list ./... | grep -v /vendor/ | grep -v /integ)
TEST_TIMEOUT?=30m
EXTENDED_TEST_TIMEOUT=45m

GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GO_VERSION_MIN=1.11
LD_FLAGS="-s -w"

fmt:
	gofmt -w $(GOFMT_FILES)

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| grep -v '.*github.com/hashicorp/vault$$' \
		| xargs go vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# prepare the build and test environment downloading external (not vendored dependencies)
prep:
	# gox simplifies building for multiple architectures
	echo => Install gox
	go get github.com/mitchellh/gox
	go get github.com/aktau/github-release

build: prep
	rm -rf bin/ pkg/
	echo => Build the binaries for the follownig OS Windows, Linux and Darwin on x64
	gox -os="linux darwin windows" -arch="amd64" -output="pkg/{{.OS}}_{{.Arch}}/vauth" -ldflags "-s -w -X main.version=$(shell git describe --tags || git rev-parse --short HEAD || echo dev)" -verbose ./...

test:
	VAULT_ADDR= \
	VAULT_TOKEN= \
	VAULT_DEV_ROOT_TOKEN_ID= \
	VAULT_ACC= \
	go test -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=$(TEST_TIMEOUT) -parallel=20 -v

release: prep
	./release.sh