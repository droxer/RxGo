package main_testable

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/droxer/RxGo/pkg/rx/streams"
)

// BoundedBufferAdvanced demonstrates advanced bounded buffer patterns
func RunBoundedBufferAdvanced() {
	fmt.Println("=== Advanced Backpressure Patterns ===")
	// Example 1: Dynamic demand adjustment
	dynamicDemand()

	// Example 2: Buffer overflow protection
	bufferOverflow()

	// Example 3: Error handling with backpressure
	errorHandling()
}

// dynamicDemand shows how to adjust demand based on processing speed
func dynamicDemand() {
	fmt.Println("1. Dynamic Demand Adjustment")
	fmt.Println("   Adjusting demand based on processing time")

	publisher := streams.RangePublisher(1, 100)

	consumer := &dynamicDemandConsumer{
		baseDemand:      10,
		minDemand:       1,
		maxDemand:       50,
		processingTimes: make([]time.Duration, 0),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	publisher.Subscribe(ctx, consumer)

	<-consumer.done
	fmt.Printf("   ✅ Completed with dynamic demand adjustment\n\n")
}

// bufferOverflow demonstrates protection against buffer overflow
func bufferOverflow() {
	fmt.Println("2. Buffer Overflow Protection")
	fmt.Println("   Protecting against memory pressure")

	publisher := streams.RangePublisher(1, 1000)

	consumer := &bufferedConsumer{
		maxBuffer:    50,
		buffer:       make([]int, 0, 50),
		memoryLimit:  100, // Simulate memory limit
		currentUsage: 0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	publisher.Subscribe(ctx, consumer)

	<-consumer.done
	fmt.Printf("   ✅ Completed with buffer overflow protection\n\n")
}

// errorHandling shows error recovery with backpressure
func errorHandling() {
	fmt.Println("3. Error Handling with Backpressure")
	fmt.Println("   Handling errors while maintaining flow control")

	publisher := streams.RangePublisher(1, 50)

	consumer := &errorHandlingConsumer{
		maxRetries:   3,
		currentRetry: 0,
		errorCount:   0,
		processed:    make([]int, 0),
		retryDelays:  []time.Duration{100 * time.Millisecond, 500 * time.Millisecond, 1 * time.Second},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	publisher.Subscribe(ctx, consumer)

	<-consumer.done
	fmt.Printf("   ✅ Completed with %d errors handled\n\n", consumer.errorCount)
}

// dynamicDemandConsumer adjusts demand based on processing speed
type dynamicDemandConsumer struct {
	baseDemand      int64
	minDemand       int64
	maxDemand       int64
	processingTimes []time.Duration
	sub             streams.Subscription
	done            chan struct{}
	mu              sync.Mutex
}

func (d *dynamicDemandConsumer) OnSubscribe(sub streams.Subscription) {
	d.sub = sub
	d.done = make(chan struct{})
	sub.Request(d.baseDemand)
}

func (d *dynamicDemandConsumer) OnNext(value int) {
	start := time.Now()

	// Simulate variable processing time
	processingTime := time.Duration(value%3+1) * 50 * time.Millisecond
	time.Sleep(processingTime)

	elapsed := time.Since(start)

	d.mu.Lock()
	d.processingTimes = append(d.processingTimes, elapsed)

	// Adjust demand based on processing time
	newDemand := d.calculateDemand()
	if newDemand != d.baseDemand {
		fmt.Printf("   📊 Adjusting demand from %d to %d (processing time: %v)\n",
			d.baseDemand, newDemand, elapsed)
		d.baseDemand = newDemand
	}
	d.mu.Unlock()

	d.sub.Request(1)
}

func (d *dynamicDemandConsumer) calculateDemand() int64 {
	if len(d.processingTimes) < 3 {
		return d.baseDemand
	}

	// Calculate average processing time
	total := time.Duration(0)
	for _, t := range d.processingTimes[len(d.processingTimes)-3:] {
		total += t
	}
	average := total / 3

	// Adjust demand based on processing speed
	if average < 100*time.Millisecond {
		return min(d.baseDemand+5, d.maxDemand)
	} else if average > 300*time.Millisecond {
		return max(d.baseDemand-2, d.minDemand)
	}
	return d.baseDemand
}

func (d *dynamicDemandConsumer) OnError(err error) {
	fmt.Printf("   ❌ Dynamic demand consumer error: %v\n", err)
	close(d.done)
}

func (d *dynamicDemandConsumer) OnComplete() {
	fmt.Printf("   ✅ Dynamic demand consumer completed\n")
	close(d.done)
}

// bufferedConsumer implements buffer overflow protection
type bufferedConsumer struct {
	maxBuffer    int64
	buffer       []int
	memoryLimit  int64
	currentUsage int64
	sub          streams.Subscription
	done         chan struct{}
	mu           sync.Mutex
}

func (b *bufferedConsumer) OnSubscribe(sub streams.Subscription) {
	b.sub = sub
	b.done = make(chan struct{})
	sub.Request(b.maxBuffer)
}

func (b *bufferedConsumer) OnNext(value int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Simulate memory usage
	b.currentUsage += int64(len(fmt.Sprint(value)))

	// Check memory pressure
	if b.currentUsage >= b.memoryLimit {
		fmt.Printf("   🚨 Memory pressure detected (usage: %d/%d)\n", b.currentUsage, b.memoryLimit)
		b.flushBuffer()
		b.currentUsage = 0
	}

	b.buffer = append(b.buffer, value)

	if len(b.buffer) >= int(b.maxBuffer) {
		b.flushBuffer()
	}
}

func (b *bufferedConsumer) flushBuffer() {
	fmt.Printf("   📦 Flushing %d items from buffer\n", len(b.buffer))
	// Process the buffer
	b.buffer = b.buffer[:0]

	// Request more items
	b.sub.Request(b.maxBuffer)
}

func (b *bufferedConsumer) OnError(err error) {
	fmt.Printf("   ❌ Buffered consumer error: %v\n", err)
	close(b.done)
}

func (b *bufferedConsumer) OnComplete() {
	if len(b.buffer) > 0 {
		b.flushBuffer()
	}
	fmt.Printf("   ✅ Buffered consumer completed\n")
	close(b.done)
}

// errorHandlingConsumer demonstrates error recovery with backpressure
type errorHandlingConsumer struct {
	maxRetries   int
	currentRetry int
	errorCount   int
	processed    []int
	retryDelays  []time.Duration
	sub          streams.Subscription
	done         chan struct{}
	mu           sync.Mutex
}

func (e *errorHandlingConsumer) OnSubscribe(sub streams.Subscription) {
	e.sub = sub
	e.done = make(chan struct{})
	sub.Request(1) // Start with one item
}

func (e *errorHandlingConsumer) OnNext(value int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Simulate occasional errors
	if value%7 == 0 && e.currentRetry < e.maxRetries {
		e.errorCount++
		e.currentRetry++

		delay := e.retryDelays[minInt(e.currentRetry-1, len(e.retryDelays)-1)]
		fmt.Printf("   ⚠️  Error processing %d, retrying in %v (attempt %d/%d)\n",
			value, delay, e.currentRetry, e.maxRetries)

		time.Sleep(delay)

		// Re-request the same item
		e.sub.Request(1)
		return
	}

	// Reset retry count on success
	e.currentRetry = 0
	e.processed = append(e.processed, value)
	fmt.Printf("   ✅ Processed %d successfully\n", value)

	// Request next item
	e.sub.Request(1)
}

func (e *errorHandlingConsumer) OnError(err error) {
	fmt.Printf("   ❌ Error handling consumer error: %v\n", err)
	close(e.done)
}

func (e *errorHandlingConsumer) OnComplete() {
	fmt.Printf("   ✅ Error handling consumer completed with %d items\n", len(e.processed))
	close(e.done)
}

// Helper functions
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
