# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

RxGo is a Reactive Extensions library for Go with full Reactive Streams 1.0.4 compliance. It provides two distinct reactive models:

1. **Observable API (Push Model)** - Simple, intuitive reactive programming without backpressure
2. **Reactive Streams API (Pull Model)** - Full Reactive Streams 1.0.4 compliance with backpressure support

## Package Structure

- `github.com/droxer/RxGo/pkg/observable` - Core Observable API with operators
- `github.com/droxer/RxGo/pkg/scheduler` - Execution context control (threading models)
- `github.com/droxer/RxGo/pkg/streams` - Full Reactive Streams 1.0.4 compliance
- `github.com/droxer/RxGo/pkg/adapters.go` - Interoperability between Observable and Reactive Streams APIs

## Common Development Commands

### Build and Test
```bash
# Run all tests
make test

# Run tests with race detection
make race

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench

# Verify all examples build and run
make verify-examples
```

### Code Quality
```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run security checks
make check-security

# Run all quality checks
make check
```

### Dependencies
```bash
# Install/update dependencies
make deps

# Install development tools
make install-tools
```

## Architecture Overview

The library implements two distinct reactive programming models:

### Observable API (Push Model)
- Located in `pkg/observable/`
- Producer controls emission rate
- Data is pushed to subscribers as soon as it's available
- No built-in backpressure mechanism
- Simple and intuitive for basic use cases
- Key files: `observable.go`, `operators.go`

### Reactive Streams API (Pull Model)
- Located in `pkg/streams/`
- Subscriber controls consumption rate through demand requests
- Full backpressure support with multiple strategies
- Reactive Streams 1.0.4 compliant
- Key files: `interfaces.go`, `compliant.go`, `processors.go`

### Key Design Patterns
1. **Generics** - Full type safety throughout the API
2. **Context Integration** - Graceful cancellation support
3. **Modular Design** - Clean separation of concerns
4. **Interoperability** - Adapters to convert between APIs

## Testing Approach

Tests are organized by package:
- `pkg/observable/*_test.go` - Tests for Observable API
- `pkg/streams/*_test.go` - Tests for Reactive Streams API

Use table-driven tests for comprehensive test coverage. Tests should verify:
- Normal operation
- Error conditions
- Context cancellation
- Edge cases

## Release Process

The library uses GoReleaser for releases:
```bash
# Check release configuration
make release-check

# Create a release (requires proper git tags)
make release
```