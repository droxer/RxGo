package RxGo_test

import (
	"fmt"
	rx "github.com/droxer/RxGo"
	"testing"
)

type SampleSubscriber struct {
	value int
}

func (s *SampleSubscriber) OnNext(next interface{}) {
	if v, ok := next.(int); ok {
		s.value += v
		fmt.Printf("++ %d\n", s.value)
	}
}

func (s *SampleSubscriber) OnCompleted() {
	fmt.Printf("total: %d\n", s.value)
}

func (s *SampleSubscriber) OnError(e error) {

}

func (s *SampleSubscriber) IsSubscribed() bool {
	return true
}

func (s *SampleSubscriber) UnSubscribe() {

}

func (s *SampleSubscriber) Add(sub rx.Subscription) {

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
