package observable

import (
	"context"
	"testing"
)

// Test subscriber for testing
type testSubscriber[T any] struct {
	values []T
	err    error
	done   bool
}

func (t *testSubscriber[T]) Start() {}
func (t *testSubscriber[T]) OnNext(value T) {
	t.values = append(t.values, value)
}
func (t *testSubscriber[T]) OnCompleted() {
	t.done = true
}
func (t *testSubscriber[T]) OnError(err error) {
	t.err = err
	t.done = true
}

func TestJust(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5)
	sub := &testSubscriber[int]{}

	obs.Subscribe(context.Background(), sub)

	if len(sub.values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(sub.values))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range sub.values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestRange(t *testing.T) {
	obs := Range(10, 5)
	sub := &testSubscriber[int]{}

	obs.Subscribe(context.Background(), sub)

	if len(sub.values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(sub.values))
	}

	expected := []int{10, 11, 12, 13, 14}
	for i, v := range sub.values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestFromSlice(t *testing.T) {
	data := []string{"hello", "world", "test"}
	obs := FromSlice(data)
	sub := &testSubscriber[string]{}

	obs.Subscribe(context.Background(), sub)

	if len(sub.values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(sub.values))
	}

	for i, v := range sub.values {
		if v != data[i] {
			t.Errorf("Expected %s at index %d, got %s", data[i], i, v)
		}
	}
}

func TestMap(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5)
	mapped := Map(obs, func(x int) int { return x * 2 })
	sub := &testSubscriber[int]{}

	mapped.Subscribe(context.Background(), sub)

	expected := []int{2, 4, 6, 8, 10}
	for i, v := range sub.values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestFilter(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5, 6)
	filtered := Filter(obs, func(x int) bool { return x%2 == 0 })
	sub := &testSubscriber[int]{}

	filtered.Subscribe(context.Background(), sub)

	expected := []int{2, 4, 6}
	for i, v := range sub.values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestEmpty(t *testing.T) {
	obs := Empty[int]()
	sub := &testSubscriber[int]{}

	obs.Subscribe(context.Background(), sub)

	if len(sub.values) != 0 {
		t.Errorf("Expected no values, got %d", len(sub.values))
	}

	if !sub.done {
		t.Error("Expected subscription to complete")
	}
}

func TestError(t *testing.T) {
	expectedErr := context.DeadlineExceeded
	obs := Error[int](expectedErr)
	sub := &testSubscriber[int]{}

	obs.Subscribe(context.Background(), sub)

	if sub.err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, sub.err)
	}

	if !sub.done {
		t.Error("Expected subscription to complete with error")
	}
}

func TestContextCancellation(t *testing.T) {
	obs := Range(1, 1000)
	sub := &testSubscriber[int]{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel immediately
	cancel()

	obs.Subscribe(ctx, sub)

	// Should have no values due to cancellation
	if len(sub.values) != 0 {
		t.Errorf("Expected no values due to cancellation, got %d", len(sub.values))
	}

	if sub.err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", sub.err)
	}
}
