package streams

import (
	"context"
	"errors"
	"sync"
)

// BufferedPublisher implements backpressure strategies with configurable overflow behavior
type BufferedPublisher[T any] struct {
	config     BackpressureConfig
	source     func(ctx context.Context, sub Subscriber[T])
	buffer     []T
	mu         sync.RWMutex
	bufferFull chan struct{}
}

func NewBufferedPublisher[T any](
	config BackpressureConfig,
	source func(ctx context.Context, sub Subscriber[T]),
) Publisher[T] {
	return &BufferedPublisher[T]{
		config:     config,
		source:     source,
		buffer:     make([]T, 0, config.BufferSize),
		bufferFull: make(chan struct{}),
	}
}

func (bp *BufferedPublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	if sub == nil {
		panic("subscriber cannot be nil")
	}

	subscription := &bufferedSubscription[T]{
		publisher: bp,
		sub:       sub,
		ctx:       ctx,
	}
	sub.OnSubscribe(subscription)
}

// bufferedSubscription implements Subscription
type bufferedSubscription[T any] struct {
	publisher *BufferedPublisher[T]
	sub       Subscriber[T]
	ctx       context.Context
	mu        sync.Mutex
	requested int64
	cancelled bool
}

func (bs *bufferedSubscription[T]) Request(n int64) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.cancelled || n <= 0 {
		return
	}

	bs.requested += n
	bs.publisher.processWithStrategy(bs.ctx, bs)
}

func (bs *bufferedSubscription[T]) Cancel() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.cancelled = true
}

func (bp *BufferedPublisher[T]) processWithStrategy(ctx context.Context, sub *bufferedSubscription[T]) {
	go func() {
		bp.source(ctx, &strategySubscriber[T]{
			publisher: bp,
			sub:       sub,
		})
	}()
}

type strategySubscriber[T any] struct {
	publisher *BufferedPublisher[T]
	sub       *bufferedSubscription[T]
}

func (ss *strategySubscriber[T]) OnSubscribe(s Subscription) {
	// Not used in this context
}

func (ss *strategySubscriber[T]) OnNext(value T) {
	ss.publisher.handleNext(value, ss.sub)
}

func (ss *strategySubscriber[T]) OnError(err error) {
	ss.sub.sub.OnError(err)
}

func (ss *strategySubscriber[T]) OnComplete() {
	ss.publisher.flushBuffer(ss.sub)
	ss.sub.sub.OnComplete()
}

func (bp *BufferedPublisher[T]) handleNext(value T, sub *bufferedSubscription[T]) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	sub.mu.Lock()
	defer sub.mu.Unlock()

	if sub.cancelled {
		return
	}

	switch bp.config.Strategy {
	case Buffer:
		bp.handleBufferStrategy(value, sub)
	case Drop:
		bp.handleDropStrategy(value, sub)
	case Latest:
		bp.handleLatestStrategy(value, sub)
	case Error:
		bp.handleErrorStrategy(value, sub)
	}
}

func (bp *BufferedPublisher[T]) handleBufferStrategy(value T, sub *bufferedSubscription[T]) {
	if len(bp.buffer) < int(bp.config.BufferSize) {
		bp.buffer = append(bp.buffer, value)
	}

	for len(bp.buffer) > 0 && sub.requested > 0 {
		next := bp.buffer[0]
		bp.buffer = bp.buffer[1:]
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) handleDropStrategy(value T, sub *bufferedSubscription[T]) {
	if int64(len(bp.buffer)) >= bp.config.BufferSize {
		// Drop the new item
		return
	}

	bp.buffer = append(bp.buffer, value)
	for len(bp.buffer) > 0 && sub.requested > 0 {
		next := bp.buffer[0]
		bp.buffer = bp.buffer[1:]
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) handleLatestStrategy(value T, sub *bufferedSubscription[T]) {
	if int64(len(bp.buffer)) >= bp.config.BufferSize {
		// Replace oldest with latest
		if len(bp.buffer) > 0 {
			bp.buffer[0] = value
		} else {
			bp.buffer = append(bp.buffer, value)
		}
	} else {
		bp.buffer = append(bp.buffer, value)
	}

	for len(bp.buffer) > 0 && sub.requested > 0 {
		next := bp.buffer[0]
		bp.buffer = bp.buffer[1:]
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) handleErrorStrategy(value T, sub *bufferedSubscription[T]) {
	if int64(len(bp.buffer)) >= bp.config.BufferSize {
		sub.sub.OnError(errors.New("buffer overflow: buffer size limit exceeded"))
		return
	}

	bp.buffer = append(bp.buffer, value)
	for len(bp.buffer) > 0 && sub.requested > 0 {
		next := bp.buffer[0]
		bp.buffer = bp.buffer[1:]
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) flushBuffer(sub *bufferedSubscription[T]) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	sub.mu.Lock()
	defer sub.mu.Unlock()

	for len(bp.buffer) > 0 && sub.requested > 0 {
		next := bp.buffer[0]
		bp.buffer = bp.buffer[1:]
		sub.requested--
		sub.sub.OnNext(next)
	}

	// Handle remaining buffer based on strategy
	if len(bp.buffer) > 0 {
		switch bp.config.Strategy {
		case Buffer, Drop, Latest:
			// Process remaining items
			for len(bp.buffer) > 0 {
				next := bp.buffer[0]
				bp.buffer = bp.buffer[1:]
				sub.sub.OnNext(next)
			}
		case Error:
			// Signal overflow for remaining items
			sub.sub.OnError(errors.New("buffer overflow during completion"))
		}
	}
}
