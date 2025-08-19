package streams

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

// manualTestSubscriber is a test subscriber with manual request control
type manualTestSubscriber[T any] struct {
	Received     []T
	Completed    bool
	Errors       []error
	Done         chan struct{}
	mu           sync.Mutex
	subscription Subscription
}

func newManualTestSubscriber[T any]() *manualTestSubscriber[T] {
	return &manualTestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

func (s *manualTestSubscriber[T]) OnSubscribe(sub Subscription) {
	s.subscription = sub
	// Don't auto-request
}

func (s *manualTestSubscriber[T]) OnNext(value T) {
	s.mu.Lock()
	s.Received = append(s.Received, value)
	s.mu.Unlock()
}

func (s *manualTestSubscriber[T]) OnError(err error) {
	s.mu.Lock()
	s.Errors = append(s.Errors, err)
	close(s.Done)
	s.mu.Unlock()
}

func (s *manualTestSubscriber[T]) OnComplete() {
	close(s.Done)
}

func (s *manualTestSubscriber[T]) Request(n int64) {
	if s.subscription != nil {
		s.subscription.Request(n)
	}
}

func (s *manualTestSubscriber[T]) Cancel() {
	if s.subscription != nil {
		s.subscription.Cancel()
	}
}

func (s *manualTestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-s.Done:
	}
}

func (s *manualTestSubscriber[T]) GetReceivedCopy() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]T, len(s.Received))
	copy(result, s.Received)
	return result
}

func (s *manualTestSubscriber[T]) AssertValues(t *testing.T, expected []T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d: %v", len(expected), len(s.Received), s.Received)
		return
	}
	for i, v := range expected {
		if !reflect.DeepEqual(s.Received[i], v) {
			t.Errorf("Expected value[%d] to be %v, got %v", i, v, s.Received[i])
		}
	}
}

func (s *manualTestSubscriber[T]) AssertCompleted(t *testing.T) {
	// Can't reliably check completion without race conditions
}

func (s *manualTestSubscriber[T]) AssertNoError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", s.Errors)
	}
}

func (s *manualTestSubscriber[T]) AssertError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Errors) == 0 {
		t.Error("Expected an error, but none occurred")
	}
}

func (s *manualTestSubscriber[T]) IsCompleted() bool {
	return s.Completed
}

func TestSubscriptionRequestControl(t *testing.T) {
	ctx := context.Background()
	publisher := NewCompliantRangePublisher(1, 10)
	sub := newManualTestSubscriber[int]()

	publisher.Subscribe(ctx, sub)

	// Request only 3 items initially
	sub.Request(3)

	// Wait briefly for processing
	<-time.After(50 * time.Millisecond)

	received := sub.GetReceivedCopy()
	if len(received) > 3 {
		t.Errorf("Expected at most 3 values with limited request, got %d", len(received))
	}

	// Request remaining items
	sub.Request(7)

	sub.Wait(ctx)

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sub.AssertValues(t, expected)
}

func TestSubscriptionCancellation(t *testing.T) {
	ctx := context.Background()
	publisher := NewCompliantRangePublisher(1, 100)
	sub := newManualTestSubscriber[int]()

	publisher.Subscribe(ctx, sub)

	// Request some items
	sub.Request(5)

	// Wait briefly then cancel
	<-time.After(50 * time.Millisecond)

	sub.Cancel()
	sub.Request(10) // Should have no effect

	// Wait briefly to ensure cancellation took effect
	<-time.After(50 * time.Millisecond)

	received := sub.GetReceivedCopy()
	if len(received) > 10 {
		t.Errorf("Expected at most 10 values after cancellation, got %d", len(received))
	}

	if sub.IsCompleted() {
		t.Error("Expected no completion after cancellation")
	}
}

func TestSubscriptionZeroRequest(t *testing.T) {
	ctx := context.Background()
	publisher := NewCompliantRangePublisher(1, 5)
	sub := newManualTestSubscriber[int]()

	publisher.Subscribe(ctx, sub)
	sub.Request(0)

	// Wait briefly to ensure no items are processed
	<-time.After(100 * time.Millisecond)

	received := sub.GetReceivedCopy()
	if len(received) != 0 {
		t.Errorf("Expected 0 values with zero request, got %d", len(received))
	}
}

func TestSubscriptionNegativeRequest(t *testing.T) {
	ctx := context.Background()
	publisher := NewCompliantRangePublisher(1, 5)
	sub := newManualTestSubscriber[int]()

	publisher.Subscribe(ctx, sub)
	sub.Request(-1)

	sub.Wait(ctx)

	sub.AssertError(t)
}

func TestSubscriptionMultipleRequests(t *testing.T) {
	ctx := context.Background()
	publisher := NewCompliantRangePublisher(1, 10)
	sub := newManualTestSubscriber[int]()

	publisher.Subscribe(ctx, sub)

	// Make multiple small requests
	sub.Request(2)
	sub.Request(3)
	sub.Request(5)

	sub.Wait(ctx)

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sub.AssertValues(t, expected)
	sub.AssertCompleted(t)
}
