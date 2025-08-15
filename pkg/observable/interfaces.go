package observable

import (
	"context"
)

// Subscriber represents a subscriber that receives values from an Observable
// This is the observer pattern interface for reactive streams
type Subscriber[T any] interface {
	// Start is called when the subscription begins
	Start()

	// OnNext is called with each emitted value
	OnNext(next T)

	// OnCompleted is called when the Observable completes
	OnCompleted()

	// OnError is called when an error occurs
	OnError(e error)
}

// OnSubscribe is the function signature for Observable creation
type OnSubscribe[T any] func(ctx context.Context, sub Subscriber[T])

// Operator represents a transformation function that takes an Observable and returns a new one
type Operator[T, R any] func(source *Observable[T]) *Observable[R]
