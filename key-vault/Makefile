# --------------------------------------------------------------------------
# Makefile for the Fantom Key Vault Server
# (c) Fantom Foundation, 2020
#
# v0.1 (2020/06/18)  - Initial version, base server build.
# v0.2 (2020/07/10)  - Use docker as main build tool
# --------------------------------------------------------------------------

# project related vars
PROJECT := $(shell basename "$(PWD)")
DOCKER_FILES := -f key-vault-server.yml
DOCKER_FILES_DEBUG := -f key-vault-server.yml -f key-vault-server.debug.yml

## run: Run with example variables
.PHONY: run
run:
	docker-compose $(DOCKER_FILES) up --build --force-recreate

## debug: Debug with example variables
.PHONY: debug
debug:
	docker-compose $(DOCKER_FILES_DEBUG) up --build --force-recreate

## clear: Clear persistent data
.PHONY: clear
clear:
	docker-compose $(DOCKER_FILES) down -v

## generate-schema: Regenerate graphql server code from schema definition
.PHONY: generate-schema
generate-schema:
	cd key-vault-server && go run github.com/99designs/gqlgen generate

all: help
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a make command in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
