# RxGo Makefile - Modern Reactive Extensions for Go
.PHONY: default dev build build-examples test test-coverage race fmt lint vet deps deps-upgrade check-security bench clean quality install-tools

# Default target
default: deps test

# Development targets
dev: fmt lint test race

# Build targets
build:
	go build ./...

build-examples:
	mkdir -p bin
	go build -o bin/basic examples/basic.go
	go build -o bin/backpressure examples/backpressure/backpressure.go
	go build -o bin/context examples/context/context.go

# Testing targets
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

race:
	go test -v -race ./...

# Linting and formatting
fmt:
	go fmt ./...

lint:
	golangci-lint run

vet:
	go vet ./...

# Dependencies
deps:
	go mod download
	go mod tidy

deps-upgrade:
	go get -u ./...
	go mod tidy

# Security
check-security:
	gosec ./...

# Benchmarking
bench:
	go test -bench=. -benchmem ./...

# Clean up
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# All-in-one quality check
quality: fmt vet lint test race

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
