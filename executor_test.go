package RxGo_test

import (
	rx "github.com/droxer/RxGo"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestExecutorDelayTask(t *testing.T) {
	var counter int32 = 0

	wg.Add(1)

	task := rx.Task{
		Call: func() {
			atomic.AddInt32(&counter, 1)
			wg.Done()
		},
		InitDelay: time.Microsecond * 200,
	}

	executor := &rx.Executor{
		Pool: make(chan rx.Task),
	}

	executor.Start()
	defer executor.Stop()

	now := time.Now()
	executor.Submit(task)

	wg.Wait()

	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("expected 1, actual is %d", counter)

	}

	since := time.Since(now).Nanoseconds()
	if since < 200000 {
		t.Errorf("expected take more than 200 microsecond, actual taken %d", since/1000)
	}
}

func TestExecutorPeriodicTask(t *testing.T) {
	var counter int32 = 0

	wg.Add(10)

	task := rx.Task{
		Call: func() {
			atomic.AddInt32(&counter, 1)
			wg.Done()
		},
		InitDelay: time.Millisecond * 100,
		Period:    time.Millisecond * 100,
	}

	executor := &rx.Executor{
		Pool: make(chan rx.Task),
	}

	now := time.Now()
	executor.Start()
	defer executor.Stop()

	executor.Submit(task)

	wg.Wait()
	since := time.Since(now).Seconds()

	if atomic.LoadInt32(&counter) != 10 {
		t.Errorf("expected 1, actual is %d", counter)
	}

	if since < 1 {
		t.Errorf("expected take more than 200 microsecond, actual taken %d", since/1000)
	}
}
