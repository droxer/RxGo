package RxGo

import (
	"time"
)

type Scheduler interface {
	LifeCycle
	CreateWorker() Worker
}

type LifeCycle interface {
	Start()
	Stop()
}

type Worker interface {
	Subscription
	Schedule(ac action0) Subscription
	ScheduleAt(ac action0, delay time.Duration) Subscription
	SchedulePeriodically(ac action0, initDelay, period time.Duration) Subscription
}

type periodicallyScheduler struct{}

func (p *periodicallyScheduler) SchedulePeriodically(ac action0, initDelay, period time.Duration) Subscription {
	return nil
}
