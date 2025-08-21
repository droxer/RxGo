# Observable API

## Quick Overview

| Function | Purpose | Example |
|----------|---------|---------|
| `Just` | Create from values | `observable.Just(1, 2, 3)` |
| `Range` | Create from range | `observable.Range(1, 5)` |
| `Create` | Custom logic | `observable.Create(...)` |
| `FromSlice` | Create from slice | `observable.FromSlice([]int{1,2,3})` |
| `Empty` | Create empty observable | `observable.Empty[int]()` |
| `Error` | Create error observable | `observable.Error[int](err)` |

## Push Model

The Observable API implements a **push-based model** without built-in backpressure, which is different from the pull-based Reactive Streams API:

- **Push Model**: Observables push data to subscribers as soon as it's available
- **No Backpressure**: Producers control emission rate without subscriber demand control
- **Simplicity**: Clean and intuitive API for basic reactive programming

## Basic Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/droxer/RxGo/pkg/observable"
)

func main() {
	// Using Just to create observable
	justObservable := observable.Just(1, 2, 3, 4, 5)
	justObservable.Subscribe(context.Background(), observable.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func() { fmt.Println("Completed") },
		func(err error) { fmt.Printf("Error: %v\n", err) },
	))
}
```

## Using Range

Create observable from range of integers:

```go
// Create observable from range of integers
rangeObservable := observable.Range(10, 5) // Emits 10, 11, 12, 13, 14
rangeObservable.Subscribe(context.Background(), observable.NewSubscriber(
	func(v int) { fmt.Printf("Received: %d\n", v) },
	func() { fmt.Println("Completed") },
	func(err error) { fmt.Printf("Error: %v\n", err) },
))
```

## Using Create with Custom Logic

Create custom observable with your own logic:

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
customObservable.Subscribe(context.Background(), observable.NewSubscriber(
	func(v int) { fmt.Printf("Received: %d\n", v) },
	func() { fmt.Println("Completed") },
	func(err error) { fmt.Printf("Error: %v\n", err) },
))
```

## Using Operators

Transform and filter data using operators:

```go
import (
    "github.com/droxer/RxGo/pkg/observable"
    "github.com/droxer/RxGo/pkg/scheduler"
)

// Transform data using operators
obs := observable.Range(1, 10)
transformed := observable.Map(obs, func(x int) int { return x * 2 })
filtered := observable.Filter(transformed, func(x int) bool { return x > 10 })
filtered.Subscribe(context.Background(), observable.NewSubscriber(
	func(v int) { fmt.Printf("Received: %d\n", v) },
	func() { fmt.Println("Completed") },
	func(err error) { fmt.Printf("Error: %v\n", err) },
))
```

## Using Schedulers

Control execution context with different schedulers:

```go
import (
    "github.com/droxer/RxGo/pkg/observable"
    "github.com/droxer/RxGo/pkg/scheduler"
)

// Use different schedulers
obs := observable.Range(1, 5)

// Computation scheduler for CPU-bound work
observable.ObserveOn(obs, scheduler.Computation()).Subscribe(
	context.Background(), 
	observable.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func() { fmt.Println("Completed") },
		func(err error) { fmt.Printf("Error: %v\n", err) },
	))

// IO scheduler for IO-bound work
observable.ObserveOn(obs, scheduler.IO()).Subscribe(
	context.Background(), 
	observable.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func() { fmt.Println("Completed") },
		func(err error) { fmt.Printf("Error: %v\n", err) },
	))

// Single thread for sequential processing
observable.ObserveOn(obs, scheduler.NewSingleThreadScheduler()).Subscribe(
	context.Background(), 
	observable.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func() { fmt.Println("Completed") },
		func(err error) { fmt.Printf("Error: %v\n", err) },
	))
```

## Context Cancellation

Use context for graceful cancellation:

```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

// Observable that respects context cancellation
contextObservable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
    for i := 0; i < 5; i++ {
        select {
        case <-ctx.Done():
            sub.OnError(ctx.Err())
            return
        default:
            time.Sleep(30 * time.Millisecond)
            sub.OnNext(i * 100)
        }
    }
    sub.OnCompleted()
})
contextObservable.Subscribe(ctx, observable.NewSubscriber(
	func(v int) { fmt.Printf("Received: %d\n", v) },
	func() { fmt.Println("Completed") },
	func(err error) { fmt.Printf("Error: %v\n", err) },
))
```

## Key Concepts

- **Observable/Subscriber Pattern**: The `observable` package provides a clean Observable API
- **Push-based Model**: Data is pushed to subscribers as soon as it's available
- **Context Support**: All observables support context cancellation for graceful shutdown
- **Type Safety**: Generic types ensure compile-time type safety
- **Simplicity**: Clean and intuitive API for basic reactive programming

## Available Operators

The following operators are available in the `observable` package:

- **Map**: Transform each item emitted by an Observable
- **Filter**: Filter items emitted by an Observable
- **ObserveOn**: Specify the scheduler on which an Observable will operate
- **Merge**: Combine multiple Observables into one by merging their emissions
- **Concat**: Concatenate multiple Observables into a single Observable
- **Take**: Emit only the first n items emitted by an Observable
- **Skip**: Skip the first n items emitted by an Observable
- **Distinct**: Suppress duplicate items emitted by an Observable

## observable vs. streams packages

This library provides two distinct packages for reactive programming:

-   **`observable`**: This package provides a simple, Rx-style `Subscriber` interface. It's easy to use and is a good choice for basic reactive programming scenarios where you don't need fine-grained control over the data flow.

-   **`streams`**: This package provides a more complex, Reactive Streams-compliant `Subscriber` interface that supports backpressure. This is a better choice for more advanced scenarios where you need to control the rate of data production to prevent overwhelming consumers.

For more details on the Reactive Streams implementation, see the [Reactive Streams documentation](./reactive-streams.md).

## Complete Example

Here's a complete example combining multiple concepts:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/scheduler"
)

func main() {
	fmt.Println("=== RxGo Basic Example ===")

	// Example 1: Basic usage with Just
	fmt.Println("\n1. Using Just():")
	justObservable := observable.Just(1, 2, 3, 4, 5)
	justObservable.Subscribe(context.Background(), observable.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func() { fmt.Println("Completed") },
		func(err error) { fmt.Printf("Error: %v\n", err) },
	))

	// Example 2: Range observable
	fmt.Println("\n2. Using Range():")
	rangeObservable := observable.Range(10, 5)
	rangeObservable.Subscribe(context.Background(), observable.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func() { fmt.Println("Completed") },
		func(err error) { fmt.Printf("Error: %v\n", err) },
	))

	// Example 3: Using operators
	fmt.Println("\n3. Using operators:")
	obs := observable.Range(1, 5)
	transformed := observable.Map(obs, func(x int) int { return x * 10 })
	filtered := observable.Filter(transformed, func(x int) bool { return x > 20 })
	observable.ObserveOn(filtered, scheduler.Computation()).Subscribe(
		context.Background(),
		observable.NewSubscriber(
			func(v int) { fmt.Printf("Received: %d\n", v) },
			func() { fmt.Println("Completed") },
			func(err error) { fmt.Printf("Error: %v\n", err) },
		))

	// Example 4: With context cancellation
	fmt.Println("\n4. With context cancellation:")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	contextObservable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
		for i := 0; i < 5; i++ {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				time.Sleep(30 * time.Millisecond)
				sub.OnNext(i * 100)
			}
		}
		sub.OnCompleted()
	})
	contextObservable.Subscribe(ctx, observable.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func() { fmt.Println("Completed") },
		func(err error) { fmt.Printf("Error: %v\n", err) },
	))

	time.Sleep(200 * time.Millisecond)
	fmt.Println("\n=== All examples completed ===")
}
```