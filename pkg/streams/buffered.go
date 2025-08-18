package streams

import (
	"context"
	"errors"
	"fmt"
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

	if len(bp.buffer) > 0 {
		switch bp.config.Strategy {
		case Buffer, Drop, Latest:
			for len(bp.buffer) > 0 {
				next := bp.buffer[0]
				bp.buffer = bp.buffer[1:]
				sub.sub.OnNext(next)
			}
		case Error:
			sub.sub.OnError(errors.New("buffer overflow during completion"))
		}
	}
}

// CompliantRangePublisher implements a Reactive Streams 1.0.4 compliant range publisher
type CompliantRangePublisher struct {
	*compliantPublisher[int]
	start int
	end   int
}

func NewCompliantRangePublisher(start, end int) *CompliantRangePublisher {
	return &CompliantRangePublisher{
		compliantPublisher: newCompliantPublisher[int](),
		start:              start,
		end:                end,
	}
}

func (rp *CompliantRangePublisher) Subscribe(ctx context.Context, sub Subscriber[int]) {
	rp.compliantPublisher.subscribe(ctx, sub)
	go rp.process(ctx)
}

func (rp *CompliantRangePublisher) process(ctx context.Context) {
	fmt.Printf("CompliantRangePublisher.process started\n")
	defer func() {
		fmt.Printf("CompliantRangePublisher.process completed\n")
		rp.complete()
	}()

	// Handle nil context
	if ctx == nil {
		ctx = context.Background()
	}

	fmt.Printf("CompliantRangePublisher.process waiting for demand signal\n")
	select {
	case <-rp.compliantPublisher.demandSignal:
		fmt.Printf("CompliantRangePublisher.process received demand signal\n")
	case <-ctx.Done():
		fmt.Printf("CompliantRangePublisher.process context done\n")
		return
	}

	fmt.Printf("CompliantRangePublisher.process emitting values %d to %d\n", rp.start, rp.end)
	for i := rp.start; i <= rp.end; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("CompliantRangePublisher.process context done during emission\n")
			return
		default:
		}

		subs := rp.getActiveSubscribers()
		fmt.Printf("CompliantRangePublisher.process subscribers: %d\n", len(subs))
		if len(subs) == 0 {
			fmt.Printf("CompliantRangePublisher.process no subscribers, returning\n")
			return
		}

		canEmit := false
		for _, sub := range subs {
			if sub.canEmit() {
				canEmit = true
				break
			}
		}

		if !canEmit {
			fmt.Printf("CompliantRangePublisher.process no subscribers can emit, waiting for demand signal\n")
			select {
			case <-rp.compliantPublisher.demandSignal:
			case <-ctx.Done():
				return
			}
		}

		fmt.Printf("CompliantRangePublisher.process emitting %d\n", i)
		if !rp.emit(i) {
			fmt.Printf("CompliantRangePublisher.process emit returned false\n")
			return
		}
	}
}

// CompliantFromSlicePublisher implements a Reactive Streams compliant slice publisher
type CompliantFromSlicePublisher[T any] struct {
	*compliantPublisher[T]
	items []T
}

func NewCompliantFromSlicePublisher[T any](items []T) *CompliantFromSlicePublisher[T] {
	return &CompliantFromSlicePublisher[T]{
		compliantPublisher: newCompliantPublisher[T](),
		items:              items,
	}
}

func (sp *CompliantFromSlicePublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	sp.compliantPublisher.subscribe(ctx, sub)
	go sp.process(ctx)
}

func (sp *CompliantFromSlicePublisher[T]) process(ctx context.Context) {
	defer sp.complete()

	if len(sp.items) == 0 {
		return
	}

	// Handle nil context
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-sp.compliantPublisher.demandSignal:
	case <-ctx.Done():
		return
	}

	for _, item := range sp.items {
		select {
		case <-ctx.Done():
			return
		default:
		}

		subs := sp.getActiveSubscribers()
		if len(subs) == 0 {
			return
		}

		canEmit := false
		for _, sub := range subs {
			if sub.canEmit() {
				canEmit = true
				break
			}
		}

		if !canEmit {
			select {
			case <-sp.compliantPublisher.demandSignal:
			case <-ctx.Done():
				return
			}
		}

		if !sp.emit(item) {
			return
		}
	}
}

// CompliantBuilder provides fluent API for compliant publishers
type CompliantBuilder[T any] struct{}

func NewCompliantBuilder[T any]() *CompliantBuilder[T] {
	return &CompliantBuilder[T]{}
}

func (b *CompliantBuilder[T]) Range(start, end int) Publisher[int] {
	return NewCompliantRangePublisher(start, end)
}

func (b *CompliantBuilder[T]) FromSlice(items []T) Publisher[T] {
	return NewCompliantFromSlicePublisher(items)
}
