package observable

import (
	"context"

	"github.com/droxer/RxGo/pkg/scheduler"
)

// Map transforms each value emitted by the Observable using the provided function
func Map[T, R any](source *Observable[T], transform func(T) R) *Observable[R] {
	return Create(func(ctx context.Context, sub Subscriber[R]) {
		source.Subscribe(ctx, &mapSubscriber[T, R]{
			sub:       sub,
			transform: transform,
		})
	})
}

// Filter filters values emitted by the Observable based on a predicate function
func Filter[T any](source *Observable[T], predicate func(T) bool) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		source.Subscribe(ctx, &filterSubscriber[T]{
			sub:       sub,
			predicate: predicate,
		})
	})
}

// ObserveOn schedules the Observable to emit its values on the specified scheduler
// This allows control over which thread/scheduler the emissions occur on
func ObserveOn[T any](source *Observable[T], sched scheduler.Scheduler) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		source.Subscribe(ctx, &observeOnSubscriber[T]{
			scheduler: sched,
			sub:       sub,
			ctx:       ctx,
		})
	})
}

// internal subscriber types for operators

type mapSubscriber[T, R any] struct {
	sub       Subscriber[R]
	transform func(T) R
}

func (m *mapSubscriber[T, R]) Start()            { m.sub.Start() }
func (m *mapSubscriber[T, R]) OnNext(t T)        { m.sub.OnNext(m.transform(t)) }
func (m *mapSubscriber[T, R]) OnError(err error) { m.sub.OnError(err) }
func (m *mapSubscriber[T, R]) OnCompleted()      { m.sub.OnCompleted() }

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
