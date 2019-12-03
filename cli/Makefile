SHELL := /bin/bash

# need to change this to VERSION=`git describe --tags`
VERSION=v0.3.0
BINARY=bin/cft
GITHUB_REPO=github.com/GoogleCloudPlatform/cloud-foundation-toolkit
PLATFORMS := linux windows darwin
BUILD_DIR=./bin
NAME=cft
BUCKET=gs://cft-cli

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X $(GITHUB_REPO)/cli/cmd.Version=$(VERSION)"

.PHONY: build
build:
	go build ${LDFLAGS} -o ${BUILD_DIR}/${NAME}

.PHONY: publish
publish:
	gsutil cp "${BUILD_DIR}/*" "${BUCKET}/${VERSION}"
	gsutil cp "${BUILD_DIR}/*" "${BUCKET}/latest"

.PHONY: release
release: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GO111MODULE=on GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o "${BUILD_DIR}/${NAME}-$@-amd64"

.PHONY: test
test:
	go test ./...
