package streams

import (
	"context"
	"testing"
)

type TestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
}

func NewTestSubscriber[T any]() *TestSubscriber[T] {
	return &TestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

func (t *TestSubscriber[T]) OnSubscribe(sub Subscription) {
	sub.Request(int64(^uint(0) >> 1))
}

func (t *TestSubscriber[T]) OnNext(value T) {
	t.Received = append(t.Received, value)
}

func (t *TestSubscriber[T]) OnError(err error) {
	t.Errors = append(t.Errors, err)
	close(t.Done)
}

func (t *TestSubscriber[T]) OnComplete() {
	t.Completed = true
	close(t.Done)
}

func (t *TestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-t.Done:
		return
	case <-ctx.Done():
		return
	}
}

func TestFromSlicePublisher(t *testing.T) {
	sub := NewTestSubscriber[int]()
	publisher := FromSlicePublisher([]int{1, 2, 3, 4, 5})
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	expected := []int{1, 2, 3, 4, 5}
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

	if len(sub.Errors) > 0 {
		t.Errorf("Expected no errors, got %v", sub.Errors)
	}
}

func TestRangePublisher(t *testing.T) {
	sub := NewTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 5)
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	expected := []int{1, 2, 3, 4, 5}
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

func TestRangePublisherZeroCount(t *testing.T) {
	sub := NewTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 0)
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestFromSlicePublisherEmpty(t *testing.T) {
	sub := NewTestSubscriber[int]()
	publisher := FromSlicePublisher([]int{})
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestStringPublisher(t *testing.T) {
	sub := NewTestSubscriber[string]()
	publisher := FromSlicePublisher([]string{"a", "b", "c"})
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	expected := []string{"a", "b", "c"}
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

func TestPublisherContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sub := NewTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 1000)
	publisher.Subscribe(ctx, sub)

	cancel()

	sub.Wait(ctx)

	if len(sub.Received) > 0 {
		t.Logf("Received %d values before cancellation", len(sub.Received))
	}
}

func TestPublisherMultipleSubscribers(t *testing.T) {
	publisher := NewCompliantRangePublisher(1, 3)

	sub1 := NewTestSubscriber[int]()
	sub2 := NewTestSubscriber[int]()

	publisher.Subscribe(context.Background(), sub1)
	publisher.Subscribe(context.Background(), sub2)

	sub1.Wait(context.Background())
	sub2.Wait(context.Background())

	expected := []int{1, 2, 3}

	if len(sub1.Received) != len(expected) {
		t.Errorf("Sub1: Expected %d values, got %d", len(expected), len(sub1.Received))
	}

	if len(sub2.Received) != len(expected) {
		t.Errorf("Sub2: Expected %d values, got %d", len(expected), len(sub2.Received))
	}

	for i, v := range expected {
		if sub1.Received[i] != v {
			t.Errorf("Sub1: Expected %v at index %d, got %v", v, i, sub1.Received[i])
		}
		if sub2.Received[i] != v {
			t.Errorf("Sub2: Expected %v at index %d, got %v", v, i, sub2.Received[i])
		}
	}
}

func TestPublisherWithSubscription(t *testing.T) {
	sub := NewTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 3)
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	expected := []int{1, 2, 3}
	if len(sub.Received) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(sub.Received))
	}

	for i, v := range expected {
		if sub.Received[i] != v {
			t.Errorf("Expected %v at index %d, got %v", v, i, sub.Received[i])
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestFromSlicePublisherEmptyString(t *testing.T) {
	sub := NewTestSubscriber[string]()
	publisher := FromSlicePublisher([]string{})
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	if len(sub.Received) != 0 {
		t.Errorf("Expected 0 values, got %d", len(sub.Received))
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestRangePublisherLarge(t *testing.T) {
	sub := NewTestSubscriber[int]()
	publisher := NewCompliantRangePublisher(1, 100)
	publisher.Subscribe(context.Background(), sub)

	sub.Wait(context.Background())

	if len(sub.Received) != 100 {
		t.Errorf("Expected 100 values, got %d", len(sub.Received))
	}

	for i, v := range sub.Received {
		if v != i+1 {
			t.Errorf("Expected %v at index %d, got %v", i+1, i, v)
		}
	}

	if !sub.Completed {
		t.Error("Expected completion")
	}
}

func TestPublisherNilSubscriber(t *testing.T) {
	publisher := NewCompliantRangePublisher(1, 5)
	publisher.Subscribe(context.Background(), nil)

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
