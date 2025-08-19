package streams

import (
	"context"
	"errors"
	"testing"
)

func TestMapProcessor(t *testing.T) {
	source := NewCompliantRangePublisher(1, 5)
	mapProcessor := NewMapProcessor(func(x int) int { return x * 2 })

	// Connect source to processor
	source.Subscribe(context.Background(), mapProcessor)

	// Subscribe to processor
	subscriber := newLocalTestSubscriber[int]()
	mapProcessor.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	expected := []int{2, 4, 6, 8, 10}
	subscriber.AssertValues(t, expected)
	subscriber.AssertCompleted(t)
	subscriber.AssertNoError(t)
}

func TestMapProcessorStringTransform(t *testing.T) {
	source := NewCompliantFromSlicePublisher([]int{1, 2, 3})
	mapProcessor := NewMapProcessor(func(x int) string {
		return map[int]string{1: "one", 2: "two", 3: "three"}[x]
	})

	source.Subscribe(context.Background(), mapProcessor)

	subscriber := newLocalTestSubscriber[string]()
	mapProcessor.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	expected := []string{"one", "two", "three"}
	subscriber.AssertValues(t, expected)
	subscriber.AssertCompleted(t)
}

func TestFilterProcessor(t *testing.T) {
	source := NewCompliantRangePublisher(1, 10)
	filterProcessor := NewFilterProcessor(func(x int) bool { return x%2 == 0 })

	source.Subscribe(context.Background(), filterProcessor)

	subscriber := newLocalTestSubscriber[int]()
	filterProcessor.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	expected := []int{2, 4, 6, 8, 10}
	subscriber.AssertValues(t, expected)
	subscriber.AssertCompleted(t)
}

func TestFilterProcessorAllFiltered(t *testing.T) {
	source := NewCompliantRangePublisher(1, 5)
	filterProcessor := NewFilterProcessor(func(x int) bool { return x > 10 })

	source.Subscribe(context.Background(), filterProcessor)

	subscriber := newLocalTestSubscriber[int]()
	filterProcessor.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	subscriber.AssertValues(t, []int{})
	subscriber.AssertCompleted(t)
}

func TestFlatMapProcessor(t *testing.T) {
	source := NewCompliantFromSlicePublisher([]int{1, 2})
	flatMapProcessor := NewFlatMapProcessor(func(x int) Publisher[int] {
		return NewCompliantFromSlicePublisher([]int{x, x + 10})
	})

	source.Subscribe(context.Background(), flatMapProcessor)

	subscriber := newLocalTestSubscriber[int]()
	flatMapProcessor.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	// FlatMap should emit: 1, 11, 2, 12 (order may vary due to concurrency)
	received := subscriber.Received()
	if len(received) != 4 {
		t.Errorf("Expected 4 values from flatmap, got %d", len(received))
	}

	// Check that all expected values are present
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

	mapProcessor := NewMapProcessor(func(x int) int { return x * 2 })
	errorPublisher.Subscribe(context.Background(), mapProcessor)

	subscriber := newLocalTestSubscriber[int]()
	mapProcessor.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	// Should receive the mapped value before error
	if len(subscriber.Received()) != 1 || subscriber.Received()[0] != 2 {
		t.Errorf("Expected to receive [2], got %v", subscriber.Received())
	}

	// Should have received the error
	subscriber.AssertError(t)
}

func TestChainedProcessors(t *testing.T) {
	source := NewCompliantRangePublisher(1, 6)

	// Chain: filter even numbers, then multiply by 3
	filterProcessor := NewFilterProcessor(func(x int) bool { return x%2 == 0 })
	mapProcessor := NewMapProcessor(func(x int) int { return x * 3 })

	// Connect the chain: source -> filter -> map -> subscriber
	source.Subscribe(context.Background(), filterProcessor)
	filterProcessor.Subscribe(context.Background(), mapProcessor)

	subscriber := newLocalTestSubscriber[int]()
	mapProcessor.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	// Should get even numbers (2,4,6) multiplied by 3 = (6,12,18)
	expected := []int{6, 12, 18}
	subscriber.AssertValues(t, expected)
	subscriber.AssertCompleted(t)
}
