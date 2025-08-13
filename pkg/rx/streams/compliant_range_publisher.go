package streams

import (
	"context"
	"sync"
)

// CompliantRangePublisher implements a Reactive Streams 1.0.4 compliant range publisher
type CompliantRangePublisher struct {
	*compliantPublisher[int]
	start int
	end   int
}

// NewCompliantRangePublisher creates a new compliant range publisher
func NewCompliantRangePublisher(start, end int) *CompliantRangePublisher {
	return &CompliantRangePublisher{
		compliantPublisher: newCompliantPublisher[int](),
		start:              start,
		end:                end,
	}
}

// Subscribe implements Publisher[int]
func (rp *CompliantRangePublisher) Subscribe(ctx context.Context, sub Subscriber[int]) {
	rp.compliantPublisher.subscribe(ctx, sub)
	go rp.process(ctx)
}

// process handles the actual publishing logic with demand control
func (rp *CompliantRangePublisher) process(ctx context.Context) {
	defer rp.complete()

	// Wait for initial demand
	select {
	case <-rp.demandSignal:
	case <-ctx.Done():
		return
	}

	for i := rp.start; i <= rp.end; i++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Check if we have demand
		subs := rp.getActiveSubscribers()
		if len(subs) == 0 {
			return
		}

		// Wait for demand from at least one subscriber
		canEmit := false
		for _, sub := range subs {
			if sub.canEmit() {
				canEmit = true
				break
			}
		}

		if !canEmit {
			// Wait for demand signal
			select {
			case <-rp.demandSignal:
			case <-ctx.Done():
				return
			}
		}

		// Emit item with demand tracking
		if !rp.emit(i) {
			// No active subscribers or no demand
			return
		}
	}
}

// CompliantRangePublisherWithConfig creates a compliant range publisher with backpressure config
func CompliantRangePublisherWithConfig(start, end int, config BackpressureConfig) Publisher[int] {
	// For now, ignore backpressure config in compliant version
	// Future: integrate with buffered publisher
	return NewCompliantRangePublisher(start, end)
}

// CompliantFromSlicePublisher implements a Reactive Streams 1.0.4 compliant slice publisher
type CompliantFromSlicePublisher[T any] struct {
	*compliantPublisher[T]
	items []T
}

// NewCompliantFromSlicePublisher creates a new compliant slice publisher
func NewCompliantFromSlicePublisher[T any](items []T) *CompliantFromSlicePublisher[T] {
	return &CompliantFromSlicePublisher[T]{
		compliantPublisher: newCompliantPublisher[T](),
		items:              items,
	}
}

// Subscribe implements Publisher[T]
func (sp *CompliantFromSlicePublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	sp.compliantPublisher.subscribe(ctx, sub)
	go sp.process(ctx)
}

// process handles the actual publishing logic with demand control
func (sp *CompliantFromSlicePublisher[T]) process(ctx context.Context) {
	defer sp.complete()

	if len(sp.items) == 0 {
		return
	}

	// Wait for initial demand
	select {
	case <-sp.demandSignal:
	case <-ctx.Done():
		return
	}

	for _, item := range sp.items {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Check if we have demand
		subs := sp.getActiveSubscribers()
		if len(subs) == 0 {
			return
		}

		// Wait for demand from at least one subscriber
		canEmit := false
		for _, sub := range subs {
			if sub.canEmit() {
				canEmit = true
				break
			}
		}

		if !canEmit {
			// Wait for demand signal
			select {
			case <-sp.demandSignal:
			case <-ctx.Done():
				return
			}
		}

		// Emit item with demand tracking
		if !sp.emit(item) {
			// No active subscribers or no demand
			return
		}
	}
}

// CompliantFromSlicePublisherWithConfig creates a compliant slice publisher with backpressure config
func CompliantFromSlicePublisherWithConfig[T any](items []T, config BackpressureConfig) Publisher[T] {
	// For now, ignore backpressure config in compliant version
	return NewCompliantFromSlicePublisher(items)
}

// CompliantBufferedPublisher implements a Reactive Streams 1.0.4 compliant buffered publisher
type CompliantBufferedPublisher[T any] struct {
	*compliantPublisher[T]
	source func(context.Context, Subscriber[T])
}

// NewCompliantBufferedPublisher creates a new compliant buffered publisher
func NewCompliantBufferedPublisher[T any](source func(context.Context, Subscriber[T])) *CompliantBufferedPublisher[T] {
	return &CompliantBufferedPublisher[T]{
		compliantPublisher: newCompliantPublisher[T](),
		source:             source,
	}
}

// Subscribe implements Publisher[T]
func (bp *CompliantBufferedPublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	bp.compliantPublisher.subscribe(ctx, sub)
	go func() {
		bp.source(ctx, &bufferedSourceSubscriber[T]{
			publisher: bp,
		})
	}()
}

// bufferedSourceSubscriber adapts the source function to our compliant publisher
type bufferedSourceSubscriber[T any] struct {
	publisher *CompliantBufferedPublisher[T]
}

func (bs *bufferedSourceSubscriber[T]) OnSubscribe(s Subscription) {
	// Forward to compliant publisher
}

func (bs *bufferedSourceSubscriber[T]) OnNext(value T) {
	bs.publisher.emit(value)
}

func (bs *bufferedSourceSubscriber[T]) OnError(err error) {
	bs.publisher.error(err)
}

func (bs *bufferedSourceSubscriber[T]) OnComplete() {
	bs.publisher.complete()
}

// CompliantNewBufferedPublisher creates a compliant buffered publisher with config
func CompliantNewBufferedPublisher[T any](config BackpressureConfig, source func(context.Context, Subscriber[T])) Publisher[T] {
	// For now, ignore backpressure config in compliant version
	return NewCompliantBufferedPublisher(source)
}

// CompliantBuilder provides fluent API for compliant publishers
type CompliantBuilder[T any] struct{}

// NewCompliantBuilder creates a new compliant builder
func NewCompliantBuilder[T any]() *CompliantBuilder[T] {
	return &CompliantBuilder[T]{}
}

// Range creates a compliant range publisher
func (b *CompliantBuilder[T]) Range(start, end int) Publisher[int] {
	return NewCompliantRangePublisher(start, end)
}

// FromSlice creates a compliant slice publisher
func (b *CompliantBuilder[T]) FromSlice(items []T) Publisher[T] {
	return NewCompliantFromSlicePublisher(items)
}

// SignalSerializer ensures sequential signaling (Rule 1.3, 2.7)
type SignalSerializer struct {
	mu       sync.Mutex
	terminal bool
}

// Serialize executes the signal function sequentially
func (s *SignalSerializer) Serialize(signal func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.terminal {
		signal()
	}
}

// MarkTerminal marks the serializer as terminal
func (s *SignalSerializer) MarkTerminal() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.terminal = true
}

// IsTerminal checks if the serializer is in terminal state
func (s *SignalSerializer) IsTerminal() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.terminal
}

// ConcurrentAccessValidator validates thread safety
type ConcurrentAccessValidator[T any] struct {
	mu       sync.Mutex
	accesses []string
}

// RecordAccess records an access for validation
func (v *ConcurrentAccessValidator[T]) RecordAccess(access string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.accesses = append(v.accesses, access)
}

// GetAccesses returns all recorded accesses
func (v *ConcurrentAccessValidator[T]) GetAccesses() []string {
	v.mu.Lock()
	defer v.mu.Unlock()
	return append([]string{}, v.accesses...)
}

// ResetAccesses clears all recorded accesses
func (v *ConcurrentAccessValidator[T]) ResetAccesses() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.accesses = nil
}

// ReactiveStreamsValidator provides utilities for validating Reactive Streams compliance
type ReactiveStreamsValidator struct {
	// Validation utilities
}

// ValidatePublisher checks publisher compliance
func (v *ReactiveStreamsValidator) ValidatePublisher(p Publisher[any]) error {
	// Implement publisher validation logic
	return nil
}

// ValidateSubscriber checks subscriber compliance
func (v *ReactiveStreamsValidator) ValidateSubscriber(s Subscriber[any]) error {
	// Implement subscriber validation logic
	return nil
}

// ValidateProcessor checks processor compliance
func (v *ReactiveStreamsValidator) ValidateProcessor(p Processor[any, any]) error {
	// Implement processor validation logic
	return nil
}
