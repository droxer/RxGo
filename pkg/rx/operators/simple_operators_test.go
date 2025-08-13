package operators

import (
	"context"
	"testing"

	"github.com/droxer/RxGo/internal/testutil"
	"github.com/droxer/RxGo/pkg/rx"
)

// (moved to pkg/rx/testutil)

func TestMap(t *testing.T) {
	// Test basic mapping
	sub := testutil.NewTestSubscriber[int]()
	source := rx.Just(1, 2, 3)
	mapped := Map(source, func(x int) int { return x * 2 })
	mapped.Subscribe(context.Background(), sub)

	expected := []int{2, 4, 6}
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

func TestFilter(t *testing.T) {
	// Test basic filtering
	sub := testutil.NewTestSubscriber[int]()
	source := rx.Just(1, 2, 3, 4, 5, 6)
	filtered := Filter(source, func(x int) bool { return x%2 == 0 })
	filtered.Subscribe(context.Background(), sub)

	expected := []int{2, 4, 6}
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

func TestMapFilterChain(t *testing.T) {
	// Test chaining map and filter
	sub := testutil.NewTestSubscriber[int]()
	source := rx.Just(1, 2, 3, 4, 5)

	// Chain: double then filter for > 4
	doubled := Map(source, func(x int) int { return x * 2 })
	filtered := Filter(doubled, func(x int) bool { return x > 4 })

	filtered.Subscribe(context.Background(), sub)

	expected := []int{6, 8, 10}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}
}

func TestEmptySource(t *testing.T) {
	// Test with empty source
	sub := testutil.NewTestSubscriber[int]()
	source := rx.Just[int]()
	mapped := Map(source, func(x int) int { return x * 2 })
	mapped.Subscribe(context.Background(), sub)

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestStringMap(t *testing.T) {
	// Test mapping to different type
	sub := testutil.NewTestSubscriber[string]()
	source := rx.Just(1, 2, 3)
	mapped := Map(source, func(x int) string { return string(rune(x + 'A')) })
	mapped.Subscribe(context.Background(), sub)

	expected := []string{"B", "C", "D"}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}
}
