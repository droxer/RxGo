package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type goWorker struct {
}

func (t *goWorker) UnSubscribe() {

}

func (t *goWorker) IsSubscribed() bool {
	return false
}

func (t *goWorker) Schedule(ac rx.Action0) rx.Subscription {
	return nil
}

func (t *goWorker) ScheduleAt(ac rx.Action0, delay time.Duration) rx.Subscription {
	return nil
}
