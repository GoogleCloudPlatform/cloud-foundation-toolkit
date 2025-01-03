SHELL := /bin/bash

# Changing this value will trigger a new release
VERSION=v1.5.12
BINARY=bin/cft
GITHUB_REPO=github.com/GoogleCloudPlatform/cloud-foundation-toolkit
PLATFORMS := linux windows darwin
BUILD_DIR=./bin
NAME=cft
BUCKET=gs://cft-cli
INT_TEST_DIR=./bpmetadata/int-test
SCHEMA_DIR=./bpmetadata/schema
SRC_PROTO_DIR=./bpmetadata/proto
PROTOC_DIR=./bpmetadata
DOCKER_TAG_VERSION_DEVELOPER_TOOLS := 1.23
DOCKER_IMAGE_DEVELOPER_TOOLS := cft/developer-tools
REGISTRY_URL := gcr.io/cloud-foundation-cicd
GOLANGCI_VERSION := 1.63.3

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X $(GITHUB_REPO)/cli/cmd.Version=$(VERSION)"

.PHONY: build
build: protoc-gen build-schema
	go build ${LDFLAGS} -o ${BUILD_DIR}/${NAME}

.PHONY: protoc-gen
protoc-gen:
	docker run --rm \
		-v "$(CURDIR)":/workspace \
		$(REGISTRY_URL)/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
		/bin/bash -c "protoc -I=${SRC_PROTO_DIR} --go_opt=paths=source_relative --go_out=${PROTOC_DIR} ${SRC_PROTO_DIR}/*.proto && \
		protoc-go-inject-tag -input='${PROTOC_DIR}/*.pb.go'"

.PHONY: build-schema
build-schema:
	go run ./${SCHEMA_DIR} -output=${SCHEMA_DIR}

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
go_test: build
	go test ./...

.PHONY: test
test: build go_test int_test

.PHONY: docker_go_lint
docker_go_lint:
	docker run --rm -v $(PWD):/cli -w /cli  golangci/golangci-lint:v${GOLANGCI_VERSION} golangci-lint --timeout=5m -v run
