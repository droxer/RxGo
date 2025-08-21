# Qwen Code Context for RxGo

This document provides essential context for Qwen Code to effectively assist with the RxGo project.

## Project Type
This is a **Go code project** implementing Reactive Extensions for the Go programming language, aiming for full Reactive Streams 1.0.4 compliance.

## Project Overview
RxGo is a library that brings the power of Reactive Programming to Go. It provides two main APIs:

1.  **Observable API (Push Model)**: A simpler, intuitive API for basic reactive programming tasks.
2.  **Reactive Streams API (Pull Model)**: A fully compliant implementation of the Reactive Streams 1.0.4 specification, providing robust backpressure handling for production systems.

Key features include:
*   Full Go generics support for type safety.
*   Multiple backpressure strategies (Buffer, Drop, Latest, Error).
*   Context-based cancellation.
*   Thread-safe operations.
*   Retry and backoff mechanisms.

## Package Structure
*   `pkg/observable`: Core Observable API and basic operators (e.g., `Just`, `Range`, `Map`, `Filter`).
*   `pkg/streams`: Full Reactive Streams 1.0.4 implementation with `Publisher`, `Subscriber`, `Subscription`, and `Processor`.
*   `pkg/scheduler`: Execution context control (e.g., trampoline, new thread, computation pool).
*   Adapters in `pkg/`: Utilities to bridge the Observable and Reactive Streams APIs.

## Building and Running
Common development commands are managed via the `Makefile`:
*   `make deps`: Download and tidy Go modules.
*   `make build`: Build the project.
*   `make test`: Run unit tests.
*   `make test-coverage`: Run tests with coverage reporting.
*   `make race`: Run tests with the race detector.
*   `make bench`: Run benchmarks.
*   `make check`: Run a comprehensive suite including formatting, vetting, linting, tests, race detection, and security checks.
*   `make verify-examples`: Verify that all examples build and run successfully.
*   `make install-tools`: Install required development tools like `golangci-lint` and `gosec`.

## Development Conventions
*   **Code Style**: Follow standard Go conventions and idioms. Use `make fmt` and `make lint` to ensure consistency.
*   **Testing**: Write comprehensive unit tests for new functionality. Use table-driven tests where appropriate. Tests are located alongside the code they test, typically ending in `_test.go`.
*   **Documentation**: Add documentation comments for all public APIs.
*   **Contributions**: Follow the guidelines in `CONTRIBUTING.md`: fork, create a feature branch, add tests, ensure all checks pass (`make check`), and submit a PR.