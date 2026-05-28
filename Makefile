BINARY  := grn
GO_DIR  := go
VERSION := $(shell grep 'const cliVersion' $(GO_DIR)/cmd/root.go | sed 's/.*"\(.*\)".*/\1/')

.PHONY: build build-all clean install fmt vet test lint help

## build: build binary for current OS/arch → ./grn
build:
	cd $(GO_DIR) && CGO_ENABLED=0 go build -trimpath -o ../$(BINARY) .

## build-all: cross-compile all release targets
build-all:
	cd $(GO_DIR) && CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -trimpath -o ../$(BINARY)-linux-amd64 .
	cd $(GO_DIR) && CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -trimpath -o ../$(BINARY)-linux-arm64 .
	cd $(GO_DIR) && CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -trimpath -o ../$(BINARY)-darwin-amd64 .
	cd $(GO_DIR) && CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -trimpath -o ../$(BINARY)-darwin-arm64 .
	cd $(GO_DIR) && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -o ../$(BINARY)-windows-amd64.exe .

## install: build and install to /usr/local/bin
install: build
	sudo mv $(BINARY) /usr/local/bin/$(BINARY)
	@echo "Installed $(BINARY) $(VERSION) → /usr/local/bin/$(BINARY)"

## fmt: format Go source
fmt:
	cd $(GO_DIR) && go fmt ./...

## vet: run go vet
vet:
	cd $(GO_DIR) && go vet ./...

## test: run tests
test:
	cd $(GO_DIR) && go test ./...

## clean: remove built binaries
clean:
	rm -f $(BINARY) $(BINARY)-linux-amd64 $(BINARY)-linux-arm64 \
	      $(BINARY)-darwin-amd64 $(BINARY)-darwin-arm64 \
	      $(BINARY)-windows-amd64.exe

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
