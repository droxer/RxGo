package observable

import (
	"context"
)

type Observable[T any] struct {
	onSubscribe OnSubscribe[T]
}

func Create[T any](on OnSubscribe[T]) *Observable[T] {
	return &Observable[T]{onSubscribe: on}
}

func (o *Observable[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	if sub == nil {
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sub.Start()
	o.onSubscribe(ctx, sub)
}

func Just[T any](values ...T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnCompleted()
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

func Range(start, count int) *Observable[int] {
	return Create(func(ctx context.Context, sub Subscriber[int]) {
		defer sub.OnCompleted()
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

func FromSlice[T any](slice []T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnCompleted()
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

func Empty[T any]() *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnCompleted()
	})
}

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

func (s *functionalSubscriber[T]) OnCompleted() {
	if s.onCompleted != nil {
		s.onCompleted()
	}
}

func (s *functionalSubscriber[T]) OnError(e error) {
	if s.onError != nil {
		s.onError(e)
	}
}

func NewSubscriber[T any](
	onNext func(T),
	onCompleted func(),
	onError func(error),
) Subscriber[T] {
	return &functionalSubscriber[T]{
		onStart:     func() {}, // Default no-op for Start
		onNext:      onNext,
		onCompleted: onCompleted,
		onError:     onError,
	}
}
