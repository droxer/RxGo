# RxGo

Reactive Extensions for Go with full Reactive Streams 1.0.4 compliance.

[![Go Report Card](https://goreportcard.com/badge/github.com/droxer/RxGo)](https://goreportcard.com/report/github.com/droxer/RxGo)
[![GoDoc](https://godoc.org/github.com/droxer/RxGo?status.svg)](https://godoc.org/github.com/droxer/RxGo)

## Features

- **Type-safe generics** - Full Go generics support
- **Reactive Streams 1.0.4** - Complete specification compliance
- **Backpressure strategies** - Buffer, Drop, Latest, Error
- **Push & Pull Models** - Observable API (push) and Reactive Streams (pull with backpressure)
- **Retry and backoff** - Fixed, Linear, Exponential backoff with configurable retry limits
- **Context cancellation** - Graceful shutdown
- **Thread-safe** - All APIs are safe for concurrent access with proper synchronization

## Quick Start

Install the library:

```bash
go get github.com/droxer/RxGo@latest
```

### Observable API (Push Model)

Simple and intuitive API for basic reactive programming:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/observable"
)

func main() {
    // Basic usage
    obs := observable.Just(1, 2, 3, 4, 5)
    obs.Subscribe(context.Background(), observable.NewSubscriber(
        func(v int) { fmt.Printf("Got %d\n", v) },
        func() { fmt.Println("Done") },
        func(err error) { fmt.Printf("Error: %v\n", err) },
    ))
    
    // Using new operators
    numbers := observable.Range(1, 10)
    firstFive := observable.Take(numbers, 5)
    fmt.Println("\nFirst five numbers:")
    firstFive.Subscribe(context.Background(), observable.NewSubscriber(
        func(v int) { fmt.Printf("%d ", v) },
        func() { fmt.Println("\nCompleted") },
        func(err error) { fmt.Printf("Error: %v\n", err) },
    ))
}
```

**Output:**
```
Got 1
Got 2
Got 3
Got 4
Got 5
Done

First five numbers:
1 2 3 4 5 
Completed
```

### Reactive Streams API (Pull Model with Backpressure)

Full Reactive Streams 1.0.4 compliance with backpressure support:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/streams"
)

func main() {
    publisher := streams.NewCompliantRangePublisher(1, 5)
    publisher.Subscribe(context.Background(), streams.NewSubscriber(
        func(v int) { fmt.Printf("Received: %d\n", v) },
        func(err error) { fmt.Printf("Error: %v\n", err) },
        func() { fmt.Println("Completed") },
    ))
    
    // Using new processors
    fmt.Println("\nUsing TakeProcessor:")
    numbers := streams.NewCompliantRangePublisher(1, 10)
    takeProcessor := streams.NewTakeProcessor[int](5)
    numbers.Subscribe(context.Background(), takeProcessor)
    takeProcessor.Subscribe(context.Background(), streams.NewSubscriber(
        func(v int) { fmt.Printf("%d ", v) },
        func(err error) { fmt.Printf("Error: %v\n", err) },
        func() { fmt.Println("\nTake completed") },
    ))
}
```


**Output:**
```
Received: 1
Received: 2
Received: 3
Received: 4
Received: 5
Completed

Using TakeProcessor:
1 2 3 4 5 
Take completed
```


### Backpressure Strategies

Handle producer/consumer speed mismatches with four strategies:

```go
import "github.com/droxer/RxGo/pkg/streams"

// Buffer - keep all items in bounded buffer
publisher := streams.NewBufferedPublisher[int](
    streams.WithBufferStrategy(streams.Buffer),
    streams.WithBufferSize(100),
)

// Drop - discard new items when full
publisher := streams.NewBufferedPublisher[int](
    streams.WithBufferStrategy(streams.Drop),
    streams.WithBufferSize(50),
)

// Latest - keep only latest item
publisher := streams.NewBufferedPublisher[int](
    streams.WithBufferStrategy(streams.Latest),
    streams.WithBufferSize(1),
)

// Error - signal error on overflow
publisher := streams.NewBufferedPublisher[int](
    streams.WithBufferStrategy(streams.Error),
    streams.WithBufferSize(10),
)
```

### Bridging APIs

You can easily convert between the `Observable` and `Publisher` APIs using adapters. This is useful when you need to combine the simplicity of the `observable` package with the backpressure support of the `streams` package.

```go
import (
    "github.com/droxer/RxGo/pkg/adapters"
    "github.com/droxer/RxGo/pkg/observable"
    "github.com/droxer/RxGo/pkg/streams"
)

// Convert an Observable to a Publisher
obs := observable.Just(1, 2, 3)
publisher := adapters.ObservablePublisherAdapter(obs)

// Convert a Publisher to an Observable
pub := streams.NewCompliantRangePublisher(1, 5)
observable := adapters.PublisherToObservableAdapter(pub)
```

## Documentation

- [Architecture](./docs/architecture.md) - Package structure and design decisions
- [Observable API](./docs/observable.md) - Simple Observable API examples
- [Reactive Streams](./docs/reactive-streams.md) - Full Reactive Streams 1.0.4 compliance
- [Backpressure](./docs/backpressure.md) - Handle producer/consumer speed mismatches
- [Push vs Pull Models](./docs/push-pull-models.md) - Understanding push and pull models with backpressure
- [Retry and Backoff](./docs/retry-backoff.md) - Configurable retry with backoff strategies
- [Transformations](./docs/transformations.md) - Transform and process data streams with both Reactive Streams and Observable API
- [Context Cancellation](./docs/context-cancellation.md) - Graceful cancellation using Go context
- [Schedulers](./docs/schedulers.md) - Execution context control

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

[MIT License](./LICENSE)
