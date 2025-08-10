# Context Cancellation

This document demonstrates how to use Go context for graceful cancellation of reactive streams.

## 1. Timeout-Based Cancellation

Cancel streams after a specified timeout:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// ContextSubscriber demonstrates context cancellation
type ContextSubscriber struct{}

func (s *ContextSubscriber) OnSubscribe(sub rxgo.Subscription) {
    sub.Request(rxgo.Unlimited)
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
    
    publisher := rxgo.RangePublisher(1, 1000)
    subscriber := &ContextSubscriber{}
    
    publisher.Subscribe(ctx, subscriber)
    time.Sleep(time.Second)
}
```

## 2. Manual Cancellation

Cancel streams programmatically using context:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// ManualCancelSubscriber demonstrates manual cancellation
type ManualCancelSubscriber struct {
    cancel context.CancelFunc
}

func (s *ManualCancelSubscriber) OnSubscribe(sub rxgo.Subscription) {
    sub.Request(10) // Request initial batch
}

func (s *ManualCancelSubscriber) OnNext(value int) {
    fmt.Printf("Processing: %d\n", value)
    if value >= 5 {
        fmt.Println("Cancelling after receiving 5 items")
        s.cancel() // Cancel the context
    }
}

func (s *ManualCancelSubscriber) OnError(err error) {
    fmt.Printf("Stream cancelled: %v\n", err)
}

func (s *ManualCancelSubscriber) OnComplete() {
    fmt.Println("Stream completed normally")
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    publisher := rxgo.RangePublisher(1, 100)
    subscriber := &ManualCancelSubscriber{cancel: cancel}
    
    publisher.Subscribe(ctx, subscriber)
    time.Sleep(500 * time.Millisecond)
}
```

## 3. Context with Deadline

Use specific deadlines for stream processing:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// DeadlineSubscriber shows deadline-based cancellation
type DeadlineSubscriber struct {
    startTime time.Time
}

func (s *DeadlineSubscriber) OnSubscribe(sub rxgo.Subscription) {
    s.startTime = time.Now()
    sub.Request(rxgo.Unlimited)
}

func (s *DeadlineSubscriber) OnNext(value int) {
    elapsed := time.Since(s.startTime)
    fmt.Printf("[%v] Received: %d\n", elapsed, value)
    time.Sleep(100 * time.Millisecond) // Simulate processing
}

func (s *DeadlineSubscriber) OnError(err error) {
    fmt.Printf("Deadline exceeded: %v\n", err)
}

func (s *DeadlineSubscriber) OnComplete() {
    fmt.Printf("Completed processing in %v\n", time.Since(s.startTime))
}

func main() {
    // Set deadline 300ms from now
    deadline := time.Now().Add(300 * time.Millisecond)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    publisher := rxgo.RangePublisher(1, 20)
    subscriber := &DeadlineSubscriber{}
    
    publisher.Subscribe(ctx, subscriber)
    time.Sleep(500 * time.Millisecond)
}
```

## 4. Parent Context Cancellation

Handle cancellation from parent contexts:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// ParentContextSubscriber handles parent context cancellation
type ParentContextSubscriber struct {
    id string
}

func (s *ParentContextSubscriber) OnSubscribe(sub rxgo.Subscription) {
    fmt.Printf("[%s] Starting subscription\n", s.id)
    sub.Request(rxgo.Unlimited)
}

func (s *ParentContextSubscriber) OnNext(value int) {
    fmt.Printf("[%s] Processing: %d\n", s.id, value)
    time.Sleep(50 * time.Millisecond) // Simulate work
}

func (s *ParentContextSubscriber) OnError(err error) {
    fmt.Printf("[%s] Cancelled: %v\n", s.id, err)
}

func (s *ParentContextSubscriber) OnComplete() {
    fmt.Printf("[%s] Completed normally\n", s.id)
}

func processWithContext(ctx context.Context, id string) {
    publisher := rxgo.RangePublisher(1, 50)
    subscriber := &ParentContextSubscriber{id: id}
    publisher.Subscribe(ctx, subscriber)
}

func main() {
    // Create parent context with timeout
    parentCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel()
    
    // Create child contexts for different processing tasks
    child1, _ := context.WithCancel(parentCtx)
    child2, _ := context.WithCancel(parentCtx)
    
    go processWithContext(child1, "Task-1")
    go processWithContext(child2, "Task-2")
    
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Parent context cancelled - all child streams should stop")
}
```