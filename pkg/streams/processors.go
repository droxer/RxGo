package streams

import (
	"context"
	"fmt"
	"sync"
)

// MapProcessor transforms items using a mapping function
type MapProcessor[T any, R any] struct {
	transform  func(T) R
	upstream   Subscription
	downstream Subscriber[R]
	mu         sync.Mutex
	terminated bool
}

func NewMapProcessor[T any, R any](transform func(T) R) *MapProcessor[T, R] {
	return &MapProcessor[T, R]{
		transform: transform,
	}
}

func (p *MapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &mapProcessorSubscription[T, R]{processor: p}
	sub.OnSubscribe(processorSub)
}

// OnSubscribe implements Subscriber[T]
func (p *MapProcessor[T, R]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	// If we already have downstream, request all items
	if p.downstream != nil {
		s.Request(1<<63 - 1)
	}
}

// OnNext implements Subscriber[T]
func (p *MapProcessor[T, R]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	transformed := p.transform(value)
	p.downstream.OnNext(transformed)
}

// OnError implements Subscriber[T]
func (p *MapProcessor[T, R]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

// OnComplete implements Subscriber[T]
func (p *MapProcessor[T, R]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

// mapProcessorSubscription implements Subscription for map processors
type mapProcessorSubscription[T any, R any] struct {
	processor *MapProcessor[T, R]
}

func (s *mapProcessorSubscription[T, R]) Request(n int64) {
	fmt.Printf("mapProcessorSubscription.Request: %d\n", n)
	// Forward the request to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	downstream := s.processor.downstream
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.mu.Unlock()

	if upstream != nil {
		fmt.Printf("Forwarding request to upstream: %d\n", n)
		upstream.Request(n)
	} else {
		fmt.Printf("No upstream yet, will forward request when upstream is available\n")
		// If we don't have upstream yet, store the request and forward it when we get the upstream
		s.processor.mu.Lock()
		if s.processor.upstream != nil {
			fmt.Printf("Upstream became available, forwarding request: %d\n", n)
			s.processor.upstream.Request(n)
		}
		s.processor.mu.Unlock()
	}
}

func (s *mapProcessorSubscription[T, R]) Cancel() {
	fmt.Printf("mapProcessorSubscription.Cancel\n")
	// Forward the cancel to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

// FilterProcessor filters items based on a predicate
type FilterProcessor[T any] struct {
	predicate  func(T) bool
	upstream   Subscription
	downstream Subscriber[T]
	mu         sync.Mutex
	terminated bool
}

func NewFilterProcessor[T any](predicate func(T) bool) *FilterProcessor[T] {
	return &FilterProcessor[T]{
		predicate: predicate,
	}
}

func (p *FilterProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &filterProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)
}

// OnSubscribe implements Subscriber[T]
func (p *FilterProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	// If we already have downstream, request all items
	if p.downstream != nil {
		s.Request(1<<63 - 1)
	}
}

// OnNext implements Subscriber[T]
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

// OnError implements Subscriber[T]
func (p *FilterProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

// OnComplete implements Subscriber[T]
func (p *FilterProcessor[T]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

// filterProcessorSubscription implements Subscription for filter processors
type filterProcessorSubscription[T any] struct {
	processor *FilterProcessor[T]
}

func (s *filterProcessorSubscription[T]) Request(n int64) {
	// Forward the request to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	downstream := s.processor.downstream
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	} else {
		// If we don't have upstream yet, store the request and forward it when we get the upstream
		s.processor.mu.Lock()
		if s.processor.upstream != nil {
			s.processor.upstream.Request(n)
		}
		s.processor.mu.Unlock()
	}
}

func (s *filterProcessorSubscription[T]) Cancel() {
	// Forward the cancel to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

// FlatMapProcessor transforms items and flattens the results
type FlatMapProcessor[T any, R any] struct {
	transform  func(T) Publisher[R]
	upstream   Subscription
	downstream Subscriber[R]
	mu         sync.Mutex
	terminated bool
}

func NewFlatMapProcessor[T any, R any](transform func(T) Publisher[R]) *FlatMapProcessor[T, R] {
	return &FlatMapProcessor[T, R]{
		transform: transform,
	}
}

func (p *FlatMapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &flatMapProcessorSubscription[T, R]{processor: p}
	sub.OnSubscribe(processorSub)
}

func (p *FlatMapProcessor[T, R]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	// If we already have downstream, request all items
	if p.downstream != nil {
		s.Request(1<<63 - 1)
	}
}

// OnNext implements Subscriber[T]
func (p *FlatMapProcessor[T, R]) OnNext(value T) {
	p.mu.Lock()

	if p.terminated || p.downstream == nil {
		p.mu.Unlock()
		return
	}

	publisher := p.transform(value)
	// We need to unlock before subscribing to avoid deadlock
	p.mu.Unlock()

	publisher.Subscribe(context.Background(), &innerSubscriber[T, R]{processor: p})
}

// OnError implements Subscriber[T]
func (p *FlatMapProcessor[T, R]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

// OnComplete implements Subscriber[T]
func (p *FlatMapProcessor[T, R]) OnComplete() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnComplete()
}

type flatMapProcessorSubscription[T any, R any] struct {
	processor *FlatMapProcessor[T, R]
}

func (s *flatMapProcessorSubscription[T, R]) Request(n int64) {
	// Forward the request to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	downstream := s.processor.downstream
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	} else {
		// If we don't have upstream yet, store the request and forward it when we get the upstream
		s.processor.mu.Lock()
		if s.processor.upstream != nil {
			s.processor.upstream.Request(n)
		}
		s.processor.mu.Unlock()
	}
}

func (s *flatMapProcessorSubscription[T, R]) Cancel() {
	// Forward the cancel to the upstream
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
	// Inner publisher completed, but the main processor continues
}

// MergeProcessor combines multiple publishers into one by merging their emissions
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

func (p *MergeProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &mergeProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)
}

func (p *MergeProcessor[T]) OnSubscribe(s Subscription) {
	// Not directly used for merge processor
}

// OnNext implements Subscriber[T] for inner subscribers
func (p *MergeProcessor[T]) OnNext(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.downstream.OnNext(value)
}

// OnError implements Subscriber[T] for inner subscribers
func (p *MergeProcessor[T]) OnError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.terminated || p.downstream == nil {
		return
	}

	p.terminated = true
	p.downstream.OnError(err)
}

// OnComplete implements Subscriber[T] for inner subscribers
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
	ctx := s.processor.ctx
	downstream := s.processor.downstream
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.ctx = ctx
	s.processor.mu.Unlock()

	// Subscribe to all sources immediately
	s.processor.mu.Lock()
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
		s.processor.mu.Unlock()
		return
	}

	// Subscribe to all sources
	for _, source := range s.processor.sources {
		source.Subscribe(ctx, &mergeInnerSubscriber[T]{processor: s.processor})
	}
	s.processor.mu.Unlock()
}

func (s *mergeProcessorSubscription[T]) Cancel() {
	// For merge processor, cancellation affects all sources
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

// ConcatProcessor emits values from publishers sequentially
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

func (p *ConcatProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &concatProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)
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
		p.sources[p.currentIndex].Subscribe(p.ctx, &concatInnerSubscriber[T]{processor: p})
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
	ctx := s.processor.ctx
	downstream := s.processor.downstream
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.ctx = ctx
	s.processor.mu.Unlock()

	// Start with the first source
	s.processor.mu.Lock()
	if len(s.processor.sources) > 0 {
		s.processor.sources[0].Subscribe(ctx, &concatInnerSubscriber[T]{processor: s.processor})
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
	s.processor.mu.Unlock()
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

func (p *TakeProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &takeProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)
}

func (p *TakeProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	// If we already have downstream, request up to n items
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

	// If we've reached our limit, complete and cancel upstream
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
	// Forward the request to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	downstream := s.processor.downstream
	nValue := s.processor.n
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.mu.Unlock()

	if upstream != nil && nValue > 0 {
		// Request up to n items
		requestAmount := nValue
		if requestAmount > n {
			requestAmount = n
		}
		upstream.Request(requestAmount)
	} else {
		// If we don't have upstream yet, store the request and forward it when we get the upstream
		s.processor.mu.Lock()
		if s.processor.upstream != nil && s.processor.n > 0 {
			requestAmount := s.processor.n
			if requestAmount > n {
				requestAmount = n
			}
			s.processor.upstream.Request(requestAmount)
		}
		s.processor.mu.Unlock()
	}
}

func (s *takeProcessorSubscription[T]) Cancel() {
	// Forward the cancel to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

// SkipProcessor skips the first n values
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

func (p *SkipProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &skipProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)
}

func (p *SkipProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	// If we already have downstream, request all items
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
	// Forward the request to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	downstream := s.processor.downstream
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	} else {
		// If we don't have upstream yet, store the request and forward it when we get the upstream
		s.processor.mu.Lock()
		if s.processor.upstream != nil {
			s.processor.upstream.Request(n)
		}
		s.processor.mu.Unlock()
	}
}

func (s *skipProcessorSubscription[T]) Cancel() {
	// Forward the cancel to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

// DistinctProcessor suppresses duplicate items
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

func (p *DistinctProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Create a subscription for this processor and give it to the subscriber
	processorSub := &distinctProcessorSubscription[T]{processor: p}
	sub.OnSubscribe(processorSub)
}

func (p *DistinctProcessor[T]) OnSubscribe(s Subscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upstream = s

	// If we already have downstream, request all items
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
	// Forward the request to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	downstream := s.processor.downstream
	s.processor.mu.Unlock()

	// Store the downstream reference in the processor
	s.processor.mu.Lock()
	s.processor.downstream = downstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Request(n)
	} else {
		// If we don't have upstream yet, store the request and forward it when we get the upstream
		s.processor.mu.Lock()
		if s.processor.upstream != nil {
			s.processor.upstream.Request(n)
		}
		s.processor.mu.Unlock()
	}
}

func (s *distinctProcessorSubscription[T]) Cancel() {
	// Forward the cancel to the upstream
	s.processor.mu.Lock()
	upstream := s.processor.upstream
	s.processor.mu.Unlock()

	if upstream != nil {
		upstream.Cancel()
	}
}

// ProcessorBuilder provides fluent API for building processors
type ProcessorBuilder[T any, R any] struct{}

func NewProcessorBuilder[T any, R any]() *ProcessorBuilder[T, R] {
	return &ProcessorBuilder[T, R]{}
}

func (b *ProcessorBuilder[T, R]) Map(transform func(T) R) Processor[T, R] {
	return NewMapProcessor(transform)
}

func (b *ProcessorBuilder[T, R]) Filter(predicate func(T) bool) Processor[T, T] {
	return NewFilterProcessor(predicate)
}

func (b *ProcessorBuilder[T, R]) FlatMap(transform func(T) Publisher[R]) Processor[T, R] {
	return NewFlatMapProcessor(transform)
}
