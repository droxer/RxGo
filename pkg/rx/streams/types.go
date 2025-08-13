package streams

import (
	"context"
)

// Publisher represents a Reactive Streams Publisher
type Publisher[T any] interface {
	Subscribe(ctx context.Context, s Subscriber[T])
}

// Subscriber represents a Reactive Streams Subscriber
type Subscriber[T any] interface {
	OnSubscribe(s Subscription)
	OnNext(t T)
	OnError(err error)
	OnComplete()
}

// Subscription represents the link between Publisher and Subscriber
type Subscription interface {
	Request(n int64)
	Cancel()
}

// Processor represents a processing stage that is both a Subscriber and Publisher
type Processor[T any, R any] interface {
	Subscriber[T]
	Publisher[R]
}
