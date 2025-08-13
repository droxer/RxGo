package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/droxer/RxGo/pkg/rx/streams"
)

func main() {
	fmt.Println("RxGo Backpressure Examples")
	fmt.Println("==========================")

	// Run individual examples
	RunBufferStrategy()
	RunDropStrategy()
	RunLatestStrategy()
	RunErrorStrategy()
	RunHighVolumeBackpressure()

	fmt.Println("\nAll examples completed!")
}

// RunBufferStrategy demonstrates the Buffer backpressure strategy
func RunBufferStrategy() {
	fmt.Println("\n--- Buffer Strategy Example ---")

	publisher := streams.RangePublisherWithConfig(1, 10, streams.BackpressureConfig{
		Strategy:   streams.Buffer,
		BufferSize: 5,
	})

	subscriber := &simpleSubscriber{name: "BufferSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()
	fmt.Printf("Buffer strategy processed %d items\n", subscriber.processed)
}

// RunDropStrategy demonstrates the Drop backpressure strategy
func RunDropStrategy() {
	fmt.Println("\n--- Drop Strategy Example ---")

	publisher := streams.RangePublisherWithConfig(1, 10, streams.BackpressureConfig{
		Strategy:   streams.Drop,
		BufferSize: 3,
	})

	subscriber := &simpleSubscriber{name: "DropSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()
	fmt.Printf("Drop strategy processed %d items\n", subscriber.processed)
}

// RunLatestStrategy demonstrates the Latest backpressure strategy
func RunLatestStrategy() {
	fmt.Println("\n--- Latest Strategy Example ---")

	publisher := streams.RangePublisherWithConfig(1, 10, streams.BackpressureConfig{
		Strategy:   streams.Latest,
		BufferSize: 2,
	})

	subscriber := &simpleSubscriber{name: "LatestSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()
	fmt.Printf("Latest strategy processed %d items\n", subscriber.processed)
}

// RunErrorStrategy demonstrates the Error backpressure strategy
func RunErrorStrategy() {
	fmt.Println("\n--- Error Strategy Example ---")

	publisher := streams.RangePublisherWithConfig(1, 10, streams.BackpressureConfig{
		Strategy:   streams.Error,
		BufferSize: 2,
	})

	subscriber := &errorSubscriber{name: "ErrorSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()

	if subscriber.err != nil {
		fmt.Printf("Error strategy triggered: %v\n", subscriber.err)
	} else {
		fmt.Printf("Error strategy processed %d items\n", subscriber.processed)
	}
}

// RunHighVolumeBackpressure demonstrates all strategies with high-volume data
func RunHighVolumeBackpressure() {
	fmt.Println("\n--- High-Volume Backpressure Demo ---")

	strategies := []struct {
		name   string
		value  streams.BackpressureStrategy
		buffer int64
	}{
		{"Buffer", streams.Buffer, 50},
		{"Drop", streams.Drop, 10},
		{"Latest", streams.Latest, 5},
		{"Error", streams.Error, 20},
	}

	for _, strategy := range strategies {
		fmt.Printf("\nTesting %s Strategy...\n", strategy.name)

		publisher := streams.RangePublisherWithConfig(1, 50, streams.BackpressureConfig{
			Strategy:   strategy.value,
			BufferSize: strategy.buffer,
		})

		counter := &volumeCounter{name: strategy.name}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		publisher.Subscribe(ctx, counter)
		counter.wait()
		cancel()

		fmt.Printf("  Strategy: %s, Processed: %d, Status: %s\n",
			strategy.name, counter.processed, counter.status)
	}
}

// simpleSubscriber implements a basic subscriber for examples
type simpleSubscriber struct {
	name      string
	processed int
	wg        sync.WaitGroup
}

func (s *simpleSubscriber) OnSubscribe(sub streams.Subscription) {
	s.wg.Add(1)
	sub.Request(100)
}

func (s *simpleSubscriber) OnNext(value int) {
	fmt.Printf("[%s] Received: %d\n", s.name, value)
	s.processed++
	time.Sleep(100 * time.Millisecond)
}

func (s *simpleSubscriber) OnError(err error) {
	fmt.Printf("[%s] Error: %v\n", s.name, err)
	s.wg.Done()
}

func (s *simpleSubscriber) OnComplete() {
	fmt.Printf("[%s] Completed\n", s.name)
	s.wg.Done()
}

func (s *simpleSubscriber) wait() {
	s.wg.Wait()
}

// errorSubscriber implements a subscriber that handles errors
type errorSubscriber struct {
	name      string
	processed int
	err       error
	wg        sync.WaitGroup
}

func (s *errorSubscriber) OnSubscribe(sub streams.Subscription) {
	s.wg.Add(1)
	sub.Request(100)
}

func (s *errorSubscriber) OnNext(value int) {
	fmt.Printf("[%s] Received: %d\n", s.name, value)
	s.processed++
	time.Sleep(300 * time.Millisecond)
}

func (s *errorSubscriber) OnError(err error) {
	fmt.Printf("[%s] Error: %v\n", s.name, err)
	s.err = err
	s.wg.Done()
}

func (s *errorSubscriber) OnComplete() {
	fmt.Printf("[%s] Completed normally\n", s.name)
	s.wg.Done()
}

func (s *errorSubscriber) wait() {
	s.wg.Wait()
}

// volumeCounter counts processed items for high-volume demo
type volumeCounter struct {
	name      string
	processed int
	status    string
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func (c *volumeCounter) OnSubscribe(sub streams.Subscription) {
	c.wg.Add(1)
	sub.Request(100)
}

func (c *volumeCounter) OnNext(value int) {
	c.mu.Lock()
	c.processed++
	c.mu.Unlock()
	time.Sleep(20 * time.Millisecond)
}

func (c *volumeCounter) OnError(err error) {
	c.mu.Lock()
	c.status = fmt.Sprintf("Error: %v", err)
	c.mu.Unlock()
	c.wg.Done()
}

func (c *volumeCounter) OnComplete() {
	c.mu.Lock()
	c.status = "Completed"
	c.mu.Unlock()
	c.wg.Done()
}

func (c *volumeCounter) wait() {
	c.wg.Wait()
}
