package scheduler

import "time"

// Runnable represents a function that can be scheduled
type Runnable func()

// Scheduler provides an abstraction for executing work on different threads
// This can be used to control concurrency and execution context
type Scheduler interface {
	// Start starts the scheduler
	Start()
	// Stop stops the scheduler
	Stop()
	// Schedule schedules the given runnable to be executed
	Schedule(run Runnable)
	// ScheduleAt schedules the given runnable to be executed after a delay
	ScheduleAt(run Runnable, delay time.Duration)
}
