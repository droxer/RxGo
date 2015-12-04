package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type eventLoopScheduler struct {
}

func (e *eventLoopScheduler) CreateWorker() rx.Worker {
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

func (e *eventLoopWorker) Schedule(ac rx.Action0) rx.Subscription {
	return nil
}

func (e *eventLoopWorker) ScheduleAt(ac rx.Action0, delay time.Duration) rx.Subscription {
	return nil
}
