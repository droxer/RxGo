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

func (e *eventLoopWorker) Schedule(ac action0) Subscription {
	go ac()
	return nil
}

func (e *eventLoopWorker) ScheduleAt(ac action0, delay time.Duration) Subscription {
	return nil
}

func (e *eventLoopWorker) SchedulePeriodically(ac action0, initDelay, period time.Duration) Subscription {
	return nil
}
