# Data Transformation

This document demonstrates data transformation concepts using both the Observable API (push model) and Reactive Streams API (pull model with backpressure).

## Overview

RxGo provides two distinct APIs for data transformation:

- **Observable API**: Simple push-based model for basic reactive programming
- **Reactive Streams API**: Pull-based model with full backpressure support

Both APIs offer similar transformation capabilities but with different usage patterns and features.

## Observable API Transformations

### Basic Transformations

Transform data using operators from the `pkg/observable` package:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/observable"
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
    obs := observable.Range(1, 5)
    
    // Transform using Map operator
    mapped := observable.Map(obs, func(x int) int { return x * 2 })
    
    // Filter values
    filtered := observable.Filter(mapped, func(x int) bool { return x > 5 })
    
    // Subscribe to receive transformed values
    if err := filtered.Subscribe(context.Background(), &IntSubscriber{name: "Observable"}); err != nil {
        fmt.Printf("Subscription failed: %v\n", err)
    }
}
```

### Available Observable Operators

| Operator | Purpose | Example |
|----------|---------|---------|
| `Map` | Transform values | `observable.Map(obs, func(x int) int { return x * 2 })` |
| `Filter` | Filter values | `observable.Filter(obs, func(x int) bool { return x > 5 })` |
| `ObserveOn` | Control execution context | `observable.ObserveOn(obs, scheduler.Computation)` |
| `Merge` | Combine multiple observables | `observable.Merge(obs1, obs2, obs3)` |
| `Concat` | Sequentially emit from multiple observables | `observable.Concat(obs1, obs2, obs3)` |
| `Take` | Emit only first N values | `observable.Take(obs, 5)` |
| `Skip` | Skip first N values | `observable.Skip(obs, 3)` |
| `Distinct` | Emit only unique values | `observable.Distinct(obs)` |

## Reactive Streams API Transformations

### Basic Transformations

Transform data using processors from the `pkg/streams` package:

```go
package main

import (
    "context"
    "fmt"
    "math"
    
    "github.com/droxer/RxGo/pkg/streams"
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
    publisher := streams.NewCompliantRangePublisher(1, 5)
    
    // Transform using Map processor
    mapper := streams.NewMapProcessor(func(x int) int { return x * 2 })
    publisher.Subscribe(context.Background(), mapper)
    
    // Filter values using Filter processor
    filter := streams.NewFilterProcessor(func(x int) bool { return x > 5 })
    mapper.Subscribe(context.Background(), filter)
    
    // Subscribe to receive transformed values
    filter.Subscribe(context.Background(), &IntSubscriber{name: "ReactiveStreams"})
}
```

### Available Stream Processors

| Processor | Purpose | Example |
|-----------|---------|---------|
| `MapProcessor` | Transform values | `streams.NewMapProcessor(func(x int) int { return x * 2 })` |
| `FilterProcessor` | Filter values | `streams.NewFilterProcessor(func(x int) bool { return x > 5 })` |
| `FlatMapProcessor` | Transform and flatten | `streams.NewFlatMapProcessor(transformFunc)` |
| `MergeProcessor` | Combine multiple publishers | `streams.NewMergeProcessor(pub1, pub2, pub3)` |
| `ConcatProcessor` | Sequentially emit from multiple publishers | `streams.NewConcatProcessor(pub1, pub2, pub3)` |
| `TakeProcessor` | Limit number of emissions | `streams.NewTakeProcessor[int](5)` |
| `SkipProcessor` | Skip initial emissions | `streams.NewSkipProcessor[int](3)` |
| `DistinctProcessor` | Filter duplicate values | `streams.NewDistinctProcessor[int]()` |

## Comparison of Transformation Approaches

### Observable API (Push Model)
- **Usage**: Simple function calls on observables
- **Syntax**: `observable.Map(observable, transformFunc)`
- **Backpressure**: No built-in backpressure support
- **Best for**: Simple transformations, prototyping, predictable data rates

### Reactive Streams API (Pull Model)
- **Usage**: Chain processors using Subscribe method
- **Syntax**: `processor := streams.NewMapProcessor(transformFunc); publisher.Subscribe(ctx, processor)`
- **Backpressure**: Full Reactive Streams 1.0.4 compliance with demand control
- **Best for**: Production systems, unbounded data streams, backpressure requirements

## Complete Example: Parallel Transformations

Here's a complete example showing both APIs side by side:

```go
package main

import (
    "context"
    "fmt"
    "math"
    "time"
    
    "github.com/droxer/RxGo/pkg/observable"
    "github.com/droxer/RxGo/pkg/streams"
)

// Observable API subscriber
type ObservableSubscriber struct {
    name string
}

func (s *ObservableSubscriber) Start() {
    fmt.Printf("[Observable %s] Starting subscription\n", s.name)
}
func (s *ObservableSubscriber) OnNext(value int) {
    fmt.Printf("[Observable %s] Received: %d\n", s.name, value)
}
func (s *ObservableSubscriber) OnError(err error) {
    fmt.Printf("[Observable %s] Error: %v\n", s.name, err)
}
func (s *ObservableSubscriber) OnCompleted() {
    fmt.Printf("[Observable %s] Completed\n", s.name)
}

// Reactive Streams subscriber
type StreamsSubscriber struct {
    name string
}

func (s *StreamsSubscriber) OnSubscribe(sub streams.Subscription) {
    fmt.Printf("[Streams %s] Starting subscription\n", s.name)
    sub.Request(math.MaxInt64) // Request all items
}

func (s *StreamsSubscriber) OnNext(value int) {
    fmt.Printf("[Streams %s] Received: %d\n", s.name, value)
}

func (s *StreamsSubscriber) OnError(err error) {
    fmt.Printf("[Streams %s] Error: %v\n", s.name, err)
}

func (s *StreamsSubscriber) OnComplete() {
    fmt.Printf("[Streams %s] Completed\n", s.name)
}

func main() {
    fmt.Println("=== Observable API Transformation ===")
    
    // Observable API transformation
    obs := observable.Range(1, 5)
    mapped := observable.Map(obs, func(x int) int { return x * 3 })
    filtered := observable.Filter(mapped, func(x int) bool { return x > 6 })
    if err := filtered.Subscribe(context.Background(), &ObservableSubscriber{name: "Transform"}); err != nil {
        fmt.Printf("Subscription failed: %v\n", err)
    }
    
    fmt.Println("\n=== Reactive Streams API Transformation ===")
    
    // Reactive Streams API transformation
    publisher := streams.NewCompliantRangePublisher(1, 5)
    mapper := streams.NewMapProcessor(func(x int) int { return x * 3 })
    publisher.Subscribe(context.Background(), mapper)
    filter := streams.NewFilterProcessor(func(x int) bool { return x > 6 })
    mapper.Subscribe(context.Background(), filter)
    filter.Subscribe(context.Background(), &StreamsSubscriber{name: "Transform"})
    
    time.Sleep(100 * time.Millisecond)
    fmt.Println("\n=== All transformations completed ===")
}
```

## Advanced Examples

### Merging and Concatenating Streams

Both APIs support combining multiple streams:

```go
// Observable API - Merge
obs1 := observable.Just(1, 2, 3)
obs2 := observable.Just(4, 5, 6)
merged := observable.Merge(obs1, obs2)

// Observable API - Concat
concatenated := observable.Concat(obs1, obs2)

// Reactive Streams API - Merge
pub1 := streams.NewCompliantFromSlicePublisher([]int{1, 2, 3})
pub2 := streams.NewCompliantFromSlicePublisher([]int{4, 5, 6})
mergeProcessor := streams.NewMergeProcessor(pub1, pub2)

// Reactive Streams API - Concat
concatProcessor := streams.NewConcatProcessor(pub1, pub2)
```

### Stream Control Operators

Control the flow of your streams with take, skip, and distinct operations:

```go
// Observable API
numbers := observable.Range(1, 10)
firstFive := observable.Take(numbers, 5)
skipped := observable.Skip(numbers, 3)
unique := observable.Distinct(observable.FromSlice([]int{1, 2, 2, 3, 3, 4}))

// Reactive Streams API
publisher := streams.NewCompliantRangePublisher(1, 10)
takeProcessor := streams.NewTakeProcessor[int](5)
skipProcessor := streams.NewSkipProcessor[int](3)
distinctProcessor := streams.NewDistinctProcessor[int]()
```

## Key Concepts

### 1. Transformation Composition
- **Observable API**: Chain operators using function composition
- **Reactive Streams API**: Chain processors using Subscribe method calls

### 2. Type Safety
- Both APIs use Go generics for compile-time type safety
- Transformations maintain type information throughout the chain

### 3. Error Handling
- Both APIs propagate errors through the transformation chain
- Subscribers handle errors in the OnError method

### 4. Context Support
- Both APIs support context cancellation for graceful shutdown
- Context is passed through the Subscribe method

## When to Use Each API

### Use Observable API When:
- Building simple reactive applications
- Prototyping or learning reactive programming
- Working with predictable data rates
- Need clean, intuitive API

### Use Reactive Streams API When:
- Building production systems with potentially unbounded data streams
- Need for backpressure control
- Requiring Reactive Streams 1.0.4 compliance
- Working with external data sources that may outpace consumers