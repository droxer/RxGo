package streams

import (
	"context"
	"testing"
)

func TestNewCompliantFromSlicePublisher(t *testing.T) {
	data := []string{"hello", "world", "reactive", "streams"}
	publisher := NewCompliantFromSlicePublisher(data)

	subscriber := newLocalTestSubscriber[string]()
	publisher.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	subscriber.AssertValues(t, data)
	subscriber.AssertCompleted(t)
	subscriber.AssertNoError(t)
}

func TestNewCompliantFromSlicePublisherEmpty(t *testing.T) {
	publisher := NewCompliantFromSlicePublisher([]int{})

	subscriber := newLocalTestSubscriber[int]()
	publisher.Subscribe(context.Background(), subscriber)

	subscriber.Wait(context.Background())

	subscriber.AssertValues(t, []int{})
	subscriber.AssertCompleted(t)
	subscriber.AssertNoError(t)
}

func TestCompliantPublisherDemandControl(t *testing.T) {
	publisher := NewCompliantRangePublisher(1, 5)

	subscriber := newLocalManualRequestSubscriber[int]()
	publisher.Subscribe(context.Background(), subscriber)

	// Request all items to ensure test completes
	subscriber.RequestItems(5)
	subscriber.Wait(context.Background())

	if len(subscriber.Received()) != 5 {
		t.Errorf("Expected 5 values with compliant publisher, got %d", len(subscriber.Received()))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range subscriber.Received() {
		if v != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}
}

func TestCompliantPublisherMultipleSubscribers(t *testing.T) {
	publisher := NewCompliantRangePublisher(10, 3)

	sub1 := newLocalTestSubscriber[int]()
	sub2 := newLocalTestSubscriber[int]()

	publisher.Subscribe(context.Background(), sub1)
	publisher.Subscribe(context.Background(), sub2)

	sub1.Wait(context.Background())
	sub2.Wait(context.Background())

	expected := []int{10, 11, 12}
	sub1.AssertValues(t, expected)
	sub2.AssertValues(t, expected)
	sub1.AssertCompleted(t)
	sub2.AssertCompleted(t)
}

func TestCompliantPublisherNegativeRequest(t *testing.T) {
	publisher := NewCompliantRangePublisher(1, 3)

	subscriber := newLocalManualRequestSubscriber[int]()
	publisher.Subscribe(context.Background(), subscriber)

	// Request negative items - should trigger error
	subscriber.RequestItems(-1)
	subscriber.Wait(context.Background())

	if len(subscriber.Errors()) == 0 {
		t.Error("Expected error for negative request, got none")
	}
}
