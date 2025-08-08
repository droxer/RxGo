package observable_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/droxer/RxGo/internal/scheduler"
	rx "github.com/droxer/RxGo/pkg/rxgo"
)

var wg sync.WaitGroup

type SampleSubscriber[T any] struct {
	doNext func(next T)
}

func (s *SampleSubscriber[T]) Start() {}

func (s *SampleSubscriber[T]) OnNext(next T) {
	s.doNext(next)
}

func (s *SampleSubscriber[T]) OnCompleted() {}

func (s *SampleSubscriber[T]) OnError(e error) {}

func TestObservable(t *testing.T) {
	var counter atomic.Int64
	sub := &SampleSubscriber[int]{
		doNext: func(p int) {
			counter.Add(int64(p))
		},
	}

	observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)
		}
		sub.OnCompleted()
	})

	observable.Subscribe(context.Background(), sub)

	if counter.Load() != 45 {
		t.Errorf("expected 45, got %d", counter.Load())
	}
}

func TestObservableSchedule(t *testing.T) {
	var counter atomic.Int64
	sub := &SampleSubscriber[int]{
		doNext: func(p int) {
			counter.Add(int64(p))
			wg.Done()
		},
	}

	wg.Add(10)

	observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)
		}
		sub.OnCompleted()
	})

	observable.ObserveOn(scheduler.Computation).Subscribe(context.Background(), sub)

	wg.Wait()
	if counter.Load() != 45 {
		t.Errorf("expected 45, got %d", counter.Load())
	}
}
