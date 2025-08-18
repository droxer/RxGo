package testutil

import (
	"context"
	"reflect"
	"testing"

	"github.com/droxer/RxGo/pkg/streams"
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

func (t *TestSubscriber[T]) Start() {}

func (t *TestSubscriber[T]) OnNext(value T) {
	t.Received = append(t.Received, value)
}

func (t *TestSubscriber[T]) OnError(err error) {
	t.Errors = append(t.Errors, err)
	if t.Done != nil {
		close(t.Done)
	}
}

func (t *TestSubscriber[T]) OnComplete() {
	t.Completed = true
	if t.Done != nil {
		close(t.Done)
	}
}

func (t *TestSubscriber[T]) OnSubscribe(s streams.Subscription) {
	s.Request(int64(^uint(0) >> 1))
}

func (t *TestSubscriber[T]) AssertNoError(tst *testing.T) {
	if len(t.Errors) > 0 {
		tst.Errorf("Expected no errors, got %v", t.Errors)
	}
}

func (t *TestSubscriber[T]) AssertCompleted(tst *testing.T) {
	if !t.Completed {
		tst.Error("Expected subscriber to complete")
	}
}

func (t *TestSubscriber[T]) AssertNotCompleted(tst *testing.T) {
	if t.Completed {
		tst.Error("Expected subscriber to not complete")
	}
}

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

func (t *TestSubscriber[T]) AssertError(tst *testing.T) {
	if len(t.Errors) == 0 {
		tst.Error("Expected an error, but none occurred")
	}
}

func (t *TestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-t.Done:
		return
	case <-ctx.Done():
		return
	}
}
