package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/droxer/RxGo/internal/publisher"
)

type ContextAwareSubscriber struct {
	received int
	wg       *sync.WaitGroup
}

func (s *ContextAwareSubscriber) OnSubscribe(sub publisher.Subscription) {
	fmt.Println("Context-aware subscriber subscribed")
	sub.Request(100) // Request many items
}

func (s *ContextAwareSubscriber) OnNext(value int) {
	s.received++
	fmt.Printf("Received: %d\n", value)
}

func (s *ContextAwareSubscriber) OnError(err error) {
	fmt.Printf("Context cancelled: %v\n", err)
	s.wg.Done()
}

func (s *ContextAwareSubscriber) OnComplete() {
	fmt.Printf("Completed, total received: %d\n", s.received)
	s.wg.Done()
}

func main() {
	fmt.Println("=== Context Cancellation Example ===")
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	subscriber := &ContextAwareSubscriber{wg: &wg}
	publisher := publisher.NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[int]) {
		defer sub.OnComplete()
		for i := 1; i <= 100; i++ {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(i)
			}
		}
	})

	publisher.Subscribe(ctx, subscriber)
	wg.Wait()
	fmt.Println("Context cancellation example completed!")
}
