.PHONY: help
help: ## List of available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the extauth Docker container
	 docker build -t extauth:v1 .

.PHONY: run
run: ## Run the container locally to listen on port 5050
	docker run --rm --network host extauth:v1

