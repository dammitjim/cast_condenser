build: ## Runs go build on the project
	go build && mv condenser ./bin/condenser

dev: ## Builds and runs the service with local environment
	go build && mv condenser ./bin/condenser && ./bin/condenser

run: ## Runs the service with local environment unless overridden
	./bin/condenser

test: ## Runs gb test with the -v verbose flag
	go test -v

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
