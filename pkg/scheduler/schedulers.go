package scheduler

// Predefined schedulers for common use cases
var (
	// Computation scheduler for CPU-bound work with fixed thread pool
	Computation Scheduler = ComputationScheduler()
	// IO scheduler for IO-bound work with cached thread pool
	IO Scheduler = IOScheduler()
	// Trampoline scheduler for immediate execution on current thread
	Trampoline Scheduler = TrampolineScheduler()
	// NewThread creates a new thread for each task
	NewThread Scheduler = NewThreadScheduler()
	// SingleThread uses a single dedicated thread
	SingleThread Scheduler = SingleThreadScheduler()
)
