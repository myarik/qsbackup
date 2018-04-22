PKGS := $(shell go list ./... | grep -v /vendor)

GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=qsbackup
MAIN_PKG=github.com/myarik/qsbackup/cmd/qsbackup
VERSION ?= vlatest
PLATFORMS := linux darwin
os = $(word 1, $@)
CONF ?= /Users/yaroslavmuravskiy/Documents/test_backup/config.yml

help: _help_

_help_:
	@echo make test -- run test
	@echo make build -- build binar pkg
	@echo make clean -- clean
	@echo make VERSION=v0.0.1 release -j2 -- build a release pkg with a version
	@echo make release -j2 -- build a release pkg
	@echo make CONF=/tmp/config.yml run


.PHONY: test
test:
	$(GOTEST) -v $(PKGS) && go vet && golint

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PKG)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f rm -rf release/$(BINARY_NAME)-vlatest-*

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY_NAME)-$(VERSION)-$(os)-amd64 $(MAIN_PKG)

.PHONY: release
release: linux darwin

run:
	go run cmd/qsbackup/main.go -c $(CONF) -debug