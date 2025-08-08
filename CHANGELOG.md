# Changelog

All notable changes to the RxGo project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-08-08

### Added
- Initial release of RxGo with full Reactive Streams 1.0.3 compliance
- Type-safe generics support throughout the API (Go 1.23+)
- Observable API for simple, callback-based reactive programming
- Reactive Streams API with backpressure support
- Context-based cancellation support
- Backpressure strategies: Buffer, Drop, Latest, Error
- Thread-safe signal delivery
- Memory efficient with bounded buffers
- Full interoperability between Observable and Reactive Streams APIs
- Comprehensive documentation with examples
- Makefile with build, test, and documentation targets
- Package documentation for godoc.org

### Features
- **Publisher[T]** - Type-safe data source with demand control
- **ReactiveSubscriber[T]** - Complete subscriber interface with lifecycle
- **Subscription** - Request/cancel control with backpressure
- **Processor[T,R]** - Transforming publisher
- **Backpressure support** - Full demand-based flow control
- **Non-blocking guarantees** - Async processing with context cancellation

### Utility Functions
- `observable.Just[T](values...)` - Create from literal values
- `publisher.NewRangePublisher(start, count)` - Integer sequence publisher
- `publisher.FromSlice[T](slice)` - Create from slice
- `publisher.NewReactivePublisher[T](fn)` - Custom publisher creation

### Documentation
- Comprehensive README with examples
- Package-level documentation for all public packages
- Go doc comments for all exported types and functions
- Examples directory with practical usage patterns

### Compatibility
- Go 1.23 or higher required (for generics support)
- Full backward compatibility between Observable and Reactive Streams APIs
- Context-aware cancellation to prevent goroutine leaks