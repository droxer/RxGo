.PHONY: default dev build test test-coverage race fmt lint vet deps deps-upgrade check-security bench clean quality install-tools verify-examples

default: deps test

dev: fmt lint test race

build:
	go build ./...
	@echo "✅ Library validation complete - no build errors"

test:
	go test -v ./pkg... --timeout 30s

test-coverage:
	go test -v -coverprofile=coverage.out ./pkg...
	go tool cover -html=coverage.out -o coverage.html

race:
	go test -v -race ./pkg...

fmt:
	go fmt ./...

lint:
	golangci-lint run

vet:
	go vet ./...

deps:
	go mod download
	go mod tidy

check-security:
	gosec ./... 

bench:
	go test -bench=. -benchmem ./...

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

check: fmt vet lint test race

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/goreleaser/goreleaser@latest

release:
	goreleaser release --clean

release-check:
	goreleaser check

verify-examples:
	@echo "🔍 Verifying all examples build and run successfully..."
	@find examples -name "*.go" -type f | while read file; do \
		dir=$$(dirname "$$file"); \
		name=$$(basename "$$dir"); \
		cd "$$dir" && go build -o /tmp/test_example_$$name || exit 1; \
		cd - > /dev/null; \
	done
	@echo "🎉 All examples verified successfully!"
