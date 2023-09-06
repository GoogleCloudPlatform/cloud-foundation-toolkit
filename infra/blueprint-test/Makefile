SHELL := /bin/bash

.PHONY: docker_go_lint
docker_go_lint:
	docker run --rm -v $(PWD):/blueprint-test -w /blueprint-test  golangci/golangci-lint:v1.54.1 golangci-lint --timeout=5m -v run
