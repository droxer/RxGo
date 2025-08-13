# Quick Start Guide

This guide provides a simple example to get you started with RxGo.

## Basic Example

Here's a minimal working example that demonstrates the core concepts:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/rx"
)

type IntSubscriber struct{}

func (s *IntSubscriber) Start() {
    fmt.Println("Starting subscription")
}
func (s *IntSubscriber) OnNext(value int) { fmt.Println(value) }
func (s *IntSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *IntSubscriber) OnCompleted() { fmt.Println("Completed!") }

func main() {
    // Using Just to create observable
    obs := rx.Just(1, 2, 3, 4, 5)
    obs.Subscribe(context.Background(), &IntSubscriber{})
}
```

## Expected Output

```
Starting subscription
1
2
3
4
5
Completed!
```

## Backpressure Strategies

RxGo provides four strategies to handle producer/consumer speed mismatches:

```go
import (
    "context"
    "github.com/droxer/RxGo/pkg/rx/streams"
)

// Buffer strategy - keep all items in bounded buffer
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 100,
})

// Drop strategy - discard new items when full
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Drop,
    BufferSize: 50,
})

// Latest strategy - keep only latest item
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Latest,
    BufferSize: 1,
})

// Error strategy - signal error on overflow
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Error,
    BufferSize: 10,
})
```

## Next Steps

For more detailed examples and advanced usage patterns, see:

- **[Basic Usage](./basic-usage.md)** - Simple Observable API examples
- **[Reactive Streams API](./reactive-streams.md)** - Full Reactive Streams 1.0.4 compliance
- **[Backpressure Control](./backpressure.md)** - Handle producer/consumer speed mismatches
- **[Context Cancellation](./context-cancellation.md)** - Graceful cancellation using Go context
- **[Data Transformation](./data-transformation.md)** - Transform and process data streams

## Running Examples

You can run the basic example:

```bash
go run examples/basic/basic.go
```

Or run the backpressure examples:

```bash
go run examples/backpressure/main.go
```

Or create your own Go file with the code above and run:

```bash
go mod init my-project
go get github.com/droxer/RxGo@latest
go run main.go
```