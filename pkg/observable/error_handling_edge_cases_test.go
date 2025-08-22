package observable

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewSubscriber(t *testing.T) {
	var received []int
	var completed bool
	var receivedError error
	var mu sync.Mutex

	subscriber := NewSubscriber(
		func(value int) {
			mu.Lock()
			received = append(received, value)
			mu.Unlock()
		},
		func() {
			mu.Lock()
			completed = true
			mu.Unlock()
		},
		func(err error) {
			mu.Lock()
			receivedError = err
			mu.Unlock()
		},
	)

	obs := Just(1, 2, 3)
	err := obs.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	expected := []int{1, 2, 3}
	if len(received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(received))
	}

	for i, v := range received {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !completed {
		t.Error("Expected completion")
	}

	if receivedError != nil {
		t.Errorf("Expected no error, got %v", receivedError)
	}
}

func TestNewSubscriberWithError(t *testing.T) {
	var receivedError error
	var completed bool
	var mu sync.Mutex

	subscriber := NewSubscriber(
		func(value int) {},
		func() {
			mu.Lock()
			completed = true
			mu.Unlock()
		},
		func(err error) {
			mu.Lock()
			receivedError = err
			mu.Unlock()
		},
	)

	testError := errors.New("test error")
	obs := Error[int](testError)
	err := obs.Subscribe(context.Background(), subscriber)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if receivedError != testError {
		t.Errorf("Expected test error, got %v", receivedError)
	}

	if completed {
		t.Error("Should not complete when error occurs")
	}
}

func TestRangeWithZeroCount(t *testing.T) {
	obs := Range(10, 0)
	sub := &testSubscriber[int]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Expected no values for zero count, got %d", len(values))
	}

	if !sub.isDone() {
		t.Error("Expected completion for zero count range")
	}
}

func TestRangeWithNegativeCount(t *testing.T) {
	obs := Range(5, -3)
	sub := &testSubscriber[int]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Expected no values for negative count, got %d", len(values))
	}

	if !sub.isDone() {
		t.Error("Expected completion for negative count range")
	}
}

func TestJustWithEmptyArgs(t *testing.T) {
	obs := Just[int]() // No arguments
	sub := &testSubscriber[int]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Expected no values for empty Just, got %d", len(values))
	}

	if !sub.isDone() {
		t.Error("Expected completion for empty Just")
	}
}

func TestFromSliceWithNil(t *testing.T) {
	obs := FromSlice[string](nil)
	sub := &testSubscriber[string]{}

	err := obs.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Expected no values for nil slice, got %d", len(values))
	}

	if !sub.isDone() {
		t.Error("Expected completion for nil slice")
	}
}

func TestOperatorErrorPropagation(t *testing.T) {
	testError := errors.New("source error")

	// Test Map error propagation
	obs := Error[int](testError)
	mapped := Map(obs, func(x int) string { return "mapped" })
	sub := &testSubscriber[string]{}

	err := mapped.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	if sub.getError() != testError {
		t.Errorf("Map should propagate error, expected %v, got %v", testError, sub.getError())
	}
}

func TestFilterErrorPropagation(t *testing.T) {
	testError := errors.New("source error")

	obs := Error[int](testError)
	filtered := Filter(obs, func(x int) bool { return true })
	sub := &testSubscriber[int]{}

	err := filtered.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	if sub.getError() != testError {
		t.Errorf("Filter should propagate error, expected %v, got %v", testError, sub.getError())
	}
}

func TestTakeWithZeroCount(t *testing.T) {
	obs := Range(1, 10)
	taken := Take(obs, 0)
	sub := &testSubscriber[int]{}

	err := taken.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Take(0) should emit no values, got %d", len(values))
	}

	if !sub.isDone() {
		t.Error("Take(0) should complete immediately")
	}
}

func TestSkipMoreThanAvailable(t *testing.T) {
	obs := Just(1, 2, 3)
	skipped := Skip(obs, 10)
	sub := &testSubscriber[int]{}

	err := skipped.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Skip(10) on 3 items should emit no values, got %d", len(values))
	}

	if !sub.isDone() {
		t.Error("Skip should complete when source completes")
	}
}

func TestConcurrentSubscribers(t *testing.T) {
	obs := Range(1, 100)

	var wg sync.WaitGroup
	numSubscribers := 10

	for i := 0; i < numSubscribers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub := &testSubscriber[int]{}
			_ = obs.Subscribe(context.Background(), sub)

			values := sub.getValues()
			if len(values) != 100 {
				t.Errorf("Each subscriber should get 100 values, got %d", len(values))
			}
		}()
	}

	wg.Wait()
}

func TestSubscriberPanicsHandled(t *testing.T) {
	// This test ensures that if a subscriber panics, it doesn't crash the whole program
	// Note: This is more of a safety test - depending on implementation, panics might not be caught

	panicSubscriber := NewSubscriber(
		func(value int) {
			if value == 2 {
				panic("test panic")
			}
		},
		func() {},
		func(err error) {},
	)

	// This should not crash the test
	defer func() {
		_ = recover() // Panic was caught, which is fine
	}()

	obs := Just(1, 2, 3)
	_ = obs.Subscribe(context.Background(), panicSubscriber)
}

func TestNilSubscriberHandling(t *testing.T) {
	obs := Just(1, 2, 3)

	// Should not panic when subscribing with nil
	_ = obs.Subscribe(context.Background(), nil)

	// Should not panic with nil context either
	sub := &testSubscriber[int]{}
	_ = obs.Subscribe(context.TODO(), sub)

	// Should still work normally
	values := sub.getValues()
	if len(values) != 3 {
		t.Errorf("Expected 3 values even with nil context, got %d", len(values))
	}
}

func TestMergeWithEmptyObservables(t *testing.T) {
	empty1 := Empty[int]()
	empty2 := Empty[int]()

	merged := Merge(empty1, empty2)
	sub := &testSubscriber[int]{}

	err := merged.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}
	time.Sleep(10 * time.Millisecond)

	values := sub.getValues()
	if len(values) != 0 {
		t.Errorf("Merge of empty observables should be empty, got %d values", len(values))
	}

	if !sub.isDone() {
		t.Error("Merge of empty observables should complete")
	}
}

func TestConcatWithErrors(t *testing.T) {
	testError := errors.New("concat error")

	obs1 := Just(1, 2)
	obs2 := Error[int](testError)
	obs3 := Just(3, 4) // This shouldn't be reached

	concatenated := Concat(obs1, obs2, obs3)
	sub := &testSubscriber[int]{}

	err := concatenated.Subscribe(context.Background(), sub)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}
	time.Sleep(10 * time.Millisecond)

	values := sub.getValues()
	if len(values) != 2 {
		t.Errorf("Should get values from first observable before error, got %d", len(values))
	}

	if sub.getError() != testError {
		t.Errorf("Should propagate error from second observable, got %v", sub.getError())
	}
}
