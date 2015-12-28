package schedulers

import (
	"runtime"
	"time"
)

type Runnable func()

type job struct {
	run   Runnable
	delay time.Duration
}

type Scheduler interface {
	Start()
	Stop()
	Schedule(run Runnable)
	ScheduleAt(run Runnable, delay time.Duration)
}

var (
	Computation Scheduler
	IO          Scheduler
)

func init() {
	Computation = newEventLoopScheduler(maxParallelism())
	IO = newThreadPoolScheduler(time.Second * 120)
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCpu := runtime.NumCPU()
	if maxProcs < numCpu {
		return maxProcs
	}
	return numCpu
}
