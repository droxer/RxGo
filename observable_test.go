package RxGo_test

import (
	"fmt"
	rx "github.com/droxer/RxGo"
	"testing"
)

var counter = 0
var sub = &rx.Subscriber{
	OnNext: func(next interface{}) {
		if v, ok := next.(int); ok {
			counter += v
			fmt.Printf("++ %d\n", counter)
		}
	},
	OnCompleted: func() {
		fmt.Printf("total: %d\n", counter)
	},
}

func TestCreateObservable(t *testing.T) {
	observable := rx.Create(func(sub *rx.Subscriber) {
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
