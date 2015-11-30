package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type eventLoopScheduler struct {
}

func (e *eventLoopScheduler) CreateWorker() Worker {
	return &eventLoopWorker{}
}

type eventLoopWorker struct {
	periodicallyScheduler
}

func (e *eventLoopWorker) UnSubscribe() {
}

func (e *eventLoopWorker) IsSubscribed() bool {
	return false
}

func (e *eventLoopWorker) Schedule(ac action0) rx.Subscription {
	return nil
}

func (e *eventLoopWorker) ScheduleAt(ac action0, delay int, unit time.Duration) rx.Subscription {
	return nil
}
