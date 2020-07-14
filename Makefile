# Disable echoing of commands
MAKEFLAGS += --silent

# Add build-time variables
PREFIX := $(shell go list ./version)
VERSION := 1.0.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
GO_VERSION := $(shell go version)
COMPILE_TIME := $(shell LC_ALL=en_US date)

BUILD_VARIABLES := -X "$(PREFIX).Version=$(VERSION)" -X "$(PREFIX).Commit=$(COMMIT)" -X "$(PREFIX).GoVersion=$(GO_VERSION)" -X "$(PREFIX).CompileTime=$(COMPILE_TIME)"
BUILD_FLAGS := -ldflags '$(BUILD_VARIABLES)'

source := $(shell find . -type f -name '*.go')

# Configure build output
binary = systembolaget
build = GOOS=$(1) GOARCH=$(2) go build $(BUILD_FLAGS) -o build/$(binary)$(3) ./cmd/systembolaget/systembolaget.go
tar = cd build && tar -czf $(1)_$(2).tar.gz $(binary)$(3) && rm $(binary)$(3)
zip = cd build && zip $(1)_$(2).zip $(binary)$(3) && rm $(binary)$(3)

.PHONY: build package format lint clean

build: build/systembolaget

# Build for the native platform. For cross-platform builds, see "package" below
build/systembolaget: $(source) Makefile
	go build $(BUILD_FLAGS) -o $@ cmd/systembolaget/systembolaget.go

format: $(source) Makefile
	gofmt -l -s -w .

lint: $(source) Makefile
	golint .

clean:
	rm -rf build

package: windows darwin linux

##
# Linux
linux: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(sources)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

##
# Windows
windows: build/windows_386.zip build/windows_amd64.zip

build/windows_386.zip: $(sources)
	$(call build,windows,386,.exe)
	$(call zip,windows,386,.exe)

build/windows_amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)

##
# Darwin
darwin: build/darwin_amd64.tar.gz

build/darwin_amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)
