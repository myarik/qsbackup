PKGS := $(shell go list ./... | grep -v /vendor)

GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
BINARY_NAME=qsbackup

help: _help_

_help_:
	@echo make test -- run test

.PHONY: test
test:
	$(GOTEST) -v $(PKGS)