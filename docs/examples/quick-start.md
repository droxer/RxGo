# Quick Start Guide

This guide provides a simple example to get you started with RxGo.

## Basic Example

Here's a minimal working example that demonstrates the core concepts:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

type IntSubscriber struct{}

func (s *IntSubscriber) Start() {}
func (s *IntSubscriber) OnNext(value int) { fmt.Println(value) }
func (s *IntSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *IntSubscriber) OnCompleted() { fmt.Println("Completed!") }

func main() {
    rxgo.Create(func(ctx context.Context, sub rxgo.Subscriber[int]) {
        for i := 0; i < 5; i++ {
            sub.OnNext(i)
        }
        sub.OnCompleted()
    }).Subscribe(context.Background(), &IntSubscriber{})
}
```

## Expected Output

```
0
1
2
3
4
Completed!
```

## Next Steps

For more detailed examples and advanced usage patterns, see:

- **[Basic Usage](./basic-usage.md)** - Simple Observable API examples
- **[Reactive Streams API](./reactive-streams.md)** - Full Reactive Streams 1.0.3 compliance
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
go get github.com/droxer/RxGo@v0.1.1
go run main.go
```