// Package observable provides the main Observable API for RxGo
package observable

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/droxer/RxGo/internal/scheduler"
)

// Subscriber defines the interface for receiving values from an Observable
type Subscriber[T any] interface {
	Start()
	OnNext(next T)
	OnCompleted()
	OnError(e error)
}

// OnSubscribe defines the function signature for creating an Observable
type OnSubscribe[T any] func(ctx context.Context, sub Subscriber[T])

// Observable represents a stream of values that can be observed
// This is the primary API for reactive programming in RxGo
type Observable[T any] struct {
	onSubscribe OnSubscribe[T]
}

// Create creates a new Observable with the given OnSubscribe function
func Create[T any](on OnSubscribe[T]) *Observable[T] {
	return &Observable[T]{onSubscribe: on}
}

// Subscribe starts the Observable and begins emitting values to the Subscriber
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

// Just creates an Observable that emits the provided values
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

// Range creates an Observable that emits integers in the specified range
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

// Empty creates an Observable that completes without emitting any items
func Empty[T any]() *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnCompleted()
	})
}

// Error creates an Observable that immediately signals an error
func Error[T any](err error) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.OnError(err)
	})
}

// Never creates an Observable that never signals any event
func Never[T any]() *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		<-ctx.Done()
	})
}

// FromSlice creates an Observable from a slice of values
func FromSlice[T any](items []T) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		defer sub.OnCompleted()
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

// ObserveOn schedules the Observable to emit its values on the specified scheduler
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

// Map transforms each value emitted by the Observable
func Map[T, R any](source *Observable[T], transform func(T) R) *Observable[R] {
	return Create(func(ctx context.Context, sub Subscriber[R]) {
		source.Subscribe(ctx, &mapSubscriber[T, R]{
			sub:       sub,
			transform: transform,
		})
	})
}

// Filter filters values emitted by the Observable
func Filter[T any](source *Observable[T], predicate func(T) bool) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		source.Subscribe(ctx, &filterSubscriber[T]{
			sub:       sub,
			predicate: predicate,
		})
	})
}

// mapSubscriber implements the mapping transformation
type mapSubscriber[T, R any] struct {
	sub       Subscriber[R]
	transform func(T) R
}

func (m *mapSubscriber[T, R]) Start()            { m.sub.Start() }
func (m *mapSubscriber[T, R]) OnNext(t T)        { m.sub.OnNext(m.transform(t)) }
func (m *mapSubscriber[T, R]) OnError(err error) { m.sub.OnError(err) }
func (m *mapSubscriber[T, R]) OnCompleted()      { m.sub.OnCompleted() }

// filterSubscriber implements the filtering transformation
type filterSubscriber[T any] struct {
	sub       Subscriber[T]
	predicate func(T) bool
}

func (f *filterSubscriber[T]) Start() { f.sub.Start() }
func (f *filterSubscriber[T]) OnNext(t T) {
	if f.predicate(t) {
		f.sub.OnNext(t)
	}
}
func (f *filterSubscriber[T]) OnError(err error) { f.sub.OnError(err) }
func (f *filterSubscriber[T]) OnCompleted()      { f.sub.OnCompleted() }
