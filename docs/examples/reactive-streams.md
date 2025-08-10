# Reactive Streams API (Advanced)

This document contains advanced usage examples using the full Reactive Streams 1.0.3 compliant API with backpressure support.

## 1. Reactive Streams Publisher

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

## 2. Custom Publisher Creation

Create custom publishers with full backpressure support:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// CustomPublisher demonstrates creating a reactive publisher
func CustomPublisher() rxgo.Publisher[string] {
    return rxgo.NewReactivePublisher(func(ctx context.Context, sub rxgo.ReactiveSubscriber[string]) {
        messages := []string{"Hello", "World", "from", "RxGo"}
        
        for _, msg := range messages {
            select {
            case <-ctx.Done():
                sub.OnError(ctx.Err())
                return
            default:
                sub.OnNext(msg)
                time.Sleep(50 * time.Millisecond)
            }
        }
        sub.OnComplete()
    })
}

func main() {
    publisher := CustomPublisher()
    subscriber := &LoggingSubscriber[string]{name: "Custom"}
    
    publisher.Subscribe(context.Background(), subscriber)
    time.Sleep(500 * time.Millisecond)
}
```

## 3. FromSlice Publisher

Create publishers from existing slices:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

func main() {
    data := []int{100, 200, 300, 400, 500}
    publisher := rxgo.FromSlice(data)
    
    subscriber := &LoggingSubscriber[int]{name: "FromSlice"}
    publisher.Subscribe(context.Background(), subscriber)
}
```