package pkg

import (
	"context"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/streams"
)

// ObservablePublisherAdapter adapts an Observable to Publisher interface
// This allows interoperability between Observable and Publisher types
func ObservablePublisherAdapter[T any](obs *observable.Observable[T]) streams.Publisher[T] {
	return streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[T]) {
		obs.Subscribe(ctx, &observableSubscriberAdapter[T]{sub: sub})
	})
}

// observableSubscriberAdapter adapts Subscriber to observable.Subscriber
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

// observableSubscription implements Subscription for ObservablePublisherAdapter
type observableSubscription struct{}

func (r *observableSubscription) Request(n int64) {
	// Observable pattern doesn't support backpressure
}

func (r *observableSubscription) Cancel() {
	// Observable pattern uses context cancellation
}

// PublisherToObservableAdapter adapts a Publisher to Observable interface
// This allows interoperability between Publisher and Observable types
func PublisherToObservableAdapter[T any](publisher streams.Publisher[T]) *observable.Observable[T] {
	return observable.Create(func(ctx context.Context, sub observable.Subscriber[T]) {
		publisher.Subscribe(ctx, &publisherSubscriberAdapter[T]{sub: sub})
	})
}

// publisherSubscriberAdapter adapts observable.Subscriber to streams.Subscriber
type publisherSubscriberAdapter[T any] struct {
	sub observable.Subscriber[T]
}

func (a *publisherSubscriberAdapter[T]) OnSubscribe(s streams.Subscription) {
	// Start the observable subscription
	a.sub.Start()

	// Request all items (Observable pattern)
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
