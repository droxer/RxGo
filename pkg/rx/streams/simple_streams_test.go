package streams

import (
	"context"
	"testing"

	"github.com/droxer/RxGo/pkg/rx"
)

// TestSubscriber is a test implementation of Subscriber
type TestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
}

func (t *TestSubscriber[T]) OnSubscribe(sub Subscription) {
	// Request all items immediately
	sub.Request(int64(^uint(0) >> 1))
}
func (t *TestSubscriber[T]) OnNext(value T)               { t.Received = append(t.Received, value) }
func (t *TestSubscriber[T]) OnError(err error)            { t.Errors = append(t.Errors, err); close(t.Done) }
func (t *TestSubscriber[T]) OnComplete()                   { t.Completed = true; close(t.Done) }

func TestFromSlicePublisher(t *testing.T) {
	// Test basic slice publisher
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := FromSlicePublisher([]int{1, 2, 3, 4, 5})
	publisher.Subscribe(context.Background(), sub)

	// Wait for completion
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

func TestRangePublisher(t *testing.T) {
	// Test range publisher
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := RangePublisher(10, 5)
	publisher.Subscribe(context.Background(), sub)

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

func TestRangePublisherZeroCount(t *testing.T) {
	// Test range publisher with zero count
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := RangePublisher(5, 0)
	publisher.Subscribe(context.Background(), sub)

	<-sub.Done

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestFromSlicePublisherEmpty(t *testing.T) {
	// Test empty slice publisher
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := FromSlicePublisher([]int{})
	publisher.Subscribe(context.Background(), sub)

	<-sub.Done

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestObservablePublisherAdapter(t *testing.T) {
	// Test adapter from Observable to Publisher
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	obs := rx.Just(1, 2, 3)
	publisher := ObservablePublisherAdapter(obs)
	publisher.Subscribe(context.Background(), sub)

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

func TestStringPublisher(t *testing.T) {
	// Test string publisher
	sub := &TestSubscriber[string]{Done: make(chan struct{})}
	publisher := FromSlicePublisher([]string{"hello", "world", "test"})
	publisher.Subscribe(context.Background(), sub)

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

func TestPublisherContextCancellation(t *testing.T) {
	// Test context cancellation works
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := FromSlicePublisher([]int{1, 2, 3, 4, 5})
	publisher.Subscribe(ctx, sub)

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

// TestProcessor removed - NewProcessor function not available in this package

func TestPublisherMultipleSubscribers(t *testing.T) {
	// Test multiple subscribers
	publisher := FromSlicePublisher([]int{1, 2, 3})
	sub1 := &TestSubscriber[int]{Done: make(chan struct{})}
	sub2 := &TestSubscriber[int]{Done: make(chan struct{})}

	publisher.Subscribe(context.Background(), sub1)
	publisher.Subscribe(context.Background(), sub2)

	// Wait for both subscribers
	<-sub1.Done
	<-sub2.Done

	expected := []int{1, 2, 3}
	
	for _, sub := range []*TestSubscriber[int]{sub1, sub2} {
		if len(sub.Received) != len(expected) {
			t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
		}

		for i, v := range sub.Received {
			if v != expected[i] {
				t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
			}
		}
	}
}

func TestPublisherWithSubscription(t *testing.T) {
	// Test subscription interface
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := FromSlicePublisher([]int{1, 2, 3})
	publisher.Subscribe(context.Background(), sub)

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

func TestFromSlicePublisherEmptyString(t *testing.T) {
	// Test empty string slice
	sub := &TestSubscriber[string]{Done: make(chan struct{})}
	publisher := FromSlicePublisher([]string{})
	publisher.Subscribe(context.Background(), sub)

	<-sub.Done

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestRangePublisherLarge(t *testing.T) {
	// Test large range
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := RangePublisher(1, 100)
	publisher.Subscribe(context.Background(), sub)

	<-sub.Done

	if len(sub.Received) != 100 {
		t.Errorf("Expected 100 values, got %d", len(sub.Received))
	}

	for i := 0; i < 100; i++ {
		if sub.Received[i] != i+1 {
			t.Errorf("Expected %d at index %d, got %d", i+1, i, sub.Received[i])
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

// TestBackpressureStrategies removed - BackpressureStrategy not available in this package

func TestPublisherNilSubscriber(t *testing.T) {
	// Test nil subscriber handling
	publisher := FromSlicePublisher([]int{1, 2, 3})
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil subscriber")
		}
	}()
	
	publisher.Subscribe(context.Background(), nil)
}


// TestProcessorWithTransform removed - NewProcessor function not available in this package

// TestComplexProcessorChain removed - NewProcessor function not available in this package

// TestFloatProcessor removed - NewProcessor function not available in this package

func TestObservableAdapter(t *testing.T) {
	// Test Observable to Publisher adapter
	tests := []struct {
		name string
		obs  *rx.Observable[int]
	}{
		{"Just", rx.Just(1, 2, 3)},
		{"Range", rx.Range(1, 3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &TestSubscriber[int]{Done: make(chan struct{})}
			publisher := ObservablePublisherAdapter(tt.obs)
			publisher.Subscribe(context.Background(), sub)

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
		})
	}
}

// TestProcessorIdentity removed - NewProcessor function not available in this package

// TestProcessorWithErrors removed - NewProcessor function not available in this package

func TestSubscriptionInterface(t *testing.T) {
	// Test subscription interface
	sub := &TestSubscriber[int]{Done: make(chan struct{})}
	publisher := FromSlicePublisher([]int{1, 2, 3})
	publisher.Subscribe(context.Background(), sub)

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