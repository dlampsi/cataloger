.DEFAULT_GOAL := help

PROJECT_NAME=cataloger
BINARY_NAME?=cataloger
RELEASE?=0.0.0
BUILD_TIME?=$(shell date '+%Y-%m-%d_%H:%M:%S')

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: tests build

clean: ## Go clean
	@go clean

.PHONY: tests
tests: ## Run tests
	go test -race -coverprofile=coverage.out ./...

.PHONY: build
build: clean ## Build binary file
	@go build \
	-ldflags "-s -w -X ${PROJECT_NAME}/info.Version=${RELEASE} \
	-X ${PROJECT_NAME}/info.BuildNumber=${BUILD_NUMBER} \
	-X ${PROJECT_NAME}/info.BuildTime=${BUILD_TIME} \
	-X ${PROJECT_NAME}/info.CommitHash=${COMMIT_HASH}" -o ${BINARY_NAME}
