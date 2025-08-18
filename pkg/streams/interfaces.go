package streams

import (
	"context"
)

type Publisher[T any] interface {
	Subscribe(ctx context.Context, s Subscriber[T])
}

type Subscriber[T any] interface {
	OnSubscribe(s Subscription)

	OnNext(t T)

	OnError(err error)

	OnComplete()
}

type Subscription interface {
	Request(n int64)

	Cancel()
}

type Processor[T any, R any] interface {
	Subscriber[T]
	Publisher[R]
}
