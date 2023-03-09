.PHONY: help
help: ## List of available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the Envoy Docker container
	 docker build -t envoy-local:v1 .

.PHONY: run
run: ## Run the Envoy container locally to listen on port 10000
	docker run --rm --network host envoy-local:v1

