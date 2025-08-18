package scheduler

import "runtime"

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCpu := runtime.NumCPU()
	if maxProcs < numCpu {
		return maxProcs
	}
	return numCpu
}
