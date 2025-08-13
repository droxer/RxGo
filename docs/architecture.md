# Architecture & Package Structure

This document describes the clean, modular architecture of RxGo, including the package structure and API design choices.

## Overview

RxGo provides a modular, clean API that combines the best of both worlds:

- **Simple Observable API** - Clean, intuitive reactive programming
- **Full Reactive Streams 1.0.4 Compliance** - Built-in backpressure support
- **Modular Package Structure** - Organized into logical subpackages
- **Flexible Imports** - Use what you need, when you need it

## Package Structure

The library is organized into clean, focused subpackages:

### Core Packages

#### `github.com/droxer/RxGo/pkg/rx`
- **Purpose**: Core Observable API and basic reactive programming
- **Key Features**:
  - `Observable[T]` - Main observable type with generics support
  - `Subscriber[T]` - Functional subscriber interface
  - `Just()`, `Range()`, `Create()` - Observable creation functions
  - Context-based cancellation support
  - Type-safe throughout with Go generics

#### `github.com/droxer/RxGo/pkg/rx/scheduler`
- **Purpose**: Advanced threading and execution control
- **Key Features**:
  - `Scheduler` interface - Unified scheduling abstraction
  - **Trampoline** - Immediate execution for lightweight operations
  - **NewThread** - New goroutine for each task
  - **SingleThread** - Sequential processing on dedicated goroutine
  - **Computation** - Fixed thread pool for CPU-bound work
  - **IO** - Cached thread pool for I/O-bound work

#### `github.com/droxer/RxGo/pkg/rx/operators`
- **Purpose**: Data transformation and processing operators
- **Key Features**:
  - `Map()` - Transform values using functions
  - `Filter()` - Filter values based on predicates
  - `ObserveOn()` - Schedule observable emissions on specific schedulers
  - Type-safe operators with generics

#### `github.com/droxer/RxGo/pkg/rx/streams`
- **Purpose**: Full Reactive Streams 1.0.4 compliance
- **Key Features**:
  - `Publisher[T]` - Type-safe data source with demand control
  - `ReactiveSubscriber[T]` - Complete subscriber interface with lifecycle
  - `Subscription` - Request/cancel control with backpressure
  - `Processor[T,R]` - Transforming publisher
  - Full backpressure support with multiple strategies

## Design Principles

### 1. Modularity
Each package has a single, well-defined responsibility:
- `rx/` - Core reactive concepts
- `scheduler/` - Execution control
- `operators/` - Data transformation
- `streams/` - Reactive Streams specification

### 2. Type Safety
Full generics support throughout the API:
- Type-safe Observable creation and processing
- Compile-time type checking
- No interface{} or reflection usage

### 3. Context Integration
Native Go context support:
- Graceful cancellation using `context.Context`
- Goroutine leak prevention
- Timeout and deadline support

### 4. Performance
Optimized for Go's concurrency model:
- Zero-allocation signal delivery where possible
- Lock-free data structures
- Efficient goroutine management
- Bounded buffers to prevent memory issues

## Usage Patterns

### Simple Usage
```go
import "github.com/droxer/RxGo/pkg/rx"

obs := rx.Just(1, 2, 3, 4, 5)
obs.Subscribe(ctx, subscriber)
```

### With Schedulers
```go
import (
    "github.com/droxer/RxGo/pkg/rx"
    "github.com/droxer/RxGo/pkg/rx/scheduler"
    "github.com/droxer/RxGo/pkg/rx/operators"
)

obs := rx.Range(1, 10)
transformed := operators.ObserveOn(obs, scheduler.Computation)
transformed.Subscribe(ctx, subscriber)
```

### Full Reactive Streams
```go
import "github.com/droxer/RxGo/pkg/rx/streams"

publisher := streams.NewPublisher[int]()
subscriber := streams.NewReactiveSubscriber[int](
    func(v int) { fmt.Println("Received:", v) },
    func() { fmt.Println("Completed") },
    func(err error) { fmt.Println("Error:", err) },
)
publisher.Subscribe(subscriber)
```


## Best Practices

### Package Selection
- Use `pkg/rx` for simple reactive programming
- Use `pkg/rx/scheduler` for advanced thread control
- Use `pkg/rx/operators` for data transformation
- Use `pkg/rx/streams` for production systems with backpressure

### Performance Guidelines
- Always use context cancellation for long-running streams
- Implement proper backpressure with `Request(n)` calls
- Use bounded buffers to prevent memory exhaustion
- Prefer the Reactive Streams API for production use

### Thread Safety
All packages are designed to be thread-safe:
- Concurrent subscription support
- Safe scheduler usage across goroutines
- Atomic operations where appropriate
- No shared mutable state