SHELL := /bin/bash

# Changing this value will trigger a new release
VERSION=v1.1.9
BINARY=bin/cft
GITHUB_REPO=github.com/GoogleCloudPlatform/cloud-foundation-toolkit
PLATFORMS := linux windows darwin
BUILD_DIR=./bin
NAME=cft
BUCKET=gs://cft-cli
INT_TEST_DIR=./bpmetadata/int-test
SCHEMA_DIR=./bpmetadata/schema
PROTO_DIR=./bpmetadata/proto

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X $(GITHUB_REPO)/cli/cmd.Version=$(VERSION)"

# need an extra build command to embed an most up to date schema in the binary
.PHONY: build
build:
	go run ./${SCHEMA_DIR} -output=${SCHEMA_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${NAME}

.PHONY: protoc
protoc:
	protoc --go_out=${PROTO_DIR}/out --go_opt=paths=source_relative ${PROTO_DIR}/*.proto
	go build ${PROTO_DIR}/out/...

.PHONY: publish
publish:
	gcloud alpha storage cp "${BUILD_DIR}/*" "${BUCKET}/${VERSION}"
	gcloud alpha storage cp "${BUILD_DIR}/*" "${BUCKET}/latest"

.PHONY: release
release: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GO111MODULE=on GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o "${BUILD_DIR}/${NAME}-$@-amd64"

.PHONY: int_test
int_test:
	${INT_TEST_DIR}/workflow.sh ${INT_TEST_DIR}

.PHONY: go_test
go_test:
	go test ./...

.PHONY: test
test: build go_test int_test

.PHONY: docker_go_lint
docker_go_lint:
	docker run --rm -v $(PWD):/cli -w /cli  golangci/golangci-lint:v1.52.2 golangci-lint --timeout=5m -v run
