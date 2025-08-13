# RxGo

Modern Reactive Extensions for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/droxer/RxGo)](https://goreportcard.com/report/github.com/droxer/RxGo)
[![GoDoc](https://godoc.org/github.com/droxer/RxGo?status.svg)](https://godoc.org/github.com/droxer/RxGo)

## Overview

RxGo is a reactive programming library for Go that provides both the original Observable API and full Reactive Streams 1.0.4 compliance with backpressure support. It enables you to compose asynchronous and event-based programs using observable sequences.

## Features

### **✅ Core Features**
- **Type-safe generics** throughout the API
- **Context-based cancellation** support
- **Backpressure strategies**: Buffer, Drop, Latest, Error
- **Thread-safe** signal delivery
- **Memory efficient** with bounded buffers
- **Full interoperability** between old and new APIs

### **✅ Reactive Streams 1.0.4 Compliance**
- **Publisher[T]** - Type-safe data source with demand control
- **ReactiveSubscriber[T]** - Complete subscriber interface with lifecycle
- **Subscription** - Request/cancel control with backpressure
- **Processor[T,R]** - Transforming publisher
- **Backpressure support** - Full demand-based flow control
- **Non-blocking guarantees** - Async processing with context cancellation
- **Binary compatible** with Reactive Streams 1.0.4

### **✅ Backpressure Strategies**
- **Buffer** - Queue items when producer is faster than consumer
- **Drop** - Drop new items when buffer is full
- **Latest** - Keep only the latest item when buffer is full
- **Error** - Signal error when buffer overflows

## Quick Start

### Simple Usage

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/rx"
)

func main() {
    // Create an observable
    obs := rx.Just(1, 2, 3, 4, 5)
    
    // Subscribe with a functional subscriber
    obs.Subscribe(context.Background(), rx.NewSubscriber(
        func(v int) {
            fmt.Printf("Received: %d\n", v)
        },
        func() {
            fmt.Println("Completed!")
        },
        func(err error) {
            fmt.Printf("Error: %v\n", err)
        },
    ))
}
```

### With Schedulers

```go
import (
	"github.com/droxer/RxGo/pkg/rx"
	"github.com/droxer/RxGo/pkg/rx/operators"
	"github.com/droxer/RxGo/pkg/rx/scheduler"
)

// Use with different schedulers
obs := rx.Range(1, 10)

// Process on computation scheduler
operators.ObserveOn(obs, scheduler.Computation).Subscribe(ctx, subscriber)

// Process on IO scheduler
operators.ObserveOn(obs, scheduler.IO).Subscribe(ctx, subscriber)
```

### Data Transformation

```go
import (
	"github.com/droxer/RxGo/pkg/rx"
	"github.com/droxer/RxGo/pkg/rx/operators"
)

// Map and Filter operations
obs := rx.Range(1, 10)
transformed := operators.Map(obs, func(x int) int { return x * x })
filtered := operators.Filter(transformed, func(x int) bool { return x > 10 })
filtered.Subscribe(context.Background(), subscriber)
```

### Installation

### Latest Version
```bash
go get github.com/droxer/RxGo@latest
```

### Go Modules
Add to your `go.mod`:
```go
require github.com/droxer/RxGo latest
```

### Requirements
- Go 1.23 or higher (for generics support)

## Architecture

RxGo provides a **modular, clean architecture** that combines intuitive reactive programming with full Reactive Streams 1.0.4 compliance. The library is organized into focused subpackages for maximum flexibility and clarity.

### Package Overview
- **`pkg/rx`** - Core Observable API with type-safe generics
- **`pkg/rx/scheduler`** - Advanced scheduling with 5 scheduler types
- **`pkg/rx/operators`** - Data transformation operators
- **`pkg/rx/streams`** - Full Reactive Streams 1.0.4 compliance

For complete architectural details, see [Architecture Documentation](./docs/architecture.md).

## Documentation

Comprehensive documentation is available in the [docs](./docs/) directory:

- **[Quick Start](./docs/quick-start.md)** - Get started in 5 minutes
- **[Architecture](./docs/architecture.md)** - Package structure and design decisions
- **[Basic Usage](./docs/basic-usage.md)** - Simple Observable API examples
- **[Reactive Streams API](./docs/reactive-streams.md)** - Full Reactive Streams 1.0.4 compliance
- **[Backpressure Control](./docs/backpressure.md)** - Handle producer/consumer speed mismatches
- **[Context Cancellation](./docs/context-cancellation.md)** - Graceful cancellation using Go context
- **[Data Transformation](./docs/data-transformation.md)** - Transform and process data streams
- **[Schedulers](./docs/schedulers.md)** - Execution context control with different scheduler types
- **[API Reference](./docs/api-reference.md)** - Complete API documentation
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
