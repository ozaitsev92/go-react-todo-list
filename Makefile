.PHONY: test
test: ## Run all the tests
	go test -v -race -timeout 30s ./...
