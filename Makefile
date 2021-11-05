.PHONY: all build run gotool clean help
.DEFAULT_GOAL := help
BINARY_NAME=firewall-api
GOCMD=go
GOBUILD=$(GOCMD) build
GOBUILD_DIR=cmd
OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin

build:
	scripts/build.sh $(BINARY_NAME)

clean:
    ifneq ($(wildcard) $(OUT_DIR),)
    	$(shell rm -fr $(OUT_DIR))
    endif

help:
	@ echo "all build run gotool clean help"
	@ echo ""
	@ echo "Example:"
	@ echo "	make build"
	@ echo "	make clean"
	@ echo "	make help"
	@ echo "	make rpm"
