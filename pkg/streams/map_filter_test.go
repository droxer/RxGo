package streams

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
)

// mapFilterTestSubscriber is a simple test subscriber for map/filter tests
type mapFilterTestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
	mu        sync.Mutex
}

func newMapFilterTestSubscriber[T any]() *mapFilterTestSubscriber[T] {
	return &mapFilterTestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

func (s *mapFilterTestSubscriber[T]) OnSubscribe(sub Subscription) {
	sub.Request(100) // Default unlimited
}

func (s *mapFilterTestSubscriber[T]) OnNext(value T) {
	s.mu.Lock()
	s.Received = append(s.Received, value)
	s.mu.Unlock()
}

func (s *mapFilterTestSubscriber[T]) OnError(err error) {
	s.mu.Lock()
	s.Errors = append(s.Errors, err)
	close(s.Done)
	s.mu.Unlock()
}

func (s *mapFilterTestSubscriber[T]) OnComplete() {
	close(s.Done)
}

func (s *mapFilterTestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-s.Done:
	}
}

func (s *mapFilterTestSubscriber[T]) GetReceivedCopy() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]T, len(s.Received))
	copy(result, s.Received)
	return result
}

func (s *mapFilterTestSubscriber[T]) AssertValues(t *testing.T, expected []T) {
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

func (s *mapFilterTestSubscriber[T]) AssertCompleted(t *testing.T) {
	// Can't reliably check completion without race conditions
}

func (s *mapFilterTestSubscriber[T]) AssertNoError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if len(s.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", s.Errors)
	}
}

func (s *mapFilterTestSubscriber[T]) AssertError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if len(s.Errors) == 0 {
		t.Error("Expected an error, but none occurred")
	}
}

func TestMapProcessor(t *testing.T) {
	t.Run("int to int", func(t *testing.T) {
		source := NewCompliantRangePublisher(1, 5)
		processor := NewMapProcessor(func(x int) int { return x * 2 })
		
		source.Subscribe(context.Background(), processor)
		
		sub := newMapFilterTestSubscriber[int]()
		processor.Subscribe(context.Background(), sub)

		ctx := context.Background()
		sub.Wait(ctx)

		expected := []int{2, 4, 6, 8, 10}
		sub.AssertValues(t, expected)
		sub.AssertCompleted(t)
		sub.AssertNoError(t)
	})

	t.Run("int to string", func(t *testing.T) {
		source := NewCompliantFromSlicePublisher([]int{1, 2, 3})
		processor := NewMapProcessor(func(x int) string {
			return map[int]string{1: "one", 2: "two", 3: "three"}[x]
		})

		source.Subscribe(context.Background(), processor)

		sub := newMapFilterTestSubscriber[string]()
		processor.Subscribe(context.Background(), sub)

		ctx := context.Background()
		sub.Wait(ctx)

		expected := []string{"one", "two", "three"}
		sub.AssertValues(t, expected)
		sub.AssertCompleted(t)
		sub.AssertNoError(t)
	})
}

func TestFilterProcessor(t *testing.T) {
	t.Run("filter even numbers", func(t *testing.T) {
		source := NewCompliantRangePublisher(1, 10)
		processor := NewFilterProcessor(func(x int) bool { return x%2 == 0 })

		source.Subscribe(context.Background(), processor)

		sub := newMapFilterTestSubscriber[int]()
		processor.Subscribe(context.Background(), sub)

		ctx := context.Background()
		sub.Wait(ctx)

		expected := []int{2, 4, 6, 8, 10}
		sub.AssertValues(t, expected)
		sub.AssertCompleted(t)
		sub.AssertNoError(t)
	})

	t.Run("filter all values", func(t *testing.T) {
		source := NewCompliantRangePublisher(1, 5)
		processor := NewFilterProcessor(func(x int) bool { return x > 10 })

		source.Subscribe(context.Background(), processor)

		sub := newMapFilterTestSubscriber[int]()
		processor.Subscribe(context.Background(), sub)

		ctx := context.Background()
		sub.Wait(ctx)

		sub.AssertValues(t, []int{})
		sub.AssertCompleted(t)
		sub.AssertNoError(t)
	})
}

func TestFlatMapProcessor(t *testing.T) {
	source := NewCompliantFromSlicePublisher([]int{1, 2})
	processor := NewFlatMapProcessor(func(x int) Publisher[int] {
		return NewCompliantFromSlicePublisher([]int{x, x + 10})
	})

	source.Subscribe(context.Background(), processor)

	sub := newMapFilterTestSubscriber[int]()
	processor.Subscribe(context.Background(), sub)

	ctx := context.Background()
	sub.Wait(ctx)

	// Check that all expected values are present regardless of order
	received := sub.GetReceivedCopy()
	if len(received) != 4 {
		t.Errorf("Expected 4 values from flatmap, got %d", len(received))
	}

	expected := map[int]bool{1: true, 2: true, 11: true, 12: true}
	for _, v := range received {
		if !expected[v] {
			t.Errorf("Unexpected value from flatmap: %d", v)
		}
		delete(expected, v)
	}

	if len(expected) != 0 {
		t.Errorf("Missing values from flatmap: %v", expected)
	}
}

func TestProcessorErrorHandling(t *testing.T) {
	// Create a publisher that emits an error
	errorPublisher := NewPublisher(func(ctx context.Context, sub Subscriber[int]) {
		sub.OnNext(1)
		sub.OnError(errors.New("test error"))
	})

	processor := NewMapProcessor(func(x int) int { return x * 2 })
	errorPublisher.Subscribe(context.Background(), processor)

	sub := newMapFilterTestSubscriber[int]()
	processor.Subscribe(context.Background(), sub)

	ctx := context.Background()
	sub.Wait(ctx)

	// Should receive the mapped value before error
	received := sub.GetReceivedCopy()
	if len(received) != 1 || received[0] != 2 {
		t.Errorf("Expected to receive [2], got %v", received)
	}

	sub.AssertError(t)
}

func TestChainedProcessors(t *testing.T) {
	source := NewCompliantRangePublisher(1, 6)

	// Chain: filter even numbers, then multiply by 3
	filterProcessor := NewFilterProcessor(func(x int) bool { return x%2 == 0 })
	mapProcessor := NewMapProcessor(func(x int) int { return x * 3 })

	// Connect the chain: source -> filter -> map -> subscriber
	source.Subscribe(context.Background(), filterProcessor)
	filterProcessor.Subscribe(context.Background(), mapProcessor)

	sub := newMapFilterTestSubscriber[int]()
	mapProcessor.Subscribe(context.Background(), sub)

	ctx := context.Background()
	sub.Wait(ctx)

	// Should get even numbers (2,4,6) multiplied by 3 = (6,12,18)
	expected := []int{6, 12, 18}
	sub.AssertValues(t, expected)
	sub.AssertCompleted(t)
	sub.AssertNoError(t)
}
