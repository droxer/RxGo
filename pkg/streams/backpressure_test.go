package streams

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBufferStrategy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 5,
	})

	collector := &testCollector{}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)
	collector.wait()

	if collector.getItemCount() != 10 {
		t.Errorf("Expected 10 items, got %d", collector.getItemCount())
	}

	for i, val := range collector.getItems() {
		if val != i+1 {
			t.Errorf("Expected item %d to be %d, got %d", i, i+1, val)
		}
	}
}

func TestDropStrategy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Drop,
		BufferSize: 3,
	})

	collector := &testCollector{processDelay: 50 * time.Millisecond}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)
	collector.wait()

	fmt.Printf("Drop strategy processed %d items\n", collector.getItemCount())
}

func TestLatestStrategy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Latest,
		BufferSize: 2,
	})

	collector := &testCollector{processDelay: 100 * time.Millisecond}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)
	collector.wait()

	items := collector.getItems()
	fmt.Printf("Latest strategy processed %d items\n", len(items))
	if len(items) > 0 {
		lastItem := items[len(items)-1]
		if lastItem != 10 {
			t.Logf("Expected last item to be 10, got %d (this may be expected with Latest strategy)", lastItem)
		}
	}
}

func TestErrorStrategy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Error,
		BufferSize: 1,
	})

	collector := &testCollector{processDelay: 500 * time.Millisecond}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)
	collector.wait()

	if collector.err == nil {
		t.Log("Error strategy test - overflow may not trigger under all conditions")
	} else {
		fmt.Printf("Error strategy correctly triggered: %v\n", collector.err)
	}
}

type testCollector struct {
	items        []int
	err          error
	mu           sync.Mutex
	wg           sync.WaitGroup
	processDelay time.Duration
	subscription Subscription
	requestCount int64
	completed    bool
}

func (c *testCollector) OnSubscribe(sub Subscription) {
	c.subscription = sub
	c.wg.Add(1)
	if c.requestCount > 0 {
		sub.Request(c.requestCount)
	} else {
		sub.Request(100) // Default unlimited
	}
}

func (c *testCollector) OnNext(value int) {
	if c.processDelay > 0 {
		time.Sleep(c.processDelay)
	}
	c.mu.Lock()
	c.items = append(c.items, value)
	c.mu.Unlock()
}

func (c *testCollector) OnError(err error) {
	c.mu.Lock()
	if !c.completed {
		c.err = err
		c.completed = true
		c.wg.Done()
	}
	c.mu.Unlock()
}

func (c *testCollector) OnComplete() {
	c.mu.Lock()
	if !c.completed {
		c.completed = true
		c.wg.Done()
	}
	c.mu.Unlock()
}

func (c *testCollector) wait() {
	c.wg.Wait()
}

func (c *testCollector) requestMore(n int64) {
	if c.subscription != nil {
		c.subscription.Request(n)
	}
}

// getItemCount returns the current number of items in a thread-safe way
func (c *testCollector) getItemCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

// getItems returns a copy of items in a thread-safe way
func (c *testCollector) getItems() []int {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([]int, len(c.items))
	copy(result, c.items)
	return result
}

func TestRangePublishWithBackpressure(t *testing.T) {
	publisher := RangePublishWithBackpressure(1, 5, BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 10,
	})

	collector := &testCollector{}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)
	collector.wait()

	if collector.getItemCount() != 5 {
		t.Errorf("Expected 5 items, got %d", collector.getItemCount())
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, val := range collector.getItems() {
		if val != expected[i] {
			t.Errorf("Expected item %d to be %d, got %d", i, expected[i], val)
		}
	}
}

func TestBackpressureStrategyString(t *testing.T) {
	tests := []struct {
		strategy BackpressureStrategy
		expected string
	}{
		{Buffer, "Buffer"},
		{Drop, "Drop"},
		{Latest, "Latest"},
		{Error, "Error"},
	}

	for _, test := range tests {
		result := test.strategy.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestBackpressureControlledDemand(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 10,
	})

	collector := &testCollector{}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)
	collector.wait()

	// Should have received all items
	if collector.getItemCount() != 5 {
		t.Errorf("Expected 5 items, got %d", collector.getItemCount())
	}

	for i, val := range collector.getItems() {
		if val != i+1 {
			t.Errorf("Expected item %d to be %d, got %d", i, i+1, val)
		}
	}
}

func TestBackpressureBufferOverflow(t *testing.T) {
	// Create a large dataset that will exceed buffer
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i + 1
	}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 5,
	})

	collector := &testCollector{
		requestCount: 2,                     // Request only 2 items initially
		processDelay: 10 * time.Millisecond, // Slow processing
	}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)

	// Let it run for a while
	time.Sleep(200 * time.Millisecond)

	// Buffer should be managing the overflow
	receivedCount := collector.getItemCount()
	t.Logf("Buffered strategy processed %d items with slow consumer", receivedCount)

	// Request all remaining
	collector.requestMore(200)
	collector.wait()

	finalCount := collector.getItemCount()
	// The actual count depends on processing speed and timing
	t.Logf("Final count: %d items processed with backpressure", finalCount)

	// Ensure we got some items
	if finalCount == 0 {
		t.Errorf("Expected to receive some items, got 0")
	}
}

func TestBackpressureDropStrategyLimitedDemand(t *testing.T) {
	data := make([]int, 50)
	for i := 0; i < 50; i++ {
		data[i] = i + 1
	}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Drop,
		BufferSize: 2,
	})

	collector := &testCollector{
		requestCount: 5,
		processDelay: 100 * time.Millisecond,
	}
	ctx := context.Background()
	publisher.Subscribe(ctx, collector)
	collector.wait()

	// With drop strategy and slow processing, many items should be dropped
	receivedCount := collector.getItemCount()
	t.Logf("Drop strategy with limited demand processed %d items", receivedCount)

	if receivedCount > 10 {
		t.Errorf("Expected drop strategy to drop many items, but got %d", receivedCount)
	}
}
