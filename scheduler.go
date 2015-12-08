package RxGo

import (
	"time"
)

type Scheduler interface {
	CreateWorker() Worker
}

type Worker interface {
	Subscription
	Schedule(ac runnable) Subscription
	ScheduleAt(ac runnable, delay time.Duration) Subscription
	SchedulePeriodically(ac runnable, initDelay, period time.Duration) Subscription
}
