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

## Usage Examples

Comprehensive usage examples are available in the [docs/examples](./docs/examples/) directory:

- **[Quick Start](./docs/examples/quick-start.md)** - Simple example to get you started
- **[Basic Usage](./docs/examples/basic-usage.md)** - Simple Observable API examples
- **[Reactive Streams API](./docs/examples/reactive-streams.md)** - Full Reactive Streams 1.0.3 compliance
- **[Backpressure Control](./docs/examples/backpressure.md)** - Handle producer/consumer speed mismatches
- **[Context Cancellation](./docs/examples/context-cancellation.md)** - Graceful cancellation using Go context
- **[Data Transformation](./docs/examples/data-transformation.md)** - Transform and process data streams

## API Reference

### Core Interfaces

#### Observable API (Simple)
```go
type Observable[T any] struct {
    // Contains filtered or unexported fields
}

func (o *Observable[T]) Subscribe(ctx context.Context, sub Subscriber[T])
```

#### Subscriber[T] Interface (Observable API)
```go
type Subscriber[T any] interface {
    Start()
    OnNext(next T)
    OnCompleted()
    OnError(e error)
}
```

#### Publisher[T] Interface (Reactive Streams API)
```go
type Publisher[T any] interface {
    Subscribe(ctx context.Context, s SubscriberReactive[T])
}
```

#### SubscriberReactive[T] Interface (Reactive Streams API)
```go
type SubscriberReactive[T any] interface {
    OnSubscribe(s Subscription)
    OnNext(t T)
    OnError(err error)
    OnComplete()
}
```

#### Subscription Interface
```go
type Subscription interface {
    Request(n int64)
    Cancel()
}
```

### Constants
```go
const (
    rxgo.Unlimited = 1 << 62 // Maximum request value
)
```

### Utility Functions

| Function | Description | Example |
|----------|-------------|---------|
| `rxgo.Just[T](values...)` | Create from literal values | `rxgo.Just(1, 2, 3)` |
| `rxgo.Range(start, count)` | Integer sequence observable | `rxgo.Range(1, 10)` |
| `rxgo.Create[T](fn)` | Custom observable creation | `rxgo.Create(customProducer)` |
| `rxgo.RangePublisher(start, count)` | Integer sequence publisher | `rxgo.RangePublisher(1, 10)` |
| `rxgo.FromSlicePublisher[T](slice)` | Create publisher from slice | `rxgo.FromSlicePublisher([]int{1, 2, 3})` |
| `rxgo.NewPublisher[T](fn)` | Custom publisher creation | `rxgo.NewPublisher(customProducer)` |

## Running Examples and Tests

### Basic Examples
```bash
# Run basic examples
go run examples/basic/basic.go

# Run backpressure examples
go run examples/backpressure/backpressure.go

# Run context examples
go run examples/context/context.go

```

### Testing
```bash
# Run all tests
make test

# Run tests with race detection
make race

# Run benchmarks
make bench
```

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
