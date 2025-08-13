package main_testable

import (
	"context"
	"fmt"
	"time"

	"github.com/droxer/RxGo/pkg/rx/streams"
)

// BoundedBufferBasic demonstrates basic bounded buffer with controlled consumption
func RunBoundedBufferBasic() {
	fmt.Println("=== Bounded Buffer Basic Examples ===")
	// Example 1: Basic backpressure with controlled demand
	basicBackpressure()

	// Example 2: Batch processing with demand-based flow
	batchProcessing()

	// Example 3: Rate limiting with backpressure
	rateLimiting()
}

// basicBackpressure shows how to control flow with demand
func basicBackpressure() {
	fmt.Println("1. Basic Backpressure Example")
	fmt.Println("   Processing 100 items with demand of 5 items at a time")

	publisher := streams.RangePublisher(1, 100)

	consumer := &controlledConsumer{
		name:      "Consumer1",
		demand:    5,
		processed: make([]int, 0),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publisher.Subscribe(ctx, consumer)

	// Wait for completion
	<-consumer.done

	fmt.Printf("   ✅ Processed %d items in batches of %d\n", len(consumer.processed), consumer.demand)
	fmt.Printf("   First batch: %v\n\n", consumer.processed[:5])
}

// batchProcessing demonstrates demand-based batch processing
func batchProcessing() {
	fmt.Println("2. Batch Processing Example")
	fmt.Println("   Processing items in batches of 10")

	publisher := streams.RangePublisher(1, 50)

	consumer := &batchConsumer{
		batchSize: 10,
		buffer:    make([]int, 0, 10),
		processed: 0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publisher.Subscribe(ctx, consumer)

	// Wait for completion
	<-consumer.done

	fmt.Printf("   ✅ Processed %d items in %d batches\n\n", 50, 50/10)
}

// rateLimiting shows how to implement rate limiting with backpressure
func rateLimiting() {
	fmt.Println("3. Rate Limiting Example")
	fmt.Println("   Processing 1 item per second")

	publisher := streams.RangePublisher(1, 10)

	consumer := &rateLimitedConsumer{
		rate:      time.Second,
		processed: make([]int, 0),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	publisher.Subscribe(ctx, consumer)

	// Wait for completion
	<-consumer.done

	fmt.Printf("   ✅ Processed %d items at 1 per second\n\n", len(consumer.processed))
}

// controlledConsumer implements backpressure with controlled demand
type controlledConsumer struct {
	name      string
	demand    int64
	processed []int
	sub       streams.Subscription
	done      chan struct{}
}

func (c *controlledConsumer) OnSubscribe(sub streams.Subscription) {
	c.sub = sub
	c.done = make(chan struct{})
	// Request initial batch
	sub.Request(c.demand)
}

func (c *controlledConsumer) OnNext(value int) {
	c.processed = append(c.processed, value)
	fmt.Printf("   📥 %s received: %d (total: %d)\n", c.name, value, len(c.processed))

	// Request next batch when current one is processed
	if len(c.processed)%int(c.demand) == 0 {
		c.sub.Request(c.demand)
	}
}

func (c *controlledConsumer) OnError(err error) {
	fmt.Printf("   ❌ %s error: %v\n", c.name, err)
	close(c.done)
}

func (c *controlledConsumer) OnComplete() {
	fmt.Printf("   ✅ %s completed\n", c.name)
	close(c.done)
}

// batchConsumer implements batch processing with backpressure
type batchConsumer struct {
	batchSize int64
	buffer    []int
	processed int
	sub       streams.Subscription
	done      chan struct{}
}

func (b *batchConsumer) OnSubscribe(sub streams.Subscription) {
	b.sub = sub
	b.done = make(chan struct{})
	sub.Request(b.batchSize)
}

func (b *batchConsumer) OnNext(value int) {
	b.buffer = append(b.buffer, value)

	if len(b.buffer) >= int(b.batchSize) {
		b.processBatch()
		b.sub.Request(b.batchSize)
	}
}

func (b *batchConsumer) processBatch() {
	fmt.Printf("   📦 Processing batch: %v\n", b.buffer)
	b.processed += len(b.buffer)
	b.buffer = b.buffer[:0] // Clear buffer
}

func (b *batchConsumer) OnError(err error) {
	fmt.Printf("   ❌ Batch consumer error: %v\n", err)
	close(b.done)
}

func (b *batchConsumer) OnComplete() {
	if len(b.buffer) > 0 {
		b.processBatch()
	}
	close(b.done)
}

// rateLimitedConsumer implements rate limiting with backpressure
type rateLimitedConsumer struct {
	rate      time.Duration
	processed []int
	lastTime  time.Time
	sub       streams.Subscription
	done      chan struct{}
}

func (r *rateLimitedConsumer) OnSubscribe(sub streams.Subscription) {
	r.sub = sub
	r.done = make(chan struct{})
	r.lastTime = time.Now()
	sub.Request(1) // Request one item initially
}

func (r *rateLimitedConsumer) OnNext(value int) {
	// Rate limiting
	sinceLast := time.Since(r.lastTime)
	if sinceLast < r.rate {
		time.Sleep(r.rate - sinceLast)
	}

	r.processed = append(r.processed, value)
	r.lastTime = time.Now()
	fmt.Printf("   ⏱️  Processed: %d at %v\n", value, r.lastTime.Format("15:04:05"))

	// Request next item
	r.sub.Request(1)
}

func (r *rateLimitedConsumer) OnError(err error) {
	fmt.Printf("   ❌ Rate limited consumer error: %v\n", err)
	close(r.done)
}

func (r *rateLimitedConsumer) OnComplete() {
	fmt.Printf("   ✅ Rate limited consumer completed\n")
	close(r.done)
}
