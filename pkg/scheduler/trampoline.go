package scheduler

import "time"

// trampolineScheduler executes tasks immediately on the calling thread
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

// TrampolineScheduler executes tasks immediately on the calling thread
func TrampolineScheduler() Scheduler {
	return &trampolineScheduler{}
}
