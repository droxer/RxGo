package streams

import (
	"errors"
	"sync"
	"sync/atomic"
)

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

func (p *compliantPublisher[T]) emit(item T) bool {
	if p.terminal.Load() {
		return false
	}

	// Get a snapshot of subscribers to avoid deadlocks
	p.mu.RLock()
	subs := make([]*compliantSubscription[T], 0, len(p.subscribers))
	for sub := range p.subscribers {
		subs = append(subs, sub)
	}
	p.mu.RUnlock()

	emitted := false

	for _, sub := range subs {
		// Atomically check if we can emit and decrement demand in one operation
		if sub.tryEmitAndDecrement() {
			sub.mu.Lock()
			if !sub.isCancelled() {
				// Release lock before calling subscriber to avoid deadlocks
				subscriber := sub.subscriber
				sub.mu.Unlock()
				subscriber.OnNext(item)
				emitted = true
			} else {
				sub.mu.Unlock()
			}
		}
	}

	return emitted
}

func (p *compliantPublisher[T]) complete() {
	if p.terminal.CompareAndSwap(false, true) {
		p.mu.Lock()
		defer p.mu.Unlock()

		subs := make([]*compliantSubscription[T], 0, len(p.subscribers))
		for sub := range p.subscribers {
			subs = append(subs, sub)
		}

		// Clear the subscribers map immediately to prevent new operations
		p.subscribers = make(map[*compliantSubscription[T]]struct{})

		// Notify subscribers outside the lock to avoid deadlocks
		for _, sub := range subs {
			if sub.markCompleted() {
				sub.mu.Lock()
				if !sub.isCancelled() {
					subscriber := sub.subscriber
					sub.mu.Unlock()
					subscriber.OnComplete()
				} else {
					sub.mu.Unlock()
				}
			}
		}
	}
}

func (p *compliantPublisher[T]) removeSubscription(sub *compliantSubscription[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.subscribers, sub)
}

func (p *compliantPublisher[T]) getActiveSubscribers() []*compliantSubscription[T] {
	p.mu.RLock()
	defer p.mu.RUnlock()

	subs := make([]*compliantSubscription[T], 0, len(p.subscribers))
	for sub := range p.subscribers {
		subs = append(subs, sub)
	}
	return subs
}

type compliantSubscription[T any] struct {
	requested  atomic.Int64
	cancelled  atomic.Bool
	completed  atomic.Bool
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

func (s *compliantSubscription[T]) Request(n int64) {
	if n <= 0 {
		s.subscriber.OnError(errors.New("non-positive subscription request"))
		return
	}

	if s.cancelled.Load() {
		return
	}

	old := s.requested.Add(n)
	if old == 0 {
		select {
		case s.publisher.demandSignal <- struct{}{}:
		default:
		}
	}
}

func (s *compliantSubscription[T]) Cancel() {
	if s.cancelled.CompareAndSwap(false, true) {
		s.publisher.removeSubscription(s)
	}
}

func (s *compliantSubscription[T]) isCancelled() bool {
	return s.cancelled.Load()
}

func (s *compliantSubscription[T]) markCompleted() bool {
	return s.completed.CompareAndSwap(false, true)
}

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

// tryEmitAndDecrement atomically checks if the subscription can emit and decrements demand
func (s *compliantSubscription[T]) tryEmitAndDecrement() bool {
	// First check if the subscription is in a state where it can emit
	// without acquiring the mutex to avoid unnecessary locking
	if s.cancelled.Load() || s.completed.Load() {
		return false
	}

	// Try to decrement demand atomically
	return s.decrementDemand()
}
