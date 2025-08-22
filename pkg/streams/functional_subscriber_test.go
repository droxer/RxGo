package streams

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewSubscriber(t *testing.T) {
	var received []int
	var completed bool
	var receivedError error
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

	publisher := NewCompliantRangePublisher(1, 3)
	err := publisher.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Give some time for async processing
	// Note: This is a bit fragile, but needed for testing
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	expected := []int{1, 2, 3}
	if len(received) != len(expected) {
		t.Errorf("Expected %d values, got %d: %v", len(expected), len(received), received)
	}

	for i, v := range received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !completed {
		t.Error("Expected completion")
	}

	if receivedError != nil {
		t.Errorf("Expected no error, got %v", receivedError)
	}
}

func TestNewSubscriberWithError(t *testing.T) {
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

	testError := errors.New("test error")
	publisher := NewPublisher(func(ctx context.Context, sub Subscriber[int]) {
		sub.OnNext(1)
		sub.OnNext(2)
		sub.OnError(testError)
	})

	err := publisher.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Give some time for processing
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

func TestNewSubscriberNilHandlers(t *testing.T) {
	// Test that nil handlers don't cause panics
	subscriber := NewSubscriber[int](nil, nil, nil)

	publisher := NewCompliantRangePublisher(1, 2)
	err := publisher.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Should not panic - this test passes if no panic occurs
}

func TestFunctionalSubscriberInterface(t *testing.T) {
	subscriber := NewSubscriber(
		func(value int) {},
		func(err error) {},
		func() {},
	)

	// Verify it implements the Subscriber interface
	var _ Subscriber[int] = subscriber
}

func TestNewSubscriberAutoRequest(t *testing.T) {
	var received []int
	var mu sync.Mutex

	subscriber := NewSubscriber(
		func(value int) {
			mu.Lock()
			received = append(received, value)
			mu.Unlock()
		},
		func(err error) {},
		func() {},
	)

	// The functional subscriber should auto-request unlimited items
	publisher := NewCompliantRangePublisher(1, 5)
	err := publisher.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Give time for processing
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(received) != 5 {
		t.Errorf("Expected all 5 values to be received, got %d", len(received))
	}
}
