package observable

import (
	"context"
	"fmt"
	"sync"

	"github.com/droxer/RxGo/pkg/scheduler"
)

func Map[T, R any](source *Observable[T], transform func(T) R) *Observable[R] {
	return Create(func(ctx context.Context, sub Subscriber[R]) {
		err := source.Subscribe(ctx, &mapSubscriber[T, R]{
			sub:       sub,
			transform: transform,
		})
		if err != nil {
			sub.OnError(err)
		}
	})
}

func Filter[T any](source *Observable[T], predicate func(T) bool) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		err := source.Subscribe(ctx, &filterSubscriber[T]{
			sub:       sub,
			predicate: predicate,
		})
		if err != nil {
			sub.OnError(err)
		}
	})
}

func ObserveOn[T any](source *Observable[T], sched scheduler.Scheduler) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		err := source.Subscribe(ctx, &observeOnSubscriber[T]{
			scheduler: sched,
			sub:       sub,
			ctx:       ctx,
		})
		if err != nil {
			sub.OnError(err)
		}
	})
}

func Merge[T any](sources ...*Observable[T]) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var firstError error
		allCompleted := make(chan struct{})
		terminated := false

		sub.Start()

		// Create a thread-safe subscriber that coordinates with the parent
		mergeSub := &mergeSubscriber[T]{
			sub:          sub,
			wg:           &wg,
			mu:           &mu,
			firstError:   &firstError,
			allCompleted: allCompleted,
			terminated:   &terminated,
			totalSources: len(sources),
		}

		// Subscribe to all sources
		for _, source := range sources {
			wg.Add(1)
			go func(s *Observable[T]) {
				defer wg.Done()
				err := s.Subscribe(ctx, mergeSub)
				if err != nil {
					mergeSub.OnError(err)
				}
			}(source)
		}

		// Wait for all sources to complete and then signal parent completion
		go func() {
			wg.Wait()
			mu.Lock()
			err := firstError
			alreadyTerminated := terminated
			// Check if the channel was already closed by an error
			select {
			case <-allCompleted:
				// Channel already closed by an error, nothing to do
				mu.Unlock()
				return
			default:
				close(allCompleted)
				terminated = true
			}
			mu.Unlock()

			// Only signal completion if not already terminated by an error
			if !alreadyTerminated {
				if err != nil {
					sub.OnError(err)
				} else {
					sub.OnComplete()
				}
			}
		}()
	})
}

func Concat[T any](sources ...*Observable[T]) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		sub.Start()
		go func() {
			defer sub.OnComplete()
			for _, source := range sources {
				done := make(chan struct{})
				errorChan := make(chan error, 1)
				concatSub := &concatSubscriber[T]{
					sub:       sub,
					done:      done,
					errorChan: errorChan,
				}
				err := source.Subscribe(ctx, concatSub)
				if err != nil {
					sub.OnError(err)
					return
				}
				select {
				case <-done:
				case err := <-errorChan:
					// Propagate error and stop processing subsequent sources
					sub.OnError(err)
					return
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				}
			}
		}()
	})
}

func Take[T any](source *Observable[T], n int) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		err := source.Subscribe(ctx, &takeSubscriber[T]{
			sub: sub,
			n:   n,
		})
		if err != nil {
			sub.OnError(err)
		}
	})
}

func Skip[T any](source *Observable[T], n int) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		err := source.Subscribe(ctx, &skipSubscriber[T]{
			sub: sub,
			n:   n,
		})
		if err != nil {
			sub.OnError(err)
		}
	})
}

func Distinct[T comparable](source *Observable[T]) *Observable[T] {
	return DistinctWithLimit(source, 10000)
}

func DistinctWithLimit[T comparable](source *Observable[T], maxSize int) *Observable[T] {
	return Create(func(ctx context.Context, sub Subscriber[T]) {
		err := source.Subscribe(ctx, &distinctSubscriber[T]{
			sub:      sub,
			seen:     make(map[T]struct{}),
			maxSize:  maxSize,
			overflow: false,
		})
		if err != nil {
			sub.OnError(err)
		}
	})
}

type mapSubscriber[T, R any] struct {
	sub       Subscriber[R]
	transform func(T) R
}

func (m *mapSubscriber[T, R]) Start()            { m.sub.Start() }
func (m *mapSubscriber[T, R]) OnNext(t T)        { m.sub.OnNext(m.transform(t)) }
func (m *mapSubscriber[T, R]) OnError(err error) { m.sub.OnError(err) }
func (m *mapSubscriber[T, R]) OnComplete()       { m.sub.OnComplete() }

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
func (f *filterSubscriber[T]) OnComplete()       { f.sub.OnComplete() }

type observeOnSubscriber[T any] struct {
	scheduler scheduler.Scheduler
	sub       Subscriber[T]
	ctx       context.Context
	started   bool
	mu        sync.Mutex
}

func (o *observeOnSubscriber[T]) Start() {
	o.mu.Lock()
	defer o.mu.Unlock()

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

type mergeSubscriber[T any] struct {
	sub          Subscriber[T]
	wg           *sync.WaitGroup
	mu           *sync.Mutex
	firstError   *error
	allCompleted chan struct{}
	terminated   *bool
	totalSources int
}

func (m *mergeSubscriber[T]) Start() { /* No need to start sub */ }

func (m *mergeSubscriber[T]) OnNext(t T) {
	m.sub.OnNext(t)
}

func (m *mergeSubscriber[T]) OnError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Only propagate the first error
	if *m.firstError == nil {
		*m.firstError = err
		*m.terminated = true
		// Cancel all other operations
		select {
		case <-m.allCompleted:
			// Already completed
		default:
			// Close the channel to signal completion
			close(m.allCompleted)
		}
		m.sub.OnError(err)
	}
}

func (m *mergeSubscriber[T]) OnComplete() {
	// Completion is handled by the parent goroutine that waits for all sources
	// The individual mergeSubscriber doesn't need to do anything special here
	// The parent goroutine will close the allCompleted channel and signal completion
}

type concatSubscriber[T any] struct {
	sub       Subscriber[T]
	done      chan struct{}
	errorChan chan error
}

func (c *concatSubscriber[T]) Start() { /* No need to start sub */ }
func (c *concatSubscriber[T]) OnNext(t T) {
	c.sub.OnNext(t)
}
func (c *concatSubscriber[T]) OnError(err error) {
	// Send error to error channel - don't close done to avoid race condition
	c.errorChan <- err
}
func (c *concatSubscriber[T]) OnComplete() {
	close(c.done)
}

type takeSubscriber[T any] struct {
	sub   Subscriber[T]
	n     int
	count int
	mu    sync.Mutex
}

func (t *takeSubscriber[T]) Start() { t.sub.Start() }
func (t *takeSubscriber[T]) OnNext(value T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.count < t.n {
		t.sub.OnNext(value)
		t.count++
	}
	if t.count == t.n {
		t.sub.OnComplete()
	}
}
func (t *takeSubscriber[T]) OnError(err error) { t.sub.OnError(err) }
func (t *takeSubscriber[T]) OnComplete() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.count < t.n {
		t.sub.OnComplete()
	}
}

type skipSubscriber[T any] struct {
	sub   Subscriber[T]
	n     int
	count int
	mu    sync.Mutex
}

func (s *skipSubscriber[T]) Start() { s.sub.Start() }
func (s *skipSubscriber[T]) OnNext(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count >= s.n {
		s.sub.OnNext(value)
	} else {
		s.count++
	}
}
func (s *skipSubscriber[T]) OnError(err error) { s.sub.OnError(err) }
func (s *skipSubscriber[T]) OnComplete()       { s.sub.OnComplete() }

type distinctSubscriber[T comparable] struct {
	sub      Subscriber[T]
	seen     map[T]struct{}
	maxSize  int
	overflow bool
	mu       sync.Mutex
}

func (d *distinctSubscriber[T]) Start() { d.sub.Start() }
func (d *distinctSubscriber[T]) OnNext(value T) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Check if we've already hit the memory limit
	if d.overflow {
		// In overflow mode, we can't guarantee distinctness, so we skip all items
		// to prevent memory issues
		return
	}

	// Check if we're about to exceed the memory limit
	if len(d.seen) >= d.maxSize {
		d.overflow = true
		// Signal an error to indicate we can't maintain distinctness
		d.sub.OnError(fmt.Errorf("distinct operator exceeded maximum memory limit of %d items", d.maxSize))
		return
	}

	if _, exists := d.seen[value]; !exists {
		d.seen[value] = struct{}{}
		d.sub.OnNext(value)
	}
}
func (d *distinctSubscriber[T]) OnError(err error) { d.sub.OnError(err) }
func (d *distinctSubscriber[T]) OnComplete()       { d.sub.OnComplete() }
