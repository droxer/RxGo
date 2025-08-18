package scheduler

import "time"

type trampolineScheduler struct{}

func (t *trampolineScheduler) Start() {}
func (t *trampolineScheduler) Stop()  {}
func (t *trampolineScheduler) Schedule(run Runnable) {
	run()
}
func (t *trampolineScheduler) ScheduleAt(run Runnable, delay time.Duration) {
	if delay > 0 {
		time.Sleep(delay)
	}
	run()
}

func TrampolineScheduler() Scheduler {
	return &trampolineScheduler{}
}
