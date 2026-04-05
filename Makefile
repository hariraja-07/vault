.PHONY: build build-all clean test install help

BINARY_NAME=vault
BUILD_DIR=.

help: ## Show this help
	@echo "Vault Build Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z_-]+:.*##.*$$/ { printf "  %-20s %s\n", $$1, $$3 }' $(MAKEFILE_LIST)

build: ## Build for current platform
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/vault

build-all: ## Build for all platforms (linux, macos, windows)
	GOOS=linux GOARCH=amd64 go build -o vault-linux ./cmd/vault
	GOOS=darwin GOARCH=amd64 go build -o vault-macos-intel ./cmd/vault
	GOOS=darwin GOARCH=arm64 go build -o vault-macos-arm ./cmd/vault
	GOOS=windows go build -o vault-windows.exe ./cmd/vault

clean: ## Remove built binaries
	rm -f vault-linux vault-macos-intel vault-macos-arm vault-windows.exe $(BINARY_NAME)

test: ## Run tests
	go test ./...

install: ## Install to /usr/local/bin
	go install ./cmd/vault
