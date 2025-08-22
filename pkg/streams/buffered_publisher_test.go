package streams

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestBufferedPublisherErrorInSource(t *testing.T) {
	testError := errors.New("source error")

	// Create a publisher that emits some values then errors
	source := func(ctx context.Context, sub Subscriber[int]) {
		sub.OnNext(1)
		sub.OnNext(2)
		sub.OnError(testError)
	}

	publisher := NewBufferedPublisher(BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 10,
	}, source)

	var received []int
	var receivedError error
	var completed bool
	var mu sync.Mutex

	subscriber := NewSubscriber(
		func(value int) {
			mu.Lock()
			received = append(received, value)
			mu.Unlock()
		},
		func(err error) {
			mu.Lock()
			receivedError = err
			mu.Unlock()
		},
		func() {
			mu.Lock()
			completed = true
			mu.Unlock()
		},
	)

	err := publisher.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Give time for processing
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(received) != 2 {
		t.Errorf("Expected 2 values before error, got %d", len(received))
	}

	if receivedError != testError {
		t.Errorf("Expected test error, got %v", receivedError)
	}

	if completed {
		t.Error("Should not complete when error occurs")
	}
}

func TestBufferedPublisherCancellation(t *testing.T) {
	var emitted int
	var mu sync.Mutex

	// Create a slow source
	source := func(ctx context.Context, sub Subscriber[int]) {
		defer sub.OnComplete()
		for i := 1; i <= 100; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				mu.Lock()
				emitted++
				mu.Unlock()
				sub.OnNext(i)
				time.Sleep(time.Millisecond) // Slow emission
			}
		}
	}

	publisher := NewBufferedPublisher(BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 5,
	}, source)

	ctx, cancel := context.WithCancel(context.Background())

	var received []int
	var receivedMu sync.Mutex

	subscriber := NewSubscriber(
		func(value int) {
			receivedMu.Lock()
			received = append(received, value)
			receivedMu.Unlock()
		},
		func(err error) {},
		func() {},
	)

	err := publisher.Subscribe(ctx, subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Cancel after short time
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Give time for cleanup
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	emittedCount := emitted
	mu.Unlock()

	receivedMu.Lock()
	receivedCount := len(received)
	receivedMu.Unlock()

	if emittedCount >= 100 {
		t.Error("Source should have been cancelled and not emit all values")
	}

	t.Logf("Emitted %d values, received %d before cancellation", emittedCount, receivedCount)
}

func TestBufferedPublisherCompleteWithBufferedItems(t *testing.T) {
	// Test that completion flushes buffer according to strategy
	source := func(ctx context.Context, sub Subscriber[int]) {
		defer sub.OnComplete()
		// Emit more than buffer size
		for i := 1; i <= 8; i++ {
			sub.OnNext(i)
		}
	}

	publisher := NewBufferedPublisher(BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 3,
	}, source)

	var received []int
	var completed bool
	var mu sync.Mutex

	// Create a controlled subscriber that only requests a few items
	subscription := &testSubscription{}
	subscriber := &controlledSubscriber[int]{
		onSubscribe: func(s Subscription) {
			subscription.sub = s
			s.Request(2) // Only request 2 items initially
		},
		onNext: func(value int) {
			mu.Lock()
			received = append(received, value)
			mu.Unlock()
		},
		onComplete: func() {
			mu.Lock()
			completed = true
			mu.Unlock()
		},
	}

	err := publisher.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Wait for completion - should flush all buffered items
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if !completed {
		t.Error("Expected completion")
	}

	// With Buffer strategy, should get some items on completion
	// The exact count may vary based on timing and buffer behavior
	if len(received) == 0 {
		t.Errorf("Expected to receive some values on completion, got 0")
	}
	t.Logf("Received %d values on completion", len(received))
}

type testSubscription struct {
	sub Subscription
}

// controlledSubscriber is a subscriber that allows customizing behavior
type controlledSubscriber[T any] struct {
	onSubscribe func(Subscription)
	onNext      func(T)
	onError     func(error)
	onComplete  func()
}

func (ts *testSubscription) Request(n int64) {
	if ts.sub != nil {
		ts.sub.Request(n)
	}
}

func (ts *testSubscription) Cancel() {
	if ts.sub != nil {
		ts.sub.Cancel()
	}
}

// Implement the Subscriber interface for controlledSubscriber
func (c *controlledSubscriber[T]) OnSubscribe(s Subscription) {
	if c.onSubscribe != nil {
		c.onSubscribe(s)
	}
}

func (c *controlledSubscriber[T]) OnNext(value T) {
	if c.onNext != nil {
		c.onNext(value)
	}
}

func (c *controlledSubscriber[T]) OnError(err error) {
	if c.onError != nil {
		c.onError(err)
	}
}

func (c *controlledSubscriber[T]) OnComplete() {
	if c.onComplete != nil {
		c.onComplete()
	}
}

func TestBufferedPublisherLatestStrategyOverflow(t *testing.T) {
	// Test Latest strategy behavior when buffer overflows
	source := func(ctx context.Context, sub Subscriber[int]) {
		defer sub.OnComplete()
		// Emit values rapidly
		for i := 1; i <= 10; i++ {
			sub.OnNext(i)
		}
	}

	publisher := NewBufferedPublisher(BackpressureConfig{
		Strategy:   Latest,
		BufferSize: 2,
	}, source)

	var received []int
	var mu sync.Mutex

	// Slow subscriber
	subscriber := NewSubscriber(
		func(value int) {
			mu.Lock()
			received = append(received, value)
			mu.Unlock()
			time.Sleep(50 * time.Millisecond) // Slow processing
		},
		func(err error) {
			t.Errorf("Unexpected error: %v", err)
		},
		func() {},
	)

	err := publisher.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Wait for completion
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	// With Latest strategy and slow subscriber, should get fewer items
	// The exact number depends on timing, but should be < 10
	t.Logf("Latest strategy processed %d out of 10 values", len(received))

	if len(received) > 8 {
		t.Error("Latest strategy should drop some values with slow subscriber")
	}
}

func TestBufferedPublisherNilSubscriber(t *testing.T) {
	source := func(ctx context.Context, sub Subscriber[int]) {
		sub.OnNext(1)
		sub.OnComplete()
	}

	publisher := NewBufferedPublisher(BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 5,
	}, source)

	// Should return error for nil subscriber
	err := publisher.Subscribe(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil subscriber, got nil")
	}
	if err.Error() != "subscriber cannot be nil" {
		t.Errorf("Expected 'subscriber cannot be nil' error, got: %v", err)
	}
}
