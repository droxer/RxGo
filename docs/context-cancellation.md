# Context Cancellation

This document demonstrates how to use Go context for graceful cancellation of reactive streams, consistent with the actual context example.

## Context Cancellation Example

Cancel streams using context timeout for graceful shutdown:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/droxer/RxGo/pkg/rx"
)

type ContextAwareSubscriber struct {
    received int
}

func (s *ContextAwareSubscriber) Start() {
    fmt.Println("Context-aware subscriber started")
}

func (s *ContextAwareSubscriber) OnNext(value int) {
    s.received++
    fmt.Printf("Received: %d\n", value)
}

func (s *ContextAwareSubscriber) OnError(err error) {
    fmt.Printf("Context cancelled: %v\n", err)
}

func (s *ContextAwareSubscriber) OnCompleted() {
    fmt.Printf("Completed, total received: %d\n", s.received)
}

func main() {
    fmt.Println("=== Context Cancellation Example ===")

    ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel()

    subscriber := &ContextAwareSubscriber{}

    // Create observable that respects context cancellation
    observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
        defer sub.OnCompleted()
        for i := 1; i <= 100; i++ {
            select {
            case <-ctx.Done():
                sub.OnError(ctx.Err())
                return
            default:
                sub.OnNext(i)
                time.Sleep(10 * time.Millisecond) // Small delay to show cancellation
            }
        }
    })

    observable.Subscribe(ctx, subscriber)
    time.Sleep(500 * time.Millisecond) // Wait for completion
    fmt.Println("Context cancellation example completed!")
}
```

## Expected Output

```
=== Context Cancellation Example ===
Context-aware subscriber started
Received: 1
Received: 2
Received: 3
...
Context cancelled: context deadline exceeded
Context cancellation example completed!
```

## Key Concepts

- **Context Timeout**: Uses `context.WithTimeout` to automatically cancel after specified duration
- **Graceful Shutdown**: The observable respects context cancellation and exits cleanly
- **Real-world Usage**: Essential for production systems with timeout requirements
- **Select Statement**: Uses Go's select statement to handle context cancellation
