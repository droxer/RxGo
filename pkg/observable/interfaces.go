package observable

import (
	"context"
)

// Subscriber is the interface that receives values from an Observable.
// It follows the push-based model where values are pushed to the subscriber
// as they become available. This interface is designed for simplicity and
// ease of use, without built-in backpressure support.
//
// For applications requiring backpressure control, consider using the
// Reactive Streams API (pkg/streams) which provides full backpressure support.
type Subscriber[T any] interface {
	// Start is called when the subscription is established.
	// This method is called before any OnNext calls.
	Start()

	// OnNext is called for each emitted item.
	// This method may be called zero or more times.
	OnNext(next T)

	// OnComplete is called when the Observable completes successfully.
	// No further calls to OnNext or OnComplete will be made after this.
	// OnComplete and OnError are mutually exclusive.
	OnComplete()

	// OnError is called when the Observable terminates with an error.
	// No further calls to OnNext, OnComplete, or OnError will be made after this.
	// OnComplete and OnError are mutually exclusive.
	OnError(e error)
}

// OnSubscribe is a function that defines how an Observable produces values.
// It is called when a Subscriber subscribes to the Observable.
//
// Parameters:
//   - ctx: The context for the subscription, used for cancellation
//   - sub: The subscriber that will receive values
type OnSubscribe[T any] func(ctx context.Context, sub Subscriber[T])

// Operator is a function that transforms one Observable into another.
// Operators are composable and can be chained together to create
// complex data processing pipelines.
//
// Example:
//
//	obs := observable.Just(1, 2, 3, 4, 5)
//	obs = obs.Map(func(x int) int { return x * 2 })
//	obs = obs.Filter(func(x int) bool { return x > 4 })
//	// Result: emits 6, 8, 10
type Operator[T, R any] func(source *Observable[T]) *Observable[R]
