default: build

test:
	go test ./... -v

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-blueprint ~/.tflint.d/plugins
