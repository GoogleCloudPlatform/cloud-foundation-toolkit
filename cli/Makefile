SHELL := /bin/bash

# Changing this value will trigger a new release
VERSION=v1.1.11
BINARY=bin/cft
GITHUB_REPO=github.com/GoogleCloudPlatform/cloud-foundation-toolkit
PLATFORMS := linux windows darwin
BUILD_DIR=./bin
NAME=cft
BUCKET=gs://cft-cli
INT_TEST_DIR=./bpmetadata/int-test
SCHEMA_DIR=./bpmetadata/schema
PROTO_DIR=./bpmetadata/proto
DOCKER_TAG_VERSION_DEVELOPER_TOOLS := 1.14
DOCKER_IMAGE_DEVELOPER_TOOLS := cft/developer-tools
REGISTRY_URL := gcr.io/cloud-foundation-cicd

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X $(GITHUB_REPO)/cli/cmd.Version=$(VERSION)"

.PHONY: build
build: protoc build-go

.PHONY: build-go
build-go:
	go run ./${SCHEMA_DIR} -output=${SCHEMA_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${NAME}

.PHONY: protoc
protoc:
	docker run --rm -it \
		-v "$(CURDIR)":/workspace \
		$(REGISTRY_URL)/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
		mkdir -p ${PROTO_DIR}/out && protoc --go_out=${PROTO_DIR}/out --go_opt=paths=source_relative ${PROTO_DIR}/*.proto

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
