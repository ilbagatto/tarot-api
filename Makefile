# Makefile for managing the Tarot API project

APP_NAME=tarot-api
GO_BIN=$(shell go env GOPATH)/bin

.PHONY: all build run test clean docs

# Build the project
build: docs
	go build -o $(APP_NAME) ./cmd/main.go

# Build the project for linux platform
build-linux: docs
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) ./cmd/main.go

# Run the application (depends on build)
run: build
	./$(APP_NAME)

# Run unit-tests
test-unit:
	go test ./internal/...

# Generate Swagger documentation (depends on build)
docs:
	rm -rf docs/
	$(GO_BIN)/swag init -g ./cmd/main.go

# Clean up build artifacts
clean:
	rm -rf $(APP_NAME) docs/

# Run integration tests
test-integration:
	./scripts/run_integration_tests.sh



