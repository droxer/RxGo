package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type Scheduler interface {
	CreateWorker() Worker
}

type periodicallyScheduler struct{}

func (p *periodicallyScheduler) SchedulePeriodically(ac action0, delay int, period int, unit time.Duration) rx.Subscription {
	return nil
}

type Worker interface {
	rx.Subscription
	Schedule(ac action0) rx.Subscription
	ScheduleAt(ac action0, delay int, unit time.Duration) rx.Subscription
	SchedulePeriodically(ac action0, delay int, period int, unit time.Duration) rx.Subscription
}

func Computation() Scheduler {
	return &eventLoopScheduler{}
}

func IO() Scheduler {
	return &threadPoolScheduler{}
}
