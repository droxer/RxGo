package schedulers

import (
	rx "github.com/droxer/RxGo"
	"time"
)

type periodicallyScheduler struct{}

func (p *periodicallyScheduler) SchedulePeriodically(ac rx.Action0, initDelay, period time.Duration) rx.Subscription {
	return nil
}

func Computation() rx.Scheduler {
	return &eventLoopScheduler{}
}
