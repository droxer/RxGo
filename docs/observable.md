# Observable API

## Quick Overview

| Function | Purpose | Example |
|----------|---------|---------|
| `Just` | Create from values | `rx.Just(1, 2, 3)` |
| `Range` | Create from range | `rx.Range(1, 5)` |
| `Create` | Custom logic | `rx.Create(...)` |

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
rangeObservable.Subscribe(context.Background(), &IntSubscriber{name: "Range"})
```

## Using Create with Custom Logic

Create custom observable with your own logic:

```go
// Create custom observable
import "github.com/droxer/RxGo/pkg/observable"

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

## Using Operators

Transform and filter data using operators:

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

Control execution context with different schedulers:

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

## Key Concepts

- **Observable/Subscriber Pattern**: The `rx` package provides a clean Observable API
- **Push-based Model**: Data is pushed to subscribers as soon as it's available
- **Context Support**: All observables support context cancellation for graceful shutdown
- **Type Safety**: Generic types ensure compile-time type safety
- **Simplicity**: Clean and intuitive API for basic reactive programming

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
	"github.com/droxer/RxGo/pkg/operators"
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
	transformed := operators.Map(obs, func(x int) int { return x * 10 })
	filtered := operators.Filter(transformed, func(x int) bool { return x > 20 })
	operators.ObserveOn(filtered, scheduler.Computation()).Subscribe(
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