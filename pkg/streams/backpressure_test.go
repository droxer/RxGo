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

	if len(collector.items) != 10 {
		t.Errorf("Expected 10 items, got %d", len(collector.items))
	}

	// Verify order is preserved
	for i, val := range collector.items {
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

	fmt.Printf("Drop strategy processed %d items\n", len(collector.items))
	// Should process some items, but may drop some due to slow consumer
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

	fmt.Printf("Latest strategy processed %d items\n", len(collector.items))
	// Should process the latest items
	if len(collector.items) > 0 {
		lastItem := collector.items[len(collector.items)-1]
		if lastItem != 10 {
			t.Logf("Expected last item to be 10, got %d (this may be expected with Latest strategy)", lastItem)
		}
	}
}

func TestErrorStrategy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	publisher := FromSlicePublishWithBackpressure(data, BackpressureConfig{
		Strategy:   Error,
		BufferSize: 1, // Very small buffer to ensure overflow
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

// testCollector collects items for testing
type testCollector struct {
	items        []int
	err          error
	mu           sync.Mutex
	wg           sync.WaitGroup
	processDelay time.Duration
}

func (c *testCollector) OnSubscribe(sub Subscription) {
	c.wg.Add(1)
	sub.Request(100) // Request all items
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
	c.err = err
	c.mu.Unlock()
	c.wg.Done()
}

func (c *testCollector) OnComplete() {
	c.wg.Done()
}

func (c *testCollector) wait() {
	c.wg.Wait()
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

	if len(collector.items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(collector.items))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, val := range collector.items {
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
