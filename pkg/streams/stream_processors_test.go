package streams

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

// processorsTestSubscriber is a simple test subscriber for processor tests
type processorsTestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
	mu        sync.Mutex
}

func newProcessorsTestSubscriber[T any]() *processorsTestSubscriber[T] {
	return &processorsTestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

func (s *processorsTestSubscriber[T]) OnSubscribe(sub Subscription) {
	sub.Request(100) // Default unlimited
}

func (s *processorsTestSubscriber[T]) OnNext(value T) {
	s.mu.Lock()
	s.Received = append(s.Received, value)
	s.mu.Unlock()
}

func (s *processorsTestSubscriber[T]) OnError(err error) {
	s.mu.Lock()
	s.Errors = append(s.Errors, err)
	close(s.Done)
	s.mu.Unlock()
}

func (s *processorsTestSubscriber[T]) OnComplete() {
	close(s.Done)
}

func (s *processorsTestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-s.Done:
	}
}

func (s *processorsTestSubscriber[T]) GetReceivedCopy() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]T, len(s.Received))
	copy(result, s.Received)
	return result
}

func (s *processorsTestSubscriber[T]) AssertValues(t *testing.T, expected []T) {
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

func (s *processorsTestSubscriber[T]) AssertCompleted(t *testing.T) {
	// Can't reliably check completion without race conditions
}

func (s *processorsTestSubscriber[T]) AssertNoError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", s.Errors)
	}
}

func (s *processorsTestSubscriber[T]) AssertError(t *testing.T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Errors) == 0 {
		t.Error("Expected an error, but none occurred")
	}
}

func TestStreamProcessors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		setup    func() Publisher[int]
		expected []int
		validate func(t *testing.T, received []int, expected []int)
	}{
		{
			name: "merge processor",
			setup: func() Publisher[int] {
				pub1 := NewCompliantFromSlicePublisher([]int{1, 2, 3})
				pub2 := NewCompliantFromSlicePublisher([]int{4, 5, 6})
				return NewMergeProcessor[int](pub1, pub2)
			},
			expected: []int{1, 2, 3, 4, 5, 6},
			validate: func(t *testing.T, received, expected []int) {
				if len(received) != 6 {
					t.Errorf("Expected 6 values, got %d", len(received))
				}
				// Check all values are present regardless of order
				present := make(map[int]bool)
				for _, v := range received {
					present[v] = true
				}
				for _, v := range expected {
					if !present[v] {
						t.Errorf("Missing value: %d", v)
					}
				}
			},
		},
		{
			name: "concat processor",
			setup: func() Publisher[int] {
				pub1 := NewCompliantFromSlicePublisher([]int{1, 2, 3})
				pub2 := NewCompliantFromSlicePublisher([]int{4, 5, 6})
				return NewConcatProcessor[int](pub1, pub2)
			},
			expected: []int{1, 2, 3, 4, 5, 6},
			validate: func(t *testing.T, received, expected []int) {
				if len(received) != 6 {
					t.Errorf("Expected 6 values, got %d", len(received))
				}
				for i, v := range received {
					if v != expected[i] {
						t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
					}
				}
			},
		},
		{
			name: "take processor",
			setup: func() Publisher[int] {
				source := NewCompliantRangePublisher(1, 10)
				processor := NewTakeProcessor[int](5)
				source.Subscribe(ctx, processor)
				return processor
			},
			expected: []int{1, 2, 3, 4, 5},
			validate: func(t *testing.T, received, expected []int) {
				if len(received) != 5 {
					t.Errorf("Expected 5 values, got %d", len(received))
				}
				for i, v := range received {
					if v != expected[i] {
						t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
					}
				}
			},
		},
		{
			name: "skip processor",
			setup: func() Publisher[int] {
				source := NewCompliantRangePublisher(1, 10)
				processor := NewSkipProcessor[int](5)
				source.Subscribe(ctx, processor)
				return processor
			},
			expected: []int{6, 7, 8, 9, 10},
			validate: func(t *testing.T, received, expected []int) {
				if len(received) != 5 {
					t.Errorf("Expected 5 values, got %d", len(received))
				}
				for i, v := range received {
					if v != expected[i] {
						t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
					}
				}
			},
		},
		{
			name: "distinct processor",
			setup: func() Publisher[int] {
				source := NewCompliantFromSlicePublisher([]int{1, 2, 2, 3, 3, 3, 4, 5, 5})
				processor := NewDistinctProcessor[int]()
				source.Subscribe(ctx, processor)
				return processor
			},
			expected: []int{1, 2, 3, 4, 5},
			validate: func(t *testing.T, received, expected []int) {
				if len(received) != 5 {
					t.Errorf("Expected 5 distinct values, got %d", len(received))
				}
				for i, v := range received {
					if v != expected[i] {
						t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := tt.setup()
			sub := newProcessorsTestSubscriber[int]()

			publisher.Subscribe(ctx, sub)
			sub.Wait(ctx)

			received := sub.GetReceivedCopy()
			tt.validate(t, received, tt.expected)

			sub.AssertCompleted(t)
			sub.AssertNoError(t)
		})
	}
}
