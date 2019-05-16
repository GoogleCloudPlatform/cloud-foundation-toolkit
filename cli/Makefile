SHELL   := /bin/bash

# need to change this to VERSION=`git describe --tags`
VERSION=1.0.0
BINARY=bin/cft


# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=v1.0.0"

build:
	go build ${LDFLAGS} -o bin/cft