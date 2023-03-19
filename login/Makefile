.PHONY: help
help: ## List of available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the login Docker container
	 docker build -t login:v1 .

.PHONY: run
run: ## Run the container locally to listen on port 9040
	docker run --rm --network host login:v1

