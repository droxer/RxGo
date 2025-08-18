package streams

import (
	"context"
	"testing"
)

func TestMergeProcessor(t *testing.T) {
	ctx := context.Background()

	pub1 := NewCompliantFromSlicePublisher([]int{1, 2, 3})
	pub2 := NewCompliantFromSlicePublisher([]int{4, 5, 6})

	mergeProcessor := NewMergeProcessor[int](pub1, pub2)

	sub := NewTestSubscriber[int]()
	mergeProcessor.Subscribe(ctx, sub)

	sub.Wait(ctx)

	if len(sub.Received) != 6 {
		t.Errorf("Expected 6 values, got %d", len(sub.Received))
	}

	expected := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true}
	for _, v := range sub.Received {
		if !expected[v] {
			t.Errorf("Unexpected value: %d", v)
		}
		delete(expected, v)
	}

	if len(expected) != 0 {
		t.Errorf("Missing values: %v", expected)
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestConcatProcessor(t *testing.T) {
	ctx := context.Background()

	pub1 := NewCompliantFromSlicePublisher([]int{1, 2, 3})
	pub2 := NewCompliantFromSlicePublisher([]int{4, 5, 6})

	concatProcessor := NewConcatProcessor[int](pub1, pub2)

	sub := NewTestSubscriber[int]()
	concatProcessor.Subscribe(ctx, sub)

	// Wait for completion
	sub.Wait(ctx)

	if len(sub.Received) != 6 {
		t.Errorf("Expected 6 values, got %d", len(sub.Received))
	}

	// Check that values are in the correct order
	expected := []int{1, 2, 3, 4, 5, 6}
	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestTakeProcessor(t *testing.T) {
	ctx := context.Background()

	source := NewCompliantRangePublisher(1, 10)

	takeProcessor := NewTakeProcessor[int](5)
	source.Subscribe(ctx, takeProcessor)

	sub := NewTestSubscriber[int]()
	takeProcessor.Subscribe(ctx, sub)

	sub.Wait(ctx)

	if len(sub.Received) != 5 {
		t.Errorf("Expected 5 values, got %d", len(sub.Received))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestSkipProcessor(t *testing.T) {
	ctx := context.Background()

	source := NewCompliantRangePublisher(1, 10)

	skipProcessor := NewSkipProcessor[int](5)
	source.Subscribe(ctx, skipProcessor)

	sub := NewTestSubscriber[int]()
	skipProcessor.Subscribe(ctx, sub)

	sub.Wait(ctx)

	if len(sub.Received) != 5 {
		t.Errorf("Expected 5 values, got %d", len(sub.Received))
	}

	expected := []int{6, 7, 8, 9, 10}
	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestDistinctProcessor(t *testing.T) {
	ctx := context.Background()

	source := NewCompliantFromSlicePublisher([]int{1, 2, 2, 3, 3, 3, 4, 5, 5})

	distinctProcessor := NewDistinctProcessor[int]()
	source.Subscribe(ctx, distinctProcessor)

	sub := NewTestSubscriber[int]()
	distinctProcessor.Subscribe(ctx, sub)

	sub.Wait(ctx)

	if len(sub.Received) != 5 {
		t.Errorf("Expected 5 distinct values, got %d", len(sub.Received))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range sub.Received {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}
