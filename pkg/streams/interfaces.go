package streams

import (
	"context"
)

// Publisher is the interface that emits a stream of values according to the
// Reactive Streams specification. Publishers are pull-based and support
// backpressure through the Subscription interface.
//
// Publishers follow the Reactive Streams specification which provides
// fully compliant backpressure support. For simpler use cases without
// backpressure requirements, consider using the Observable API (pkg/observable).
type Publisher[T any] interface {
	// Subscribe subscribes a Subscriber to this Publisher.
	// This method establishes the subscription and begins the flow of values
	// from the Publisher to the Subscriber.
	//
	// Parameters:
	//   ctx: The context for the subscription, used for cancellation
	//   s: The subscriber that will receive values
	//
	// Returns:
	//   error if subscription fails (e.g., nil subscriber)
	Subscribe(ctx context.Context, s Subscriber[T]) error
}

// Subscriber is the interface that receives values from a Publisher.
// It follows the Reactive Streams specification with full backpressure support.
//
// Subscribers are part of the pull-based model where they control the rate
// of consumption through demand requests.
type Subscriber[T any] interface {
	// OnSubscribe is called when the subscription is established.
	// The subscriber should request data via the Subscription parameter.
	// This method is called before any OnNext calls.
	OnSubscribe(s Subscription)

	// OnNext is called for each emitted item.
	// This method may be called zero or more times.
	OnNext(t T)

	// OnError is called when the Publisher terminates with an error.
	// No further calls to OnNext, OnComplete, or OnError will be made after this.
	// OnComplete and OnError are mutually exclusive.
	OnError(err error)

	// OnComplete is called when the Publisher completes successfully.
	// No further calls to OnNext or OnComplete will be made after this.
	// OnComplete and OnError are mutually exclusive.
	OnComplete()
}

// Subscription represents a one-to-one lifecycle of a Subscriber subscribing
// to a Publisher. It is used for requesting data and cancelling the subscription.
type Subscription interface {
	// Request requests up to n additional items from the Publisher.
	// This method can be called multiple times to request more items.
	// A Publisher will emit items up to the requested amount.
	//
	// Parameters:
	//   n: The number of additional items to request
	Request(n int64)

	// Cancel cancels the subscription and instructs the Publisher to stop
	// emitting items. No further calls to the Subscriber will be made.
	Cancel()
}

// Processor represents a processing stage that is both a Subscriber and Publisher.
// It consumes values from an upstream Publisher and emits processed values
// to downstream Subscribers.
type Processor[T any, R any] interface {
	Subscriber[T]
	Publisher[R]
}
