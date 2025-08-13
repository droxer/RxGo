package streams

import (
	"context"
)

// MapProcessor transforms items using a mapping function
type MapProcessor[T any, R any] struct {
	*compliantPublisher[R]
	transform func(T) R
}

func NewMapProcessor[T any, R any](transform func(T) R) *MapProcessor[T, R] {
	return &MapProcessor[T, R]{
		compliantPublisher: newCompliantPublisher[R](),
		transform:          transform,
	}
}

func (p *MapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) {
	p.compliantPublisher.subscribe(ctx, sub)
}

// OnSubscribe implements Subscriber[T]
func (p *MapProcessor[T, R]) OnSubscribe(s Subscription) {
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

// FilterProcessor filters items based on a predicate
type FilterProcessor[T any] struct {
	*compliantPublisher[T]
	predicate func(T) bool
}

func NewFilterProcessor[T any](predicate func(T) bool) *FilterProcessor[T] {
	return &FilterProcessor[T]{
		compliantPublisher: newCompliantPublisher[T](),
		predicate:          predicate,
	}
}

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

// FlatMapProcessor transforms items and flattens the results
type FlatMapProcessor[T any, R any] struct {
	*compliantPublisher[R]
	transform func(T) Publisher[R]
}

func NewFlatMapProcessor[T any, R any](transform func(T) Publisher[R]) *FlatMapProcessor[T, R] {
	return &FlatMapProcessor[T, R]{
		compliantPublisher: newCompliantPublisher[R](),
		transform:          transform,
	}
}

func (p *FlatMapProcessor[T, R]) Subscribe(ctx context.Context, sub Subscriber[R]) {
	p.compliantPublisher.subscribe(ctx, sub)
}

func (p *FlatMapProcessor[T, R]) OnSubscribe(s Subscription) {
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
	p.compliantPublisher.complete()
}

func (p *FlatMapProcessor[T, R]) newInnerSubscriber() Subscriber[R] {
	return &innerSubscriber[T, R]{
		processor: p,
	}
}

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
