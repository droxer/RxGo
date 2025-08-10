// Package rxgo provides the main entry point for the RxGo library.
// This package provides a unified API for both Observable and Reactive Streams patterns.
package rxgo

import (
	"context"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/reactive"
)

// Observable API

// Create creates a new Observable with the given OnSubscribe function
func Create[T any](on observable.OnSubscribe[T]) *observable.Observable[T] {
	return observable.Create(on)
}

// Just creates an Observable that emits the provided values
func Just[T any](values ...T) *observable.Observable[T] {
	return observable.Just(values...)
}

// Range creates an Observable that emits integers in the specified range
func Range(start, count int) *observable.Observable[int] {
	return observable.Range(start, count)
}

// FromSlice creates an Observable from a slice of values
func FromSlice[T any](items []T) *observable.Observable[T] {
	return observable.FromSlice(items)
}

// Empty creates an Observable that completes without emitting any items
func Empty[T any]() *observable.Observable[T] {
	return observable.Empty[T]()
}

// Error creates an Observable that immediately signals an error
func Error[T any](err error) *observable.Observable[T] {
	return observable.Error[T](err)
}

// Never creates an Observable that never signals any event
func Never[T any]() *observable.Observable[T] {
	return observable.Never[T]()
}

// Reactive Streams API

// NewPublisher creates a new Publisher with Reactive Streams compliance
func NewPublisher[T any](onSubscribe func(ctx context.Context, sub reactive.Subscriber[T])) reactive.Publisher[T] {
	return reactive.NewPublisher(onSubscribe)
}

// RangePublisher creates a Publisher that emits a range of integers
func RangePublisher(start, count int) reactive.Publisher[int] {
	return reactive.Range(start, count)
}

// FromSlicePublisher creates a Publisher from a slice of values
func FromSlicePublisher[T any](items []T) reactive.Publisher[T] {
	return reactive.FromSlice(items)
}

// EmptyPublisher creates a Publisher that completes without emitting any items
func EmptyPublisher[T any]() reactive.Publisher[T] {
	return reactive.Empty[T]()
}

// ErrorPublisher creates a Publisher that immediately signals an error
func ErrorPublisher[T any](err error) reactive.Publisher[T] {
	return reactive.Error[T](err)
}

// NeverPublisher creates a Publisher that never signals any event
func NeverPublisher[T any]() reactive.Publisher[T] {
	return reactive.Never[T]()
}

// JustPublisher creates a Publisher that emits the provided values
func JustPublisher[T any](values ...T) reactive.Publisher[T] {
	return reactive.Just(values...)
}