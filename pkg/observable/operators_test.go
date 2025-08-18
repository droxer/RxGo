package observable

import (
	"context"
	"testing"
	"time"

	"github.com/droxer/RxGo/pkg/scheduler"
)

func TestMerge(t *testing.T) {
	ctx := context.Background()

	obs1 := Just(1, 2, 3)
	obs2 := Just(4, 5, 6)

	merged := Merge(obs1, obs2)

	sub := &testSubscriber[int]{}
	merged.Subscribe(ctx, sub)

	// Wait a bit for async operations
	time.Sleep(10 * time.Millisecond)

	values := sub.getValues()
	if len(values) != 6 {
		t.Errorf("Expected 6 values, got %d", len(values))
	}

	// Check that all values are present (order may vary due to concurrency)
	expected := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true}
	for _, v := range values {
		if !expected[v] {
			t.Errorf("Unexpected value: %d", v)
		}
		delete(expected, v)
	}

	if len(expected) != 0 {
		t.Errorf("Missing values: %v", expected)
	}
}

func TestConcat(t *testing.T) {
	ctx := context.Background()

	obs1 := Just(1, 2, 3)
	obs2 := Just(4, 5, 6)

	concatenated := Concat(obs1, obs2)

	sub := &testSubscriber[int]{}
	concatenated.Subscribe(ctx, sub)

	// Wait a bit for async operations
	time.Sleep(10 * time.Millisecond)

	values := sub.getValues()
	if len(values) != 6 {
		t.Errorf("Expected 6 values, got %d", len(values))
	}

	// Check that values are in the correct order
	expected := []int{1, 2, 3, 4, 5, 6}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestTake(t *testing.T) {
	ctx := context.Background()

	obs := Range(1, 10) // 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
	taken := Take(obs, 5)

	sub := &testSubscriber[int]{}
	taken.Subscribe(ctx, sub)

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

func TestSkip(t *testing.T) {
	ctx := context.Background()

	obs := Range(1, 10) // 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
	skipped := Skip(obs, 5)

	sub := &testSubscriber[int]{}
	skipped.Subscribe(ctx, sub)

	values := sub.getValues()
	if len(values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(values))
	}

	expected := []int{6, 7, 8, 9, 10}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestDistinct(t *testing.T) {
	ctx := context.Background()

	obs := FromSlice([]int{1, 2, 2, 3, 3, 3, 4, 5, 5})
	distinct := Distinct(obs)

	sub := &testSubscriber[int]{}
	distinct.Subscribe(ctx, sub)

	values := sub.getValues()
	if len(values) != 5 {
		t.Errorf("Expected 5 distinct values, got %d", len(values))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestObserveOn(t *testing.T) {
	ctx := context.Background()

	obs := Just(1, 2, 3)
	scheduled := ObserveOn(obs, scheduler.Trampoline)

	sub := &testSubscriber[int]{}
	scheduled.Subscribe(ctx, sub)

	// Wait a bit for async operations
	time.Sleep(10 * time.Millisecond)

	values := sub.getValues()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	expected := []int{1, 2, 3}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}
