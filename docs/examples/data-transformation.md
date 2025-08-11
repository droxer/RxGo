# Data Transformation

This document demonstrates data transformation concepts using the actual RxGo examples and APIs.

## Available Transformations

Based on the actual examples, here are the core transformation patterns:

### Range Transformation

Transform a range of values:

```go
package main

import (
    "context"
    "fmt"

    "github.com/droxer/RxGo/pkg/rxgo"
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
    // Transform using Range
    obs := rxgo.Range(1, 5)
    obs.Subscribe(context.Background(), &IntSubscriber{name: "Range"})
}
```

### Just Transformation

Create observable from literal values:

```go
// Create observable from literal values and transform
obs := rxgo.Just(10, 20, 30, 40, 50)
obs.Subscribe(context.Background(), &IntSubscriber{name: "Just"})
```

### Create with Custom Logic

Transform data using custom creation logic:

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/droxer/RxGo/pkg/observable"
    "github.com/droxer/RxGo/pkg/rxgo"
)

// Custom transformation using Create
customObservable := rxgo.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
    defer sub.OnCompleted()
    for i := 1; i <= 100; i++ {
        select {
        case <-ctx.Done():
            sub.OnError(ctx.Err())
            return
        default:
            sub.OnNext(i)
            time.Sleep(10 * time.Millisecond) // Transform with delay
        }
    }
})
```

## Practical Transformation Examples

### Processing Numbers

Transform and process a sequence of numbers:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/droxer/RxGo/pkg/rxgo"
)

type ProcessingSubscriber struct {
    name string
}

func (s *ProcessingSubscriber) Start() {
    fmt.Printf("[%s] Starting processing\n", s.name)
}

func (s *ProcessingSubscriber) OnNext(value int) {
    // Transform: multiply by 2
    transformed := value * 2
    fmt.Printf("[%s] Original: %d, Transformed: %d\n", s.name, value, transformed)
}

func (s *ProcessingSubscriber) OnError(err error) {
    fmt.Printf("[%s] Processing error: %v\n", s.name, err)
}

func (s *ProcessingSubscriber) OnCompleted() {
    fmt.Printf("[%s] Processing completed\n", s.name)
}

func main() {
    fmt.Println("=== Data Transformation Example ===")

    // Transform numbers 1-10
    numbers := rxgo.Range(1, 10)
    numbers.Subscribe(context.Background(), &ProcessingSubscriber{name: "Transformer"})

    time.Sleep(100 * time.Millisecond)
    fmt.Println("\n=== Transformation completed ===")
}
```

### Context-Aware Transformations

Transform data with context cancellation support:

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/droxer/RxGo/pkg/observable"
    "github.com/droxer/RxGo/pkg/rxgo"
)

type ContextTransformer struct {
    received int
}

func (s *ContextTransformer) Start() {
    fmt.Println("Context-aware transformation started")
}

func (s *ContextTransformer) OnNext(value int) {
    s.received++
    // Transform: square the number
    squared := value * value
    fmt.Printf("Input: %d, Output: %d\n", value, squared)
}

func (s *ContextTransformer) OnError(err error) {
    fmt.Printf("Transformation cancelled: %v\n", err)
}

func (s *ContextTransformer) OnCompleted() {
    fmt.Printf("Transformation completed, processed: %d items\n", s.received)
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel()

    transformer := &ContextTransformer{}

    // Transform with context support
    obs := rxgo.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        defer sub.OnCompleted()
        for i := 1; i <= 100; i++ {
            select {
            case <-ctx.Done():
                sub.OnError(ctx.Err())
                return
            default:
                sub.OnNext(i)
                time.Sleep(10 * time.Millisecond)
            }
        }
    })

    obs.Subscribe(ctx, transformer)
    time.Sleep(500 * time.Millisecond)
}
```

## Key Transformation Patterns

### 1. Simple Value Transformation
- Use `rxgo.Range()` for sequential numbers
- Use `rxgo.Just()` for literal values
- Use `observable.Create()` for custom sequences

### 2. Context-Aware Processing
- Always use context cancellation for production systems
- Handle `ctx.Done()` for graceful shutdown
- Use select statements for non-blocking operations

### 3. Subscriber-Based Processing
- Implement the subscriber interface with `Start()`, `OnNext()`, `OnError()`, `OnCompleted()`
- Use meaningful subscriber names for debugging
- Handle errors appropriately in the `OnError` method

### 4. Synchronous vs Asynchronous
- The actual examples show synchronous processing patterns
- Use appropriate delays with `time.Sleep()` for demonstration
- Consider using goroutines for concurrent processing in production

## Running Transformation Examples

```bash
# Run basic transformation examples
go run examples/basic/basic.go

# Run context-aware transformations
go run examples/context/context.go

# Run with custom transformations
./bin/basic
./bin/context
```