.PHONY: build
build: ## Build a version
	go build -o ./build/apiserver -v ./cmd/apiserver

.PHONY: test
test: ## Run all the tests
	go test -v -race -timeout 30s ./...

.PHONY: clean
clean: ## Remove temporary files
	go clean

#TODO get db connection string from ENV variables
.PHONY: migrate
migrate: ## Run migrations
	migrate -path migrations -database "postgres://localhost:5432/todo_db?sslmode=disable&user=user&password=password" up

#TODO get db connection string from ENV variables
.PHONY: rollback
rollback: ## Rollback migrations
	migrate -path migrations -database "postgres://localhost:5432/todo_db?sslmode=disable&user=user&password=password" down

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
