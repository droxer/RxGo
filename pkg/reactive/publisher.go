package reactive

import (
	"context"

	"github.com/droxer/RxGo/internal/publisher"
)

// Publisher represents a provider of a potentially unbounded number of sequenced elements,
// publishing them according to the demand received from its Subscriber(s).
type Publisher[T any] interface {
	Subscribe(ctx context.Context, s Subscriber[T])
}

// Subscriber is a consumer of elements published by a Publisher.
type Subscriber[T any] interface {
	OnSubscribe(s Subscription)
	OnNext(t T)
	OnError(err error)
	OnComplete()
}

// Subscription represents a one-to-one lifecycle of a Subscriber subscribing to a Publisher.
type Subscription interface {
	Request(n int64)
	Cancel()
}

// Processor represents a processing stage that is both a Subscriber and a Publisher.
type Processor[T, R any] interface {
	Publisher[R]
	Subscriber[T]
}

// NewPublisher creates a new Publisher with Reactive Streams compliance
func NewPublisher[T any](onSubscribe func(ctx context.Context, sub Subscriber[T])) Publisher[T] {
	return &publisherAdapter[T]{
		onSubscribe: onSubscribe,
	}
}

// FromSlice creates a Publisher from a slice of values
func FromSlice[T any](items []T) Publisher[T] {
	return &publisherAdapter[T]{
		onSubscribe: func(ctx context.Context, sub Subscriber[T]) {
			bridge := &subscriberBridge[T]{sub: sub}
			internalPub := publisher.FromSlice(items)
			internalPub.Subscribe(ctx, bridge)
		},
	}
}

// Range creates a Publisher that emits a range of integers
func Range(start, count int) Publisher[int] {
	return &publisherAdapter[int]{
		onSubscribe: func(ctx context.Context, sub Subscriber[int]) {
			bridge := &subscriberBridge[int]{sub: sub}
			internalPub := publisher.RangePublisher(start, count)
			internalPub.Subscribe(ctx, bridge)
		},
	}
}

// Empty creates a Publisher that completes without emitting any items
func Empty[T any]() Publisher[T] {
	return NewPublisher(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnSubscribe(&emptySubscription{})
		sub.OnComplete()
	})
}

// Error creates a Publisher that immediately signals an error
func Error[T any](err error) Publisher[T] {
	return NewPublisher(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnSubscribe(&emptySubscription{})
		sub.OnError(err)
	})
}

// Never creates a Publisher that never signals any event
func Never[T any]() Publisher[T] {
	return NewPublisher(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnSubscribe(&emptySubscription{})
		<-ctx.Done()
	})
}

// Just creates a Publisher that emits the provided values
func Just[T any](values ...T) Publisher[T] {
	return FromSlice(values)
}


type publisherAdapter[T any] struct {
	onSubscribe func(ctx context.Context, sub Subscriber[T])
}

func (a *publisherAdapter[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	a.onSubscribe(ctx, sub)
}


type subscriberBridge[T any] struct {
	sub Subscriber[T]
}

func (b *subscriberBridge[T]) OnSubscribe(s publisher.Subscription) {
	b.sub.OnSubscribe(&subscriptionBridge{internal: s})
}

func (b *subscriberBridge[T]) OnNext(t T) {
	b.sub.OnNext(t)
}

func (b *subscriberBridge[T]) OnError(err error) {
	b.sub.OnError(err)
}

func (b *subscriberBridge[T]) OnComplete() {
	b.sub.OnComplete()
}


type subscriptionBridge struct {
	internal publisher.Subscription
}

func (b *subscriptionBridge) Request(n int64) {
	b.internal.Request(n)
}

func (b *subscriptionBridge) Cancel() {
	b.internal.Cancel()
}


type emptySubscription struct{}

func (e *emptySubscription) Request(n int64) {}
func (e *emptySubscription) Cancel()         {}
