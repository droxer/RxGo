# RxGo

Modern Reactive Extensions for Go

## Examples (Go 1.21+)

RxGo now supports both the original Observable API and the new Reactive Streams API with full backpressure support.

### **1. Original Observable API (Backward Compatible)**

```go
package main

import (
    "context"
    "fmt"
    
    rx "github.com/droxer/RxGo"
)

type IntSubscriber struct{}

func (s *IntSubscriber) Start() {}
func (s *IntSubscriber) OnNext(value int) { fmt.Println(value) }
func (s *IntSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *IntSubscriber) OnCompleted() { fmt.Println("Completed!") }

func main() {
    rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
        for i := 0; i < 5; i++ {
            sub.OnNext(i)
        }
        sub.OnCompleted()
    }).Subscribe(context.Background(), &IntSubscriber{})
}
```

### **2. Reactive Streams API (New - Full Compliance)**

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    rx "github.com/droxer/RxGo"
)

// ReactiveSubscriber implementation
type LoggingSubscriber[T any] struct {
    name string
}

func (s *LoggingSubscriber[T]) OnSubscribe(sub rx.Subscription) {
    fmt.Printf("[%s] Subscribed, requesting 3 items\n", s.name)
    sub.Request(3) // Backpressure control
}

func (s *LoggingSubscriber[T]) OnNext(value T) {
    fmt.Printf("[%s] Received: %v\n", s.name, value)
}

func (s *LoggingSubscriber[T]) OnError(err error) {
    fmt.Printf("[%s] Error: %v\n", s.name, err)
}

func (s *LoggingSubscriber[T]) OnComplete() {
    fmt.Printf("[%s] Completed\n", s.name)
}

func main() {
    // Create a publisher with backpressure
    publisher := rx.RangePublisher(1, 10)
    
    // Subscribe with backpressure control
    subscriber := &LoggingSubscriber[int]{name: "Demo"}
    publisher.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```

### **3. Backpressure Strategies**

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    rx "github.com/droxer/RxGo"
)

type SlowSubscriber struct {
    received []int
}

func (s *SlowSubscriber) OnSubscribe(sub rx.Subscription) {
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
    // Create publisher with buffer strategy
    publisher := rx.RangePublisher(1, 100).
        WithBackpressureStrategy(rx.Buffer).
        WithBufferSize(10)
    
    subscriber := &SlowSubscriber{}
    publisher.Subscribe(context.Background(), subscriber)
    
    time.Sleep(time.Second)
}
```

### **4. Context Cancellation**

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    rx "github.com/droxer/RxGo"
)

type ContextSubscriber struct{}

func (s *ContextSubscriber) OnSubscribe(sub rx.Subscription) {
    sub.Request(rx.Unlimited)
}

func (s *ContextSubscriber) OnNext(value int) {
    fmt.Printf("Received: %d\n", value)
}

func (s *ContextSubscriber) OnError(err error) {
    fmt.Printf("Cancelled: %v\n", err)
}

func (s *ContextSubscriber) OnComplete() {
    fmt.Println("Completed")
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    publisher := rx.RangePublisher(1, 1000)
    subscriber := &ContextSubscriber{}
    
    publisher.Subscribe(ctx, subscriber)
    
    time.Sleep(time.Second)
}
```

### **5. Processor (Transforming Publisher)**

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    rx "github.com/droxer/RxGo"
)

// Simple processor that doubles values
type DoubleProcessor struct{}

func (p *DoubleProcessor) OnSubscribe(sub rx.Subscription) {
    sub.Request(rx.Unlimited)
}

func (p *DoubleProcessor) OnNext(value int) {
    p.OnNext(value * 2) // Forward doubled value
}

func (p *DoubleProcessor) OnError(err error) {
    p.OnError(err)
}

func (p *DoubleProcessor) OnComplete() {
    p.OnComplete()
}

func main() {
    source := rx.RangePublisher(1, 5)
    
    // Create a processor (both Subscriber and Publisher)
    processor := &DoubleProcessor{}
    
    // Chain them together
    source.Subscribe(context.Background(), processor)
    processor.Subscribe(context.Background(), &LoggingSubscriber[int]{name: "Output"})
    
    time.Sleep(100 * time.Millisecond)
}
```

## Features

### **✅ Core Features**
- **Type-safe generics** throughout the API
- **Context-based cancellation** support
- **Backpressure strategies**: Buffer, Drop, Latest, Error
- **Thread-safe** signal delivery
- **Memory efficient** with bounded buffers
- **Full interoperability** between old and new APIs

### **✅ Reactive Streams 1.0.3 Compliance**
- **Publisher[T]** - Type-safe data source with demand control
- **ReactiveSubscriber[T]** - Complete subscriber interface with lifecycle
- **Subscription** - Request/cancel control with backpressure
- **Processor[T,R]** - Transforming publisher
- **Backpressure support** - Full demand-based flow control
- **Non-blocking guarantees** - Async processing with context cancellation

### **✅ Backpressure Strategies**
- **Buffer** - Queue items when producer is faster than consumer
- **Drop** - Drop new items when buffer is full
- **Latest** - Keep only the latest item when buffer is full
- **Error** - Signal error when buffer overflows

### **✅ Utility Functions**
- `rx.Just[T](values...)` - Create from literal values
- `rx.RangePublisher(start, count)` - Integer sequence publisher
- `rx.FromSlice[T](slice)` - Create from slice
- `rx.NewReactivePublisher[T](fn)` - Custom publisher creation

## API Documentation

### **Publisher[T] Interface**
```go
type Publisher[T any] interface {
    Subscribe(ctx context.Context, s ReactiveSubscriber[T])
}
```

### **ReactiveSubscriber[T] Interface**
```go
type ReactiveSubscriber[T any] interface {
    OnSubscribe(s Subscription)
    OnNext(t T)
    OnError(err error)
    OnComplete()
}
```

### **Subscription Interface**
```go
type Subscription interface {
    Request(n int64)
    Cancel()
}
```

### **Constants**
```go
const (
    rx.Unlimited = 1 << 62 // Maximum request value
    
    // Backpressure strategies
    Buffer BackpressureStrategy = iota
    Drop
    Latest
    Error
)
```

## Installation

```bash
go get github.com/droxer/RxGo
```

## Running Examples

```bash
# Run basic examples
go run examples/basic.go

# Run tests
go test -v ./...
```

## Migration Guide

### **From Old API to Reactive Streams**

**Before:**
```go
obs := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
    sub.OnNext(42)
    sub.OnCompleted()
})
obs.Subscribe(ctx, subscriber)
```

**After:**
```go
publisher := rx.NewReactivePublisher(func(ctx context.Context, sub rx.ReactiveSubscriber[int]) {
    sub.OnSubscribe(rx.NewSubscription())
    sub.OnNext(42)
    sub.OnComplete()
})
publisher.Subscribe(ctx, reactiveSubscriber)
```

## Performance Considerations

- **Zero-allocation** signal delivery for common cases
- **Lock-free** data structures where possible
- **Context-aware** cancellation to prevent goroutine leaks
- **Bounded buffers** to prevent memory issues
- **Backpressure** to handle producer/consumer speed mismatches

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

Copyright 2015

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
