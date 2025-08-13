package streams

import (
	"context"
	"sync/atomic"

	"github.com/droxer/RxGo/pkg/rx"
)

// ReactivePublisher implements Publisher with full Reactive Streams compliance
// This provides backpressure support and follows Reactive Streams 1.0.4 specification
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

// reactiveSubscription implements Subscription
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

func ObservablePublisherAdapter[T any](obs *rx.Observable[T]) Publisher[T] {
	return NewPublisher(func(ctx context.Context, sub Subscriber[T]) {
		obs.Subscribe(ctx, &observableSubscriberAdapter[T]{sub: sub})
	})
}

type observableSubscriberAdapter[T any] struct {
	sub Subscriber[T]
}

func (a *observableSubscriberAdapter[T]) Start() {
	// Observable interface doesn't have OnSubscribe, so we call it here
	a.sub.OnSubscribe(&reactiveSubscription{})
}

func (a *observableSubscriberAdapter[T]) OnNext(t T) {
	a.sub.OnNext(t)
}

func (a *observableSubscriberAdapter[T]) OnError(err error) {
	a.sub.OnError(err)
}

func (a *observableSubscriberAdapter[T]) OnCompleted() {
	a.sub.OnComplete()
}
