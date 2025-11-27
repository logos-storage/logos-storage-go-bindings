# Makefile for Codex Go Bindings

NIM_CODEX_DIR := vendor/nim-codex
NIM_CODEX_LIB_DIR   := $(abspath $(NIM_CODEX_DIR)/library)
NIM_CODEX_BUILD_DIR := $(abspath $(NIM_CODEX_DIR)/build)

CGO_CFLAGS  := -I$(NIM_CODEX_LIB_DIR)
CGO_LDFLAGS := -L$(NIM_CODEX_BUILD_DIR) -lcodex -Wl,-rpath,$(NIM_CODEX_BUILD_DIR)

.PHONY: all clean update libcodex build test

all: build

submodules:
	@echo "Fetching submodules..."
	@git submodule update --init --recursive

update: | submodules
	@echo "Updating nim-codex..."
	@$(MAKE) -C $(NIM_CODEX_DIR) update

libcodex:
	@echo "Building libcodex..."
	@$(MAKE) -C $(NIM_CODEX_DIR) libcodex

build:
	@echo "Building Codex Go Bindings..."
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o codex-go ./codex

test:
	@echo "Running tests..."
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOTESTFLAGS="-timeout=2m" gotestsum --packages="./..." -f testname -- $(if $(filter-out test,$(MAKECMDGOALS)),-run "$(filter-out test,$(MAKECMDGOALS))")

%:
	@:

clean:
	@echo "Cleaning up..."
	@git submodule deinit -f $(NIM_CODEX_DIR)
	@rm -f codex-go