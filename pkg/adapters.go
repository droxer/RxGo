package pkg

import (
	"context"
	"fmt"
	"sync"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/streams"
)

func ObservablePublisherAdapter[T any](obs *observable.Observable[T]) streams.Publisher[T] {
	return streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[T]) {
		err := obs.Subscribe(ctx, &observableSubscriberAdapter[T]{sub: sub})
		if err != nil {
			sub.OnError(err)
		}
	})
}

type observableSubscriberAdapter[T any] struct {
	sub streams.Subscriber[T]
}

func (a *observableSubscriberAdapter[T]) Start() {
	a.sub.OnSubscribe(&observableSubscription{})
}

func (a *observableSubscriberAdapter[T]) OnNext(t T) {
	a.sub.OnNext(t)
}

func (a *observableSubscriberAdapter[T]) OnError(err error) {
	a.sub.OnError(err)
}

func (a *observableSubscriberAdapter[T]) OnComplete() {
	a.sub.OnComplete()
}

type observableSubscription struct{}

func (r *observableSubscription) Request(n int64) {
}

func (r *observableSubscription) Cancel() {
}

func PublisherToObservableAdapter[T any](publisher streams.Publisher[T]) *observable.Observable[T] {
	return observable.Create(func(ctx context.Context, sub observable.Subscriber[T]) {
		adapter := &publisherSubscriberAdapter[T]{
			sub:              sub,
			initialBatchSize: 128,  // Default batch size
			lowWaterMark:     64,   // Request more when outstanding drops below this
			highWaterMark:    256,  // Maximum outstanding requests
			maxBufferSize:    1024, // Maximum buffer size to prevent unbounded growth
		}
		err := publisher.Subscribe(ctx, adapter)
		if err != nil {
			sub.OnError(err)
		}
	})
}

type publisherSubscriberAdapter[T any] struct {
	sub          observable.Subscriber[T]
	subscription streams.Subscription
	outstanding  int64
	buffered     int64
	mu           sync.Mutex
	closed       bool
	// Configuration for backpressure handling
	initialBatchSize int64
	lowWaterMark     int64
	highWaterMark    int64
	maxBufferSize    int64
}

func (a *publisherSubscriberAdapter[T]) OnSubscribe(s streams.Subscription) {
	a.subscription = s
	a.sub.Start()
	// Request an initial batch based on configuration
	a.requestMore(a.initialBatchSize)
}

func (a *publisherSubscriberAdapter[T]) OnNext(t T) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	// Check for buffer overflow
	if a.buffered >= a.maxBufferSize {
		// Buffer overflow - cancel subscription and notify error
		a.closed = true
		if a.subscription != nil {
			a.subscription.Cancel()
		}
		a.sub.OnError(fmt.Errorf("buffer overflow: maximum buffer size of %d exceeded", a.maxBufferSize))
		return
	}

	a.sub.OnNext(t)
	a.outstanding--
	a.buffered++

	// Request more if we're running low but haven't exceeded high water mark
	if a.outstanding < a.lowWaterMark && !a.closed && (a.outstanding+a.initialBatchSize) <= a.highWaterMark {
		a.requestMore(a.initialBatchSize)
	}
}

func (a *publisherSubscriberAdapter[T]) OnError(err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	a.closed = true
	a.buffered = 0 // Reset buffer count on error
	a.sub.OnError(err)
}

func (a *publisherSubscriberAdapter[T]) OnComplete() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	a.closed = true
	a.buffered = 0 // Reset buffer count on completion
	a.sub.OnComplete()
}

func (a *publisherSubscriberAdapter[T]) requestMore(n int64) {
	if a.subscription != nil && !a.closed {
		a.outstanding += n
		a.subscription.Request(n)
	}
}
