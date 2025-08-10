# Basic Usage

This document contains basic usage examples for RxGo.

## 1. Observable API (Simple)

The simplest way to get started with RxGo using the backward-compatible Observable API:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/observable"
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
    observable := observable.Just(1, 2, 3, 4, 5)
    observable.Subscribe(context.Background(), &IntSubscriber{name: "Just"})
}
```

## 2. Range Observable

Create observables from ranges:

```go
rangeObservable := observable.Range(10, 5)
rangeObservable.Subscribe(context.Background(), &IntSubscriber{name: "Range"})
```

## 3. Create with Custom Logic

Create observables with custom logic:

```go
customObservable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
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