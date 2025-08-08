package RxGo_test

import (
	"context"
	"testing"
	"sync"

	rx "github.com/droxer/RxGo"
	"github.com/droxer/RxGo/schedulers"
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
	var counter = 0
	sub := &SampleSubscriber[int]{
		doNext: func(p int) {
			counter += p
		},
	}

	observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)
		}
		sub.OnCompleted()
	})

	observable.Subscribe(context.Background(), sub)

	if counter != 45 {
		t.Errorf("expected 45, got %d", counter)
	}
}

func TestObservableSchedule(t *testing.T) {
	var counter = 0
	sub := &SampleSubscriber[int]{
		doNext: func(p int) {
			counter += p
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

	observable.ObserveOn(schedulers.Computation).Subscribe(context.Background(), sub)

	wg.Wait()
	if counter != 45 {
		t.Errorf("expected 45, got %d", counter)
	}
}
