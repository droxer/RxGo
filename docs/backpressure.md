# Backpressure Control

This document demonstrates how to handle producer/consumer speed mismatches using backpressure strategies, consistent with the actual backpressure example.

## Basic Backpressure Example

Handle slow consumers using controlled processing:

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/droxer/RxGo/pkg/rxgo"
)

type BackpressureSubscriber struct {
    name     string
    received []int
    mu       sync.Mutex
    wg       *sync.WaitGroup
}

func (s *BackpressureSubscriber) Start() {
    fmt.Printf("[%s] Starting subscription\n", s.name)
}

func (s *BackpressureSubscriber) OnNext(value int) {
    s.mu.Lock()
    s.received = append(s.received, value)
    s.mu.Unlock()
    fmt.Printf("[%s] Processing: %d (total: %d)\n", s.name, value, len(s.received))
    time.Sleep(100 * time.Millisecond) // Simulate slow processing
}

func (s *BackpressureSubscriber) OnError(err error) {
    fmt.Printf("[%s] Error: %v\n", s.name, err)
    s.wg.Done()
}

func (s *BackpressureSubscriber) OnCompleted() {
    fmt.Printf("[%s] Completed, received %d items\n", s.name, len(s.received))
    s.wg.Done()
}

func main() {
    fmt.Println("=== Backpressure Example ===")
    var wg sync.WaitGroup
    wg.Add(1)

    subscriber := &BackpressureSubscriber{
        name: "Backpressure",
        wg:   &wg,
    }

    // Create observable
    obs := rxgo.Range(1, 10)
    obs.Subscribe(context.Background(), subscriber)
    wg.Wait()
    fmt.Println("Backpressure example completed!")
}
```

## Expected Output

When you run this example, you'll see output like:

```
=== Backpressure Example ===
[Backpressure] Starting subscription
[Backpressure] Processing: 1 (total: 1)
[Backpressure] Processing: 2 (total: 2)
[Backpressure] Processing: 3 (total: 3)
...
[Backpressure] Processing: 10 (total: 10)
[Backpressure] Completed, received 10 items
Backpressure example completed!
```

## Key Concepts

- **Controlled Processing**: The subscriber processes items at its own pace using simulated delays
- **Thread Safety**: Uses mutex for safe concurrent access to shared state
- **WaitGroup**: Ensures proper completion handling
- **Real-world Application**: This pattern is useful when consumers are slower than producers