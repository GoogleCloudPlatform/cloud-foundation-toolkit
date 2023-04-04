SHELL := /bin/bash

# Changing this value will trigger a new release
VERSION=v0.9.1
BINARY=bin/cft
GITHUB_REPO=github.com/GoogleCloudPlatform/cloud-foundation-toolkit
PLATFORMS := linux windows darwin
BUILD_DIR=./bin
SCHEMA_DIR=./bpmetadata/schema
NAME=cft
BUCKET=gs://cft-cli

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X $(GITHUB_REPO)/cli/cmd.Version=$(VERSION)"

# need an extra build command to embed an most up to date schema in the binary
.PHONY: build
build:
	go run ./${SCHEMA_DIR} -output=${SCHEMA_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${NAME}
	
.PHONY: publish
publish:
	gcloud alpha storage cp "${BUILD_DIR}/*" "${BUCKET}/${VERSION}"
	gcloud alpha storage cp "${BUILD_DIR}/*" "${BUCKET}/latest"

.PHONY: release
release: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GO111MODULE=on GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o "${BUILD_DIR}/${NAME}-$@-amd64"

.PHONY: go_test
go_test:
	go test ./...

.PHONY: test
test: build go_test
