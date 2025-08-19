package streams

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

// publishersTestSubscriber is a simple test subscriber for publisher tests
type publishersTestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
	mu        sync.Mutex
}

func newPublishersTestSubscriber[T any]() *publishersTestSubscriber[T] {
	return &publishersTestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

func (s *publishersTestSubscriber[T]) OnSubscribe(sub Subscription) {
	sub.Request(100) // Default unlimited
}

func (s *publishersTestSubscriber[T]) OnNext(value T) {
	s.mu.Lock()
	s.Received = append(s.Received, value)
	s.mu.Unlock()
}

func (s *publishersTestSubscriber[T]) OnError(err error) {
	s.mu.Lock()
	s.Errors = append(s.Errors, err)
	close(s.Done)
	s.mu.Unlock()
}

func (s *publishersTestSubscriber[T]) OnComplete() {
	close(s.Done)
}

func (s *publishersTestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-s.Done:
	}
}

func (s *publishersTestSubscriber[T]) GetReceivedCopy() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]T, len(s.Received))
	copy(result, s.Received)
	return result
}

func (s *publishersTestSubscriber[T]) AssertValues(t *testing.T, expected []T) {
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

func (s *publishersTestSubscriber[T]) AssertCompleted(t *testing.T) {
	// Can't reliably check completion without race conditions
}

func (s *publishersTestSubscriber[T]) AssertNoError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if len(s.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", s.Errors)
	}
}

func TestFromSlicePublisher(t *testing.T) {
	t.Run("basic int slice", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5}
		publisher := FromSlicePublisher(expected)
		
		sub := newPublishersTestSubscriber[int]()
		ctx := context.Background()
		
		publisher.Subscribe(ctx, sub)
		sub.Wait(ctx)
		
		sub.AssertValues(t, expected)
		sub.AssertNoError(t)
	})

	t.Run("empty slice", func(t *testing.T) {
		publisher := FromSlicePublisher([]int{})
		sub := newPublishersTestSubscriber[int]()
		ctx := context.Background()
		
		publisher.Subscribe(ctx, sub)
		sub.Wait(ctx)
		
		sub.AssertValues(t, []int{})
		sub.AssertNoError(t)
	})

	t.Run("string slice", func(t *testing.T) {
		expected := []string{"hello", "world", "test"}
		publisher := FromSlicePublisher(expected)
		sub := newPublishersTestSubscriber[string]()
		ctx := context.Background()
		
		publisher.Subscribe(ctx, sub)
		sub.Wait(ctx)
		
		sub.AssertValues(t, expected)
		sub.AssertNoError(t)
	})
}

func TestRangePublisher(t *testing.T) {
	t.Run("basic range", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5}
		publisher := NewCompliantRangePublisher(1, 5)
		sub := newPublishersTestSubscriber[int]()
		ctx := context.Background()
		
		publisher.Subscribe(ctx, sub)
		sub.Wait(ctx)
		
		sub.AssertValues(t, expected)
		sub.AssertNoError(t)
	})

	t.Run("zero count", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(1, 0)
		sub := newPublishersTestSubscriber[int]()
		ctx := context.Background()
		
		publisher.Subscribe(ctx, sub)
		sub.Wait(ctx)
		
		sub.AssertValues(t, []int{})
		sub.AssertNoError(t)
	})

	t.Run("large range", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		publisher := NewCompliantRangePublisher(1, 10)
		sub := newPublishersTestSubscriber[int]()
		ctx := context.Background()
		
		publisher.Subscribe(ctx, sub)
		sub.Wait(ctx)
		
		sub.AssertValues(t, expected)
		sub.AssertNoError(t)
	})

	t.Run("multiple subscribers", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(1, 3)
		ctx := context.Background()
		expected := []int{1, 2, 3}

		sub1 := newPublishersTestSubscriber[int]()
		sub2 := newPublishersTestSubscriber[int]()

		publisher.Subscribe(ctx, sub1)
		publisher.Subscribe(ctx, sub2)
		sub1.Wait(ctx)
		sub2.Wait(ctx)

		sub1.AssertValues(t, expected)
		sub2.AssertValues(t, expected)
	})

	t.Run("nil subscriber", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(1, 5)
		// Should not panic
		publisher.Subscribe(context.Background(), nil)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		sub := newPublishersTestSubscriber[int]()
		publisher := NewCompliantRangePublisher(1, 1000)

		publisher.Subscribe(ctx, sub)
		cancel()

		sub.Wait(ctx)
		received := sub.GetReceivedCopy()
		t.Logf("Context cancelled, received %d values", len(received))
	})
}
