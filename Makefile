# Makefile for Logos Storage Go Bindings

LOGOS_STORAGE_NIM_DIR := vendor/logos-storage-nim
LOGOS_STORAGE_NIM_LIB_DIR   := $(abspath $(LOGOS_STORAGE_NIM_DIR)/library)
LOGOS_STORAGE_NIM_BUILD_DIR := $(abspath $(LOGOS_STORAGE_NIM_DIR)/build)

CGO_CFLAGS  := -I$(LOGOS_STORAGE_NIM_LIB_DIR)
CGO_LDFLAGS := -L$(LOGOS_STORAGE_NIM_BUILD_DIR) -lstorage -Wl,-rpath,$(LOGOS_STORAGE_NIM_BUILD_DIR)

.PHONY: all clean update libstorage build test

all: build

submodules:
	@echo "Fetching submodules..."
	@git submodule update --init --recursive

update: | submodules
	@echo "Updating logos-storage-nim..."
	@$(MAKE) -C $(LOGOS_STORAGE_NIM_DIR) update

libstorage:
	@echo "Building libstorage..."
	@$(MAKE) -C $(LOGOS_STORAGE_NIM_DIR) libstorage

libstorage-with-debug-api:
	@echo "Building libstorage..."
	@$(MAKE) -C $(LOGOS_STORAGE_NIM_DIR) libstorage STORAGE_LIB_PARAMS="-d:storage_enable_api_debug_peers"

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
	@git submodule deinit -f $(LOGOS_STORAGE_NIM_DIR)
	@rm -f storage-go