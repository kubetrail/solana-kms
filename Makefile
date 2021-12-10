SHELL := /bin/bash

# formatting color values
RD="$(shell tput setaf 1)"
YE="$(shell tput setaf 3)"
NC="$(shell tput sgr0)"

# please define PROJECT env. var for your Google
# cloud project ID
NAME="solana-kms"
VERSION="0.0.1"
CATEGORY="services"
IMG="us-central1-docker.pkg.dev/${PROJECT}/artifacts/services/${NAME}"

.PHONY: all
all: _sanity goimports vendor
	@echo -e ${YE}▶ building and installing ${NAME} binary${NC}
	@go install

# sanity check
.PHONY: _sanity
_sanity:
	@if [[ -z "${PROJECT}" ]]; then \
	    echo "please set PROJECT env. var for your Google cloud project"; \
	    exit 1; \
	fi
	@for cmd in podman kubectl helm go goimports; do \
		if [[ -z $$(command -v $${cmd}) ]]; then \
                	echo "$${cmd} not found. pl. install."; \
                	exit 1; \
        	fi; \
	done

.PHONY: goimports
goimports:
	@echo -e ${YE}▶ goimports formatting${NC}
	@goimports -w -l main.go
	@goimports -w -l cmd
	@goimports -w -l pkg

.PHONY: vendor
vendor:
	@echo -e ${YE}▶ regenerating vendor folder${NC}
	@rm -rf vendor
	@go mod vendor

