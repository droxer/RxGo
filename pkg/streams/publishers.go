package streams

import (
	"context"
	"sync/atomic"
)

type ReactivePublisher[T any] struct {
	onSubscribe func(ctx context.Context, sub Subscriber[T])
}

func NewPublisher[T any](onSubscribe func(ctx context.Context, sub Subscriber[T])) Publisher[T] {
	return &ReactivePublisher[T]{
		onSubscribe: onSubscribe,
	}
}

func (p *ReactivePublisher[T]) Subscribe(ctx context.Context, s Subscriber[T]) {
	if s == nil {
		panic("subscriber cannot be nil")
	}

	subscription := &reactiveSubscription{}
	s.OnSubscribe(subscription)

	go func() {
		p.processWithBackpressure(ctx, s, subscription)
	}()
}

func (p *ReactivePublisher[T]) processWithBackpressure(ctx context.Context, s Subscriber[T], sub *reactiveSubscription) {
	p.onSubscribe(ctx, s)
}

func FromSlicePublisher[T any](items []T) Publisher[T] {
	return NewPublisher(func(ctx context.Context, sub Subscriber[T]) {
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

func RangePublisher(start, count int) Publisher[int] {
	return NewPublisher(func(ctx context.Context, sub Subscriber[int]) {
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

func RangePublishWithBackpressure(start, count int, config BackpressureConfig) Publisher[int] {
	return NewBufferedPublisher(BackpressureConfig{
		Strategy:   config.Strategy,
		BufferSize: config.BufferSize,
	}, func(ctx context.Context, sub Subscriber[int]) {
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

func FromSlicePublishWithBackpressure[T any](items []T, config BackpressureConfig) Publisher[T] {
	return NewBufferedPublisher(BackpressureConfig{
		Strategy:   config.Strategy,
		BufferSize: config.BufferSize,
	}, func(ctx context.Context, sub Subscriber[T]) {
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

type reactiveSubscription struct {
	cancelled atomic.Bool
	requested atomic.Int64
}

func (s *reactiveSubscription) Request(n int64) {
	if n <= 0 {
		return
	}
	if s.cancelled.Load() {
		return
	}
	s.requested.Add(n)
}

func (s *reactiveSubscription) Cancel() {
	s.cancelled.Store(true)
}

type functionalSubscriber[T any] struct {
	onSubscribe func(Subscription)
	onNext      func(T)
	onError     func(error)
	onComplete  func()
}

func (s *functionalSubscriber[T]) OnSubscribe(sub Subscription) {
	if s.onSubscribe != nil {
		s.onSubscribe(sub)
	}
}

func (s *functionalSubscriber[T]) OnNext(next T) {
	if s.onNext != nil {
		s.onNext(next)
	}
}

func (s *functionalSubscriber[T]) OnError(e error) {
	if s.onError != nil {
		s.onError(e)
	}
}

func (s *functionalSubscriber[T]) OnComplete() {
	if s.onComplete != nil {
		s.onComplete()
	}
}

func NewSubscriber[T any](
	onNext func(T),
	onError func(error),
	onComplete func(),
) Subscriber[T] {
	return &functionalSubscriber[T]{
		onSubscribe: func(s Subscription) {
			s.Request(1<<63 - 1)
		},
		onNext:     onNext,
		onError:    onError,
		onComplete: onComplete,
	}
}
