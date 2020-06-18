# --------------------------------------------------------------------------
# Makefile for the Fantom Key Vault Server
# (c) Fantom Foundation, 2020
#
# v0.1 (2020/06/18)  - Initial version, base server build.
# --------------------------------------------------------------------------

# project related vars
PROJECT := $(shell basename "$(PWD)")

# go related vars
GOBASE := $(shell pwd)
GOBIN=$(CURDIR)/build

## server: Make the API server as bin/frd
server:
	go build -o $(GOBIN)/keyvault ./cmd/keyvault

.PHONY: help
all: help
help: Makefile
	@echo
	@echo "Choose a make command in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
