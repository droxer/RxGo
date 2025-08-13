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

## Next Steps

For more detailed examples and advanced usage patterns, see:

- **[Basic Usage](./basic-usage.md)** - Simple Observable API examples
- **[Reactive Streams API](./reactive-streams.md)** - Full Reactive Streams 1.0.4 compliance
- **[Backpressure Control](./backpressure.md)** - Handle producer/consumer speed mismatches
- **[Context Cancellation](./context-cancellation.md)** - Graceful cancellation using Go context
- **[Data Transformation](./data-transformation.md)** - Transform and process data streams

## Running the Example

You can run this example directly:

```bash
go run examples/basic/basic.go
```

Or create your own Go file with the code above and run:

```bash
go mod init my-project
go get github.com/droxer/RxGo@latest
go run main.go
```