package rx

import (
	"context"
	"testing"
)

// TestSubscriber is a test implementation of Subscriber
type TestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
}

func (t *TestSubscriber[T]) Start()            {}
func (t *TestSubscriber[T]) OnNext(value T)    { t.Received = append(t.Received, value) }
func (t *TestSubscriber[T]) OnError(err error) { t.Errors = append(t.Errors, err); close(t.Done) }
func (t *TestSubscriber[T]) OnCompleted()      { t.Completed = true; close(t.Done) }

func TestJust(t *testing.T) {
	// Test basic Just creation
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	obs := Just(1, 2, 3, 4, 5)
	obs.Subscribe(context.Background(), sub)

	<-sub.Done

	expected := []int{1, 2, 3, 4, 5}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestRange(t *testing.T) {
	// Test Range creation
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	obs := Range(10, 5)
	obs.Subscribe(context.Background(), sub)

	<-sub.Done

	expected := []int{10, 11, 12, 13, 14}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestCreate(t *testing.T) {
	// Test Create with custom subscriber
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	obs := Create(func(ctx context.Context, subscriber Subscriber[int]) {
		for i := 1; i <= 3; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				subscriber.OnNext(i)
			}
		}
		subscriber.OnCompleted()
	})
	obs.Subscribe(context.Background(), sub)

	<-sub.Done

	expected := []int{1, 2, 3}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestSubscribeWithContext(t *testing.T) {
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	obs := Create(func(ctx context.Context, subscriber Subscriber[int]) {
		for i := 1; i <= 3; i++ {
			subscriber.OnNext(i)
		}
		subscriber.OnCompleted()
	})
	obs.Subscribe(ctx, sub)

	<-sub.Done

	expected := []int{1, 2, 3}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}
}

func TestSubscriberInterface(t *testing.T) {
	// Test custom subscriber implementation
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	obs := Just(1, 2, 3)
	obs.Subscribe(context.Background(), sub)

	<-sub.Done

	expected := []int{1, 2, 3}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestEmptyObservable(t *testing.T) {
	// Test empty observable
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	obs := Just[int]()
	obs.Subscribe(context.Background(), sub)

	<-sub.Done

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestStringObservable(t *testing.T) {
	// Test string observable
	sub := &TestSubscriber[string]{Done: make(chan struct{})}
	obs := Just("hello", "world", "test")
	obs.Subscribe(context.Background(), sub)

	<-sub.Done

	expected := []string{"hello", "world", "test"}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}