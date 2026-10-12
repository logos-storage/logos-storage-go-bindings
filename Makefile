# Makefile for Logos Storage Go Bindings

LIBS_DIR := $(abspath ./libs)

CGO_CFLAGS  := -I$(LIBS_DIR)
CGO_LDFLAGS := -L$(LIBS_DIR) -lstorage -Wl,-rpath,$(LIBS_DIR)

# Fetch configuration
# OS is not set with ?= because Windows defines OS=Windows_NT in the environment
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
VERSION ?= v0.4.3
DOWNLOAD_URL := https://github.com/logos-storage/logos-storage-nim/releases/download/$(VERSION)/libstorage-$(OS)-$(ARCH)-$(VERSION).zip

.PHONY: all clean fetch build test

all: build

fetch:
	@echo "Fetching libstorage from $(DOWNLOAD_URL)"
	rm -Rf $(LIBS_DIR)/*
	curl -fSL --create-dirs -o $(LIBS_DIR)/libstorage-$(OS)-$(ARCH)-$(VERSION).zip $(DOWNLOAD_URL)
	unzip -o -qq $(LIBS_DIR)/libstorage-$(OS)-$(ARCH)-$(VERSION).zip -d $(LIBS_DIR)
	rm -f $(LIBS_DIR)/*.zip

build:
	@echo "Building Logos Storage Go Bindings..."
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o storage-go ./storage

test:
	@echo "Running tests..."
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOTESTFLAGS="-timeout=2m" gotestsum --packages="./..." -f testname -- $(if $(filter-out test,$(MAKECMDGOALS)),-run "$(filter-out test,$(MAKECMDGOALS))")

test-with-params:
	@echo "Running tests..."
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOTESTFLAGS="-timeout=2m" gotestsum --packages="./..." -f testname -- $(ARGS)

%:
	@:

clean:
	@echo "Cleaning up..."
	@rm -Rf $(LIBS_DIR)/*
	@rm -f storage-go