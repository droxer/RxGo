package rx

import (
	"context"
	"time"

	"github.com/droxer/RxGo/pkg/rx/scheduler"
)

func (o *Observable[T]) ObserveOn(scheduler scheduler.Scheduler) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		o.Subscribe(ctx, &observeOnSubscriber[T]{
			scheduler: scheduler,
			sub:       sub,
			ctx:       ctx,
		})
	})
}

func Schedule[T any](value T, delay time.Duration) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		select {
		case <-time.After(delay):
			if ctx.Err() == nil {
				sub.OnNext(value)
				sub.OnComplete()
			} else {
				sub.OnError(ctx.Err())
			}
		case <-ctx.Done():
			sub.OnError(ctx.Err())
		}
	})
}

func Interval(interval time.Duration) *Observable[int] {
	return Create(func(ctx context.Context, sub Subscriber[int]) {
		defer sub.OnComplete()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case <-ticker.C:
				sub.OnNext(count)
				count++
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			}
		}
	})
}

type observeOnSubscriber[T any] struct {
	scheduler scheduler.Scheduler
	sub       Subscriber[T]
	ctx       context.Context
	started   bool
}

func (o *observeOnSubscriber[T]) Start() {
	if !o.started {
		o.started = true
		o.scheduler.Schedule(func() {
			o.sub.Start()
		})
	}
}

func (o *observeOnSubscriber[T]) OnNext(t T) {
	o.scheduler.Schedule(func() {
		if o.ctx.Err() != nil {
			return
		}
		o.sub.OnNext(t)
	})
}

func (o *observeOnSubscriber[T]) OnError(err error) {
	o.scheduler.Schedule(func() {
		if o.ctx.Err() != nil {
			return
		}
		o.sub.OnError(err)
	})
}

func (o *observeOnSubscriber[T]) OnComplete() {
	o.scheduler.Schedule(func() {
		if o.ctx.Err() != nil {
			return
		}
		o.sub.OnComplete()
	})
}
