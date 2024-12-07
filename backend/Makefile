# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean

# Build target
BINARY_NAME = app

.PHONY: all
all: test build

.PHONY: build
build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/app

.PHONY: test
test:
	$(GOTEST) -race -v ./...

.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

.PHONY: coverage
coverage:
	$(GOTEST) -v ./... -cover -coverprofile=c.out

.PHONY: hcoverage
hcoverage:
	$(GOCMD) tool cover -html=c.out

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)