APP_NAME := go-loglint
LOGLINT_CONFIG ?= ./configs/loglint.yaml

.PHONY: help run-loglint build test clean

help: ## Show available commands
	@echo "Available targets:"
	@echo "  make run-loglint - Run loglint locally"
	@echo "  make build       - Build loglint binary"
	@echo "  make build-plugin - Build golangci-lint plugin (CGO_ENABLED=1)"
	@echo "  make test        - Run all tests"
	@echo "  make clean       - Remove build artifacts"

run-loglint: ## Run loglint with local config
	LOGLINT_CONFIG_PATH=$(LOGLINT_CONFIG) go run ./cmd

build: ## Build loglint binary
	go build -o ./bin/$(APP_NAME) ./cmd

build-plugin: ## Build golangci-lint plugin
	CGO_ENABLED=1 go build -buildmode=plugin -o ./bin/loglint.so ./plugin

test: ## Run all tests
	go test ./...

clean: ## Remove build artifacts
	rm -rf ./bin