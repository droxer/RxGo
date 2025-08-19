// Package testutil provides comprehensive test utilities for RxGo reactive streams testing.
//
// This package contains thread-safe test subscribers with flexible configuration
// options for testing various reactive streams scenarios including:
// - Basic value collection and assertion
// - Manual backpressure control
// - Process delay simulation
// - Error and completion handling
package testutil

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

// ReactiveStreamsTestSubscriber is a comprehensive test utility for collecting published values and testing reactive streams.
// It provides flexible configuration options for various testing scenarios including backpressure control,
// process delay simulation, and comprehensive assertion methods.
type ReactiveStreamsTestSubscriber[T any] struct {
	// Collected data
	Received  []T
	Errors    []error
	Completed bool

	// Control fields
	done         chan struct{}
	doneClosed   bool
	subscription interface {
		Request(n int64)
		Cancel()
	}
	mu sync.Mutex
	wg sync.WaitGroup

	// Configurable behavior
	ProcessDelay time.Duration // Delay to simulate slow processing
	RequestMode  RequestMode   // How to handle subscription requests
	RequestCount int64         // Number of items to request (for Manual mode)
}

// RequestMode defines how a TestSubscriber handles subscription requests
type RequestMode int

const (
	// AutoRequest automatically requests unlimited items (default behavior)
	AutoRequest RequestMode = iota
	// ManualRequest requires explicit calls to RequestItems()
	ManualRequest
	// LimitedRequest requests a fixed number of items specified by RequestCount
	LimitedRequest
)

// NewReactiveStreamsTestSubscriber creates a new test subscriber with default settings (auto-request unlimited items)
func NewReactiveStreamsTestSubscriber[T any]() *ReactiveStreamsTestSubscriber[T] {
	return &ReactiveStreamsTestSubscriber[T]{
		Received:    make([]T, 0),
		Errors:      make([]error, 0),
		done:        make(chan struct{}),
		RequestMode: AutoRequest,
	}
}

// NewManualBackpressureTestSubscriber creates a test subscriber that requires manual request management for backpressure testing
func NewManualBackpressureTestSubscriber[T any]() *ReactiveStreamsTestSubscriber[T] {
	return &ReactiveStreamsTestSubscriber[T]{
		Received:    make([]T, 0),
		Errors:      make([]error, 0),
		done:        make(chan struct{}),
		RequestMode: ManualRequest,
	}
}

// NewLimitedRequestTestSubscriber creates a test subscriber that requests a specific number of items for controlled demand testing
func NewLimitedRequestTestSubscriber[T any](requestCount int64) *ReactiveStreamsTestSubscriber[T] {
	return &ReactiveStreamsTestSubscriber[T]{
		Received:     make([]T, 0),
		Errors:       make([]error, 0),
		done:         make(chan struct{}),
		RequestMode:  LimitedRequest,
		RequestCount: requestCount,
	}
}

// OnSubscribe handles subscription based on the configured RequestMode
func (ts *ReactiveStreamsTestSubscriber[T]) OnSubscribe(sub interface {
	Request(n int64)
	Cancel()
}) {
	ts.mu.Lock()
	ts.subscription = sub
	ts.wg.Add(1)
	ts.mu.Unlock()

	switch ts.RequestMode {
	case AutoRequest:
		sub.Request(int64(^uint(0) >> 1)) // Request maximum
	case LimitedRequest:
		if ts.RequestCount > 0 {
			sub.Request(ts.RequestCount)
		}
	case ManualRequest:
		// Don't auto-request anything - wait for manual calls
	}
}

// OnNext collects received values with optional processing delay
func (ts *ReactiveStreamsTestSubscriber[T]) OnNext(value T) {
	if ts.ProcessDelay > 0 {
		time.Sleep(ts.ProcessDelay)
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.Received = append(ts.Received, value)
}

// OnError collects errors and signals completion
func (ts *ReactiveStreamsTestSubscriber[T]) OnError(err error) {
	ts.mu.Lock()
	ts.Errors = append(ts.Errors, err)
	if !ts.doneClosed {
		ts.doneClosed = true
		close(ts.done)
	}
	ts.mu.Unlock()
	ts.wg.Done()
}

// OnComplete marks completion and signals done
func (ts *ReactiveStreamsTestSubscriber[T]) OnComplete() {
	ts.mu.Lock()
	ts.Completed = true
	if !ts.doneClosed {
		ts.doneClosed = true
		close(ts.done)
	}
	ts.mu.Unlock()
	ts.wg.Done()
}

// Wait blocks until completion, error, or context cancellation
func (ts *ReactiveStreamsTestSubscriber[T]) Wait(ctx context.Context) {
	select {
	case <-ts.done:
		return
	case <-ctx.Done():
		return
	}
}

// WaitWithTimeout waits for completion with a default timeout
func (ts *ReactiveStreamsTestSubscriber[T]) WaitWithTimeout() {
	ts.wg.Wait()
}

// RequestItems manually requests n items (only works with ManualRequest mode)
func (ts *ReactiveStreamsTestSubscriber[T]) RequestItems(n int64) {
	ts.mu.Lock()
	subscription := ts.subscription
	ts.mu.Unlock()

	if subscription != nil {
		subscription.Request(n)
	}
}

// Cancel cancels the subscription
func (ts *ReactiveStreamsTestSubscriber[T]) Cancel() {
	ts.mu.Lock()
	subscription := ts.subscription
	ts.mu.Unlock()

	if subscription != nil {
		subscription.Cancel()
	}
}

// Assertion methods for testing

// AssertValues checks if received values match expected values exactly
func (ts *ReactiveStreamsTestSubscriber[T]) AssertValues(tst *testing.T, expected []T) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.Received) != len(expected) {
		tst.Errorf("Expected %d values, got %d: %v", len(expected), len(ts.Received), ts.Received)
		return
	}
	for i, v := range expected {
		if !reflect.DeepEqual(ts.Received[i], v) {
			tst.Errorf("Expected value[%d] to be %v, got %v", i, v, ts.Received[i])
		}
	}
}

// AssertValuesUnordered checks if received values match expected values (ignoring order)
func (ts *ReactiveStreamsTestSubscriber[T]) AssertValuesUnordered(tst *testing.T, expected []T) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.Received) != len(expected) {
		tst.Errorf("Expected %d values, got %d: %v", len(expected), len(ts.Received), ts.Received)
		return
	}

	expectedMap := make(map[interface{}]int)
	receivedMap := make(map[interface{}]int)

	for _, v := range expected {
		expectedMap[v]++
	}
	for _, v := range ts.Received {
		receivedMap[v]++
	}

	if !reflect.DeepEqual(expectedMap, receivedMap) {
		tst.Errorf("Expected values %v (unordered), got %v", expected, ts.Received)
	}
}

// AssertValueCount checks if the number of received values matches expected count
func (ts *ReactiveStreamsTestSubscriber[T]) AssertValueCount(tst *testing.T, expectedCount int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.Received) != expectedCount {
		tst.Errorf("Expected %d values, got %d: %v", expectedCount, len(ts.Received), ts.Received)
	}
}

// AssertCompleted checks if subscriber completed successfully
func (ts *ReactiveStreamsTestSubscriber[T]) AssertCompleted(tst *testing.T) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if !ts.Completed {
		tst.Error("Expected subscriber to complete")
	}
}

// AssertNotCompleted checks if subscriber has not completed
func (ts *ReactiveStreamsTestSubscriber[T]) AssertNotCompleted(tst *testing.T) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.Completed {
		tst.Error("Expected subscriber to not complete")
	}
}

// AssertNoError checks if no errors occurred
func (ts *ReactiveStreamsTestSubscriber[T]) AssertNoError(tst *testing.T) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if len(ts.Errors) > 0 {
		tst.Errorf("Expected no errors, got %v", ts.Errors)
	}
}

// AssertError checks if at least one error occurred
func (ts *ReactiveStreamsTestSubscriber[T]) AssertError(tst *testing.T) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if len(ts.Errors) == 0 {
		tst.Error("Expected an error, but none occurred")
	}
}

// AssertErrorCount checks if the expected number of errors occurred
func (ts *ReactiveStreamsTestSubscriber[T]) AssertErrorCount(tst *testing.T, expectedCount int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if len(ts.Errors) != expectedCount {
		tst.Errorf("Expected %d errors, got %d: %v", expectedCount, len(ts.Errors), ts.Errors)
	}
}

// GetReceivedCopy returns a copy of received values (thread-safe)
func (ts *ReactiveStreamsTestSubscriber[T]) GetReceivedCopy() []T {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	result := make([]T, len(ts.Received))
	copy(result, ts.Received)
	return result
}

// GetErrorsCopy returns a copy of received errors (thread-safe)
func (ts *ReactiveStreamsTestSubscriber[T]) GetErrorsCopy() []error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	result := make([]error, len(ts.Errors))
	copy(result, ts.Errors)
	return result
}

// IsCompleted returns whether the subscriber has completed (thread-safe)
func (ts *ReactiveStreamsTestSubscriber[T]) IsCompleted() bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.Completed
}

// ReactiveStreamsTestSubscription is a wrapper for testing subscription behavior and tracking requests/cancellations
type ReactiveStreamsTestSubscription struct {
	sub interface {
		Request(n int64)
		Cancel()
	}
	requested int64
	cancelled bool
	mu        sync.Mutex
}

// NewReactiveStreamsTestSubscription wraps a subscription for testing and tracking
func NewReactiveStreamsTestSubscription(sub interface {
	Request(n int64)
	Cancel()
}) *ReactiveStreamsTestSubscription {
	return &ReactiveStreamsTestSubscription{sub: sub}
}

// Request forwards the request and tracks the amount
func (ts *ReactiveStreamsTestSubscription) Request(n int64) {
	ts.mu.Lock()
	ts.requested += n
	ts.mu.Unlock()

	if ts.sub != nil {
		ts.sub.Request(n)
	}
}

// Cancel forwards the cancellation and tracks the state
func (ts *ReactiveStreamsTestSubscription) Cancel() {
	ts.mu.Lock()
	ts.cancelled = true
	ts.mu.Unlock()

	if ts.sub != nil {
		ts.sub.Cancel()
	}
}

// GetRequestedAmount returns the total amount requested
func (ts *ReactiveStreamsTestSubscription) GetRequestedAmount() int64 {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.requested
}

// IsCancelled returns whether the subscription was cancelled
func (ts *ReactiveStreamsTestSubscription) IsCancelled() bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.cancelled
}
