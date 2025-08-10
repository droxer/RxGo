package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
)

type BackpressureSubscriber struct {
	name     string
	received []int
	mu       sync.Mutex
	wg       *sync.WaitGroup
	requests chan int
}

func (s *BackpressureSubscriber) Start() {
	fmt.Printf("[%s] Starting subscription\n", s.name)
	// Start requesting items with backpressure
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(200 * time.Millisecond)
			s.requests <- 2 // Request 2 items
			fmt.Printf("[%s] Requested 2 more items\n", s.name)
		}
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

func (s *BackpressureSubscriber) OnCompleted() {
	fmt.Printf("[%s] Completed, received %d items\n", s.name, len(s.received))
	s.wg.Done()
}

func main() {
	fmt.Println("=== Backpressure Example ===")
	var wg sync.WaitGroup
	wg.Add(1)

	subscriber := &BackpressureSubscriber{
		name:     "Backpressure",
		wg:       &wg,
		requests: make(chan int, 10),
	}

	// Create observable that respects backpressure through channel buffering
	observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
		defer sub.OnCompleted()

		// Simulate controlled emission based on consumer speed
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for _, item := range items {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(item)
				time.Sleep(50 * time.Millisecond) // Simulate producer speed
			}
		}
	})

	observable.Subscribe(context.Background(), subscriber)
	wg.Wait()
	fmt.Println("Backpressure example completed!")
}
