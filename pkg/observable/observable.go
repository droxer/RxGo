package observable

import (
	"context"
)

// Observable represents a stream of values that can be observed
// It follows the reactive programming pattern where values are emitted over time
type Observable[T any] struct {
	onSubscribe OnSubscribe[T]
}

// Create creates a new Observable with the given subscription function
// This is the primary way to create custom Observables
func Create[T any](on OnSubscribe[T]) *Observable[T] {
	return &Observable[T]{onSubscribe: on}
}

// Subscribe connects a Subscriber to this Observable
// The Subscriber will receive all emitted values until completion or error
func (o *Observable[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	if sub == nil {
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sub.Start()
	o.onSubscribe(ctx, sub)
}

// Just creates an Observable that emits the provided values in sequence
// This is a convenience function for creating Observables from static data
func Just[T any](values ...T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnCompleted()
		for _, value := range values {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(value)
			}
		}
	})
}

// Range creates an Observable that emits a range of integers
// Similar to a for loop but as an Observable stream
func Range(start, count int) *Observable[int] {
	return Create(func(ctx context.Context, sub Subscriber[int]) {
		defer sub.OnCompleted()
		if count <= 0 {
			return
		}
		for i := 0; i < count; i++ {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(start + i)
			}
		}
	})
}

// FromSlice creates an Observable from a slice of values
// This is useful for converting existing data structures to Observable streams
func FromSlice[T any](slice []T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnCompleted()
		for _, value := range slice {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(value)
			}
		}
	})
}

// Empty creates an Observable that completes immediately without emitting any values
// Useful for testing or as a starting point for conditional streams
func Empty[T any]() *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnCompleted()
	})
}

// Error creates an Observable that immediately errors with the given error
// Useful for error handling and testing error scenarios
func Error[T any](err error) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnError(err)
	})
}
