package rx

import (
	"context"
)

type Subscriber[T any] interface {
	Start()
	OnNext(next T)
	OnComplete()
	OnError(e error)
}

type OnSubscribe[T any] func(ctx context.Context, sub Subscriber[T])

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
		defer sub.OnComplete()
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
		defer sub.OnComplete()
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

func Empty[T any]() *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnComplete()
	})
}

func Error[T any](err error) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnError(err)
	})
}

func Never[T any]() *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		<-ctx.Done()
	})
}

func FromSlice[T any](items []T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnComplete()
		for _, item := range items {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(item)
			}
		}
	})
}

func NewSubscriber[T any](
	onNext func(T),
	onComplete func(),
	onError func(error),
) Subscriber[T] {
	return &simpleSubscriber[T]{
		onNext:     onNext,
		onComplete: onComplete,
		onError:    onError,
	}
}

type simpleSubscriber[T any] struct {
	onNext     func(T)
	onComplete func()
	onError    func(error)
	started    bool
}

func (s *simpleSubscriber[T]) Start() { s.started = true }
func (s *simpleSubscriber[T]) OnNext(t T) {
	if s.onNext != nil {
		s.onNext(t)
	}
}
func (s *simpleSubscriber[T]) OnComplete() {
	if s.onComplete != nil {
		s.onComplete()
	}
}
func (s *simpleSubscriber[T]) OnError(err error) {
	if s.onError != nil {
		s.onError(err)
	}
}
