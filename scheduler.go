package RxGo

import (
	"time"
)

type Scheduler interface {
	CreateWorker() Worker
}

type Worker interface {
	Subscription
	Schedule(ac Action0) Subscription
	ScheduleAt(ac Action0, delay int, unit time.Duration) Subscription
	SchedulePeriodically(ac Action0, delay int, period int, unit time.Duration) Subscription
}
