package streams

import (
	"context"
	"testing"
	"time"
)

func TestSubscriptionRequestControl(t *testing.T) {
	subscriber := newLocalManualRequestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 10)
	publisher.Subscribe(context.Background(), subscriber)

	// Request only 3 items
	subscriber.RequestItems(3)

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	if len(subscriber.Received()) > 3 {
		t.Errorf("Expected at most 3 values with limited request, got %d", len(subscriber.Received()))
	}

	// Request more items
	subscriber.RequestItems(7)

	subscriber.Wait(context.Background())

	if len(subscriber.Received()) != 10 {
		t.Errorf("Expected 10 values total, got %d", len(subscriber.Received()))
	}

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i, v := range subscriber.Received() {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}
}

func TestSubscriptionCancellation(t *testing.T) {
	subscriber := newLocalManualRequestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 100)
	publisher.Subscribe(context.Background(), subscriber)

	// Request a few items
	subscriber.RequestItems(5)
	time.Sleep(50 * time.Millisecond)

	// Cancel subscription
	subscriber.CancelSubscription()

	// Try to request more (should have no effect)
	subscriber.RequestItems(10)
	time.Sleep(50 * time.Millisecond)

	if len(subscriber.Received()) > 10 {
		t.Errorf("Expected at most 10 values after cancellation, got %d", len(subscriber.Received()))
	}

	// Should not be completed (was cancelled)
	if subscriber.Completed() {
		t.Error("Expected no completion after cancellation")
	}
}

func TestSubscriptionZeroRequest(t *testing.T) {
	subscriber := newLocalManualRequestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 5)
	publisher.Subscribe(context.Background(), subscriber)

	// Request 0 items - should receive nothing
	subscriber.RequestItems(0)
	time.Sleep(100 * time.Millisecond)

	if len(subscriber.Received()) != 0 {
		t.Errorf("Expected 0 values with zero request, got %d", len(subscriber.Received()))
	}
}

func TestSubscriptionNegativeRequest(t *testing.T) {
	subscriber := newLocalManualRequestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 5)
	publisher.Subscribe(context.Background(), subscriber)

	// Request negative items - should trigger error according to reactive streams spec
	subscriber.RequestItems(-1)

	subscriber.Wait(context.Background())

	if len(subscriber.Errors()) == 0 {
		t.Error("Expected error for negative request, got none")
	}
}

func TestSubscriptionMultipleRequests(t *testing.T) {
	subscriber := newLocalManualRequestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 10)
	publisher.Subscribe(context.Background(), subscriber)

	// Make multiple small requests
	subscriber.RequestItems(2)
	subscriber.RequestItems(3)
	subscriber.RequestItems(5)

	subscriber.Wait(context.Background())

	if len(subscriber.Received()) != 10 {
		t.Errorf("Expected 10 values with multiple requests, got %d", len(subscriber.Received()))
	}

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i, v := range subscriber.Received() {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}

	if !subscriber.Completed() {
		t.Error("Expected completion")
	}
}
