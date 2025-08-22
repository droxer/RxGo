package observable

import (
	"context"
	"sync"
	"testing"
)

type testSubscriber[T any] struct {
	mu     sync.Mutex
	values []T
	err    error
	done   bool
}

func (t *testSubscriber[T]) Start() {}
func (t *testSubscriber[T]) OnNext(value T) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.values = append(t.values, value)
}
func (t *testSubscriber[T]) OnComplete() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.done = true
}
func (t *testSubscriber[T]) OnError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.err = err
	t.done = true
}

func (t *testSubscriber[T]) getValues() []T {
	t.mu.Lock()
	defer t.mu.Unlock()
	result := make([]T, len(t.values))
	copy(result, t.values)
	return result
}

func (t *testSubscriber[T]) getError() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.err
}

func (t *testSubscriber[T]) isDone() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.done
}

func TestJust(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5)
	sub := &testSubscriber[int]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(values))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestRange(t *testing.T) {
	obs := Range(10, 5)
	sub := &testSubscriber[int]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(values))
	}

	expected := []int{10, 11, 12, 13, 14}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestFromSlice(t *testing.T) {
	data := []string{"hello", "world", "test"}
	obs := FromSlice(data)
	sub := &testSubscriber[string]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	for i, v := range values {
		if v != data[i] {
			t.Errorf("Expected %s at index %d, got %s", data[i], i, v)
		}
	}
}

func TestMap(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5)
	mapped := Map(obs, func(x int) int { return x * 2 })
	sub := &testSubscriber[int]{}

	err := mapped.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	expected := []int{2, 4, 6, 8, 10}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestFilter(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5, 6)
	filtered := Filter(obs, func(x int) bool { return x%2 == 0 })
	sub := &testSubscriber[int]{}

	err := filtered.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	expected := []int{2, 4, 6}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestEmpty(t *testing.T) {
	obs := Empty[int]()
	sub := &testSubscriber[int]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Expected no values, got %d", len(values))
	}

	if !sub.isDone() {
		t.Error("Expected subscription to complete")
	}
}

func TestError(t *testing.T) {
	expectedErr := context.DeadlineExceeded
	obs := Error[int](expectedErr)
	sub := &testSubscriber[int]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	if sub.getError() != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, sub.getError())
	}

	if !sub.isDone() {
		t.Error("Expected subscription to complete with error")
	}
}

func TestContextCancellation(t *testing.T) {
	obs := Range(1, 1000)
	sub := &testSubscriber[int]{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cancel()

	err := obs.Subscribe(ctx, sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Expected no values due to cancellation, got %d", len(values))
	}

	if sub.getError() != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", sub.getError())
	}
}
