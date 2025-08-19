package streams

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

// compliantTestSubscriber is a simple test subscriber for testing publishers
type compliantTestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
	mu        sync.Mutex
}

func newCompliantTestSubscriber[T any]() *compliantTestSubscriber[T] {
	return &compliantTestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

func (s *compliantTestSubscriber[T]) OnSubscribe(sub Subscription) {
	sub.Request(100) // Default unlimited
}

func (s *compliantTestSubscriber[T]) OnNext(value T) {
	s.mu.Lock()
	s.Received = append(s.Received, value)
	s.mu.Unlock()
}

func (s *compliantTestSubscriber[T]) OnError(err error) {
	s.mu.Lock()
	s.Errors = append(s.Errors, err)
	close(s.Done)
	s.mu.Unlock()
}

func (s *compliantTestSubscriber[T]) OnComplete() {
	close(s.Done)
}

func (s *compliantTestSubscriber[T]) wait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-s.Done:
	}
}

func (s *compliantTestSubscriber[T]) assertValues(t *testing.T, expected []T) {
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

func (s *compliantTestSubscriber[T]) assertNoError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", s.Errors)
	}
}

func (s *compliantManualTestSubscriber[T]) AssertError(t *testing.T) {
	s.compliantTestSubscriber.mu.Lock()
	defer s.compliantTestSubscriber.mu.Unlock()

	if len(s.Errors) == 0 {
		t.Error("Expected an error, but none occurred")
	}
}

// compliantManualTestSubscriber is a test subscriber with manual request control
type compliantManualTestSubscriber[T any] struct {
	*compliantTestSubscriber[T]
	subscription Subscription
}

func newCompliantManualTestSubscriber[T any]() *compliantManualTestSubscriber[T] {
	return &compliantManualTestSubscriber[T]{
		compliantTestSubscriber: newCompliantTestSubscriber[T](),
	}
}

func (s *compliantManualTestSubscriber[T]) OnSubscribe(sub Subscription) {
	s.subscription = sub
	// Don't auto-request
}

func (s *compliantManualTestSubscriber[T]) request(n int64) {
	if s.subscription != nil {
		s.subscription.Request(n)
	}
}

func TestCompliantFromSlicePublisher(t *testing.T) {
	t.Run("string slice", func(t *testing.T) {
		data := []string{"hello", "world", "reactive", "streams"}
		publisher := NewCompliantFromSlicePublisher(data)
		sub := newCompliantTestSubscriber[string]()
		ctx := context.Background()

		publisher.Subscribe(ctx, sub)
		sub.wait(ctx)

		sub.assertValues(t, data)
		sub.assertNoError(t)
	})

	t.Run("empty slice", func(t *testing.T) {
		publisher := NewCompliantFromSlicePublisher([]int{})
		sub := newCompliantTestSubscriber[int]()
		ctx := context.Background()

		publisher.Subscribe(ctx, sub)
		sub.wait(ctx)

		sub.assertValues(t, []int{})
		sub.assertNoError(t)
	})
}

func TestCompliantRangePublisher(t *testing.T) {
	t.Run("basic range", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5}
		publisher := NewCompliantRangePublisher(1, 5)
		sub := newCompliantTestSubscriber[int]()
		ctx := context.Background()

		publisher.Subscribe(ctx, sub)
		sub.wait(ctx)

		sub.assertValues(t, expected)
		sub.assertNoError(t)
	})

	t.Run("zero count", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(1, 0)
		sub := newCompliantTestSubscriber[int]()
		ctx := context.Background()

		publisher.Subscribe(ctx, sub)
		sub.wait(ctx)

		sub.assertValues(t, []int{})
		sub.assertNoError(t)
	})

	t.Run("demand control", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(1, 5)
		sub := newCompliantManualTestSubscriber[int]()
		ctx := context.Background()

		publisher.Subscribe(ctx, sub)
		sub.request(5)

		sub.wait(ctx)

		expected := []int{1, 2, 3, 4, 5}
		sub.assertValues(t, expected)
	})

	t.Run("negative request", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(1, 3)
		sub := newCompliantManualTestSubscriber[int]()
		ctx := context.Background()

		publisher.Subscribe(ctx, sub)
		sub.request(-1)

		sub.wait(ctx)

		sub.AssertError(t)
	})

	t.Run("multiple subscribers", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(10, 3)
		ctx := context.Background()

		sub1 := newCompliantTestSubscriber[int]()
		sub2 := newCompliantTestSubscriber[int]()

		publisher.Subscribe(ctx, sub1)
		publisher.Subscribe(ctx, sub2)
		sub1.wait(ctx)
		sub2.wait(ctx)

		expected := []int{10, 11, 12}
		sub1.assertValues(t, expected)
		sub2.assertValues(t, expected)
	})
}
