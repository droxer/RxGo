package RxGo

import (
	"time"
)

type Scheduler interface {
	CreateWorker() Worker
}

type Worker interface {
	Subscription
	Schedule(ac action0) Subscription
	ScheduleAt(ac action0, delay time.Duration) Subscription
	SchedulePeriodically(ac action0, initDelay, period time.Duration) Subscription
}
