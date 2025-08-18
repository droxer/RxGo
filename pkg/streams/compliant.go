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

	p.mu.RLock()
	defer p.mu.RUnlock()

	emitted := false
	for sub := range p.subscribers {
		if sub.canEmit() && sub.decrementDemand() {
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

func (p *compliantPublisher[T]) complete() {
	if p.terminal.CompareAndSwap(false, true) {
		p.mu.Lock()
		defer p.mu.Unlock()

		for sub := range p.subscribers {
			if sub.markCompleted() {
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
	subscription := &compliantSubscription[T]{
		subscriber: sub,
		publisher:  pub,
	}
	return subscription
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

func (s *compliantSubscription[T]) canEmit() bool {
	return !s.cancelled.Load() && !s.completed.Load() && s.requested.Load() > 0
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
