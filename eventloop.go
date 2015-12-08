package RxGo

import (
	"time"
)

type eventLoopScheduler struct {
}

func (e *eventLoopScheduler) CreateWorker() Worker {
	return &eventLoopWorker{
		executor: NewExecutor(),
	}
}

type eventLoopWorker struct {
	executor *Executor
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
