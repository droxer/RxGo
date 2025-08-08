package RxGo

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/droxer/RxGo/schedulers"
)

type Operator[T any] interface {
	Call(ctx context.Context, sub Subscriber[T]) Subscriber[T]
}

type opObserveOn[T any] struct {
	scheduler schedulers.Scheduler
}

type observeOnSubscriber[T any] struct {
	ctx       context.Context
	scheduler schedulers.Scheduler
	child     Subscriber[T]
}

func (op *opObserveOn[T]) Call(ctx context.Context, sub Subscriber[T]) Subscriber[T] {
	return &observeOnSubscriber[T]{
		ctx:       ctx,
		scheduler: op.scheduler,
		child:     sub,
	}
}

func (o *observeOnSubscriber[T]) Start() {
	o.scheduler.Start()
}

func (o *observeOnSubscriber[T]) OnNext(next T) {
	o.scheduler.Schedule(func() {
		select {
		case <-o.ctx.Done():
			// Context cancelled, skip processing
			return
		default:
			func() {
				defer func() {
					if r := recover(); r != nil {
						err := fmt.Errorf("observeOn panic: %v\n%s", r, debug.Stack())
						o.child.OnError(err)
					}
				}()
				o.child.OnNext(next)
			}()
		}
	})
}

func (o *observeOnSubscriber[T]) OnError(e error) {
	o.child.OnError(e)
}

func (o *observeOnSubscriber[T]) OnCompleted() {
	o.scheduler.Stop()
	select {
	case <-o.ctx.Done():
		// Context cancelled, skip processing
		return
	default:
		o.child.OnCompleted()
	}
}
