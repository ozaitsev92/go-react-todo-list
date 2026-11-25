GO ?= go
BINARY := bin/app

.PHONY: all
all: lint test build

.PHONY: build
build: ## Build the server binary into ./bin
	$(GO) build -trimpath -o $(BINARY) ./cmd/app

.PHONY: run
run: ## Run with in-memory storage and a throwaway dev key (never use in production)
	STORAGE=memory JWT_SECURE_COOKIE=false LOG_LEVEL=debug \
	JWT_SIGNING_KEY=$${JWT_SIGNING_KEY:-dev-only-insecure-key-0123456789abcdef} \
	$(GO) run ./cmd/app

.PHONY: test
test: ## Unit tests (no Docker needed)
	$(GO) test -race -shuffle=on ./...

.PHONY: test-integration
test-integration: ## Unit + MongoDB integration tests (needs Docker or MONGO_TEST_URI)
	$(GO) test -race -shuffle=on -tags integration ./...

.PHONY: coverage
coverage: ## Coverage across packages, including integration tests
	$(GO) test -tags integration -coverpkg=./... -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -1

.PHONY: coverage-html
coverage-html: coverage
	$(GO) tool cover -html=coverage.out

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: lint
lint: ## go vet + golangci-lint (config in .golangci.yml)
	$(GO) vet ./...
	$(GO) vet -tags integration ./...
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run --build-tags integration ./...

.PHONY: vuln
vuln: ## Report known vulnerabilities in reachable code
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: clean
clean:
	$(GO) clean
	rm -rf bin coverage.out
