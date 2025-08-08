// Package rxgo provides the main entry point for the RxGo library.
package rxgo

import (
	"context"

	"github.com/droxer/RxGo/internal/publisher"
	"github.com/droxer/RxGo/pkg/observable"
)

// NewReactivePublisher creates a new Publisher with Reactive Streams compliance
func NewReactivePublisher[T any](onSubscribe func(ctx context.Context, sub publisher.ReactiveSubscriber[T])) publisher.Publisher[T] {
	return publisher.NewReactivePublisher(onSubscribe)
}

// FromSlice creates a Publisher from a slice of values
func FromSlice[T any](items []T) publisher.Publisher[T] {
	return publisher.FromSlice(items)
}

// RangePublisher creates a Publisher that emits a range of integers
func RangePublisher(start, count int) publisher.Publisher[int] {
	return publisher.RangePublisher(start, count)
}

// Create creates an Observable using the legacy API
func Create[T any](onSubscribe func(ctx context.Context, sub observable.Subscriber[T])) *observable.Observable[T] {
	return observable.Create(onSubscribe)
}

// Just creates an Observable that emits the provided values
func Just[T any](values ...T) *observable.Observable[T] {
	return observable.Just(values...)
}

// Range creates an Observable that emits integers in the specified range
func Range(start, count int) *observable.Observable[int] {
	return observable.Range(start, count)
}
