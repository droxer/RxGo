package observable

import (
	"context"

	"github.com/droxer/RxGo/internal/publisher"
)

// BridgeSubscriber adapts the new Reactive Streams Subscriber to the old Subscriber interface
type BridgeSubscriber[T any] struct {
	oldSubscriber Subscriber[T]
	subscription  publisher.Subscription
}

func NewBridgeSubscriber[T any](oldSub Subscriber[T]) *BridgeSubscriber[T] {
	return &BridgeSubscriber[T]{
		oldSubscriber: oldSub,
	}
}

func (bs *BridgeSubscriber[T]) OnSubscribe(s publisher.Subscription) {
	bs.subscription = s
	bs.oldSubscriber.Start()
	// Auto-request unlimited for backward compatibility
	s.Request(1 << 62) // Max int64 value
}

func (bs *BridgeSubscriber[T]) OnNext(t T) {
	bs.oldSubscriber.OnNext(t)
}

func (bs *BridgeSubscriber[T]) OnError(err error) {
	bs.oldSubscriber.OnError(err)
}

func (bs *BridgeSubscriber[T]) OnComplete() {
	bs.oldSubscriber.OnCompleted()
}

// ObservablePublisher bridges old Observable to new Publisher
func ObservablePublisher[T any](obs *Observable[T]) publisher.Publisher[T] {
	return publisher.NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[T]) {
		// Create bridge subscriber to convert ReactiveSubscriber to old Subscriber
		bridge := &bridgeFromReactive[T]{
			reactiveSubscriber: sub,
			ctx:                ctx,
		}

		// Use old Observable with the bridge
		obs.onSubscribe(ctx, bridge)
	})
}

// bridgeFromReactive adapts ReactiveSubscriber to old Subscriber interface
type bridgeFromReactive[T any] struct {
	reactiveSubscriber publisher.ReactiveSubscriber[T]
	ctx                context.Context
	subscription       publisher.Subscription
}

func (br *bridgeFromReactive[T]) Start() {
	// This is called by the old Observable system
	// We need to bridge from old to new
	br.reactiveSubscriber.OnSubscribe(br.subscription)
}

func (br *bridgeFromReactive[T]) OnNext(t T) {
	br.reactiveSubscriber.OnNext(t)
}

func (br *bridgeFromReactive[T]) OnError(err error) {
	br.reactiveSubscriber.OnError(err)
}

func (br *bridgeFromReactive[T]) OnCompleted() {
	br.reactiveSubscriber.OnComplete()
}

// PublisherObservable bridges new Publisher to old Observable
type PublisherObservable[T any] struct {
	publisher publisher.Publisher[T]
}

func NewPublisherObservable[T any](p publisher.Publisher[T]) *PublisherObservable[T] {
	return &PublisherObservable[T]{publisher: p}
}

func (po *PublisherObservable[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Create bridge subscriber to convert from old Subscriber to new ReactiveSubscriber
	bridge := NewBridgeSubscriber[T](sub)
	po.publisher.Subscribe(ctx, bridge)
}
