package RxGo_test

import (
	rx "github.com/droxer/RxGo"
	"sync"
	"testing"
)

var wg sync.WaitGroup

type SampleSubscriber struct {
	doNext func(next interface{})
}

func (s *SampleSubscriber) Start() {
}

func (s *SampleSubscriber) OnNext(next interface{}) {
	s.doNext(next)
}

func (s *SampleSubscriber) OnCompleted() {
}

func (s *SampleSubscriber) OnError(e error) {
}

func TestObservable(t *testing.T) {
	var counter = 0
	sub := &SampleSubscriber{
		doNext: func(p interface{}) {
			if v, ok := p.(int); ok {
				counter += v
			}
		},
	}

	observable := rx.Create(func(sub rx.Subscriber) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)
		}
		sub.OnCompleted()
	})

	observable.Subscribe(sub)

	if counter != 45 {
		t.Errorf("expected 45, got %d", counter)
	}
}

func TestObservableSchedule(t *testing.T) {
	var counter = 0
	sub := &SampleSubscriber{
		doNext: func(p interface{}) {
			if v, ok := p.(int); ok {
				counter += v
			}
			wg.Done()
		},
	}

	wg.Add(10)

	observable := rx.Create(func(sub rx.Subscriber) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)

		}
		sub.OnCompleted()
	})

	observable.ObserveOn(rx.ComputationScheduler).Subscribe(sub)

	wg.Wait()
	if counter != 45 {
		t.Errorf("expected 45, got %d", counter)
	}
}
