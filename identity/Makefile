# --------------------------------------------------------------------------
# Makefile for the OpenID Connect API
# (c) Fantom Foundation, 2020
#
# v0.1 (2020/0/20)  - Initial version.
# --------------------------------------------------------------------------

# project related vars
PROJECT := $(shell basename "$(PWD)")
DOCKER_FILES := -f identity.yml -f database.yml
DOCKER_FILES_DEBUG := -f identity.yml -f identity.debug.yml -f database.yml

## run: Run identity manager with example variables
.PHONY: run
run:
	docker-compose $(DOCKER_FILES) up --build --force-recreate

## debug: Debug identity manager with example variables
.PHONY: debug
debug:
	docker-compose $(DOCKER_FILES_DEBUG) up --build --force-recreate

## clear: Clear persistent data
.PHONY: clear
clear:
	docker-compose $(DOCKER_FILES) down -v

all: help
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a make command in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
