.PHONY: docker-build docker-up docker-down test

docker-build: ## Build docker image
	@echo "Building Docker image..."
	docker build -t reviewer-assigner .

compose-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	docker-compose up -d

compose-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

test: ## Run tests
	@echo "Running tests..."
	go test ./...
