DB_TYPE ?= postgres
DB_NAME ?= todo_db
DB_USER ?= user
DB_PWD ?= password
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432

DB_CONNECTION_STRING ?= "$(DB_TYPE)://$(DB_USER):$(DB_PWD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable"

.PHONY: build test clean migrate rollback help

build: lint ## Build a version
	go build -o ./build/apiserver -v ./cmd/apiserver

run: build ## Run the app
	./build/apiserver

test: ## Run all the tests
	go test -v -race -timeout 30s ./...

clean: ## Remove temporary files
	go clean

postgresup: ## Start PostgreSQl server
	docker run --name todoapp-postgres -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PWD) -e POSTGRES_DB=$(DB_NAME) -p 5432:5432 -d postgres

postgresdown: ## Stop PostgreSQl server
	docker stop todoapp-postgres
	docker rm todoapp-postgres

migrate: ## Run migrations
	migrate -path migrations -database $(DB_CONNECTION_STRING) up

rollback: ## Rollback migrations
	migrate -path migrations -database $(DB_CONNECTION_STRING) down

lint: ## Lint
	golangci-lint run

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
