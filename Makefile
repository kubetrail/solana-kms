SHELL := /bin/bash

# formatting color values
RD="$(shell tput setaf 1)"
YE="$(shell tput setaf 3)"
NC="$(shell tput sgr0)"

.PHONY: all
all: goimports vendor
	@echo -e ${YE}▶ building and installing ${NAME} binary${NC}
	@go install

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

