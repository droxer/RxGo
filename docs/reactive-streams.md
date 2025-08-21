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

	"github.com/droxer/RxGo/pkg/streams"
)

func main() {
	// Create publisher from range
	publisher := streams.RangePublisher(1, 10)
	publisher.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func(err error) { fmt.Printf("Error: %v\n", err) },
		func() { fmt.Println("Completed") },
	))
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
import "github.com/droxer/RxGo/pkg/streams"

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

## streams vs. observable packages

This library provides two distinct packages for reactive programming:

-   **`streams`**: This package provides a more complex, Reactive Streams-compliant `Subscriber` interface that supports backpressure. This is a better choice for more advanced scenarios where you need to control the rate of data production to prevent overwhelming consumers.

-   **`observable`**: This package provides a simple, Rx-style `Subscriber` interface. It's easy to use and is a good choice for basic reactive programming scenarios where you don't need fine-grained control over the data flow.

For more details on the `observable` package, see the [Observable documentation](./observable.md).

## Complete Example

Here's a complete reactive streams example:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/droxer/RxGo/pkg/streams"
)

func main() {
	fmt.Println("=== Reactive Streams Example ===")

	// Example 1: Range publisher
	fmt.Println("\n1. Range Publisher:")
	rangePublisher := streams.RangePublisher(1, 5)
	rangePublisher.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func(err error) { fmt.Printf("Error: %v\n", err) },
		func() { fmt.Println("Completed") },
	))

	// Example 2: FromSlice publisher
	fmt.Println("\n2. FromSlice Publisher:")
	slicePublisher := streams.FromSlicePublisher([]int{10, 20, 30, 40, 50})
	slicePublisher.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func(err error) { fmt.Printf("Error: %v\n", err) },
		func() { fmt.Println("Completed") },
	))

	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n=== Reactive Streams completed ===")
}
```
