APP_NAME=storageserver

.PHONY: help

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

test: ## Run the tests
	cd src && go test  -v ./...

build: ## Builds the binary
	cd src && go build -o ../$(APP_NAME)

build-docker: ## Builds the docker image
	docker build -t $(APP_NAME) .

setup-local-environment: ## Sets up a local testing environment
	docker compose up
	