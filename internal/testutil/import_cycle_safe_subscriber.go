// Package testutil provides local test utilities for streams package to avoid import cycles.
//
// This file provides simplified versions of test utilities that work with
// interface{} instead of the concrete streams.Subscription type to avoid
// circular imports when the streams package needs to use test utilities.
package testutil

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

// ImportCycleSafeTestSubscriber is a simplified test subscriber for use within the streams package
// to avoid import cycles. It uses interface{} for subscriptions instead of concrete types to prevent
// circular dependencies when streams package needs test utilities.
type ImportCycleSafeTestSubscriber[T any] struct {
	Received     []T
	Completed    bool
	Errors       []error
	Done         chan struct{}
	doneClosed   bool
	mu           sync.Mutex
	ProcessDelay time.Duration // Optional delay to simulate slow processing
}

// NewImportCycleSafeTestSubscriber creates a new import cycle safe test subscriber
func NewImportCycleSafeTestSubscriber[T any]() *ImportCycleSafeTestSubscriber[T] {
	return &ImportCycleSafeTestSubscriber[T]{
		Received: make([]T, 0),
		Errors:   make([]error, 0),
		Done:     make(chan struct{}),
	}
}

// OnSubscribe handles subscription and auto-requests unlimited items
func (t *ImportCycleSafeTestSubscriber[T]) OnSubscribe(s interface {
	Request(n int64)
	Cancel()
}) {
	s.Request(int64(^uint(0) >> 1)) // Request max int64
}

// OnNext collects received values with optional processing delay
func (t *ImportCycleSafeTestSubscriber[T]) OnNext(value T) {
	if t.ProcessDelay > 0 {
		time.Sleep(t.ProcessDelay)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.Received = append(t.Received, value)
}

// OnError collects errors and signals completion
func (t *ImportCycleSafeTestSubscriber[T]) OnError(err error) {
	t.mu.Lock()
	t.Errors = append(t.Errors, err)
	if !t.doneClosed {
		t.doneClosed = true
		close(t.Done)
	}
	t.mu.Unlock()
}

// OnComplete marks completion and signals done
func (t *ImportCycleSafeTestSubscriber[T]) OnComplete() {
	t.mu.Lock()
	t.Completed = true
	if !t.doneClosed {
		t.doneClosed = true
		close(t.Done)
	}
	t.mu.Unlock()
}

// Wait blocks until completion or context cancellation
func (t *ImportCycleSafeTestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-t.Done:
		return
	case <-ctx.Done():
		return
	}
}

// Assertion methods

// AssertNoError checks if no errors occurred
func (t *ImportCycleSafeTestSubscriber[T]) AssertNoError(tst *testing.T) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.Errors) > 0 {
		tst.Errorf("Expected no errors, got %v", t.Errors)
	}
}

// AssertCompleted checks if subscriber completed successfully
func (t *ImportCycleSafeTestSubscriber[T]) AssertCompleted(tst *testing.T) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.Completed {
		tst.Error("Expected subscriber to complete")
	}
}

// AssertValues checks if received values match expected
func (t *ImportCycleSafeTestSubscriber[T]) AssertValues(tst *testing.T, expected []T) {
	t.mu.Lock()
	defer t.mu.Unlock()

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

// AssertError checks if at least one error occurred
func (t *ImportCycleSafeTestSubscriber[T]) AssertError(tst *testing.T) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.Errors) == 0 {
		tst.Error("Expected an error, but none occurred")
	}
}

// GetReceivedCopy returns a thread-safe copy of received values
func (t *ImportCycleSafeTestSubscriber[T]) GetReceivedCopy() []T {
	t.mu.Lock()
	defer t.mu.Unlock()
	result := make([]T, len(t.Received))
	copy(result, t.Received)
	return result
}

// GetErrorsCopy returns a thread-safe copy of received errors
func (t *ImportCycleSafeTestSubscriber[T]) GetErrorsCopy() []error {
	t.mu.Lock()
	defer t.mu.Unlock()
	result := make([]error, len(t.Errors))
	copy(result, t.Errors)
	return result
}

// IsCompleted returns completion status in a thread-safe way
func (t *ImportCycleSafeTestSubscriber[T]) IsCompleted() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.Completed
}

// ImportCycleSafeManualRequestSubscriber allows manual control over subscription demand while avoiding import cycles
type ImportCycleSafeManualRequestSubscriber[T any] struct {
	*ImportCycleSafeTestSubscriber[T]
	subscription interface {
		Request(n int64)
		Cancel()
	}
}

// NewImportCycleSafeManualRequestSubscriber creates a subscriber that requires explicit demand requests while avoiding import cycles
func NewImportCycleSafeManualRequestSubscriber[T any]() *ImportCycleSafeManualRequestSubscriber[T] {
	return &ImportCycleSafeManualRequestSubscriber[T]{
		ImportCycleSafeTestSubscriber: NewImportCycleSafeTestSubscriber[T](),
	}
}

// OnSubscribe stores the subscription for manual control (overrides LocalTestSubscriber)
func (s *ImportCycleSafeManualRequestSubscriber[T]) OnSubscribe(sub interface {
	Request(n int64)
	Cancel()
}) {
	s.subscription = sub
	// Don't auto-request unlimited items
}

// RequestItems manually requests n items to test backpressure control
func (s *ImportCycleSafeManualRequestSubscriber[T]) RequestItems(n int64) {
	if s.subscription != nil {
		s.subscription.Request(n)
	}
}

// CancelSubscription cancels the subscription
func (s *ImportCycleSafeManualRequestSubscriber[T]) CancelSubscription() {
	if s.subscription != nil {
		s.subscription.Cancel()
	}
}
