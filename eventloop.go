package RxGo

import (
	"time"
)

type eventLoopScheduler struct {
}

func (e *eventLoopScheduler) CreateWorker() Worker {
	return &eventLoopWorker{}
}

func (e *eventLoopScheduler) Start() {

}

func (e *eventLoopScheduler) Stop() {

}

type eventLoopWorker struct {
	periodicallyScheduler
}

func (e *eventLoopWorker) UnSubscribe() {
}

func (e *eventLoopWorker) IsSubscribed() bool {
	return false
}

func (e *eventLoopWorker) Schedule(ac action0) Subscription {
	ac()
	return nil
}

func (e *eventLoopWorker) ScheduleAt(ac action0, delay time.Duration) Subscription {
	return nil
}
