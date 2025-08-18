package observable

import (
	"context"
	"sync"

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

// Merge combines multiple Observables into one by merging their emissions
func Merge[T any](sources ...*Observable[T]) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		var wg sync.WaitGroup
		completed := make(chan bool, len(sources))

		sub.Start()

		for _, source := range sources {
			wg.Add(1)
			go func(s *Observable[T]) {
				defer wg.Done()
				s.Subscribe(ctx, &mergeSubscriber[T]{sub: sub})
				completed <- true
			}(source)
		}

		// Wait for all sources to complete
		go func() {
			wg.Wait()
			close(completed)
			sub.OnCompleted()
		}()
	})
}

// Concat emits the values from the first Observable, then the second, and so on
func Concat[T any](sources ...*Observable[T]) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.Start()
		go func() {
			defer sub.OnCompleted()
			for _, source := range sources {
				done := make(chan struct{})
				source.Subscribe(ctx, &concatSubscriber[T]{
					sub:  sub,
					done: done,
				})
				// Wait for this source to complete before moving to the next
				select {
				case <-done:
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				}
			}
		}()
	})
}

// Take emits only the first n values emitted by the source Observable
func Take[T any](source *Observable[T], n int) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		source.Subscribe(ctx, &takeSubscriber[T]{
			sub: sub,
			n:   n,
		})
	})
}

// Skip skips the first n values emitted by the source Observable
func Skip[T any](source *Observable[T], n int) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		source.Subscribe(ctx, &skipSubscriber[T]{
			sub: sub,
			n:   n,
		})
	})
}

// Distinct suppresses duplicate items emitted by the source Observable
func Distinct[T comparable](source *Observable[T]) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		source.Subscribe(ctx, &distinctSubscriber[T]{
			sub:  sub,
			seen: make(map[T]struct{}),
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

type mergeSubscriber[T any] struct {
	sub Subscriber[T]
}

func (m *mergeSubscriber[T]) Start()            { /* No need to start sub */ }
func (m *mergeSubscriber[T]) OnNext(t T)        { m.sub.OnNext(t) }
func (m *mergeSubscriber[T]) OnError(err error) { m.sub.OnError(err) }
func (m *mergeSubscriber[T]) OnCompleted()      { /* Don't complete parent yet */ }

type concatSubscriber[T any] struct {
	sub  Subscriber[T]
	done chan struct{}
}

func (c *concatSubscriber[T]) Start()            { /* No need to start sub */ }
func (c *concatSubscriber[T]) OnNext(t T)        { c.sub.OnNext(t) }
func (c *concatSubscriber[T]) OnError(err error) { c.sub.OnError(err); close(c.done) }
func (c *concatSubscriber[T]) OnCompleted()      { close(c.done) }

type takeSubscriber[T any] struct {
	sub   Subscriber[T]
	n     int
	count int
}

func (t *takeSubscriber[T]) Start() { t.sub.Start() }
func (t *takeSubscriber[T]) OnNext(value T) {
	if t.count < t.n {
		t.sub.OnNext(value)
		t.count++
	}
	if t.count == t.n {
		t.sub.OnCompleted()
	}
}
func (t *takeSubscriber[T]) OnError(err error) { t.sub.OnError(err) }
func (t *takeSubscriber[T]) OnCompleted() {
	if t.count < t.n {
		t.sub.OnCompleted()
	}
}

type skipSubscriber[T any] struct {
	sub   Subscriber[T]
	n     int
	count int
}

func (s *skipSubscriber[T]) Start() { s.sub.Start() }
func (s *skipSubscriber[T]) OnNext(value T) {
	if s.count >= s.n {
		s.sub.OnNext(value)
	} else {
		s.count++
	}
}
func (s *skipSubscriber[T]) OnError(err error) { s.sub.OnError(err) }
func (s *skipSubscriber[T]) OnCompleted()      { s.sub.OnCompleted() }

type distinctSubscriber[T comparable] struct {
	sub  Subscriber[T]
	seen map[T]struct{}
	mu   sync.Mutex
}

func (d *distinctSubscriber[T]) Start() { d.sub.Start() }
func (d *distinctSubscriber[T]) OnNext(value T) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.seen[value]; !exists {
		d.seen[value] = struct{}{}
		d.sub.OnNext(value)
	}
}
func (d *distinctSubscriber[T]) OnError(err error) { d.sub.OnError(err) }
func (d *distinctSubscriber[T]) OnCompleted()      { d.sub.OnCompleted() }
