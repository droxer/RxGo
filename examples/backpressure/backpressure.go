package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	rx "github.com/droxer/RxGo/pkg/rxgo"
)

type BackpressureSubscriber struct {
	name         string
	received     []int
	subscription rx.Subscription
	mu           sync.Mutex
	wg           *sync.WaitGroup
}

func (s *BackpressureSubscriber) OnSubscribe(sub rx.Subscription) {
	s.subscription = sub
	fmt.Printf("[%s] Subscribed, will request items slowly\n", s.name)

	// Request items slowly to demonstrate backpressure
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(200 * time.Millisecond)
			sub.Request(2)
			fmt.Printf("[%s] Requested 2 more items\n", s.name)
		}

		// Wait for processing and then complete
		time.Sleep(500 * time.Millisecond)
		s.wg.Done()
	}()
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

func (s *BackpressureSubscriber) OnComplete() {
	fmt.Printf("[%s] Completed, received %d items\n", s.name, len(s.received))
	s.wg.Done()
}

func main() {
	fmt.Println("=== Backpressure Example ===")
	var wg sync.WaitGroup
	wg.Add(1)

	subscriber := &BackpressureSubscriber{name: "Backpressure", wg: &wg}

	publisher := rx.NewReactivePublisher(func(ctx context.Context, sub rx.ReactiveSubscriber[int]) {
		defer sub.OnComplete()
		for i := 1; i <= 10; i++ {
			sub.OnNext(i)
		}
	})

	publisher.Subscribe(context.Background(), subscriber)
	wg.Wait()
	fmt.Println("Backpressure example completed!")
}
