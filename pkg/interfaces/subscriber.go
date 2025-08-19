package interfaces

import "context"

// Subscription represents a one-to-one lifecycle of a Subscriber subscribing to a Publisher
type Subscription interface {
	Request(n int64)
	Cancel()
}

// Subscriber represents a subscriber to a reactive stream
type Subscriber[T any] interface {
	OnSubscribe(s Subscription)
	OnNext(t T)
	OnError(e error)
	OnComplete()
}

// ObservableSubscriber represents a subscriber to an observable
type ObservableSubscriber[T any] interface {
	Start()
	OnNext(t T)
	OnError(e error)
	OnCompleted()
}

// Publisher represents a provider of a potentially unbounded number of sequenced elements
type Publisher[T any] interface {
	Subscribe(ctx context.Context, s Subscriber[T])
}

// Observable represents an observable stream
type Observable[T any] interface {
	Subscribe(ctx context.Context, s ObservableSubscriber[T])
}
