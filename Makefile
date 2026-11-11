.DEFAULT_GOAL := help
.PHONY: help setup install dev assets build run vet check fmt fmt-check clean

GO ?= go
NPM ?= npm
GOFMT ?= gofmt
AIR_VERSION := v1.64.5
APP_BINARY := bin/ticket-sync

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*## "; printf "\nIT Support development commands\n\n"} /^[a-z-]+:.*## / {printf "  %-12s %s\n", $$1, $$2} END {print ""}' $(MAKEFILE_LIST)

setup: install ## Install tools and create .env if missing
	@test -f .env || cp .env.example .env

install: ## Install Node dependencies, Go modules and Air
	$(NPM) ci
	$(GO) mod download
	GOBIN=$(CURDIR)/bin $(GO) install github.com/air-verse/air@$(AIR_VERSION)

dev: ## Start development with Go, template and asset reload
	AIR_BIN=$(CURDIR)/bin/air $(NPM) run dev

assets: ## Compile CSS and copy browser assets
	$(NPM) run build

build: assets ## Build the production executable
	$(GO) build -buildvcs=false -trimpath -o $(APP_BINARY) ./cmd/server

run: ## Run the built app with exported environment variables
	APP_ENV=production ./$(APP_BINARY)

vet: ## Run Go static checks
	$(GO) vet ./...

check: fmt-check vet assets ## Run formatting, static checks and asset build

fmt: ## Format Go code and HTML templates
	$(GOFMT) -w app cmd config framework public resources routes
	$(NPM) run fmt:html

fmt-check: ## Check formatting without modifying files
	@test -z "$$($(GOFMT) -l app cmd config framework public resources routes)" || { echo "Run make fmt to format Go files"; exit 1; }
	$(NPM) run fmt:html:check

clean: ## Remove generated executables and development builds
	rm -rf bin .tmp
