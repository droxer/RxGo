package streams

import (
	"context"
)

// MapProcessor implements a Processor that transforms items using a mapping function
type MapProcessor[T any, R any] struct {
	*compliantPublisher[R]
	transform func(T) R
}

// NewMapProcessor creates a new MapProcessor
func NewMapProcessor[T any, R any](transform func(T) R) *MapProcessor[T, R] {
	return &MapProcessor[T, R]{
		compliantPublisher: newCompliantPublisher[R](),
		transform:          transform,
	}
}

// Subscribe implements Publisher[R]
func (p *MapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) {
	p.compliantPublisher.subscribe(ctx, sub)
}

// OnSubscribe implements Subscriber[T]
func (p *MapProcessor[T, R]) OnSubscribe(s Subscription) {
	// Forward subscription to upstream
	// This would be connected to the upstream publisher
	// In a real implementation, this would be set by the upstream
}

// OnNext implements Subscriber[T]
func (p *MapProcessor[T, R]) OnNext(value T) {
	transformed := p.transform(value)
	p.compliantPublisher.emit(transformed)
}

// OnError implements Subscriber[T]
func (p *MapProcessor[T, R]) OnError(err error) {
	p.compliantPublisher.error(err)
}

// OnComplete implements Subscriber[T]
func (p *MapProcessor[T, R]) OnComplete() {
	p.compliantPublisher.complete()
}

// FilterProcessor implements a Processor that filters items based on a predicate
type FilterProcessor[T any] struct {
	*compliantPublisher[T]
	predicate func(T) bool
}

// NewFilterProcessor creates a new FilterProcessor
func NewFilterProcessor[T any](predicate func(T) bool) *FilterProcessor[T] {
	return &FilterProcessor[T]{
		compliantPublisher: newCompliantPublisher[T](),
		predicate:          predicate,
	}
}

// Subscribe implements Publisher[T]
func (p *FilterProcessor[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	p.compliantPublisher.subscribe(ctx, sub)
}

// OnSubscribe implements Subscriber[T]
func (p *FilterProcessor[T]) OnSubscribe(s Subscription) {
	// Forward subscription to upstream
}

// OnNext implements Subscriber[T]
func (p *FilterProcessor[T]) OnNext(value T) {
	if p.predicate(value) {
		p.compliantPublisher.emit(value)
	}
}

// OnError implements Subscriber[T]
func (p *FilterProcessor[T]) OnError(err error) {
	p.compliantPublisher.error(err)
}

// OnComplete implements Subscriber[T]
func (p *FilterProcessor[T]) OnComplete() {
	p.compliantPublisher.complete()
}

// FlatMapProcessor implements a Processor that transforms items and flattens the results
type FlatMapProcessor[T any, R any] struct {
	*compliantPublisher[R]
	transform func(T) Publisher[R]
}

// NewFlatMapProcessor creates a new FlatMapProcessor
func NewFlatMapProcessor[T any, R any](transform func(T) Publisher[R]) *FlatMapProcessor[T, R] {
	return &FlatMapProcessor[T, R]{
		compliantPublisher: newCompliantPublisher[R](),
		transform:          transform,
	}
}

// Subscribe implements Publisher[R]
func (p *FlatMapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) {
	p.compliantPublisher.subscribe(ctx, sub)
}

// OnSubscribe implements Subscriber[T]
func (p *FlatMapProcessor[T, R]) OnSubscribe(s Subscription) {
	// Forward subscription to upstream
}

// OnNext implements Subscriber[T]
func (p *FlatMapProcessor[T, R]) OnNext(value T) {
	publisher := p.transform(value)
	publisher.Subscribe(context.Background(), p.newInnerSubscriber())
}

// OnError implements Subscriber[T]
func (p *FlatMapProcessor[T, R]) OnError(err error) {
	p.compliantPublisher.error(err)
}

// OnComplete implements Subscriber[T]
func (p *FlatMapProcessor[T, R]) OnComplete() {
	// Wait for all inner publishers to complete
	// Implementation would track active publishers
	p.compliantPublisher.complete()
}

func (p *FlatMapProcessor[T, R]) newInnerSubscriber() Subscriber[R] {
	return &innerSubscriber[T, R]{
		processor: p,
	}
}

// innerSubscriber handles inner publisher subscriptions
type innerSubscriber[T any, R any] struct {
	processor *FlatMapProcessor[T, R]
}

func (s *innerSubscriber[T, R]) OnSubscribe(sub Subscription) {
	sub.Request(1)
}

func (s *innerSubscriber[T, R]) OnNext(value R) {
	s.processor.compliantPublisher.emit(value)
}

func (s *innerSubscriber[T, R]) OnError(err error) {
	s.processor.compliantPublisher.error(err)
}

func (s *innerSubscriber[T, R]) OnComplete() {
	// Inner publisher completed
}

// ProcessorBuilder provides a fluent API for building processors
type ProcessorBuilder[T any, R any] struct{}

// NewProcessorBuilder creates a new processor builder
func NewProcessorBuilder[T any, R any]() *ProcessorBuilder[T, R] {
	return &ProcessorBuilder[T, R]{}
}

// Map creates a Map processor
func (b *ProcessorBuilder[T, R]) Map(transform func(T) R) Processor[T, R] {
	return NewMapProcessor(transform)
}

// Filter creates a Filter processor
func (b *ProcessorBuilder[T, R]) Filter(predicate func(T) bool) Processor[T, T] {
	return NewFilterProcessor(predicate)
}

// FlatMap creates a FlatMap processor
func (b *ProcessorBuilder[T, R]) FlatMap(transform func(T) Publisher[R]) Processor[T, R] {
	return NewFlatMapProcessor(transform)
}
