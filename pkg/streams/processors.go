package streams

import (
	"context"
	"errors"
	"sync"
)

// MapProcessor transforms items emitted by a Publisher by applying a function to each item.
// It implements the Processor interface and supports backpressure.
//
// Thread Safety:
// MapProcessor is safe for concurrent use. All methods are synchronized using a mutex.
// Multiple goroutines can safely call methods on the same MapProcessor instance.
type MapProcessor[T any, R any] struct {
	transform  func(T) R
	upstream   Subscription
	downstream Subscriber[R]
	mu         sync.Mutex
	terminated bool
}

// NewMapProcessor creates a new MapProcessor with the specified transformation function.
//
// Parameters:
//   - transform: A function that takes a value of type T and returns a value of type R
//
// Returns:
//   - A new MapProcessor instance
//
// Example:
//
//	processor := NewMapProcessor(func(x int) string {
//	    return strconv.Itoa(x)
//	})
//	// This processor will convert integers to strings
func NewMapProcessor[T any, R any](transform func(T) R) *MapProcessor[T, R] {
	return &MapProcessor[T, R]{
		transform: transform,
	}
}

// Subscribe subscribes a Subscriber to this MapProcessor.
// This method establishes the subscription and begins the flow of values
// from the upstream Publisher to the downstream Subscriber through the processor.
//
// Parameters:
//   - ctx: The context for the subscription, used for cancellation
//   - sub: The subscriber that will receive processed values
func (p *MapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &mapProcessorSubscription[T, R]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

// OnSubscribe is called when the subscription to the upstream Publisher is established.
// It stores the subscription and requests an unlimited number of items from the upstream.
//
// Parameters:
//   - s: The subscription to the upstream Publisher
func (p *MapProcessor[T, R]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	if p.downstream != nil {
		// Request a reasonable initial batch size for backpressure
		s.Request(128)
	}
}

// OnNext is called for each item emitted by the upstream Publisher.
// It applies the transformation function to the item and emits the result to the downstream Subscriber.
//
// Parameters:
//   - value: The item emitted by the upstream Publisher
func (p *MapProcessor[T, R]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	transformed := p.transform(value)
	p.downstream.OnNext(transformed)
}

// OnError is called when the upstream Publisher terminates with an error.
// It terminates the processor and propagates the error to the downstream Subscriber.
// After this method is called, no further calls to OnNext, OnComplete, or OnError will be made.
//
// Parameters:
//   - err: The error that caused the upstream Publisher to terminate
func (p *MapProcessor[T, R]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

// OnComplete is called when the upstream Publisher completes successfully.
// It terminates the processor and signals completion to the downstream Subscriber.
// After this method is called, no further calls to OnNext or OnComplete will be made.
func (p *MapProcessor[T, R]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

// mapProcessorSubscription represents the subscription from a downstream Subscriber to a MapProcessor.
// It implements the Subscription interface and delegates request/cancel operations to the upstream Publisher.
//
// Thread Safety:
// mapProcessorSubscription is safe for concurrent use through the MapProcessor's mutex.
type mapProcessorSubscription[T any, R any] struct {
	processor *MapProcessor[T, R]
}

// Request requests up to n additional items from the upstream Publisher.
// This method delegates the request to the upstream Publisher's subscription.
//
// Parameters:
//   - n: The number of additional items to request
func (s *mapProcessorSubscription[T, R]) Request(n int64) {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	}
}

// Cancel cancels the subscription and instructs the upstream Publisher to stop emitting items.
// This method delegates the cancellation to the upstream Publisher's subscription.
func (s *mapProcessorSubscription[T, R]) Cancel() {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

// FilterProcessor filters items emitted by a Publisher by applying a predicate function to each item.
// Only items for which the predicate returns true are emitted to the downstream Subscriber.
// It implements the Processor interface and supports backpressure.
//
// Thread Safety:
// FilterProcessor is safe for concurrent use. All methods are synchronized using a mutex.
// Multiple goroutines can safely call methods on the same FilterProcessor instance.
type FilterProcessor[T any] struct {
	predicate  func(T) bool
	upstream   Subscription
	downstream Subscriber[T]
	mu         sync.Mutex
	terminated bool
}

// NewFilterProcessor creates a new FilterProcessor with the specified predicate function.
//
// Parameters:
//   - predicate: A function that takes a value of type T and returns true if the value should be emitted
//
// Returns:
//   - A new FilterProcessor instance
//
// Example:
//
//	processor := NewFilterProcessor(func(x int) bool {
//	    return x % 2 == 0
//	})
//	// This processor will only emit even numbers
func NewFilterProcessor[T any](predicate func(T) bool) *FilterProcessor[T] {
	return &FilterProcessor[T]{
		predicate: predicate,
	}
}

// Subscribe subscribes a Subscriber to this FilterProcessor.
// This method establishes the subscription and begins the flow of values
// from the upstream Publisher to the downstream Subscriber through the processor.
//
// Parameters:
//   - ctx: The context for the subscription, used for cancellation
//   - sub: The subscriber that will receive filtered values
func (p *FilterProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &filterProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

// OnSubscribe is called when the subscription to the upstream Publisher is established.
// It stores the subscription and requests an unlimited number of items from the upstream.
//
// Parameters:
//   - s: The subscription to the upstream Publisher
func (p *FilterProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	if p.downstream != nil {
		// Request a reasonable initial batch size for backpressure
		s.Request(128)
	}
}

// OnNext is called for each item emitted by the upstream Publisher.
// It applies the predicate function to the item and emits the item to the downstream Subscriber
// only if the predicate returns true.
//
// Parameters:
//   - value: The item emitted by the upstream Publisher
func (p *FilterProcessor[T]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	if p.predicate(value) {
		p.downstream.OnNext(value)
	}
}

// OnError is called when the upstream Publisher terminates with an error.
// It terminates the processor and propagates the error to the downstream Subscriber.
// After this method is called, no further calls to OnNext, OnComplete, or OnError will be made.
//
// Parameters:
//   - err: The error that caused the upstream Publisher to terminate
func (p *FilterProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

// OnComplete is called when the upstream Publisher completes successfully.
// It terminates the processor and signals completion to the downstream Subscriber.
// After this method is called, no further calls to OnNext or OnComplete will be made.
func (p *FilterProcessor[T]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

// filterProcessorSubscription represents the subscription from a downstream Subscriber to a FilterProcessor.
// It implements the Subscription interface and delegates request/cancel operations to the upstream Publisher.
//
// Thread Safety:
// filterProcessorSubscription is safe for concurrent use through the FilterProcessor's mutex.
type filterProcessorSubscription[T any] struct {
	processor *FilterProcessor[T]
}

// Request requests up to n additional items from the upstream Publisher.
// This method delegates the request to the upstream Publisher's subscription.
//
// Parameters:
//   - n: The number of additional items to request
func (s *filterProcessorSubscription[T]) Request(n int64) {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	}
}

// Cancel cancels the subscription and instructs the upstream Publisher to stop emitting items.
// This method delegates the cancellation to the upstream Publisher's subscription.
func (s *filterProcessorSubscription[T]) Cancel() {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

type FlatMapProcessor[T any, R any] struct {
	transform         func(T) Publisher[R]
	ctx               context.Context
	upstream          Subscription
	downstream        Subscriber[R]
	mu                sync.Mutex
	terminated        bool
	upstreamCompleted bool
	activeInners      int
	completedInners   int
}

func NewFlatMapProcessor[T any, R any](transform func(T) Publisher[R]) *FlatMapProcessor[T, R] {
	return &FlatMapProcessor[T, R]{
		transform: transform,
	}
}

func (p *FlatMapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	if ctx == nil {
	} else {
		p.ctx = ctx
	}
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &flatMapProcessorSubscription[T, R]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

func (p *FlatMapProcessor[T, R]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	if p.downstream != nil {
		// Request a reasonable initial batch size for backpressure
		// FlatMap typically needs smaller batch sizes due to fan-out
		s.Request(32)
	}
}

func (p *FlatMapProcessor[T, R]) OnNext(value T) {
	p.mu.Lock()

	if p.terminated || p.downstream == nil {
		p.mu.Unlock()
		return
	}

	publisher := p.transform(value)
	p.activeInners++
	p.mu.Unlock()

	err := publisher.Subscribe(p.ctx, &innerSubscriber[T, R]{processor: p})
	if err != nil {
		innerSub := &innerSubscriber[T, R]{processor: p}
		innerSub.OnError(err)
	}
}

func (p *FlatMapProcessor[T, R]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

func (p *FlatMapProcessor[T, R]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.upstreamCompleted = true

	// Complete downstream only if no active inner publishers
	if p.activeInners == p.completedInners {
		p.terminated = true
		p.downstream.OnComplete()
	}
}

type flatMapProcessorSubscription[T any, R any] struct {
	processor *FlatMapProcessor[T, R]
}

func (s *flatMapProcessorSubscription[T, R]) Request(n int64) {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	}
}

func (s *flatMapProcessorSubscription[T, R]) Cancel() {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

type innerSubscriber[T any, R any] struct {
	processor *FlatMapProcessor[T, R]
}

func (s *innerSubscriber[T, R]) OnSubscribe(sub Subscription) {
	sub.Request(1<<63 - 1) // Request all items
}

func (s *innerSubscriber[T, R]) OnNext(value R) {
	s.processor.mu.Lock()
	defer s.processor.mu.Unlock()

	if s.processor.terminated || s.processor.downstream == nil {
		return
	}

	s.processor.downstream.OnNext(value)
}

func (s *innerSubscriber[T, R]) OnError(err error) {
	s.processor.mu.Lock()
	defer s.processor.mu.Unlock()

	if s.processor.terminated || s.processor.downstream == nil {
		return
	}

	s.processor.terminated = true
	s.processor.downstream.OnError(err)
}

func (s *innerSubscriber[T, R]) OnComplete() {
	s.processor.mu.Lock()
	defer s.processor.mu.Unlock()

	if s.processor.terminated || s.processor.downstream == nil {
		return
	}

	s.processor.completedInners++

	// Complete downstream if upstream is done and all inners are complete
	if s.processor.upstreamCompleted && s.processor.activeInners == s.processor.completedInners {
		s.processor.terminated = true
		s.processor.downstream.OnComplete()
	}
}

type MergeProcessor[T any] struct {
	sources          []Publisher[T]
	downstream       Subscriber[T]
	activeSources    int
	completedSources int
	mu               sync.Mutex
	terminated       bool
	ctx              context.Context
}

func NewMergeProcessor[T any](sources ...Publisher[T]) *MergeProcessor[T] {
	return &MergeProcessor[T]{
		sources:       sources,
		activeSources: len(sources),
	}
}

func (p *MergeProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	p.ctx = ctx
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &mergeProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

func (p *MergeProcessor[T]) OnSubscribe(s Subscription) {
	// Not directly used for merge processor
}

func (p *MergeProcessor[T]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.downstream.OnNext(value)
}

func (p *MergeProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

func (p *MergeProcessor[T]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.completedSources++
	if p.completedSources == p.activeSources && p.downstream != nil {
		p.terminated = true
		p.downstream.OnComplete()
	}
}

type mergeProcessorSubscription[T any] struct {
	processor *MergeProcessor[T]
}

func (s *mergeProcessorSubscription[T]) Request(n int64) {
	s.processor.mu.Lock()
	defer s.processor.mu.Unlock()

	ctx := s.processor.ctx

	if len(s.processor.sources) == 0 {
		s.processor.terminated = true
		if s.processor.downstream != nil {
			go func() {
				s.processor.mu.Lock()
				defer s.processor.mu.Unlock()
				if s.processor.downstream != nil {
					s.processor.downstream.OnComplete()
				}
			}()
		}
		return
	}

	for _, source := range s.processor.sources {
		err := source.Subscribe(ctx, &mergeInnerSubscriber[T]{processor: s.processor})
		if err != nil {
			innerSub := &mergeInnerSubscriber[T]{processor: s.processor}
			innerSub.OnError(err)
		}
	}
}

func (s *mergeProcessorSubscription[T]) Cancel() {
}

type mergeInnerSubscriber[T any] struct {
	processor *MergeProcessor[T]
}

func (s *mergeInnerSubscriber[T]) OnSubscribe(sub Subscription) {
	sub.Request(1<<63 - 1) // Request all items
}

func (s *mergeInnerSubscriber[T]) OnNext(value T) {
	s.processor.OnNext(value)
}

func (s *mergeInnerSubscriber[T]) OnError(err error) {
	s.processor.OnError(err)
}

func (s *mergeInnerSubscriber[T]) OnComplete() {
	s.processor.OnComplete()
}

type ConcatProcessor[T any] struct {
	sources      []Publisher[T]
	downstream   Subscriber[T]
	currentIndex int
	ctx          context.Context
	mu           sync.Mutex
	terminated   bool
}

func NewConcatProcessor[T any](sources ...Publisher[T]) *ConcatProcessor[T] {
	return &ConcatProcessor[T]{
		sources: sources,
	}
}

func (p *ConcatProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	p.ctx = ctx
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &concatProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

func (p *ConcatProcessor[T]) OnSubscribe(s Subscription) {
	// Not used for this processor
}

func (p *ConcatProcessor[T]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.downstream.OnNext(value)
}

func (p *ConcatProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

func (p *ConcatProcessor[T]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.currentIndex++
	if p.currentIndex < len(p.sources) && p.ctx != nil {
		// Move to the next source
		err := p.sources[p.currentIndex].Subscribe(p.ctx, &concatInnerSubscriber[T]{processor: p})
		if err != nil && p.downstream != nil {
			p.terminated = true
			p.downstream.OnError(err)
		}
	} else if p.downstream != nil {
		// All sources completed
		p.terminated = true
		p.downstream.OnComplete()
	}
}

type concatProcessorSubscription[T any] struct {
	processor *ConcatProcessor[T]
}

func (s *concatProcessorSubscription[T]) Request(n int64) {
	s.processor.mu.Lock()
	defer s.processor.mu.Unlock()

	ctx := s.processor.ctx

	// Start with the first source
	if len(s.processor.sources) > 0 {
		err := s.processor.sources[0].Subscribe(ctx, &concatInnerSubscriber[T]{processor: s.processor})
		if err != nil {
			innerSub := &concatInnerSubscriber[T]{processor: s.processor}
			innerSub.OnError(err)
		}
	} else {
		s.processor.terminated = true
		if s.processor.downstream != nil {
			go func() {
				s.processor.mu.Lock()
				defer s.processor.mu.Unlock()
				if s.processor.downstream != nil {
					s.processor.downstream.OnComplete()
				}
			}()
		}
	}
}

func (s *concatProcessorSubscription[T]) Cancel() {
	// For processors, cancellation would typically affect the upstream
}

type concatInnerSubscriber[T any] struct {
	processor *ConcatProcessor[T]
}

func (s *concatInnerSubscriber[T]) OnSubscribe(sub Subscription) {
	sub.Request(1<<63 - 1) // Request all items
}

func (s *concatInnerSubscriber[T]) OnNext(value T) {
	s.processor.OnNext(value)
}

func (s *concatInnerSubscriber[T]) OnError(err error) {
	s.processor.OnError(err)
}

func (s *concatInnerSubscriber[T]) OnComplete() {
	s.processor.OnComplete()
}

// TakeProcessor emits only the first n values
type TakeProcessor[T any] struct {
	n          int64
	count      int64
	upstream   Subscription
	downstream Subscriber[T]
	mu         sync.Mutex
	terminated bool
}

func NewTakeProcessor[T any](n int64) *TakeProcessor[T] {
	return &TakeProcessor[T]{
		n: n,
	}
}

func (p *TakeProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &takeProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

func (p *TakeProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	if p.downstream != nil && p.n > 0 {
		s.Request(p.n)
	}
}

func (p *TakeProcessor[T]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	if p.count < p.n {
		p.downstream.OnNext(value)
		p.count++
	}

	if p.count == p.n {
		p.terminated = true
		p.downstream.OnComplete()
		if p.upstream != nil {
			p.upstream.Cancel()
		}
	}
}

func (p *TakeProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

func (p *TakeProcessor[T]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

type takeProcessorSubscription[T any] struct {
	processor *TakeProcessor[T]
}

func (s *takeProcessorSubscription[T]) Request(n int64) {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	nValue := s.processor.n
	s.processor.mu.Unlock()

	if upstream != nil && nValue > 0 {
		requestAmount := nValue
		if requestAmount > n {
			requestAmount = n
		}
		upstream.Request(requestAmount)
	}
}

func (s *takeProcessorSubscription[T]) Cancel() {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

type SkipProcessor[T any] struct {
	n          int64
	count      int64
	upstream   Subscription
	downstream Subscriber[T]
	mu         sync.Mutex
	terminated bool
}

func NewSkipProcessor[T any](n int64) *SkipProcessor[T] {
	return &SkipProcessor[T]{
		n: n,
	}
}

func (p *SkipProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &skipProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

func (p *SkipProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	if p.downstream != nil {
		s.Request(1<<63 - 1)
	}
}

func (p *SkipProcessor[T]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	if p.count >= p.n {
		p.downstream.OnNext(value)
	} else {
		p.count++
	}
}

func (p *SkipProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

func (p *SkipProcessor[T]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

type skipProcessorSubscription[T any] struct {
	processor *SkipProcessor[T]
}

func (s *skipProcessorSubscription[T]) Request(n int64) {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	}
}

func (s *skipProcessorSubscription[T]) Cancel() {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

type DistinctProcessor[T comparable] struct {
	seen       map[T]struct{}
	upstream   Subscription
	downstream Subscriber[T]
	mu         sync.Mutex
	terminated bool
}

func NewDistinctProcessor[T comparable]() *DistinctProcessor[T] {
	return &DistinctProcessor[T]{
		seen: make(map[T]struct{}),
	}
}

func (p *DistinctProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	p.mu.Lock()
	p.downstream = sub
	p.mu.Unlock()

	processorSub := &distinctProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)

	return nil
}

func (p *DistinctProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	if p.downstream != nil {
		s.Request(1<<63 - 1)
	}
}

func (p *DistinctProcessor[T]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	if _, exists := p.seen[value]; !exists {
		p.seen[value] = struct{}{}
		p.downstream.OnNext(value)
	}
}

func (p *DistinctProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

func (p *DistinctProcessor[T]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

type distinctProcessorSubscription[T comparable] struct {
	processor *DistinctProcessor[T]
}

func (s *distinctProcessorSubscription[T]) Request(n int64) {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	}
}

func (s *distinctProcessorSubscription[T]) Cancel() {
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

type ProcessorBuilder[T any, R any] struct{}

func NewProcessorBuilder[T any, R any]() *ProcessorBuilder[T, R] {
	return &ProcessorBuilder[T, R]{}
}

func (b *ProcessorBuilder[T, R]) Map(transform func(T) R) *MapProcessor[T, R] {
	return NewMapProcessor(transform)
}

func (b *ProcessorBuilder[T, R]) Filter(predicate func(T) bool) *FilterProcessor[T] {
	return NewFilterProcessor(predicate)
}

func (b *ProcessorBuilder[T, R]) FlatMap(transform func(T) Publisher[R]) *FlatMapProcessor[T, R] {
	return NewFlatMapProcessor(transform)
}
