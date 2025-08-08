package publisher

import (
	"context"
	"sync"
	"sync/atomic"
)

// SubscriptionState represents the state of a subscription
type SubscriptionState int

const (
	NoSubscription SubscriptionState = iota
	Active
	Cancelled
	Completed
)

// ReactiveSubscription implements Subscription interface
type ReactiveSubscription struct {
	ctx       context.Context
	cancel    context.CancelFunc
	requested int64        // atomic counter for requested items
	state     atomic.Value // SubscriptionState
	mu        sync.RWMutex
	onRequest func(n int64)
	onCancel  func()
}

// NewReactiveSubscription creates a new ReactiveSubscription
func NewReactiveSubscription(ctx context.Context, onRequest func(n int64), onCancel func()) Subscription {
	ctx, cancel := context.WithCancel(ctx)
	sub := &ReactiveSubscription{
		ctx:       ctx,
		cancel:    cancel,
		onRequest: onRequest,
		onCancel:  onCancel,
	}
	sub.state.Store(NoSubscription)
	return sub
}

func (s *ReactiveSubscription) Request(n int64) {
	// Rule 3.9: n must be positive
	if n <= 0 {
		// Rule 3.9: IllegalArgumentException
		return
	}

	// Rule 3.6: Subscription state check
	currentState := s.state.Load().(SubscriptionState)
	if currentState == Cancelled || currentState == Completed {
		return
	}

	// Atomic addition to requested count
	atomic.AddInt64(&s.requested, n)

	// Notify publisher of demand
	if s.onRequest != nil {
		s.onRequest(n)
	}
}

func (s *ReactiveSubscription) Cancel() {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentState := s.state.Load().(SubscriptionState)
	if currentState == Cancelled {
		return
	}

	s.state.Store(Cancelled)
	s.cancel()

	if s.onCancel != nil {
		s.onCancel()
	}
}

func (s *ReactiveSubscription) GetRequested() int64 {
	return atomic.LoadInt64(&s.requested)
}

func (s *ReactiveSubscription) DecrementRequested() int64 {
	return atomic.AddInt64(&s.requested, -1)
}
