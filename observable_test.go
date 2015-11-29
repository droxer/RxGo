package RxGo_test

import (
	"fmt"
	rx "github.com/droxer/RxGo"
	"testing"
)

type TestSubscriber struct {
	value int
}

func (s *TestSubscriber) OnNext(next interface{}) {
	if v, ok := next.(int); ok {
		s.value += v
		fmt.Printf("++ %d\n", v)
	}
}

func (s *TestSubscriber) OnCompleted() {
	fmt.Printf("total: %d\n", s.value)
}

func (s *TestSubscriber) OnError(e error) {
}

func TestCreateObservable(t *testing.T) {
	sub := &TestSubscriber{10}

	observable := rx.Create(func(sub rx.Subscriber) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)
		}
		sub.OnCompleted()
	})

	observable.Subscribe(sub)

	if sub.value != 55 {
		t.Errorf("expected 55, got %d", sub.value)
	}
}

func ExampleCreateObservable() {
	sub := &TestSubscriber{0}

	observable := rx.Create(func(sub rx.Subscriber) {
		for i := 0; i < 10; i++ {
			sub.OnNext(i)
		}
		sub.OnCompleted()
	})

	observable.Subscribe(sub)
	//Output: ++ 0
	//++ 1
	//++ 2
	//++ 3
	//++ 4
	//++ 5
	//++ 6
	//++ 7
	//++ 8
	//++ 9
	//total: 45
}
