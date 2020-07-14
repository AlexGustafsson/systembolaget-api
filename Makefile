# Disable echoing of commands
MAKEFLAGS += --silent

PREFIX := $(shell go list ./version)
VERSION := 0.1.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
GO_VERSION := $(shell go version)
COMPILE_TIME := $(shell LC_ALL=en_US date)

BUILD_VARIABLES := -X "$(PREFIX).Version=$(VERSION)" -X "$(PREFIX).Commit=$(COMMIT)" -X "$(PREFIX).GoVersion=$(GO_VERSION)" -X "$(PREFIX).CompileTime=$(COMPILE_TIME)"
BUILD_FLAGS := -ldflags '$(BUILD_VARIABLES)'

source := $(shell find . -type f -name '*.go')

.PHONY: build clean format lint

build: build/systembolaget

build/systembolaget: $(source) Makefile
	go build $(BUILD_FLAGS) -o $@ cmd/systembolaget/systembolaget.go

format: $(source) Makefile
	gofmt -l -s -w .

lint: $(source) Makefile
	golint .

clean:
	rm -rf build
