package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type threadPoolScheduler struct {
}

func (t *threadPoolScheduler) CreateWorker() rx.Worker {
	return &threadWorker{}
}

func (t *threadPoolScheduler) Start() {

}

func (t *threadPoolScheduler) Stop() {

}

type threadWorker struct {
	periodicallyScheduler
}

func (t *threadWorker) UnSubscribe() {

}

func (t *threadWorker) IsSubscribed() bool {
	return false
}

func (t *threadWorker) Schedule(ac rx.Action0) rx.Subscription {
	return nil
}

func (t *threadWorker) ScheduleAt(ac rx.Action0, delay time.Duration) rx.Subscription {
	return nil
}
