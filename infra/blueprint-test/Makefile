SHELL := /bin/bash
GOLANGCI_VERSION := 1.63.3

.PHONY: docker_go_lint
docker_go_lint:
	docker run --rm -v $(PWD):/blueprint-test -w /blueprint-test  golangci/golangci-lint:v${GOLANGCI_VERSION} golangci-lint --timeout=5m -v run
