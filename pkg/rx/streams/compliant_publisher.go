package streams

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// compliantPublisher provides Reactive Streams 1.0.4 compliant base for publishers
type compliantPublisher[T any] struct {
	mu           sync.RWMutex
	subscribers  map[*compliantSubscription[T]]struct{}
	terminal     atomic.Bool
	shutdown     atomic.Bool
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

// error signals an error to all subscribers
func (p *compliantPublisher[T]) error(err error) {
	if err == nil {
		err = errors.New("publisher error")
	}

	if p.terminal.CompareAndSwap(false, true) {
		p.mu.Lock()
		defer p.mu.Unlock()

		for sub := range p.subscribers {
			sub.markError(err)
			// Rule 1.7: onError must be the last signal
			sub.mu.Lock()
			if !sub.isCancelled() {
				sub.subscriber.OnError(err)
			}
			sub.mu.Unlock()
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

// hasSubscribers checks if there are any active subscribers
func (p *compliantPublisher[T]) hasSubscribers() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.subscribers) > 0
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