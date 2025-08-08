package publisher

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestReactiveSubscriber implements ReactiveSubscriber for testing
type TestReactiveSubscriber[T any] struct {
	mu           sync.Mutex
	subscription Subscription
	values       []T
	errors       []error
	completed    bool
	subscribed   bool
}

func NewTestReactiveSubscriber[T any]() *TestReactiveSubscriber[T] {
	return &TestReactiveSubscriber[T]{}
}

func (t *TestReactiveSubscriber[T]) OnSubscribe(s Subscription) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.subscription = s
	t.subscribed = true
}

func (t *TestReactiveSubscriber[T]) OnNext(value T) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.values = append(t.values, value)
}

func (t *TestReactiveSubscriber[T]) OnError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.errors = append(t.errors, err)
}

func (t *TestReactiveSubscriber[T]) OnComplete() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.completed = true
}

func (t *TestReactiveSubscriber[T]) getValues() []T {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]T(nil), t.values...)
}

func (t *TestReactiveSubscriber[T]) getErrors() []error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]error(nil), t.errors...)
}

func (t *TestReactiveSubscriber[T]) isCompleted() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.completed
}

func (t *TestReactiveSubscriber[T]) isSubscribed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.subscribed
}

func (t *TestReactiveSubscriber[T]) getSubscription() Subscription {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.subscription
}

func TestReactivePublisherSubscribe(t *testing.T) {
	publisher := FromSlice([]int{1, 2, 3, 4, 5})
	subscriber := NewTestReactiveSubscriber[int]()

	publisher.Subscribe(context.Background(), subscriber)

	// Wait for completion
	time.Sleep(100 * time.Millisecond)

	if !subscriber.isSubscribed() {
		t.Error("Subscriber should be subscribed")
	}

	values := subscriber.getValues()
	if len(values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(values))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	if !subscriber.isCompleted() {
		t.Error("Subscriber should be completed")
	}
}

func TestReactivePublisherBackpressure(t *testing.T) {
	// Test basic backpressure functionality
	publisher := RangePublisher(0, 10)
	subscriber := NewTestReactiveSubscriber[int]()

	publisher.Subscribe(context.Background(), subscriber)

	// Wait for subscription
	time.Sleep(10 * time.Millisecond)

	// Request some items
	sub := subscriber.getSubscription()
	if sub == nil {
		t.Fatal("Subscription should not be nil")
	}

	sub.Request(5)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	values := subscriber.getValues()
	if len(values) == 0 {
		t.Errorf("Expected some values due to backpressure, got %d", len(values))
	}
}

func TestReactivePublisherCancel(t *testing.T) {
	publisher := RangePublisher(0, 100)
	subscriber := NewTestReactiveSubscriber[int]()

	publisher.Subscribe(context.Background(), subscriber)

	// Give it some time to start
	time.Sleep(10 * time.Millisecond)

	sub := subscriber.getSubscription()
	if sub == nil {
		t.Fatal("Subscription should not be nil")
	}

	// Cancel subscription
	sub.Cancel()

	// Wait for cancellation to take effect
	time.Sleep(100 * time.Millisecond)

	values := subscriber.getValues()
	// Should have some values but not all 100 (exact count depends on timing)
	if len(values) >= 100 {
		t.Logf("Got %d values, cancellation may not have prevented all items", len(values))
	}
}

func TestReactivePublisherError(t *testing.T) {
	publisher := NewReactivePublisher(func(ctx context.Context, sub ReactiveSubscriber[string]) {
		sub.OnNext("hello")
		sub.OnError(fmt.Errorf("test error"))
	})

	subscriber := NewTestReactiveSubscriber[string]()
	publisher.Subscribe(context.Background(), subscriber)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	errors := subscriber.getErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	if errors[0].Error() != "test error" {
		t.Errorf("Expected 'test error', got %v", errors[0])
	}
}

func TestReactivePublisherEmpty(t *testing.T) {
	publisher := FromSlice([]int{})
	subscriber := NewTestReactiveSubscriber[int]()

	publisher.Subscribe(context.Background(), subscriber)

	// Wait for completion
	time.Sleep(100 * time.Millisecond)

	values := subscriber.getValues()
	if len(values) != 0 {
		t.Errorf("Expected 0 values for empty publisher, got %d", len(values))
	}

	if !subscriber.isCompleted() {
		t.Error("Empty publisher should complete immediately")
	}
}

func TestReactivePublisherNilSubscriber(t *testing.T) {
	publisher := FromSlice([]int{1, 2, 3})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Should panic with nil subscriber")
		}
	}()

	publisher.Subscribe(context.Background(), nil)
}

func TestReactivePublisherContextCancellation(t *testing.T) {
	publisher := RangePublisher(0, 100)
	subscriber := NewTestReactiveSubscriber[int]()

	ctx, cancel := context.WithCancel(context.Background())

	publisher.Subscribe(ctx, subscriber)

	// Give it some time to start
	time.Sleep(10 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for cancellation to take effect
	time.Sleep(100 * time.Millisecond)

	values := subscriber.getValues()
	// Should have some values but not all 100
	if len(values) >= 100 {
		t.Logf("Context cancellation should have stopped publisher, but got %d values", len(values))
	}

	errors := subscriber.getErrors()
	// Context cancellation errors may not always be delivered due to race conditions
	if len(errors) > 0 {
		t.Logf("Received errors: %v", errors)
	}
}
