package rx

import (
	"context"
	"testing"

	"github.com/droxer/RxGo/internal/testutil"
)

func TestJust(t *testing.T) {
	// Test basic Just creation
	sub := testutil.NewTestSubscriber[int]()
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
	sub := testutil.NewTestSubscriber[int]()
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
	sub := testutil.NewTestSubscriber[int]()
	obs := Create(func(ctx context.Context, subscriber Subscriber[int]) {
		for i := 1; i <= 3; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				subscriber.OnNext(i)
			}
		}
		subscriber.OnComplete()
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

	sub := testutil.NewTestSubscriber[int]()
	obs := Create(func(ctx context.Context, subscriber Subscriber[int]) {
		for i := 1; i <= 3; i++ {
			subscriber.OnNext(i)
		}
		subscriber.OnComplete()
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
	sub := testutil.NewTestSubscriber[int]()
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
	sub := testutil.NewTestSubscriber[int]()
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
	sub := testutil.NewTestSubscriber[string]()
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
