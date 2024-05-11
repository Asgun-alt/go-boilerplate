# Variabel
TEST?=./...

# Target
default: test

test:
	@echo "==> Running tests..."
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	open cover.html
	
fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

vet:
	@echo "==> Verifying with go vet..."
	go vet $(TEST)

lint:
	@echo "==> Linting..."
	golangci-lint run --fix

.PHONY: default test testcover fmt vet lint cover

################
# Migration
################
mig-build:
	@echo ">>> Building migration..."
	@go build -o bin/migration ./cmd/migration

mig-up: mig-build
	@echo ">>> executing migration..."
	@./bin/migration migrate up
	@echo ">>> finished executing migration..."

mig-down: mig-build
	@echo ">>> Rolling back migration 1 version..."
	@./bin/migration migrate down
	@echo ">>> finished rolling bank migration 1 version..."

mig-down-all: mig-build
	@echo ">>> resetting migration..."
	@./bin/migration migrate reset
	@echo ">>> finished resetting migration..."

mig-status: mig-build
	@echo ">>> Migration Status"
	@./bin/migration migrate status

mig-create: mig-build
	@echo ">>> Create Migration"
	@./bin/migration migrate create $(migration_name) go

mig-create-sql: mig-build
	@echo ">>> Create SQL Migration"
	@./bin/migration migrate create $(name) sql

################
# BUILD BINARY
################

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 1
GIT_URL := $(shell git config --get remote.origin.url)
GIT_COMMIT_HASH := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git branch --show-current)
GIT_TAG := $(shell git tag --points-at HEAD)
GIT_TAG := $(or $(GIT_TAG),$(GIT_BRANCH))
BUILD_OS := $(shell uname -rns | sed -e 's/ /_/g') # replace space with _
BUILD_TIME := $(shell date -u +%FT%T%Z)
# GO_MOD_NAME=$(shell go list -m)
# default to static build
GO_BUILD_FLAGS=-trimpath -a -tags "osusergo,netgo" -ldflags '-extldflags=-static -w -s -v' \
-ldflags "-X '$(GO_APP_INFO_MOD_NAME)/appinfo.GitURL=$(GIT_URL)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.GitCommitHash=$(GIT_COMMIT_HASH)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.GitTag=$(GIT_TAG)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.BuildTime=$(BUILD_TIME)' \
-X '$(GO_APP_INFO_MOD_NAME)/appinfo.BuildOS=$(BUILD_OS)'"

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) -o build/go-boilerplate main.go
