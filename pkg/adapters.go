package pkg

import (
	"context"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/streams"
)

func ObservablePublisherAdapter[T any](obs *observable.Observable[T]) streams.Publisher[T] {
	return streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[T]) {
		obs.Subscribe(ctx, &observableSubscriberAdapter[T]{sub: sub})
	})
}

type observableSubscriberAdapter[T any] struct {
	sub streams.Subscriber[T]
}

func (a *observableSubscriberAdapter[T]) Start() {
	a.sub.OnSubscribe(&observableSubscription{})
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

type observableSubscription struct{}

func (r *observableSubscription) Request(n int64) {
}

func (r *observableSubscription) Cancel() {
}

func PublisherToObservableAdapter[T any](publisher streams.Publisher[T]) *observable.Observable[T] {
	return observable.Create(func(ctx context.Context, sub observable.Subscriber[T]) {
		publisher.Subscribe(ctx, &publisherSubscriberAdapter[T]{sub: sub})
	})
}

type publisherSubscriberAdapter[T any] struct {
	sub observable.Subscriber[T]
}

func (a *publisherSubscriberAdapter[T]) OnSubscribe(s streams.Subscription) {
	a.sub.Start()
	s.Request(1<<63 - 1)
}

func (a *publisherSubscriberAdapter[T]) OnNext(t T) {
	a.sub.OnNext(t)
}

func (a *publisherSubscriberAdapter[T]) OnError(err error) {
	a.sub.OnError(err)
}

func (a *publisherSubscriberAdapter[T]) OnComplete() {
	a.sub.OnCompleted()
}
