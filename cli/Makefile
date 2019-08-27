SHELL := /bin/bash

# need to change this to VERSION=`git describe --tags`
VERSION=v0.1.0
BINARY=bin/cft
GITHUB_REPO=github.com/GoogleCloudPlatform/cloud-foundation-toolkit

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X $(GITHUB_REPO)/cli/cmd.Version=$(VERSION)"

build:
	go build ${LDFLAGS} -o bin/cft

.PHONY: test
test:
	go test ./...