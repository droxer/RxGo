# Backpressure Control

This document demonstrates how to handle producer/consumer speed mismatches using backpressure strategies.

## 1. Slow Subscriber with Backpressure

Handle slow consumers using controlled demand requests:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// SlowSubscriber demonstrates controlled consumption
type SlowSubscriber struct {
    received []int
}

func (s *SlowSubscriber) OnSubscribe(sub rxgo.Subscription) {
    // Request items slowly
    go func() {
        for i := 0; i < 5; i++ {
            time.Sleep(100 * time.Millisecond)
            sub.Request(1)
        }
    }()
}

func (s *SlowSubscriber) OnNext(value int) {
    s.received = append(s.received, value)
    fmt.Printf("Processing: %d\n", value)
    time.Sleep(50 * time.Millisecond) // Simulate slow processing
}

func (s *SlowSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *SlowSubscriber) OnComplete() { fmt.Println("Done!") }

func main() {
    publisher := rxgo.RangePublisher(1, 100)
    subscriber := &SlowSubscriber{}
    publisher.Subscribe(context.Background(), subscriber)
    
    time.Sleep(time.Second)
}
```

## 2. Buffered Backpressure Strategy

Use buffering to handle temporary speed mismatches:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// BufferedSubscriber shows buffering strategy
type BufferedSubscriber struct {
    bufferSize int
}

func (s *BufferedSubscriber) OnSubscribe(sub rxgo.Subscription) {
    // Request more than immediate need to buffer
    sub.Request(int64(s.bufferSize))
}

func (s *BufferedSubscriber) OnNext(value int) {
    fmt.Printf("Buffered processing: %d\n", value)
    time.Sleep(10 * time.Millisecond) // Simulate processing
}

func (s *BufferedSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *BufferedSubscriber) OnComplete() { fmt.Println("Buffered processing complete") }

func main() {
    publisher := rxgo.RangePublisher(1, 20)
    subscriber := &BufferedSubscriber{bufferSize: 10}
    publisher.Subscribe(context.Background(), subscriber)
    
    time.Sleep(500 * time.Millisecond)
}
```

## 3. Demand-Based Flow Control

Implement sophisticated demand management:

```go
package main

import (
    "context"
    "fmt"
    "sync/atomic"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// DemandSubscriber implements dynamic demand management
type DemandSubscriber struct {
    processed int64
    target    int64
}

func (s *DemandSubscriber) OnSubscribe(sub rxgo.Subscription) {
    // Start with initial demand
    sub.Request(5)
}

func (s *DemandSubscriber) OnNext(value int) {
    count := atomic.AddInt64(&s.processed, 1)
    fmt.Printf("Processing item %d: %d\n", count, value)
    
    // Request more when 80% of previous batch is processed
    if count%s.target == int64(float64(s.target)*0.8) {
        fmt.Printf("Requesting next batch of %d items\n", s.target)
        // Note: In real implementation, we'd need access to subscription
    }
    
    time.Sleep(20 * time.Millisecond) // Simulate processing
}

func (s *DemandSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *DemandSubscriber) OnComplete() { fmt.Println("All items processed") }

func main() {
    publisher := rxgo.RangePublisher(1, 50)
    subscriber := &DemandSubscriber{target: 10}
    publisher.Subscribe(context.Background(), subscriber)
    
    time.Sleep(2 * time.Second)
}
```