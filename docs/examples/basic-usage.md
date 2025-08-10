# Basic Usage

This document contains basic usage examples for RxGo.

## 1. Unified RxGo API (Recommended)

The unified RxGo API provides a clean, modern interface for both Observable and Reactive Streams patterns:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// IntSubscriber demonstrates type-safe subscriber with generics
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
    // Create observable from literal values
    observable := rxgo.Just(1, 2, 3, 4, 5)
    observable.Subscribe(context.Background(), &IntSubscriber{name: "Just"})
}
```

## 2. Range Observable

Create observables from ranges:

```go
rangeObservable := rxgo.Range(10, 5)
rangeObservable.Subscribe(context.Background(), &IntSubscriber{name: "Range"})
```

## 3. Create with Custom Logic

Create observables with custom logic:

```go
customObservable := rxgo.Create(func(ctx context.Context, sub rxgo.Subscriber[int]) {
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

## 4. Observable API (Direct)

For direct access to the Observable API:

```go
import "github.com/droxer/RxGo/pkg/observable"

obs := observable.Just(1, 2, 3)
obs.Subscribe(context.Background(), subscriber)
```

## 5. Reactive Streams API (Direct)

For direct access to the Reactive Streams API:

```go
import "github.com/droxer/RxGo/pkg/reactive"

publisher := reactive.Just(1, 2, 3)
publisher.Subscribe(context.Background(), reactiveSubscriber)
```