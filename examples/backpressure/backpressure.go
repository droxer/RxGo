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
