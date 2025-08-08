package observable

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/droxer/RxGo/internal/scheduler"
)

// Subscriber is an interface for backward compatibility
type Subscriber[T any] interface {
	Start()
	OnNext(next T)
	OnCompleted()
	OnError(e error)
}

// Create creates a new Observable with the given OnSubscribe function.
func Create[T any](on OnSubscribe[T]) *Observable[T] {
	return &Observable[T]{onSubscribe: on}
}

// OnSubscribe defines the function signature for creating an Observable.
type OnSubscribe[T any] func(ctx context.Context, sub Subscriber[T])

// Observable represents a stream of values that can be observed.
type Observable[T any] struct {
	onSubscribe OnSubscribe[T]
}

// Subscribe starts the Observable and begins emitting values to the Subscriber.
func (o *Observable[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	if sub == nil {
		return
	}

	// Ensure context is not nil
	if ctx == nil {
		ctx = context.Background()
	}

	// Start the subscriber
	sub.Start()

	// Ensure panic recovery
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			sub.OnError(fmt.Errorf("panic in observable: %v", r))
		}
	}()

	// Execute the subscription
	o.onSubscribe(ctx, sub)
}

// Just creates an Observable that emits the provided values.
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

// Range creates an Observable that emits integers in the specified range.
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

// ObserveOn schedules the Observable to emit its values on the specified scheduler.
func (o *Observable[T]) ObserveOn(scheduler scheduler.Scheduler) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		o.Subscribe(ctx, &observeOnSubscriber[T]{
			scheduler: scheduler,
			sub:       sub,
			ctx:       ctx,
		})
	})
}

// observeOnSubscriber wraps a subscriber for scheduling

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

func (o *observeOnSubscriber[T]) OnCompleted() {
	o.scheduler.Schedule(func() {
		if o.ctx.Err() != nil {
			return
		}
		o.sub.OnCompleted()
	})
}
