# RxGo

Modern Reactive Extensions for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/droxer/RxGo)](https://goreportcard.com/report/github.com/droxer/RxGo)
[![GoDoc](https://godoc.org/github.com/droxer/RxGo?status.svg)](https://godoc.org/github.com/droxer/RxGo)

## Overview

RxGo is a reactive programming library for Go that provides both the original Observable API and full Reactive Streams 1.0.3 compliance with backpressure support. It enables you to compose asynchronous and event-based programs using observable sequences.

## Features

### **✅ Core Features**
- **Type-safe generics** throughout the API
- **Context-based cancellation** support
- **Backpressure strategies**: Buffer, Drop, Latest, Error
- **Thread-safe** signal delivery
- **Memory efficient** with bounded buffers
- **Full interoperability** between old and new APIs

### **✅ Reactive Streams 1.0.3 Compliance**
- **Publisher[T]** - Type-safe data source with demand control
- **ReactiveSubscriber[T]** - Complete subscriber interface with lifecycle
- **Subscription** - Request/cancel control with backpressure
- **Processor[T,R]** - Transforming publisher
- **Backpressure support** - Full demand-based flow control
- **Non-blocking guarantees** - Async processing with context cancellation

### **✅ Backpressure Strategies**
- **Buffer** - Queue items when producer is faster than consumer
- **Drop** - Drop new items when buffer is full
- **Latest** - Keep only the latest item when buffer is full
- **Error** - Signal error when buffer overflows

## Quick Start

### Installation

### Latest Version
```bash
go get github.com/droxer/RxGo@v0.1.1
```

### Go Modules
Add to your `go.mod`:
```go
require github.com/droxer/RxGo v0.1.1
```

### Latest Development
```bash
go get github.com/droxer/RxGo@latest
```

### Requirements
- Go 1.23 or higher (for generics support)

## API Overview

RxGo provides two compatible APIs:

1. **Observable API** - Simple, callback-based approach
2. **Reactive Streams API** - Full specification compliance with backpressure

## Documentation

Comprehensive documentation is available in the [docs](./docs/) directory:

- **[Quick Start](./docs/quick-start.md)** - Simple example to get you started
- **[Basic Usage](./docs/basic-usage.md)** - Simple Observable API examples
- **[Reactive Streams API](./docs/reactive-streams.md)** - Full Reactive Streams 1.0.3 compliance
- **[Backpressure Control](./docs/backpressure.md)** - Handle producer/consumer speed mismatches
- **[Context Cancellation](./docs/context-cancellation.md)** - Graceful cancellation using Go context
- **[Data Transformation](./docs/data-transformation.md)** - Transform and process data streams
- **[Architecture](./docs/architecture.md)** - Package structure and API choices
- **[API Reference](./docs/api-reference.md)** - Complete API documentation
- **[Schedulers](./docs/schedulers.md)** - Execution context control with different scheduler types
## Performance Considerations

### Optimization Features
- **Zero-allocation** signal delivery for common cases
- **Lock-free** data structures where possible
- **Context-aware** cancellation to prevent goroutine leaks
- **Bounded buffers** to prevent memory issues
- **Backpressure** to handle producer/consumer speed mismatches

### Best Practices
- Always use context cancellation for long-running streams
- Implement proper backpressure with `Request(n)` calls
- Use bounded buffers to prevent memory exhaustion
- Prefer the Reactive Streams API for production use

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines on how to contribute to this project.

## License

[MIT License](./LICENSE)

## Support

- **Issues**: [GitHub Issues](https://github.com/droxer/RxGo/issues)
- **Documentation**: [GoDoc](https://godoc.org/github.com/droxer/RxGo)
