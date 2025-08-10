package main

import (
	"context"
	"fmt"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/rxgo"
)

// IntSubscriber demonstrates type-safe subscriber with generics
type IntSubscriber struct {
	name string
}

func (s *IntSubscriber) Start() {
	fmt.Printf("[%s] Starting subscription\n", s.name)
}

func (s *IntSubscriber) OnNext(value int) {
	fmt.Printf("[%s] Received: %d\n", s.name, value)
}

func (s *IntSubscriber) OnError(err error) {
	fmt.Printf("[%s] Error: %v\n", s.name, err)
}

func (s *IntSubscriber) OnCompleted() {
	fmt.Printf("[%s] Completed\n", s.name)
}

func main() {
	fmt.Println("=== RxGo Modern Example ===")

	// Example 1: Basic usage with Just
	fmt.Println("\n1. Using Just():")
	justObservable := rxgo.Just(1, 2, 3, 4, 5)
	justObservable.Subscribe(context.Background(), &IntSubscriber{name: "Just"})

	// Example 2: Range observable
	fmt.Println("\n2. Using Range():")
	rangeObservable := rxgo.Range(10, 5)
	rangeObservable.Subscribe(context.Background(), &IntSubscriber{name: "Range"})

	// Example 3: Create with custom logic
	fmt.Println("\n3. Using Create():")
	customObservable := rxgo.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
		for i := 0; i < 3; i++ {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(i * 10)
			}
		}
		sub.OnCompleted()
	})
	customObservable.Subscribe(context.Background(), &IntSubscriber{name: "Create"})

	// Example 4: With context (no scheduler due to event loop issues)
	fmt.Println("\n4. With context cancellation:")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	contextObservable := rxgo.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
		for i := 0; i < 5; i++ {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(i * 100)
			}
		}
		sub.OnCompleted()
	})
	contextObservable.Subscribe(ctx, &IntSubscriber{name: "Context"})

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== All examples completed ===")
}
