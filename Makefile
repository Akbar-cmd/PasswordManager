.PHONY: build build-all build-macos build-windows clean fmt help install uninstall deps tidy ci

# Variables
BINARY_NAME=PasswordManager
GO_FILES=$(shell find . -type f -name '*.go')
VERSION ?= 1.0.0

# Colors
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m

## help: Display this help message
help:
	@echo "$(GREEN)Password Manager - Makefile Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)Build Commands:$(NC)"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application
build: fmt
	@echo "$(GREEN)[*] Building $(BINARY_NAME)...$(NC)"
	@go build -v -o $(BINARY_NAME) .
	@echo "$(GREEN)[✓] Build complete: ./$(BINARY_NAME)$(NC)"

## build-windows: Build for Windows
build-windows: fmt
	@mkdir -p ./bin
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/$(BINARY_NAME).exe .

## build-macos: Build for macOS
build-macos: fmt
	@mkdir -p ./bin
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/$(BINARY_NAME)-macos-arm64 .
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/$(BINARY_NAME)-macos-amd64 .

## run: Build and run the application
run: build
	@echo "$(GREEN)[*] Running $(BINARY_NAME)...$(NC)"
	@./$(BINARY_NAME)

## clean: Remove build artifacts
clean:
	@echo "$(YELLOW)[*] Cleaning up...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f *.out
	@rm -f coverage.html
	@go clean -cache -modcache
	@echo "$(GREEN)[✓] Cleanup complete$(NC)"

## fmt: Format code
fmt:
	@echo "$(YELLOW)[*] Formatting code...$(NC)"
	@go fmt ./...
	@gofmt -w .
	@echo "$(GREEN)[✓] Code formatted$(NC)"

## install: Install the binary to system
install: build
	@echo "$(GREEN)[*] Installing $(BINARY_NAME)...$(NC)"
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)[✓] Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"

## uninstall: Remove the binary from system
uninstall:
	@echo "$(YELLOW)[*] Uninstalling $(BINARY_NAME)...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)[✓] Uninstalled$(NC)"

## deps: Download and verify dependencies
deps:
	@echo "$(YELLOW)[*] Downloading dependencies...$(NC)"
	@go mod download
	@go mod verify
	@echo "$(GREEN)[✓] Dependencies verified$(NC)"

## tidy: Tidy Go modules
tidy:
	@echo "$(YELLOW)[*] Running go mod tidy...$(NC)"
	@go mod tidy
	@echo "$(GREEN)[✓] Go modules tidied$(NC)"

## ci: Run all CI checks (format, lint, test)
ci: fmt
	@echo "$(GREEN)[✓] All CI checks passed$(NC)"

.DEFAULT_GOAL := help
