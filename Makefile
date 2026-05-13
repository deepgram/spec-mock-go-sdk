# Copyright Deepgram, Inc. All Rights Reserved.
# SPDX-License-Identifier: MIT

ifeq ($(OS),Windows_NT)
	build_OS := Windows
	NUL = NUL
else
	build_OS := $(shell uname -s 2>/dev/null || echo Unknown)
	NUL = /dev/null
endif

.DEFAULT_GOAL := help

ROOT_DIR := $(shell git rev-parse --show-toplevel)
GO := go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOHOSTOS ?= $(shell go env GOHOSTOS)
GOHOSTARCH ?= $(shell go env GOHOSTARCH)

help: #### Display help
	@echo ''
	@echo 'Syntax: make <target>'
	@awk 'BEGIN {FS = ":.*#### "; printf "\nTargets:\n"} /^[a-zA-Z_-]+:.*?#### / { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } ' $(MAKEFILE_LIST)
	@echo ''

.PHONY: version
version: #### Display tool versions
	@echo 'ROOT_DIR: $(ROOT_DIR)'
	@echo 'GOOS: $(GOOS)'
	@echo 'GOARCH: $(GOARCH)'
	@echo 'go version: $(shell go version)'
