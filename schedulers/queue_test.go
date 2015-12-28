package schedulers_test

import (
	queue "github.com/droxer/RxGo/schedulers"
	"testing"
)

var items = []int{9, 2, 7, 10, 88, 16}

func TestAdd(t *testing.T) {
	q := queue.NewQueue()
	for _, v := range items {
		q.Add(v)
	}

	if !q.Contains(10) {
		t.Error("expected contains 10 ")
	}
}

func TestPoll(t *testing.T) {
	q := queue.NewQueue()

	for _, v := range items {
		q.Add(v)
	}

	actual := q.Poll()

	if actual != 9 {
		t.Errorf("expected 9, actual is %d \n", actual)
	}
}

func TestEmpty(t *testing.T) {
	q := queue.NewQueue()

	if q.Poll() != nil {
		t.Error("expected return nil for polling empty queue")
	}

	if q.Peek() != nil {
		t.Error("expected return nil for peeking empty queue")
	}

	if q.Contains(10) {
		t.Error("expected always return false for empty queue")
	}
}

func TestPeek(t *testing.T) {
	q := queue.NewQueue()

	for _, v := range items {
		q.Add(v)
	}

	actual := q.Peek()

	if actual != 16 {
		t.Errorf("expected 16, actual is %d \n", actual)
	}

	if q.Len() != 6 {
		t.Errorf("expected 6, actual is %d \n", q.Len())
	}
}
