package RxGo

import (
	"time"
)

type eventLoopScheduler struct {
}

func (e *eventLoopScheduler) CreateWorker() Worker {
	return &eventLoopWorker{}
}

type eventLoopWorker struct {
	Executor
}

func (e *eventLoopWorker) Schedule(ac runnable) Subscription {
	go ac()
	return nil
}

func (e *eventLoopWorker) ScheduleAt(ac runnable, delay time.Duration) Subscription {
	return nil
}

func (e *eventLoopWorker) SchedulePeriodically(ac runnable, initDelay, period time.Duration) Subscription {
	return nil
}
