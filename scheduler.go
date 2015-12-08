package RxGo

import (
	"time"
)

type Scheduler interface {
	CreateWorker() Worker
}

type Worker interface {
	Start()
	Stop()
	Schedule(ac runnable)
	ScheduleAt(ac runnable, delay time.Duration)
	SchedulePeriodically(ac runnable, initDelay, period time.Duration)
}
