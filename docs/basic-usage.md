# Basic Usage

This guide demonstrates the basic Observable API using the reorganized package structure.

## Package Imports

```go
import (
    "context"
    "github.com/droxer/RxGo/pkg/rx"
    "github.com/droxer/RxGo/pkg/rx/operators"
    "github.com/droxer/RxGo/pkg/rx/scheduler"
)
```

## Creating Observables

### Using Just

Create observable from literal values:

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
    // Using Just to create observable
    justObservable := rx.Just(1, 2, 3, 4, 5)
    justObservable.Subscribe(context.Background(), &IntSubscriber{name: "Just"})
}
```

### Using Range

Create observable from range of integers:

```go
// Create observable from range of integers
rangeObservable := rx.Range(10, 5) // Emits 10, 11, 12, 13, 14
rangeObservable.Subscribe(context.Background(), &IntSubscriber{name: "Range"})
```

### Using Create with Custom Logic

Create custom observable with your own logic:

```go
// Create custom observable
import "github.com/droxer/RxGo/pkg/rx"

customObservable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
    for i := 0; i < 3; i++ {
        select {
        case <-ctx.Done():
            sub.OnError(ctx.Err())
            return
        default:
            sub.OnNext(i * 10)
        }
    }
    sub.OnCompleted()
})
customObservable.Subscribe(context.Background(), &IntSubscriber{name: "Create"})
```

## Using Operators

### Map and Filter Operations

```go
import (
    "github.com/droxer/RxGo/pkg/rx"
    "github.com/droxer/RxGo/pkg/rx/operators"
)

// Transform data using operators
obs := rx.Range(1, 10)
transformed := operators.Map(obs, func(x int) int { return x * 2 })
filtered := operators.Filter(transformed, func(x int) bool { return x > 10 })
filtered.Subscribe(context.Background(), &IntSubscriber{name: "Operators"})
```

## Using Schedulers

### With Different Schedulers

```go
import (
    "github.com/droxer/RxGo/pkg/rx"
    "github.com/droxer/RxGo/pkg/rx/operators"
    "github.com/droxer/RxGo/pkg/rx/scheduler"
)

// Use different schedulers
obs := rx.Range(1, 5)

// Computation scheduler for CPU-bound work
operators.ObserveOn(obs, scheduler.Computation).Subscribe(ctx, subscriber)

// IO scheduler for IO-bound work
operators.ObserveOn(obs, scheduler.IO).Subscribe(ctx, subscriber)

// Single thread for sequential processing
operators.ObserveOn(obs, scheduler.SingleThread).Subscribe(ctx, subscriber)
```

## Context Cancellation

Use context for graceful cancellation:

```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

// Observable that respects context cancellation
contextObservable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
    for i := 0; i < 5; i++ {
        select {
        case <-ctx.Done():
            sub.OnError(ctx.Err())
            return
        default:
            sub.OnNext(i * 100)
        }
    }
    sub.OnCompleted()
})
contextObservable.Subscribe(ctx, &IntSubscriber{name: "Context"})
```

## Complete Example

Here's a complete example combining multiple concepts:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rx"
    "github.com/droxer/RxGo/pkg/rx/operators"
    "github.com/droxer/RxGo/pkg/rx/scheduler"
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
    fmt.Println("=== RxGo Basic Example ===")

    // Example 1: Basic usage with Just
    fmt.Println("\n1. Using Just():")
    justObservable := rx.Just(1, 2, 3, 4, 5)
    justObservable.Subscribe(context.Background(), &IntSubscriber{name: "Just"})

    // Example 2: Range observable
    fmt.Println("\n2. Using Range():")
    rangeObservable := rx.Range(10, 5)
    rangeObservable.Subscribe(context.Background(), &IntSubscriber{name: "Range"})

    // Example 3: Using operators
    fmt.Println("\n3. Using operators:")
    obs := rx.Range(1, 5)
    transformed := operators.Map(obs, func(x int) int { return x * 10 })
    filtered := operators.Filter(transformed, func(x int) bool { return x > 20 })
    operators.ObserveOn(filtered, scheduler.Computation).Subscribe(
        context.Background(), 
        &IntSubscriber{name: "Operators"})

    // Example 4: With context cancellation
    fmt.Println("\n4. With context cancellation:")
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    contextObservable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
        for i := 0; i < 5; i++ {
            select {
            case <-ctx.Done():
                sub.OnError(ctx.Err())
                return
            default:
                sub.OnNext(i * 100)
            }
        }
        sub.OnCompleted()
    })
    contextObservable.Subscribe(ctx, &IntSubscriber{name: "Context"})

    time.Sleep(100 * time.Millisecond)
    fmt.Println("\n=== All examples completed ===")
}
```