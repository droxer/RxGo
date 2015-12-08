package RxGo

import (
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	Call      runnable
	InitDelay time.Duration
	Period    time.Duration
}

var id int32 = 0
var mutex = &sync.Mutex{}

func ID() int32 {
	mutex.Lock()
	defer mutex.Unlock()

	atomic.AddInt32(&id, 1)
	atomic.LoadInt32(&id)
	return id
}

type Executor struct {
	id      int32
	pool    chan Task
	running bool
}

func NewExecutor() *Executor {
	return &Executor{
		pool: make(chan Task),
		id:   ID(),
	}
}

func (e *Executor) Start() {
	e.running = true
	go func() {
		select {
		case t, more := <-e.pool:
			if more {
				time.Sleep(t.InitDelay)
				t.Call()
			}

			if t.Period != 0 {
				go func() {
					for e.running {
						time.Sleep(t.Period)
						t.Call()
					}
				}()
			}
		}
	}()
}

func (e *Executor) Stop() {
	close(e.pool)
	e.running = false
}

func (e *Executor) Submit(t Task) {
	e.pool <- t
}
