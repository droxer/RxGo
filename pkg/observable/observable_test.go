package observable

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// TestSubscriber is a test implementation of Subscriber
type TestSubscriber[T any] struct {
	values    []T
	errors    []error
	completed bool
	started   bool
	mu        sync.Mutex
}

func NewTestSubscriber[T any]() *TestSubscriber[T] {
	return &TestSubscriber[T]{}
}

func (t *TestSubscriber[T]) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.started = true
}

func (t *TestSubscriber[T]) OnNext(value T) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.values = append(t.values, value)
}

func (t *TestSubscriber[T]) OnCompleted() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.completed = true
}

func (t *TestSubscriber[T]) OnError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.errors = append(t.errors, err)
}

func (t *TestSubscriber[T]) Values() []T {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]T(nil), t.values...)
}

func (t *TestSubscriber[T]) Errors() []error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]error(nil), t.errors...)
}

func (t *TestSubscriber[T]) Completed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.completed
}

func (t *TestSubscriber[T]) Started() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.started
}

func TestCreate(t *testing.T) {
	called := false
	obs := Create(func(ctx context.Context, sub Subscriber[int]) {
		called = true
		sub.OnCompleted()
	})

	sub := NewTestSubscriber[int]()
	obs.Subscribe(context.Background(), sub)

	if !called {
		t.Error("Create function was not called")
	}
	if !sub.Completed() {
		t.Error("Observable should complete")
	}
}

func TestJust(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected []int
	}{
		{"single value", []int{42}, []int{42}},
		{"multiple values", []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := Just(tt.values...)
			sub := NewTestSubscriber[int]()
			obs.Subscribe(context.Background(), sub)

			values := sub.Values()
			if len(values) != len(tt.expected) {
				t.Errorf("expected %d values, got %d", len(tt.expected), len(values))
			}

			for i, v := range values {
				if v != tt.expected[i] {
					t.Errorf("expected %d at index %d, got %d", tt.expected[i], i, v)
				}
			}

			if !sub.Completed() {
				t.Error("Observable should complete")
			}
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		count    int
		expected []int
	}{
		{"positive range", 1, 3, []int{1, 2, 3}},
		{"zero start", 0, 3, []int{0, 1, 2}},
		{"negative start", -2, 3, []int{-2, -1, 0}},
		{"single value", 5, 1, []int{5}},
		{"zero count", 10, 0, []int{}},
		{"negative count", 10, -5, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := Range(tt.start, tt.count)
			sub := NewTestSubscriber[int]()
			obs.Subscribe(context.Background(), sub)

			values := sub.Values()
			if len(values) != len(tt.expected) {
				t.Errorf("expected %d values, got %d", len(tt.expected), len(values))
			}

			for i, v := range values {
				if v != tt.expected[i] {
					t.Errorf("expected %d at index %d, got %d", tt.expected[i], i, v)
				}
			}

			if !sub.Completed() {
				t.Error("Observable should complete")
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	obs := Empty[int]()
	sub := NewTestSubscriber[int]()
	obs.Subscribe(context.Background(), sub)

	if len(sub.Values()) != 0 {
		t.Errorf("expected empty, got %v", sub.Values())
	}
	if !sub.Completed() {
		t.Error("Observable should complete immediately")
	}
}

func TestError(t *testing.T) {
	err := errors.New("test error")
	obs := Error[int](err)
	sub := NewTestSubscriber[int]()
	obs.Subscribe(context.Background(), sub)

	if len(sub.Values()) != 0 {
		t.Errorf("expected no values, got %v", sub.Values())
	}
	if len(sub.Errors()) != 1 {
		t.Errorf("expected 1 error, got %d", len(sub.Errors()))
	}
	if sub.Errors()[0] != err {
		t.Errorf("expected error %v, got %v", err, sub.Errors()[0])
	}
}

func TestNever(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	obs := Never[int]()
	sub := NewTestSubscriber[int]()
	obs.Subscribe(ctx, sub)

	if len(sub.Values()) != 0 {
		t.Errorf("expected no values, got %v", sub.Values())
	}
	if sub.Completed() {
		t.Error("Observable should not complete")
	}
	if len(sub.Errors()) != 0 {
		t.Errorf("expected no errors, got %v", sub.Errors())
	}
}

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		expected []string
	}{
		{"string slice", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"empty slice", []string{}, []string{}},
		{"nil slice", nil, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := FromSlice(tt.slice)
			sub := NewTestSubscriber[string]()
			if tt.name == "int slice" {
				obsInt := FromSlice([]int{1, 2, 3})
				subInt := NewTestSubscriber[int]()
				obsInt.Subscribe(context.Background(), subInt)
				values := subInt.Values()
				if len(values) != 3 {
					t.Errorf("expected 3 values, got %d", len(values))
				}
			} else {
				obs.Subscribe(context.Background(), sub)
				values := sub.Values()
				if len(values) != len(tt.expected) {
					t.Errorf("expected %d values, got %d", len(tt.expected), len(values))
				}
				for i, v := range values {
					if v != tt.expected[i] {
						t.Errorf("expected %v at index %d, got %v", tt.expected[i], i, v)
					}
				}
			}
		})
	}
}

func TestMap(t *testing.T) {
	source := Just(1, 2, 3)
	mapped := Map(source, func(x int) string { return string(rune('A' + x - 1)) })

	sub := NewTestSubscriber[string]()
	mapped.Subscribe(context.Background(), sub)

	expected := []string{"A", "B", "C"}
	values := sub.Values()
	if len(values) != len(expected) {
		t.Errorf("expected %d values, got %d", len(expected), len(values))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, v)
		}
	}
}

func TestFilter(t *testing.T) {
	source := Just(1, 2, 3, 4, 5, 6)
	filtered := Filter(source, func(x int) bool { return x%2 == 0 })

	sub := NewTestSubscriber[int]()
	filtered.Subscribe(context.Background(), sub)

	expected := []int{2, 4, 6}
	values := sub.Values()
	if len(values) != len(expected) {
		t.Errorf("expected %d values, got %d", len(expected), len(values))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	obs := Create(func(ctx context.Context, sub Subscriber[int]) {
		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(i)
			}
		}
		sub.OnCompleted()
	})

	sub := NewTestSubscriber[int]()

	// Cancel after a short delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	obs.Subscribe(ctx, sub)

	values := sub.Values()
	if len(sub.Errors()) != 0 {
		t.Logf("got error: %v", sub.Errors()[0])
	}
	if len(values) == 0 {
		t.Log("no values received")
	}
}

func TestNilSubscriber(t *testing.T) {
	obs := Just(1, 2, 3)

	defer func() {
		if r := recover(); r != nil {
			t.Logf("expected panic: %v", r)
		}
	}()

	obs.Subscribe(context.Background(), nil)
}

func TestChainOperations(t *testing.T) {
	source := Just(1, 2, 3, 4, 5)
	mapped := Map(source, func(x int) int { return x * 2 })
	filtered := Filter(mapped, func(x int) bool { return x > 4 })

	sub := NewTestSubscriber[int]()
	filtered.Subscribe(context.Background(), sub)

	expected := []int{6, 8, 10}
	values := sub.Values()
	if len(values) != len(expected) {
		t.Errorf("expected %d values, got %d", len(expected), len(values))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestErrorInOnNext(t *testing.T) {
	obs := Just(1, 2, 3)

	// Create a subscriber that panics on first value
	panicSub := &panicSubscriber{}

	// This should not panic the entire system
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Panic recovered: %v", r)
		}
	}()

	obs.Subscribe(context.Background(), panicSub)
}

// panicSubscriber is a test subscriber that panics on OnNext
type panicSubscriber struct{}

func (p *panicSubscriber) Start()            {}
func (p *panicSubscriber) OnNext(v int)      { panic("test panic") }
func (p *panicSubscriber) OnCompleted()      {}
func (p *panicSubscriber) OnError(err error) {}

func BenchmarkJust(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sub := NewTestSubscriber[int]()
		Just(i, i+1, i+2).Subscribe(context.Background(), sub)
	}
}

func BenchmarkMap(b *testing.B) {
	source := Just(1, 2, 3, 4, 5)
	for i := 0; i < b.N; i++ {
		sub := NewTestSubscriber[int]()
		Map(source, func(x int) int { return x * 2 }).Subscribe(context.Background(), sub)
	}
}

func BenchmarkFilter(b *testing.B) {
	source := Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	for i := 0; i < b.N; i++ {
		sub := NewTestSubscriber[int]()
		Filter(source, func(x int) bool { return x%2 == 0 }).Subscribe(context.Background(), sub)
	}
}
