MAKEFLAGS += --warn-undefined-variables

NC    :=\033[0m
BLUE  :=\033[36m
GREEN :=\033[0;32m

APP       := cataloger
BINARY    ?= $(APP)
# BINARY    ?= $(APP)-$(RELEASE)_$(GOARCH)_$(GOOS)
RELEASE   ?= $(shell git rev-parse --short HEAD)
COMMIT    := $(shell git rev-parse HEAD)
DATE      ?= $(shell date '+%Y-%m-%d_%H:%M:%S')
BUILD_NUM ?= 

BIN_DIR   := bin
TESTS_DIR := tests

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

default: help

help: ## Display available make commands
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE)%-12s$(NC) %s\n", $$1, $$2}'

all: deps prep lint tests build ## All included
	@echo "$(GREEN)All done!$(NC)"

deps: ## Checks all required dependencies to be installed
	@echo "$(BLUE)• Checking dependencies$(NC)"
	@which go
	@which golangci-lint

tests: ## Run all tests
	@echo "$(BLUE)• Running unit tests$(NC)"
	go test -race -coverprofile=$(TESTS_DIR)/coverage.out ./...
.PHONY: tests

lint: ## Applies linter
	@echo "$(BLUE)• Running linter$(NC)"
	golangci-lint run -c ./.golangci.yml --timeout 3m ./...

prep: ## Performs all required build preparations.
	@mkdir -p tests
	@mkdir -p $(BIN_DIR)

clean: ## Clean before build
	@echo "$(BLUE)• Clean before build$(NC)"
	go clean
	rm -f $(BIN_DIR)/$(BINARY)

build: clean ## Build binary
	@echo "$(BLUE)• Building binary$(NC)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-X '$(APP)/info.Release=$(RELEASE)' \
			-X '$(APP)/info.BuildTime=$(DATE)' \
			-X '$(APP)/info.CommitHash=$(COMMIT)' \
			-X '$(APP)/info.BuildNumber=$(BUILD_NUM)'"\
		-o $(BIN_DIR)/$(BINARY)
	@chmod +x $(BIN_DIR)/$(BINARY)
