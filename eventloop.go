package RxGo

import (
	_ "fmt"
	"runtime"
	"time"
)

type eventLoopScheduler struct {
	worker []Worker
	rr     int
}

func newEventLoopScheduler() *eventLoopScheduler {
	var numCPUs = runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	els := &eventLoopScheduler{
		worker: make([]Worker, numCPUs),
	}

	for i := 0; i < numCPUs; i++ {
		els.worker[i] = newEventLoopWorker()
	}

	return els
}

func (e *eventLoopScheduler) CreateWorker() Worker {
	e.rr++
	return e.worker[e.rr%len(e.worker)]
}

type eventLoopWorker struct {
	executor *Executor
}

func newEventLoopWorker() *eventLoopWorker {
	return &eventLoopWorker{
		executor: NewExecutor(),
	}
}

func (e *eventLoopWorker) Start() {
	e.executor.Start()
}

func (e *eventLoopWorker) Stop() {
	e.executor.Stop()
}

func (e *eventLoopWorker) Schedule(run runnable) {
	e.executor.Submit(Task{
		Run: run,
	})
}

func (e *eventLoopWorker) ScheduleAt(run runnable, delay time.Duration) {
	e.executor.Submit(Task{
		Run:       run,
		InitDelay: delay,
	})
}

func (e *eventLoopWorker) SchedulePeriodically(run runnable, initDelay, period time.Duration) {
	e.executor.Submit(Task{
		Run:       run,
		InitDelay: initDelay,
		Period:    period,
	})
}
