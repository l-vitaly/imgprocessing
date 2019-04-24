COMMON_VERSION=v1.0.2

# Git current hash
GIT_HASH=$(shell git rev-parse --short HEAD)

# Git current tag
GIT_TAG=$(shell git tag -l --contains HEAD) 

# Git current branch 
GIT_BRANCH=$(shell git symbolic-ref HEAD | sed -e 's|^refs/heads/||' | sed -e 's|_.*||')

DEFAULT_BINARY=service

ifeq ($(CMD_PATH),)
$(error "not set CMD_PATH. please set CMD_PATH and retry now")
endif

ifeq ($(BINARY),)
BINARY:=$(DEFAULT_BINARY)
endif

ifeq ($(BINARY),)
$(error "not set BINARY. please set BINARY and retry now")
endif

.DEFAULT_GOAL := help

ifeq ($(GIT_TAG),)
TAG=$(GIT_BRANCH).$(GIT_HASH)
else
TAG=$(GIT_TAG)
endif

# Build date time
BUILD=$(shell date +%FT%T%z)

# Setup the -ldflags options
LDFLAGS=-ldflags "-extldflags "-static" -X main.Tag=$(TAG) -X main.Build=$(BUILD)"

up: build
	BINARY=${BINARY} docker-compose up -d --build
	rm ${BINARY}
	
down:	
	docker-compose down

build: ## Build the binary
	@echo "Build - $(TAG)"
	CGO_ENABLED=0 go build ${LDFLAGS} -a -installsuffix cgo -o ${BINARY} ${CMD_PATH}

check: lint vet ## Runs all tests

test: ## Run the unit tests
	go test -race -v $(shell go list ./... | grep -v /vendor/)

lint: ## Lint all files
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

vet: ## Run the vet tool
	go vet $(shell go list ./... | grep -v /vendor/)

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.SILENT: build test lint vet clean docker-build docker-push help
