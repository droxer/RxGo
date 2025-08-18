package scheduler

var (
	Computation  Scheduler = ComputationScheduler()
	IO           Scheduler = IOScheduler()
	Trampoline   Scheduler = TrampolineScheduler()
	NewThread    Scheduler = NewThreadScheduler()
	SingleThread Scheduler = SingleThreadScheduler()
)
