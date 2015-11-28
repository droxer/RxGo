package RxGo_test

import (
	"fmt"
	rx "github.com/droxer/RxGo"
	"testing"
)

type simpleOnSubscribe struct {
}

func (s *simpleOnSubscribe) Call(sub rx.Subscriber) {
	for i := 0; i < 10; i++ {
		sub.OnNext(i)
	}
	sub.OnCompleted()
}

type simpleSubscriber struct {
	value int
}

func (s *simpleSubscriber) OnNext(next interface{}) {
	if v, ok := next.(int); ok {
		s.value += v
		fmt.Printf("++ %d\n", v)
	}
}

func (s *simpleSubscriber) OnCompleted() {
	fmt.Printf("total: %d\n", s.value)
}

func (s *simpleSubscriber) OnError(e error) {
}

func TestSimpleObservable(t *testing.T) {
	sub := &simpleSubscriber{10}

	observable := rx.Create(&simpleOnSubscribe{})
	observable.Subscribe(sub)

	if sub.value != 55 {
		t.Errorf("expected 55, got %d", sub.value)
	}
}

func ExampleSimpleObservable() {
	sub := &simpleSubscriber{0}

	observable := rx.Create(&simpleOnSubscribe{})
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
