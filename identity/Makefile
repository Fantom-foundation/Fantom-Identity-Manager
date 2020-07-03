# --------------------------------------------------------------------------
# Makefile for the OpenID Connect API
# (c) Fantom Foundation, 2020
#
# v0.1 (2020/0/20)  - Initial version.
# --------------------------------------------------------------------------

# project related vars
PROJECT := $(shell basename "$(PWD)")

# go related vars
GOBASE := $(shell pwd)
GOBIN=$(CURDIR)/build

## build: Build identity manager
.PHONY: build
build:
	cd identity-app && make build

## run: Run identity manager with example variables
.PHONY: run
run:
	cd identity-app && make run


## clear: Clear persistent data
.PHONY: clear
clear:
	cd identity-app && make clear

all: help
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a make command in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
