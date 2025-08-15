package testutil

import (
	"context"
	"reflect"
	"testing"
)

// TestSubscriber is a test implementation of Subscriber for testing purposes
type TestSubscriber[T any] struct {
	Received  []T
	Completed bool
	Errors    []error
	Done      chan struct{}
}

// NewTestSubscriber creates a new TestSubscriber instance
func NewTestSubscriber[T any]() *TestSubscriber[T] {
	return &TestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

// Start implements Subscriber interface
func (t *TestSubscriber[T]) Start() {}

// OnNext implements Subscriber interface
func (t *TestSubscriber[T]) OnNext(value T) {
	t.Received = append(t.Received, value)
}

// OnError implements Subscriber interface
func (t *TestSubscriber[T]) OnError(err error) {
	t.Errors = append(t.Errors, err)
	if t.Done != nil {
		close(t.Done)
	}
}

// OnCompleted implements Subscriber interface
func (t *TestSubscriber[T]) OnCompleted() {
	t.Completed = true
	if t.Done != nil {
		close(t.Done)
	}
}

// OnSubscribe implements a Subscriber interface method
func (t *TestSubscriber[T]) OnSubscribe(s interface{}) {
	// Auto-request all items
	if req, ok := s.(interface{ Request(int64) }); ok {
		req.Request(int64(^uint(0) >> 1))
	}
}

// AssertNoError checks that no errors occurred
func (t *TestSubscriber[T]) AssertNoError(tst *testing.T) {
	if len(t.Errors) > 0 {
		tst.Errorf("Expected no errors, got %v", t.Errors)
	}
}

// AssertCompleted checks that the subscriber completed
func (t *TestSubscriber[T]) AssertCompleted(tst *testing.T) {
	if !t.Completed {
		tst.Error("Expected subscriber to complete")
	}
}

// AssertNotCompleted checks that the subscriber did not complete
func (t *TestSubscriber[T]) AssertNotCompleted(tst *testing.T) {
	if t.Completed {
		tst.Error("Expected subscriber to not complete")
	}
}

// AssertValues checks that the received values match expected
func (t *TestSubscriber[T]) AssertValues(tst *testing.T, expected []T) {
	if len(t.Received) != len(expected) {
		tst.Errorf("Expected %d values, got %d: %v", len(expected), len(t.Received), t.Received)
		return
	}
	for i, v := range expected {
		if !reflect.DeepEqual(t.Received[i], v) {
			tst.Errorf("Expected value[%d] to be %v, got %v", i, v, t.Received[i])
		}
	}
}

// AssertError checks that an error occurred
func (t *TestSubscriber[T]) AssertError(tst *testing.T) {
	if len(t.Errors) == 0 {
		tst.Error("Expected an error, but none occurred")
	}
}

// Wait waits for completion or error
func (t *TestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-t.Done:
		return
	case <-ctx.Done():
		return
	}
}
