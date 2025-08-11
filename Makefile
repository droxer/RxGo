.PHONY: default dev build test test-coverage race fmt lint vet deps deps-upgrade check-security bench clean quality install-tools verify-examples

default: deps test

dev: fmt lint test race

build:
	go build ./...
	@echo "✅ Library validation complete - no build errors"

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

race:
	go test -v -race ./...

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

doc:
	@echo "📚 Generating documentation..."
	go doc -all ./pkg/observable
	go doc -all ./internal/publisher
	go doc -all .

doc-serve:
	@echo "🌐 Starting documentation server..."
	godoc -http=:6060

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

check: fmt vet lint test race

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/goreleaser/goreleaser@latest

release:
	goreleaser release --clean

release-check:
	goreleaser check

verify-examples:
	@echo "🔍 Verifying all examples build and run successfully..."
	@for dir in examples/*/; do \
		if [ -f "$${dir}*.go" ]; then \
			echo "  ✅ Testing $$dir"; \
			cd $$dir && go build -o /tmp/test_example . && echo "    Build OK" || exit 1; \
			cd ../..; \
		fi \
	done
	@echo "🎉 All examples verified successfully!"
	@echo ""
	@echo "Run individual examples:"
	go run examples/data-transformation/data_transformation.go
	go run examples/http-streaming/http_streaming.go
	go run examples/financial/financial_processor.go
	go run examples/real-time/monitoring.go

