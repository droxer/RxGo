package scheduler

import "time"

type Runnable func()

type Scheduler interface {
	Start()
	Stop()
	Schedule(run Runnable)
	ScheduleAt(run Runnable, delay time.Duration)
}
