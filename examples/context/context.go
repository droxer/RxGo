package main

import (
	"context"
	"fmt"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
)

type ContextAwareSubscriber struct {
	received int
}

func (s *ContextAwareSubscriber) Start() {
	fmt.Println("Context-aware subscriber started")
}

func (s *ContextAwareSubscriber) OnNext(value int) {
	s.received++
	fmt.Printf("Received: %d\n", value)
}

func (s *ContextAwareSubscriber) OnError(err error) {
	fmt.Printf("Context cancelled: %v\n", err)
}

func (s *ContextAwareSubscriber) OnCompleted() {
	fmt.Printf("Completed, total received: %d\n", s.received)
}

func main() {
	fmt.Println("=== Context Cancellation Example ===")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	subscriber := &ContextAwareSubscriber{}

	// Create observable that respects context cancellation
	observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
		defer sub.OnCompleted()
		for i := 1; i <= 100; i++ {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(i)
				time.Sleep(10 * time.Millisecond) // Small delay to show cancellation
			}
		}
	})

	observable.Subscribe(ctx, subscriber)
	time.Sleep(500 * time.Millisecond) // Wait for completion
	fmt.Println("Context cancellation example completed!")
}
