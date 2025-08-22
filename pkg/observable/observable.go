package observable

import (
	"context"
	"errors"
)

// Observable represents a push-based stream of values.
// It follows the Observable pattern where values are pushed to subscribers
// as they become available. This API is designed for simplicity and ease of use.
//
// For applications requiring backpressure control, consider using the
// Reactive Streams API (pkg/streams) which provides full backpressure support.
//
// Example usage:
//
//	obs := observable.Just(1, 2, 3, 4, 5)
//	obs.Subscribe(context.Background(), subscriber)
type Observable[T any] struct {
	onSubscribe OnSubscribe[T]
}

// Create creates a new Observable from an OnSubscribe function.
// The OnSubscribe function defines how the Observable produces values
// when a Subscriber subscribes to it.
//
// Parameters:
//
//	on: The function that defines how the Observable produces values
//
// Returns:
//
//	A new Observable that will call the provided function when subscribed to
//
// Example:
//
//	obs := observable.Create(func(ctx context.Context, sub Subscriber[int]) {
//	    go func() {
//	        select {
//	        case <-ctx.Done():
//	            sub.OnError(ctx.Err())
//	            return
//	        default:
//	            sub.OnNext(42)
//	            sub.OnComplete()
//	        }
//	    }()
//	})
func Create[T any](on OnSubscribe[T]) *Observable[T] {
	return &Observable[T]{onSubscribe: on}
}

// Subscribe subscribes a Subscriber to this Observable.
// This method establishes the subscription and begins the flow of values
// from the Observable to the Subscriber.
//
// Parameters:
//
//	ctx: The context for the subscription, used for cancellation
//	sub: The subscriber that will receive values
//
// Example:
//
//	obs := observable.Just(1, 2, 3)
//	obs.Subscribe(context.Background(), subscriber)
func (o *Observable[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sub.Start()
	o.onSubscribe(ctx, sub)
	return nil
}

// Just creates an Observable that emits the specified values and then completes.
// This is a convenience function for creating an Observable from a fixed set of values.
//
// Parameters:
//
//	values: The values to emit
//
// Returns:
//
//	An Observable that emits the provided values and then completes
//
// Example:
//
//	obs := observable.Just(1, 2, 3, 4, 5)
//	// Emits: 1, 2, 3, 4, 5 then completes
func Just[T any](values ...T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnComplete()
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

// Range creates an Observable that emits a range of sequential integers.
// The range starts at the specified start value and emits count sequential integers.
//
// Parameters:
//
//	start: The first integer to emit
//	count: The number of integers to emit
//
// Returns:
//
//	An Observable that emits the specified range of integers and then completes
//
// Example:
//
//	obs := observable.Range(1, 5)
//	// Emits: 1, 2, 3, 4, 5 then completes
func Range(start, count int) *Observable[int] {
	return Create(func(ctx context.Context, sub Subscriber[int]) {
		defer sub.OnComplete()
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

// FromSlice creates an Observable that emits the values from a slice and then completes.
// This is a convenience function for creating an Observable from an existing slice.
//
// Parameters:
//
//	slice: The slice of values to emit
//
// Returns:
//
//	An Observable that emits the values from the slice and then completes
//
// Example:
//
//	values := []int{1, 2, 3, 4, 5}
//	obs := observable.FromSlice(values)
//	// Emits: 1, 2, 3, 4, 5 then completes
func FromSlice[T any](slice []T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnComplete()
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

// Empty creates an Observable that completes immediately without emitting any values.
// This is useful for representing empty streams or as a starting point for operators.
//
// Returns:
//
//	An Observable that completes immediately
//
// Example:
//
//	obs := observable.Empty[int]()
//	// Completes immediately without emitting any values
func Empty[T any]() *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnComplete()
	})
}

// Error creates an Observable that terminates with the specified error.
// This is useful for representing streams that start in an error state.
//
// Parameters:
//
//	err: The error to emit
//
// Returns:
//
//	An Observable that terminates with the specified error
//
// Example:
//
//	obs := observable.Error[int](errors.New("something went wrong"))
//	// Terminates immediately with the specified error
func Error[T any](err error) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnError(err)
	})
}

type functionalSubscriber[T any] struct {
	onStart     func()
	onNext      func(T)
	onCompleted func()
	onError     func(error)
}

func (s *functionalSubscriber[T]) Start() {
	if s.onStart != nil {
		s.onStart()
	}
}

func (s *functionalSubscriber[T]) OnNext(next T) {
	if s.onNext != nil {
		s.onNext(next)
	}
}

func (s *functionalSubscriber[T]) OnComplete() {
	if s.onCompleted != nil {
		s.onCompleted()
	}
}

func (s *functionalSubscriber[T]) OnError(e error) {
	if s.onError != nil {
		s.onError(e)
	}
}

// NewSubscriber creates a new Subscriber with the specified callback functions.
// This is a convenience function for creating subscribers without implementing
// the full Subscriber interface.
//
// Parameters:
//
//	onNext: Function called for each emitted value
//	onCompleted: Function called when the Observable completes successfully
//	onError: Function called when the Observable terminates with an error
//
// Returns:
//
//	A new Subscriber that calls the provided functions
//
// Example:
//
//	subscriber := observable.NewSubscriber(
//	    func(value int) { fmt.Println("Received:", value) },
//	    func() { fmt.Println("Completed") },
//	    func(err error) { fmt.Printf("Error: %v\n", err) },
//	)
func NewSubscriber[T any](
	onNext func(T),
	onCompleted func(),
	onError func(error),
) Subscriber[T] {
	// Provide default no-op functions for nil callbacks to prevent nil pointer dereferences
	if onNext == nil {
		onNext = func(T) {}
	}
	if onCompleted == nil {
		onCompleted = func() {}
	}
	if onError == nil {
		onError = func(error) {}
	}

	return &functionalSubscriber[T]{
		onStart:     func() {}, // Default no-op for Start
		onNext:      onNext,
		onCompleted: onCompleted,
		onError:     onError,
	}
}
