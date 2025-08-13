# Reactive Streams API

This document demonstrates the Reactive Streams API using both `rxgo` and `observable` packages, consistent with actual examples.

## Basic Reactive Streams Usage

The actual examples show a simplified approach using the Observable API. Here's how the Reactive Streams concepts map to the actual implementations:

## Using Range

Create and process a range of values:

```go
package main

import (
    "context"
    "fmt"

    "github.com/droxer/RxGo/pkg/rx"
)

type IntSubscriber struct {
    name string
}

func (s *IntSubscriber) Start() {
    fmt.Printf("[%s] Starting subscription\n", s.name)
}
func (s *IntSubscriber) OnNext(value int) {
    fmt.Printf("[%s] Received: %d\n", s.name, value)
}
func (s *IntSubscriber) OnError(err error) {
    fmt.Printf("[%s] Error: %v\n", s.name, err)
}
func (s *IntSubscriber) OnCompleted() {
    fmt.Printf("[%s] Completed\n", s.name)
}

func main() {
    // Create observable from range
    obs := rx.Range(1, 10)
    obs.Subscribe(context.Background(), &IntSubscriber{name: "Range"})
}
```

## Using Just

Create observable from literal values:

```go
// Create observable from literal values
obs := rx.Just(100, 200, 300, 400, 500)
obs.Subscribe(context.Background(), &IntSubscriber{name: "Just"})
```

## Custom Observable Creation

Create custom observables with context support:

```go
import "github.com/droxer/RxGo/pkg/rx"

// Custom observable creation
customObservable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
    defer sub.OnCompleted()
    
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

- **Observable as Publisher**: The `rxgo.Range()` and `rxgo.Just()` functions act as publishers
- **Subscriber Interface**: The actual examples use a simple subscriber interface with `Start()`, `OnNext()`, `OnError()`, and `OnCompleted()` methods
- **Context Support**: All observables support context cancellation for graceful shutdown
- **Synchronous Processing**: The actual examples show synchronous processing patterns

## Complete Example

Here's a complete reactive streams example:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/droxer/RxGo/pkg/rx"
)

type IntSubscriber struct {
    name string
}

func (s *IntSubscriber) Start() {
    fmt.Printf("[%s] Starting subscription\n", s.name)
}
func (s *IntSubscriber) OnNext(value int) {
    fmt.Printf("[%s] Received: %d\n", s.name, value)
}
func (s *IntSubscriber) OnError(err error) {
    fmt.Printf("[%s] Error: %v\n", s.name, err)
}
func (s *IntSubscriber) OnCompleted() {
    fmt.Printf("[%s] Completed\n", s.name)
}

func main() {
    fmt.Println("=== Reactive Streams Example ===")

    // Example 1: Range publisher
    fmt.Println("\n1. Range Publisher:")
    rangeObs := rx.Range(1, 5)
    rangeObs.Subscribe(context.Background(), &IntSubscriber{name: "Range"})

    // Example 2: Just publisher
    fmt.Println("\n2. Just Publisher:")
    justObs := rx.Just(10, 20, 30, 40, 50)
    justObs.Subscribe(context.Background(), &IntSubscriber{name: "Just"})

    time.Sleep(100 * time.Millisecond)
    fmt.Println("\n=== Reactive Streams completed ===")
}
```