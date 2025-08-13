package streams

import (
	"errors"
	"sync"
	"sync/atomic"
)

// compliantSubscription implements Reactive Streams 1.0.4 compliant subscription
type compliantSubscription[T any] struct {
	requested atomic.Int64
	cancelled atomic.Bool
	completed atomic.Bool
	error     atomic.Value // error
	subscriber Subscriber[T]
	publisher  *compliantPublisher[T]
	mu         sync.Mutex
}

func newCompliantSubscription[T any](sub Subscriber[T], pub *compliantPublisher[T]) *compliantSubscription[T] {
	return &compliantSubscription[T]{
		subscriber: sub,
		publisher:  pub,
	}
}

// Request implements Subscription.Request with Reactive Streams 1.0.4 compliance
func (s *compliantSubscription[T]) Request(n int64) {
	if n <= 0 {
		// Rule 3.9: Must signal IllegalArgumentException for non-positive requests
		s.subscriber.OnError(errors.New("non-positive subscription request"))
		return
	}

	if s.cancelled.Load() {
		// Rule 3.6: Must be a NOP if already cancelled
		return
	}

	// Atomically add to requested demand
	old := s.requested.Add(n)
	if old == 0 {
		// Signal publisher to start/resume processing
		select {
		case s.publisher.demandSignal <- struct{}{}:
		default:
		}
	}
}

// Cancel implements Subscription.Cancel with Reactive Streams 1.0.4 compliance
func (s *compliantSubscription[T]) Cancel() {
	if s.cancelled.CompareAndSwap(false, true) {
		// Rule 3.12-3.13: Cleanup resources
		s.publisher.removeSubscription(s)
	}
}

// isCancelled checks if the subscription has been cancelled
func (s *compliantSubscription[T]) isCancelled() bool {
	return s.cancelled.Load()
}

// isCompleted checks if the subscription has completed
func (s *compliantSubscription[T]) isCompleted() bool {
	return s.completed.Load()
}

// markCompleted marks the subscription as completed
func (s *compliantSubscription[T]) markCompleted() bool {
	return s.completed.CompareAndSwap(false, true)
}

// markError marks the subscription with an error
func (s *compliantSubscription[T]) markError(err error) {
	s.error.Store(err)
}

// canEmit checks if we can emit an item based on demand
func (s *compliantSubscription[T]) canEmit() bool {
	return !s.cancelled.Load() && !s.completed.Load() && s.requested.Load() > 0
}

// decrementDemand safely decrements the requested count
func (s *compliantSubscription[T]) decrementDemand() bool {
	for {
		current := s.requested.Load()
		if current <= 0 {
			return false
		}
		if s.requested.CompareAndSwap(current, current-1) {
			return true
		}
	}
}