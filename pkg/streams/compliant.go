package streams

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// compliantPublisher provides Reactive Streams 1.0.4 compliant base for publishers
type compliantPublisher[T any] struct {
	mu           sync.RWMutex
	subscribers  map[*compliantSubscription[T]]struct{}
	terminal     atomic.Bool
	demandSignal chan struct{}
}

func newCompliantPublisher[T any]() *compliantPublisher[T] {
	return &compliantPublisher[T]{
		subscribers:  make(map[*compliantSubscription[T]]struct{}),
		demandSignal: make(chan struct{}, 1),
	}
}

// subscribe adds a compliant subscriber
func (p *compliantPublisher[T]) subscribe(ctx context.Context, sub Subscriber[T]) {
	if sub == nil {
		return
	}

	// Rule 1.9: Must reject null subscribers
	if ctx == nil {
		ctx = context.Background()
	}

	// Check if already terminated
	if p.terminal.Load() {
		return
	}

	subscription := newCompliantSubscription(sub, p)

	p.mu.Lock()
	if p.terminal.Load() {
		p.mu.Unlock()
		return
	}
	p.subscribers[subscription] = struct{}{}
	p.mu.Unlock()

	// Signal onSubscribe asynchronously
	go func() {
		select {
		case <-ctx.Done():
			// Context cancelled before subscription
			return
		default:
			// Rule 1.5: Must call onSubscribe before any other signals
			sub.OnSubscribe(subscription)
		}
	}()
}

// emit safely emits an item to all subscribers with demand
func (p *compliantPublisher[T]) emit(item T) bool {
	if p.terminal.Load() {
		return false
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	// Rule 1.1: Must not emit more than requested
	emitted := false
	for sub := range p.subscribers {
		if sub.canEmit() && sub.decrementDemand() {
			// Rule 1.3: Sequential signaling
			sub.mu.Lock()
			if !sub.isCancelled() {
				sub.subscriber.OnNext(item)
				emitted = true
			}
			sub.mu.Unlock()
		}
	}

	return emitted
}

// complete signals completion to all subscribers
func (p *compliantPublisher[T]) complete() {
	if p.terminal.CompareAndSwap(false, true) {
		p.mu.Lock()
		defer p.mu.Unlock()

		for sub := range p.subscribers {
			if sub.markCompleted() {
				// Rule 1.7: onComplete must be the last signal
				sub.mu.Lock()
				if !sub.isCancelled() {
					sub.subscriber.OnComplete()
				}
				sub.mu.Unlock()
			}
		}
		p.subscribers = make(map[*compliantSubscription[T]]struct{})
	}
}

// removeSubscription removes a subscription
func (p *compliantPublisher[T]) removeSubscription(sub *compliantSubscription[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.subscribers, sub)
}

// getActiveSubscribers returns a snapshot of active subscriptions
func (p *compliantPublisher[T]) getActiveSubscribers() []*compliantSubscription[T] {
	p.mu.RLock()
	defer p.mu.RUnlock()

	subs := make([]*compliantSubscription[T], 0, len(p.subscribers))
	for sub := range p.subscribers {
		subs = append(subs, sub)
	}
	return subs
}

// compliantSubscription implements Reactive Streams 1.0.4 compliant subscription
type compliantSubscription[T any] struct {
	requested atomic.Int64
	cancelled atomic.Bool
	completed atomic.Bool
	// error field is reserved for future use
	subscriber Subscriber[T]
	publisher  *compliantPublisher[T]
	mu         sync.Mutex
}

func newCompliantSubscription[T any](sub Subscriber[T], pub *compliantPublisher[T]) *compliantSubscription[T] {
	fmt.Printf("newCompliantSubscription called\n")
	subscription := &compliantSubscription[T]{
		subscriber: sub,
		publisher:  pub,
	}
	fmt.Printf("newCompliantSubscription: initial requested=%d\n", subscription.requested.Load())
	return subscription
}

// Request implements Subscription.Request with Reactive Streams 1.0.4 compliance
func (s *compliantSubscription[T]) Request(n int64) {
	fmt.Printf("compliantSubscription.Request: %d\n", n)
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
	fmt.Printf("compliantSubscription.Request: before Add, requested=%d\n", s.requested.Load())
	old := s.requested.Add(n)
	fmt.Printf("compliantSubscription.Request: after Add, old=%d, requested=%d\n", old, s.requested.Load())
	if old == 0 {
		// Signal publisher to start/resume processing
		fmt.Printf("compliantSubscription.Request: sending demand signal\n")
		select {
		case s.publisher.demandSignal <- struct{}{}:
			fmt.Printf("compliantSubscription.Request: demand signal sent successfully\n")
		default:
			fmt.Printf("compliantSubscription.Request: demand signal send failed (channel full)\n")
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

// markCompleted marks the subscription as completed
func (s *compliantSubscription[T]) markCompleted() bool {
	return s.completed.CompareAndSwap(false, true)
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
