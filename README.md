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

```bash
go get github.com/droxer/RxGo
```

### Requirements
- Go 1.23 or higher (for generics support)

## API Overview

RxGo provides two compatible APIs:

1. **Observable API** - Simple, callback-based approach
2. **Reactive Streams API** - Full specification compliance with backpressure

## Usage Examples

### 1. Observable API (Simple)

The simplest way to get started with RxGo using the backward-compatible Observable API:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/observable"
)

type IntSubscriber struct{}

func (s *IntSubscriber) Start() {}
func (s *IntSubscriber) OnNext(value int) { fmt.Println(value) }
func (s *IntSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *IntSubscriber) OnCompleted() { fmt.Println("Completed!") }

func main() {
    observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        for i := 0; i < 5; i++ {
            sub.OnNext(i)
        }
        sub.OnCompleted()
    }).Subscribe(context.Background(), &IntSubscriber{})
}
```

### 2. Reactive Streams API (Advanced)

Full Reactive Streams compliance with backpressure and demand control:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/internal/publisher"
)

// ReactiveSubscriber implementation
type LoggingSubscriber[T any] struct {
    name string
}

func (s *LoggingSubscriber[T]) OnSubscribe(sub publisher.Subscription) {
    fmt.Printf("[%s] Subscribed, requesting 3 items\n", s.name)
    sub.Request(3) // Backpressure control
}

func (s *LoggingSubscriber[T]) OnNext(value T) {
    fmt.Printf("[%s] Received: %v\n", s.name, value)
}

func (s *LoggingSubscriber[T]) OnError(err error) {
    fmt.Printf("[%s] Error: %v\n", s.name, err)
}

func (s *LoggingSubscriber[T]) OnComplete() {
    fmt.Printf("[%s] Completed\n", s.name)
}

func main() {
    publisher := publisher.NewRangePublisher(1, 10)
    subscriber := &LoggingSubscriber[int]{name: "Demo"}
    publisher.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```

### 3. Backpressure Control

Handle producer/consumer speed mismatches with backpressure:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/internal/publisher"
)

type SlowSubscriber struct {
    received []int
}

func (s *SlowSubscriber) OnSubscribe(sub publisher.Subscription) {
    // Request items slowly
    go func() {
        for i := 0; i < 5; i++ {
            time.Sleep(100 * time.Millisecond)
            sub.Request(1)
        }
    }()
}

func (s *SlowSubscriber) OnNext(value int) {
    s.received = append(s.received, value)
    fmt.Printf("Processing: %d\n", value)
    time.Sleep(50 * time.Millisecond) // Simulate slow processing
}

func (s *SlowSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *SlowSubscriber) OnComplete() { fmt.Println("Done!") }

func main() {
    publisher := publisher.NewRangePublisher(1, 100)
    subscriber := &SlowSubscriber{}
    publisher.Subscribe(context.Background(), subscriber)
    
    time.Sleep(time.Second)
}
```

### 4. Context Cancellation

Graceful cancellation using Go context:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/internal/publisher"
)

type ContextSubscriber struct{}

func (s *ContextSubscriber) OnSubscribe(sub publisher.Subscription) {
    sub.Request(publisher.Unlimited)
}

func (s *ContextSubscriber) OnNext(value int) {
    fmt.Printf("Received: %d\n", value)
}

func (s *ContextSubscriber) OnError(err error) {
    fmt.Printf("Cancelled: %v\n", err)
}

func (s *ContextSubscriber) OnComplete() {
    fmt.Println("Completed")
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    publisher := publisher.NewRangePublisher(1, 1000)
    subscriber := &ContextSubscriber{}
    
    publisher.Subscribe(ctx, subscriber)
    time.Sleep(time.Second)
}
```

### 5. Data Transformation

Create processors to transform data streams:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/internal/publisher"
)

// Simple processor that doubles values
type DoubleProcessor struct {
    publisher publisher.ReactivePublisher[int]
}

func (p *DoubleProcessor) OnSubscribe(sub publisher.Subscription) {
    sub.Request(publisher.Unlimited)
}

func (p *DoubleProcessor) OnNext(value int) {
    p.publisher.OnNext(value * 2) // Forward doubled value
}

func (p *DoubleProcessor) OnError(err error) {
    p.publisher.OnError(err)
}

func (p *DoubleProcessor) OnComplete() {
    p.publisher.OnComplete()
}

func main() {
    source := publisher.NewRangePublisher(1, 5)
    
    processor := publisher.NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[int]) {
        source.Subscribe(ctx, &DoubleProcessor{processor: publisher.NewReactivePublisher(func(ctx context.Context, innerSub publisher.ReactiveSubscriber[int]) {
            sub = innerSub
        })})})
    
    subscriber := publisher.NewBenchmarkSubscriber[int]()
    processor.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```

## API Reference

### Core Interfaces

#### Publisher[T] Interface
```go
type Publisher[T any] interface {
    Subscribe(ctx context.Context, s ReactiveSubscriber[T])
}
```

#### ReactiveSubscriber[T] Interface
```go
type ReactiveSubscriber[T any] interface {
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
    publisher.Unlimited = 1 << 62 // Maximum request value
)
```

### Utility Functions

| Function | Description | Example |
|----------|-------------|---------|
| `observable.Just[T](values...)` | Create from literal values | `observable.Just(1, 2, 3)` |
| `publisher.NewRangePublisher(start, count)` | Integer sequence publisher | `publisher.NewRangePublisher(1, 10)` |
| `publisher.FromSlice[T](slice)` | Create from slice | `publisher.FromSlice([]int{1, 2, 3})` |
| `publisher.NewReactivePublisher[T](fn)` | Custom publisher creation | `publisher.NewReactivePublisher(customProducer)` |

## Running Examples and Tests

### Basic Examples
```bash
# Run basic examples
go run examples/basic.go

# Run reactive streams examples
go run examples/reactive_streams.go

# Run backpressure examples
go run examples/backpressure.go
```

### Testing
```bash
# Run all tests
go test -v ./...

# Run tests with race detection
go test -race -v ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

## Migration Guide

### From Observable API to Reactive Streams API

The Reactive Streams API provides more control and backpressure support. Here's how to migrate:

#### Before (Observable API)
```go
obs := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
    sub.OnNext(42)
    sub.OnCompleted()
})
obs.Subscribe(ctx, subscriber)
```

#### After (Reactive Streams API)
```go
publisher := publisher.NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[int]) {
    sub.OnSubscribe(publisher.NewSubscription())
    sub.OnNext(42)
    sub.OnComplete()
})
publisher.Subscribe(ctx, reactiveSubscriber)
```

### Key Changes
- Replace `Subscriber` with `ReactiveSubscriber`
- Add `OnSubscribe` method to handle subscription
- Use `Request(n)` for backpressure control
- Replace `OnCompleted` with `OnComplete`

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

We welcome contributions! Please follow these steps:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Add** tests for new functionality
4. **Ensure** all tests pass (`go test -v ./...`)
5. **Commit** your changes (`git commit -m 'Add amazing feature'`)
6. **Push** to the branch (`git push origin feature/amazing-feature`)
7. **Open** a pull request

### Development Setup
```bash
# Clone the repository
git clone https://github.com/droxer/RxGo.git
cd RxGo

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## License

MIT License

Copyright (c) 2025 RxGo Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Support

- **Issues**: [GitHub Issues](https://github.com/droxer/RxGo/issues)
- **Documentation**: [GoDoc](https://godoc.org/github.com/droxer/RxGo)
