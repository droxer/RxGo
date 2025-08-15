package streams

import (
	"context"
)

// Publisher represents a Reactive Streams Publisher
// Publishers may support multi-threaded access, depending on implementation
type Publisher[T any] interface {
	// Subscribe subscribes a Subscriber to this Publisher
	// The Publisher will signal the Subscriber via its methods
	Subscribe(ctx context.Context, s Subscriber[T])
}

// Subscriber represents a Reactive Streams Subscriber
// Subscriber implementations should be thread-safe
// Subscriber methods must be called sequentially (Rule 2.7)
type Subscriber[T any] interface {
	// OnSubscribe is called when a subscription is established
	// Must be the first signal to the Subscriber
	OnSubscribe(s Subscription)

	// OnNext is called with the next data item
	// Must be preceded by OnSubscribe and followed by OnComplete or OnError
	OnNext(t T)

	// OnError is called with an error if the Publisher fails
	// Must be preceded by OnSubscribe and must be the final signal
	OnError(err error)

	// OnComplete is called when the Publisher has finished
	// Must be preceded by OnSubscribe and must be the final signal
	OnComplete()
}

// Subscription represents the link between Publisher and Subscriber
// It controls the flow of data through backpressure
// Subscription methods must be thread-safe
type Subscription interface {
	// Request requests n items from the Publisher
	// n must be positive (Rule 3.9)
	Request(n int64)

	// Cancel cancels the subscription
	// No further signals will be sent to the Subscriber
	Cancel()
}

// Processor represents a processing stage that is both a Subscriber and Publisher
// This enables building reactive pipelines
// Processors must be thread-safe and respect both Publisher and Subscriber rules
type Processor[T any, R any] interface {
	Subscriber[T]
	Publisher[R]
}
