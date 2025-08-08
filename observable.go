package RxGo

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/droxer/RxGo/schedulers"
)

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
	
	// Recover from panics and forward as errors
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("observable panic: %v\n%s", r, debug.Stack())
			sub.OnError(err)
		}
	}()
	
	o.onSubscribe(ctx, sub)
}

// ObserveOn schedules emissions on the specified scheduler.
func (o *Observable[T]) ObserveOn(sch schedulers.Scheduler) *Observable[T] {
	return o.lift(&opObserveOn[T]{scheduler: sch})
}

// lift applies an operator to the Observable.
func (o *Observable[T]) lift(op Operator[T]) *Observable[T] {
	return &Observable[T]{
		onSubscribe: func(ctx context.Context, sub Subscriber[T]) {
			transformed := op.Call(ctx, sub)
			transformed.Start()
			o.onSubscribe(ctx, transformed)
		},
	}
}

// Just creates an Observable that emits the provided values.
func Just[T any](values ...T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnCompleted()
		
		for _, v := range values {
			select {
			case <-ctx.Done():
				sub.OnError(ctx.Err())
				return
			default:
				sub.OnNext(v)
			}
		}
	})
}

// Range creates an Observable that emits a sequence of integers.
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
