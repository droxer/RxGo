// Package rxgo provides the main entry point for the RxGo library.
package rxgo

import (
	"context"

	"github.com/droxer/RxGo/internal/publisher"
	"github.com/droxer/RxGo/pkg/observable"
)

// Re-export main types for easy access
type (
	// Publisher represents a source of data that can be subscribed to
	Publisher[T any] = publisher.Publisher[T]

	// ReactiveSubscriber represents a subscriber in the Reactive Streams specification
	ReactiveSubscriber[T any] = publisher.ReactiveSubscriber[T]

	// Subscription represents the link between Publisher and ReactiveSubscriber
	Subscription = publisher.Subscription

	// Observable represents the legacy observable interface
	Observable[T any] = *observable.Observable[T]

	// Subscriber represents the legacy subscriber interface
	Subscriber[T any] = observable.Subscriber[T]
)

// NewReactivePublisher creates a new Publisher with Reactive Streams compliance
func NewReactivePublisher[T any](onSubscribe func(ctx context.Context, sub ReactiveSubscriber[T])) Publisher[T] {
	return publisher.NewReactivePublisher(onSubscribe)
}

// FromSlice creates a Publisher from a slice of values
func FromSlice[T any](items []T) Publisher[T] {
	return publisher.FromSlice(items)
}

// RangePublisher creates a Publisher that emits a range of integers
func RangePublisher(start, count int) Publisher[int] {
	return publisher.RangePublisher(start, count)
}

// Create creates an Observable using the legacy API
func Create[T any](onSubscribe func(ctx context.Context, sub Subscriber[T])) Observable[T] {
	return observable.Create(onSubscribe)
}

// Just creates an Observable that emits the provided values
func Just[T any](values ...T) Observable[T] {
	return observable.Just(values...)
}

// Range creates an Observable that emits integers in the specified range
func Range(start, count int) Observable[int] {
	return observable.Range(start, count)
}
