# paths
makefile := $(realpath $(lastword $(MAKEFILE_LIST)))
cmd_dir  := ./cmd/unifi-mcp
dist_dir := ./dist

# parallel jobs for build-release (can be overridden)
JOBS ?= 8

# executables
GO    := go
GZIP  := gzip --best
MKDIR := mkdir -p
ZIP   := zip -m

# build flags
export CGO_ENABLED := 0
BUILD_FLAGS := -ldflags="-s -w" -trimpath

# release binaries
bin := unifi-mcp
releases :=                            \
	$(dist_dir)/$(bin)-darwin-amd64    \
	$(dist_dir)/$(bin)-darwin-arm64    \
	$(dist_dir)/$(bin)-linux-386       \
	$(dist_dir)/$(bin)-linux-amd64     \
	$(dist_dir)/$(bin)-linux-arm5      \
	$(dist_dir)/$(bin)-linux-arm6      \
	$(dist_dir)/$(bin)-linux-arm7      \
	$(dist_dir)/$(bin)-linux-arm64     \
	$(dist_dir)/$(bin)-netbsd-amd64    \
	$(dist_dir)/$(bin)-openbsd-amd64   \
	$(dist_dir)/$(bin)-solaris-amd64   \
	$(dist_dir)/$(bin)-windows-amd64.exe

## build: build an executable for your architecture
.PHONY: build
build: | clean $(dist_dir) fmt lint vet
	$(GO) build $(BUILD_FLAGS) -o $(dist_dir)/unifi-mcp $(cmd_dir)

## build-release: build release executables
# Runs prepare once, then builds all binaries in parallel
# Override jobs with: make build-release JOBS=16
.PHONY: build-release
build-release: prepare
	$(MAKE) -j$(JOBS) $(releases)

# unifi-mcp-darwin-amd64
$(dist_dir)/$(bin)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-darwin-arm64
$(dist_dir)/$(bin)-darwin-arm64:
	GOARCH=arm64 GOOS=darwin \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-linux-386
$(dist_dir)/$(bin)-linux-386:
	GOARCH=386 GOOS=linux \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-linux-amd64
$(dist_dir)/$(bin)-linux-amd64:
	GOARCH=amd64 GOOS=linux \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-linux-arm5
$(dist_dir)/$(bin)-linux-arm5:
	GOARCH=arm GOOS=linux GOARM=5 \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-linux-arm6
$(dist_dir)/$(bin)-linux-arm6:
	GOARCH=arm GOOS=linux GOARM=6 \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-linux-arm7
$(dist_dir)/$(bin)-linux-arm7:
	GOARCH=arm GOOS=linux GOARM=7 \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-linux-arm64
$(dist_dir)/$(bin)-linux-arm64:
	GOARCH=arm64 GOOS=linux \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-netbsd-amd64
$(dist_dir)/$(bin)-netbsd-amd64:
	GOARCH=amd64 GOOS=netbsd \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-openbsd-amd64
$(dist_dir)/$(bin)-openbsd-amd64:
	GOARCH=amd64 GOOS=openbsd \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-solaris-amd64
$(dist_dir)/$(bin)-solaris-amd64:
	GOARCH=amd64 GOOS=solaris \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# unifi-mcp-windows-amd64
$(dist_dir)/$(bin)-windows-amd64.exe:
	GOARCH=amd64 GOOS=windows \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(ZIP) $@.zip $@ -j

.PHONY: prepare
prepare: | clean $(dist_dir) fmt lint vet test

## install: build and install unifi-mcp on your PATH
.PHONY: install
install: build
	$(GO) install $(BUILD_FLAGS) $(cmd_dir)

## clean: remove compiled executables
.PHONY: clean
clean:
	rm -f $(dist_dir)/*

## generate: generate types and client from OpenAPI spec
.PHONY: generate
generate:
	oapi-codegen -generate types -package unifi -o internal/unifi/types.gen.go api/unifi-network.json
	oapi-codegen -generate client -package unifi -o internal/unifi/client.gen.go api/unifi-network.json

## fmt: format code with 80-column wrapping
.PHONY: fmt
fmt:
	$(GO) run github.com/segmentio/golines@latest -w --max-len=80 --ignore-generated .
	$(GO) run mvdan.cc/gofumpt@latest -w .

## lint: lint go source files
.PHONY: lint
lint:
	$(GO) run github.com/mgechev/revive@latest --exclude=./internal/unifi/... ./...

## vet: vet go source files
.PHONY: vet
vet:
	$(GO) vet ./...

## test: run tests
.PHONY: test
test:
	$(GO) test ./...

## coverage: generate test coverage report
.PHONY: coverage
coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out
	@echo ""
	@echo "HTML coverage report: coverage.html"
	$(GO) tool cover -html=coverage.out -o coverage.html

## check: format, lint, vet, and test
.PHONY: check
check: | fmt lint vet test

# ./dist
$(dist_dir):
	$(MKDIR) $(dist_dir)

## sloc: count source lines of code
.PHONY: sloc
sloc:
	scc --exclude-dir vendor .

## help: display this help text
.PHONY: help
help:
	@cat $(makefile) | \
	sort             | \
	grep "^##"       | \
	sed 's/## //g'   | \
	column -t -s ':'
