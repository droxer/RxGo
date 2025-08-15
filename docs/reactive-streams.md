# Reactive Streams API

Full Reactive Streams 1.0.4 compliance with backpressure support.

## Quick Overview

| Component | Purpose | Example |
|-----------|---------|---------|
| `Publisher[T]` | Data source | `streams.RangePublisher` |
| `Subscriber[T]` | Data consumer | Custom struct implementing Subscriber interface |
| `Subscription` | Demand control | `Request(n)` |

## Pull Model

The Reactive Streams API implements a **pull-based model** with backpressure, which is different from the push-based Observable API:

- **Pull Model**: Subscribers request data using `subscription.Request(n)` 
- **Backpressure Support**: Publishers must respect subscriber demand
- **Compliance**: Full Reactive Streams 1.0.4 specification compliance

## Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "math"

    "github.com/droxer/RxGo/pkg/rx/streams"
)

type IntSubscriber struct {
    name string
}

func (s *IntSubscriber) OnSubscribe(sub streams.Subscription) {
    fmt.Printf("[%s] Starting subscription\n", s.name)
    sub.Request(math.MaxInt64) // Request all items
}

func (s *IntSubscriber) OnNext(value int) {
    fmt.Printf("[%s] Received: %d\n", s.name, value)
}

func (s *IntSubscriber) OnError(err error) {
    fmt.Printf("[%s] Error: %v\n", s.name, err)
}

func (s *IntSubscriber) OnComplete() {
    fmt.Printf("[%s] Completed\n", s.name)
}

func main() {
    // Create publisher from range
    publisher := streams.RangePublisher(1, 10)
    publisher.Subscribe(context.Background(), &IntSubscriber{name: "Range"})
}
```

## Using FromSlice

Create publisher from slice values:

```go
// Create publisher from slice values
publisher := streams.FromSlicePublisher([]int{100, 200, 300, 400, 500})
publisher.Subscribe(context.Background(), &IntSubscriber{name: "FromSlice"})
```

## Custom Publisher Creation

Create custom publishers with context support:

```go
import "github.com/droxer/RxGo/pkg/rx/streams"

// Custom publisher creation
customPublisher := streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[int]) {
    defer sub.OnComplete()
    
    values := []int{1, 2, 3, 4, 5}
    for _, v := range values {
        select {
        case <-ctx.Done():
            sub.OnError(ctx.Err())
            return
        default:
            sub.OnNext(v * 10)
        }
    }
})
```

## Key Concepts

- **Publisher/Subscriber Pattern**: The `streams` package implements the Reactive Streams specification
- **Backpressure Support**: All publishers support backpressure through the Subscription interface
- **Context Support**: All publishers support context cancellation for graceful shutdown
- **Type Safety**: Generic types ensure compile-time type safety
- **Compliance**: Full Reactive Streams 1.0.4 specification compliance

## Complete Example

Here's a complete reactive streams example:

```go
package main

import (
    "context"
    "fmt"
    "math"
    "time"

    "github.com/droxer/RxGo/pkg/rx/streams"
)

type IntSubscriber struct {
    name string
}

func (s *IntSubscriber) OnSubscribe(sub streams.Subscription) {
    fmt.Printf("[%s] Starting subscription\n", s.name)
    sub.Request(math.MaxInt64) // Request all items
}

func (s *IntSubscriber) OnNext(value int) {
    fmt.Printf("[%s] Received: %d\n", s.name, value)
}

func (s *IntSubscriber) OnError(err error) {
    fmt.Printf("[%s] Error: %v\n", s.name, err)
}

func (s *IntSubscriber) OnComplete() {
    fmt.Printf("[%s] Completed\n", s.name)
}

func main() {
    fmt.Println("=== Reactive Streams Example ===")

    // Example 1: Range publisher
    fmt.Println("\n1. Range Publisher:")
    rangePublisher := streams.RangePublisher(1, 5)
    rangePublisher.Subscribe(context.Background(), &IntSubscriber{name: "Range"})

    // Example 2: FromSlice publisher
    fmt.Println("\n2. FromSlice Publisher:")
    slicePublisher := streams.FromSlicePublisher([]int{10, 20, 30, 40, 50})
    slicePublisher.Subscribe(context.Background(), &IntSubscriber{name: "FromSlice"})

    time.Sleep(100 * time.Millisecond)
    fmt.Println("\n=== Reactive Streams completed ===")
}
```