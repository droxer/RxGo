package streams

import (
	"context"
	"errors"
	"sync"
)

type BufferedPublisher[T any] struct {
	config     BackpressureConfig
	source     func(ctx context.Context, sub Subscriber[T])
	buffer     []T
	head       int // Index of the first element
	tail       int // Index where the next element will be inserted
	count      int // Number of elements in the buffer
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
		buffer:     make([]T, config.BufferSize),
		head:       0,
		tail:       0,
		count:      0,
		bufferFull: make(chan struct{}),
	}
}

func (bp *BufferedPublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	subscription := &bufferedSubscription[T]{
		publisher: bp,
		sub:       sub,
		ctx:       ctx,
	}
	sub.OnSubscribe(subscription)

	return nil
}

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
	// Check if buffer is full
	if bp.count >= int(bp.config.BufferSize) {
		// Buffer is full, drop the value according to Buffer strategy
		// In Buffer strategy, we should block or drop based on implementation
		// For now, we'll just return (drop the value)
		return
	}

	// Add the value to the buffer
	bp.buffer[bp.tail] = value
	bp.tail = (bp.tail + 1) % int(bp.config.BufferSize)
	bp.count++

	// Process buffered items if there's demand
	for bp.count > 0 && sub.requested > 0 {
		next := bp.buffer[bp.head]
		bp.head = (bp.head + 1) % int(bp.config.BufferSize)
		bp.count--
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) handleDropStrategy(value T, sub *bufferedSubscription[T]) {
	// Check if buffer is full
	if bp.count >= int(bp.config.BufferSize) {
		// Buffer is full, drop the value according to Drop strategy
		return
	}

	// Add the value to the buffer
	bp.buffer[bp.tail] = value
	bp.tail = (bp.tail + 1) % int(bp.config.BufferSize)
	bp.count++

	// Process buffered items if there's demand
	for bp.count > 0 && sub.requested > 0 {
		next := bp.buffer[bp.head]
		bp.head = (bp.head + 1) % int(bp.config.BufferSize)
		bp.count--
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) handleLatestStrategy(value T, sub *bufferedSubscription[T]) {
	if bp.count >= int(bp.config.BufferSize) {
		// Buffer is full, replace the oldest item (at head position) with the new value
		bp.buffer[bp.head] = value
		// Move head forward since we're replacing the oldest item
		bp.head = (bp.head + 1) % int(bp.config.BufferSize)
		// tail stays the same, count stays the same
	} else {
		// Buffer is not full, add the value normally
		bp.buffer[bp.tail] = value
		bp.tail = (bp.tail + 1) % int(bp.config.BufferSize)
		bp.count++
	}

	// Process buffered items if there's demand
	for bp.count > 0 && sub.requested > 0 {
		next := bp.buffer[bp.head]
		bp.head = (bp.head + 1) % int(bp.config.BufferSize)
		bp.count--
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) handleErrorStrategy(value T, sub *bufferedSubscription[T]) {
	if bp.count >= int(bp.config.BufferSize) {
		sub.sub.OnError(errors.New("buffer overflow: buffer size limit exceeded"))
		return
	}

	// Add the value to the buffer
	bp.buffer[bp.tail] = value
	bp.tail = (bp.tail + 1) % int(bp.config.BufferSize)
	bp.count++

	// Process buffered items if there's demand
	for bp.count > 0 && sub.requested > 0 {
		next := bp.buffer[bp.head]
		bp.head = (bp.head + 1) % int(bp.config.BufferSize)
		bp.count--
		sub.requested--
		sub.sub.OnNext(next)
	}
}

func (bp *BufferedPublisher[T]) flushBuffer(sub *bufferedSubscription[T]) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	sub.mu.Lock()
	defer sub.mu.Unlock()

	// Process all remaining buffered items
	for bp.count > 0 && sub.requested > 0 {
		next := bp.buffer[bp.head]
		bp.head = (bp.head + 1) % int(bp.config.BufferSize)
		bp.count--
		sub.requested--
		sub.sub.OnNext(next)
	}

	// Handle any remaining items based on strategy
	if bp.count > 0 {
		switch bp.config.Strategy {
		case Buffer, Drop, Latest:
			// For these strategies, emit all remaining items regardless of demand
			for bp.count > 0 {
				next := bp.buffer[bp.head]
				bp.head = (bp.head + 1) % int(bp.config.BufferSize)
				bp.count--
				sub.sub.OnNext(next)
			}
		case Error:
			sub.sub.OnError(errors.New("buffer overflow during completion"))
		}
	}
}

type CompliantRangePublisher struct {
	*compliantPublisher[int]
	start   int
	count   int
	started bool
	mu      sync.Mutex
}

func NewCompliantRangePublisher(start, count int) *CompliantRangePublisher {
	return &CompliantRangePublisher{
		compliantPublisher: newCompliantPublisher[int](),
		start:              start,
		count:              count,
	}
}

func (rp *CompliantRangePublisher) Subscribe(ctx context.Context, sub Subscriber[int]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// Create subscription and call OnSubscribe synchronously
	subscription := newCompliantSubscription(sub, rp.compliantPublisher)

	rp.compliantPublisher.mu.Lock()
	if rp.compliantPublisher.terminal.Load() {
		rp.compliantPublisher.mu.Unlock()
		return nil
	}
	rp.compliantPublisher.subscribers[subscription] = struct{}{}
	rp.compliantPublisher.mu.Unlock()

	// Call OnSubscribe synchronously to establish demand before processing
	sub.OnSubscribe(subscription)

	// Only start processing if we're the first subscriber
	rp.mu.Lock()
	if !rp.started {
		rp.started = true
		go rp.process(ctx)
	}
	rp.mu.Unlock()

	return nil
}

func (rp *CompliantRangePublisher) process(ctx context.Context) {
	defer rp.complete()

	if ctx == nil {
		ctx = context.Background()
	}

	if rp.count <= 0 {
		return
	}

	for i := 0; i < rp.count; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Wait until we can emit this item
		for {
			subs := rp.getActiveSubscribers()
			if len(subs) == 0 {
				return
			}

			// Check if any subscriber can accept the item
			if rp.emit(rp.start + i) {
				break // Successfully emitted, move to next item
			}

			// No subscriber could accept the item, wait for demand
			select {
			case <-rp.compliantPublisher.demandSignal:
				// Try again after receiving demand signal
			case <-ctx.Done():
				return
			}
		}
	}
}

type CompliantFromSlicePublisher[T any] struct {
	*compliantPublisher[T]
	items   []T
	started bool
	mu      sync.Mutex
}

func NewCompliantFromSlicePublisher[T any](items []T) *CompliantFromSlicePublisher[T] {
	return &CompliantFromSlicePublisher[T]{
		compliantPublisher: newCompliantPublisher[T](),
		items:              items,
	}
}

func (sp *CompliantFromSlicePublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// Create subscription and call OnSubscribe synchronously
	subscription := newCompliantSubscription(sub, sp.compliantPublisher)

	sp.compliantPublisher.mu.Lock()
	if sp.compliantPublisher.terminal.Load() {
		sp.compliantPublisher.mu.Unlock()
		return nil
	}
	sp.compliantPublisher.subscribers[subscription] = struct{}{}
	sp.compliantPublisher.mu.Unlock()

	// Call OnSubscribe synchronously to establish demand before processing
	sub.OnSubscribe(subscription)

	// Only start processing if we're the first subscriber
	sp.mu.Lock()
	if !sp.started {
		sp.started = true
		go sp.process(ctx)
	}
	sp.mu.Unlock()

	return nil
}

func (sp *CompliantFromSlicePublisher[T]) process(ctx context.Context) {
	defer sp.complete()

	if len(sp.items) == 0 {
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	for _, item := range sp.items {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Wait until we can emit this item
		for {
			subs := sp.getActiveSubscribers()
			if len(subs) == 0 {
				return
			}

			// Check if any subscriber can accept the item
			if sp.emit(item) {
				break // Successfully emitted, move to next item
			}

			// No subscriber could accept the item, wait for demand
			select {
			case <-sp.compliantPublisher.demandSignal:
				// Try again after receiving demand signal
			case <-ctx.Done():
				return
			}
		}
	}
}

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
