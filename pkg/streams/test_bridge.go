package streams

import (
	"context"
	"testing"

	"github.com/droxer/RxGo/internal/testutil"
)

// localTestSubscriber wraps testutil.ImportCycleSafeTestSubscriber to work with concrete Subscription types
type localTestSubscriber[T any] struct {
	inner *testutil.ImportCycleSafeTestSubscriber[T]
}

// newLocalTestSubscriber creates a local test subscriber for use in streams tests
func newLocalTestSubscriber[T any]() *localTestSubscriber[T] {
	return &localTestSubscriber[T]{
		inner: testutil.NewImportCycleSafeTestSubscriber[T](),
	}
}

// OnSubscribe converts the concrete Subscription to interface and delegates
func (t *localTestSubscriber[T]) OnSubscribe(sub Subscription) {
	t.inner.OnSubscribe(sub)
}

// OnNext delegates to inner
func (t *localTestSubscriber[T]) OnNext(value T) {
	t.inner.OnNext(value)
}

// OnError delegates to inner
func (t *localTestSubscriber[T]) OnError(err error) {
	t.inner.OnError(err)
}

// OnComplete delegates to inner
func (t *localTestSubscriber[T]) OnComplete() {
	t.inner.OnComplete()
}

// Wait delegates to inner
func (t *localTestSubscriber[T]) Wait(ctx context.Context) {
	t.inner.Wait(ctx)
}

// AssertValues delegates to inner
func (t *localTestSubscriber[T]) AssertValues(tst *testing.T, expected []T) {
	t.inner.AssertValues(tst, expected)
}

// AssertCompleted delegates to inner
func (t *localTestSubscriber[T]) AssertCompleted(tst *testing.T) {
	t.inner.AssertCompleted(tst)
}

// AssertNoError delegates to inner
func (t *localTestSubscriber[T]) AssertNoError(tst *testing.T) {
	t.inner.AssertNoError(tst)
}

// AssertError delegates to inner
func (t *localTestSubscriber[T]) AssertError(tst *testing.T) {
	t.inner.AssertError(tst)
}

// Received exposes the received values
func (t *localTestSubscriber[T]) Received() []T {
	return t.inner.GetReceivedCopy()
}

// Errors exposes the received errors
func (t *localTestSubscriber[T]) Errors() []error {
	return t.inner.GetErrorsCopy()
}

// Completed exposes the completion status
func (t *localTestSubscriber[T]) Completed() bool {
	return t.inner.IsCompleted()
}

// localManualRequestSubscriber wraps testutil.ImportCycleSafeManualRequestSubscriber
type localManualRequestSubscriber[T any] struct {
	inner *testutil.ImportCycleSafeManualRequestSubscriber[T]
}

// newLocalManualRequestSubscriber creates a manual request subscriber
func newLocalManualRequestSubscriber[T any]() *localManualRequestSubscriber[T] {
	return &localManualRequestSubscriber[T]{
		inner: testutil.NewImportCycleSafeManualRequestSubscriber[T](),
	}
}

// OnSubscribe converts the concrete Subscription to interface and delegates
func (s *localManualRequestSubscriber[T]) OnSubscribe(sub Subscription) {
	s.inner.OnSubscribe(sub)
}

// OnNext delegates to inner
func (s *localManualRequestSubscriber[T]) OnNext(value T) {
	s.inner.OnNext(value)
}

// OnError delegates to inner
func (s *localManualRequestSubscriber[T]) OnError(err error) {
	s.inner.OnError(err)
}

// OnComplete delegates to inner
func (s *localManualRequestSubscriber[T]) OnComplete() {
	s.inner.OnComplete()
}

// Wait delegates to inner
func (s *localManualRequestSubscriber[T]) Wait(ctx context.Context) {
	s.inner.Wait(ctx)
}

// RequestItems delegates to inner
func (s *localManualRequestSubscriber[T]) RequestItems(n int64) {
	s.inner.RequestItems(n)
}

// CancelSubscription delegates to inner
func (s *localManualRequestSubscriber[T]) CancelSubscription() {
	s.inner.CancelSubscription()
}

// AssertValues delegates to inner
func (s *localManualRequestSubscriber[T]) AssertValues(tst *testing.T, expected []T) {
	s.inner.AssertValues(tst, expected)
}

// AssertCompleted delegates to inner
func (s *localManualRequestSubscriber[T]) AssertCompleted(tst *testing.T) {
	s.inner.AssertCompleted(tst)
}

// AssertNoError delegates to inner
func (s *localManualRequestSubscriber[T]) AssertNoError(tst *testing.T) {
	s.inner.AssertNoError(tst)
}

// AssertError delegates to inner
func (s *localManualRequestSubscriber[T]) AssertError(tst *testing.T) {
	s.inner.AssertError(tst)
}

// Received exposes the received values
func (s *localManualRequestSubscriber[T]) Received() []T {
	return s.inner.GetReceivedCopy()
}

// Errors exposes the received errors
func (s *localManualRequestSubscriber[T]) Errors() []error {
	return s.inner.GetErrorsCopy()
}

// Completed exposes the completion status
func (s *localManualRequestSubscriber[T]) Completed() bool {
	return s.inner.IsCompleted()
}
