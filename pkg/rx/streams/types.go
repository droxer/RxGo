package streams

import (
	"context"
)

// Publisher represents a Reactive Streams Publisher that can be subscribed to
// This provides full Reactive Streams 1.0.4 compliance with backpressure support
type Publisher[T any] interface {
	Subscribe(ctx context.Context, s Subscriber[T])
}

// Subscriber represents a Reactive Streams Subscriber
// This is the modern, type-safe interface for subscribers
// It follows the Reactive Streams specification 1.0.4
type Subscriber[T any] interface {
	OnSubscribe(s Subscription)
	OnNext(t T)
	OnError(err error)
	OnComplete()
}

// Subscription represents the link between Publisher and Subscriber
// This interface is part of the Reactive Streams specification
type Subscription interface {
	Request(n int64)
	Cancel()
}
