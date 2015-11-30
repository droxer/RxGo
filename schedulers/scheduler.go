package schedulers

import (
	rx "github.com/droxer/RxGo"
)

type Scheduler interface {
	CreateWorker() Worker
}

type Worker interface {
	rx.Subscription
	Schedule(ac action0)
}

func Computation() Scheduler {
	return &eventLoopScheduler{}
}

func IO() Scheduler {
	return &threadPoolScheduler{}
}
