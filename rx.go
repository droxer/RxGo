package RxGo

import (
	"runtime"
)

var (
	ComputationScheduler Scheduler
)

func init() {
	ComputationScheduler = newEventLoopScheduler(maxParallelism())
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCpu := runtime.NumCPU()
	if maxProcs < numCpu {
		return maxProcs
	}
	return numCpu
}
