// Package streams_test contains comprehensive tests for RxGo reactive streams publishers.
//
// This file tests basic publisher functionality including:
// - FromSlicePublisher: Creating publishers from Go slices
// - RangePublisher: Creating publishers that emit sequential integers
// - Publisher lifecycle: Subscription, emission, completion, and error handling
// - Context cancellation and nil subscriber handling
//
// Test Categories:
// - Basic functionality tests
// - Empty/zero value tests
// - Large data set tests
// - Error handling tests
// - Concurrency tests
package streams

import (
	"context"
	"testing"
)

// TestFromSlicePublisher tests the FromSlicePublisher functionality
func TestFromSlicePublisher(t *testing.T) {
	t.Run("BasicIntSlice", func(t *testing.T) {
		sub := newLocalTestSubscriber[int]()
		publisher := FromSlicePublisher([]int{1, 2, 3, 4, 5})
		publisher.Subscribe(context.Background(), sub)

		sub.Wait(context.Background())

		expected := []int{1, 2, 3, 4, 5}
		sub.AssertValues(t, expected)
		sub.AssertCompleted(t)
		sub.AssertNoError(t)
	})

	t.Run("EmptySlice", func(t *testing.T) {
		sub := newLocalTestSubscriber[int]()
		publisher := FromSlicePublisher([]int{})
		publisher.Subscribe(context.Background(), sub)

		sub.Wait(context.Background())

		sub.AssertValues(t, []int{})
		sub.AssertCompleted(t)
		sub.AssertNoError(t)
	})

	t.Run("StringSlice", func(t *testing.T) {
		sub := newLocalTestSubscriber[string]()
		publisher := FromSlicePublisher([]string{"a", "b", "c"})
		publisher.Subscribe(context.Background(), sub)

		sub.Wait(context.Background())

		expected := []string{"a", "b", "c"}
		sub.AssertValues(t, expected)
		sub.AssertCompleted(t)
		sub.AssertNoError(t)
	})
}

func TestRangePublisher(t *testing.T) {
	sub := newLocalTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 5)
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	expected := []int{1, 2, 3, 4, 5}
	sub.AssertValues(t, expected)
	sub.AssertCompleted(t)
}

func TestRangePublisherZeroCount(t *testing.T) {
	sub := newLocalTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 0)
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	sub.AssertValues(t, []int{})
	sub.AssertCompleted(t)
}

func TestRangePublisherLarge(t *testing.T) {
	sub := newLocalTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 100)
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	expected := make([]int, 100)
	for i := 0; i < 100; i++ {
		expected[i] = i + 1
	}
	sub.AssertValues(t, expected)
	sub.AssertCompleted(t)
}

func TestPublisherNilSubscriber(t *testing.T) {
	publisher := NewCompliantRangePublisher(1, 5)
	// Should not panic
	publisher.Subscribe(context.Background(), nil)
}

func TestPublisherContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sub := newLocalTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 1000)
	publisher.Subscribe(ctx, sub)

	cancel()
	sub.Wait(ctx)

	if len(sub.Received()) > 0 {
		t.Logf("Received %d values before cancellation", len(sub.Received()))
	}
}

func TestPublisherMultipleSubscribers(t *testing.T) {
	publisher := NewCompliantRangePublisher(1, 3)

	sub1 := newLocalTestSubscriber[int]()
	sub2 := newLocalTestSubscriber[int]()

	publisher.Subscribe(context.Background(), sub1)
	publisher.Subscribe(context.Background(), sub2)

	sub1.Wait(context.Background())
	sub2.Wait(context.Background())

	expected := []int{1, 2, 3}
	sub1.AssertValues(t, expected)
	sub2.AssertValues(t, expected)
}

func TestSubscriptionInterface(t *testing.T) {
	publisher := NewCompliantRangePublisher(1, 5)
	if publisher == nil {
		t.Error("NewCompliantRangePublisher does not implement Publisher interface")
	}

	slicePublisher := FromSlicePublisher([]string{"test"})
	if slicePublisher == nil {
		t.Error("FromSlicePublisher does not implement Publisher interface")
	}
}
