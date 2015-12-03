package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type periodicallyScheduler struct{}

func (p *periodicallyScheduler) SchedulePeriodically(ac rx.Action0, delay int, period int, unit time.Duration) rx.Subscription {
	return nil
}

func Computation() rx.Scheduler {
	return &eventLoopScheduler{}
}

func IO() rx.Scheduler {
	return &threadPoolScheduler{}
}
