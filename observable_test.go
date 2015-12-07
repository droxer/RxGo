package RxGo_test

import (
	rx "github.com/droxer/RxGo"
	"testing"
)

type SampleSubscriber struct {
	value int
}

func (s *SampleSubscriber) OnNext(next interface{}) {
	if v, ok := next.(int); ok {
		s.value += v
	}
}

func (s *SampleSubscriber) OnCompleted() {
}

func (s *SampleSubscriber) OnError(e error) {

}

func TestCreateObservable(t *testing.T) {

	observable := rx.Create(func(sub rx.Subscriber) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)
		}
		sub.OnCompleted()
	})

	sub := &SampleSubscriber{0}
	observable.Subscribe(sub)

	if sub.value != 45 {
		t.Errorf("expected 45, got %d", sub.value)
	}
}
