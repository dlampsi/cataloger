PROJECT_NAME=cataloger
RELEASE?=0.0.0
BUILD_TIME?=$(shell date '+%Y-%m-%d_%H:%M:%S')
OUTPUT?=${PROJECT_NAME}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build

clean:
	@go clean

tests: clean ## Run tests
	go test -v

build: clean ## Build package
	@go build \
	-ldflags "-s -w -X ${PROJECT_NAME}/info.Version=${RELEASE} \
	-X ${PROJECT_NAME}/info.BuildTime=${BUILD_TIME}" -o ${OUTPUT}
