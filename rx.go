package RxGo

var ComputationScheduler Scheduler

func init() {
	ComputationScheduler = &eventLoopScheduler{}
}
