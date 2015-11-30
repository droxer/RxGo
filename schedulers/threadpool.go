package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type threadPoolScheduler struct {
}

func (t *threadPoolScheduler) CreateWorker() Worker {
	return &threadWorker{}
}

type threadWorker struct {
	periodicallyScheduler
}

func (t *threadWorker) UnSubscribe() {

}

func (t *threadWorker) IsSubscribed() bool {
	return false
}

func (t *threadWorker) Schedule(ac action0) rx.Subscription {
	return nil
}

func (t *threadWorker) ScheduleAt(ac action0, delay int, unit time.Duration) rx.Subscription {
	return nil
}
