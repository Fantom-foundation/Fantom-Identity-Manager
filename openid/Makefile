# --------------------------------------------------------------------------
# Makefile for the OpenID Connect API
# (c) Fantom Foundation, 2020
#
# v0.1 (2020/06/20)  - Initial version.
# --------------------------------------------------------------------------

# project related vars
PROJECT := $(shell basename "$(PWD)")
DOCKER_FILES := -f openid.yml -f database.yml

## run: Run OpenID Connect API
.PHONY: run
run:
	docker-compose $(DOCKER_FILES) up --build --force-recreate

## register-client: Create Client ${CLIENT_ID} with secret ${CLIENT_SECRET} for callbacks ${CLIENT_CALLBACKS}
.PHONY: register-client
register-client:
	docker-compose $(DOCKER_FILES) exec hydra \
		hydra clients create \
		--skip-tls-verify \
    	--endpoint http://127.0.0.1:4445/ \
		--id ${CLIENT_ID} \
		--secret ${CLIENT_SECRET} \
		--grant-types authorization_code,refresh_token,client_credentials,implicit \
		--response-types token,code,id_token \
		--scope openid,offline \
		--callbacks ${CLIENT_CALLBACKS}

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
