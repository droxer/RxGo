package observable

import (
	"context"
)

type Subscriber[T any] interface {
	Start()

	OnNext(next T)

	OnCompleted()

	OnError(e error)
}

type OnSubscribe[T any] func(ctx context.Context, sub Subscriber[T])

type Operator[T, R any] func(source *Observable[T]) *Observable[R]
