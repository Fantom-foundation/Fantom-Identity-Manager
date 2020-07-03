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
	go build -o $(GOBIN)/identity-app .

## run: Run identity manager with example variables
.PHONY: run
run:
	COOKIE_STORE_KEY=01234567890123456789012345678901 \
	SESSION_STORE_KEY=01234567890123456789012345678901 \
	HYDRA_ADMIN_URL=http://localhost:4445 \
	PORT=3000 \
	IMPORT_USERS=users.example.json \
	$(GOBIN)/identity-app


## clear: Clear persistent data
.PHONY: clear
clear:
	rm -f $(GOBIN)/identity-app

all: help
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a make command in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
