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
go get github.com/droxer/RxGo@v0.1.0
```

### Go Modules
Add to your `go.mod`:
```go
require github.com/droxer/RxGo v0.1.0
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
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// ReactiveSubscriber implementation
type LoggingSubscriber[T any] struct {
    name string
}

func (s *LoggingSubscriber[T]) OnSubscribe(sub rxgo.Subscription) {
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
    publisher := rxgo.RangePublisher(1, 10)
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
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

type SlowSubscriber struct {
    received []int
}

func (s *SlowSubscriber) OnSubscribe(sub rxgo.Subscription) {
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
    publisher := rxgo.RangePublisher(1, 100)
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
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

type ContextSubscriber struct{}

func (s *ContextSubscriber) OnSubscribe(sub rxgo.Subscription) {
    sub.Request(rxgo.Unlimited)
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
    
    publisher := rxgo.RangePublisher(1, 1000)
    subscriber := &ContextSubscriber{}
    
    publisher.Subscribe(ctx, subscriber)
    time.Sleep(time.Second)
}
```

### 5. Data Transformation

Create observables that transform data streams:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/observable"
)

func main() {
    // Create source observable
    source := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        defer sub.OnCompleted()
        for i := 1; i <= 5; i++ {
            sub.OnNext(i)
        }
    })
    
    // Transform the stream by doubling values
    transformed := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        source.Subscribe(ctx, observable.Subscriber[int]{
            Start: func() {
                sub.Start()
            },
            OnNext: func(value int) {
                sub.OnNext(value * 2) // Double the value
            },
            OnError: func(err error) {
                sub.OnError(err)
            },
            OnCompleted: func() {
                sub.OnCompleted()
            },
        })
    })
    
    // Subscribe to the transformed stream
    type IntSubscriber struct{}
    
    func (s *IntSubscriber) Start() {
        fmt.Println("Starting transformation subscriber")
    }
    
    func (s *IntSubscriber) OnNext(value int) {
        fmt.Printf("Received: %d\n", value)
    }
    
    func (s *IntSubscriber) OnError(err error) {
        fmt.Printf("Error: %v\n", err)
    }
    
    func (s *IntSubscriber) OnCompleted() {
        fmt.Println("Transformation completed")
    }
    
    subscriber := &IntSubscriber{}
    transformed.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```

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
    Subscribe(ctx context.Context, s ReactiveSubscriber[T])
}
```

#### ReactiveSubscriber[T] Interface (Reactive Streams API)
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
    rxgo.Unlimited = 1 << 62 // Maximum request value
)
```

### Utility Functions

| Function | Description | Example |
|----------|-------------|---------|
| `observable.Just[T](values...)` | Create from literal values | `observable.Just(1, 2, 3)` |
| `observable.Range(start, count)` | Integer sequence observable | `observable.Range(1, 10)` |
| `observable.Create[T](fn)` | Custom observable creation | `observable.Create(customProducer)` |
| `rxgo.RangePublisher(start, count)` | Integer sequence publisher | `rxgo.RangePublisher(1, 10)` |
| `rxgo.FromSlice[T](slice)` | Create publisher from slice | `rxgo.FromSlice([]int{1, 2, 3})` |
| `rxgo.NewReactivePublisher[T](fn)` | Custom publisher creation | `rxgo.NewReactivePublisher(customProducer)` |

## Running Examples and Tests

### Basic Examples
```bash
# Run basic examples
go run examples/basic/basic.go

# Run backpressure examples
go run examples/backpressure/backpressure.go

# Run context examples
go run examples/context/context.go

# Run compiled binaries
./bin/basic
./bin/backpressure
./bin/context
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
