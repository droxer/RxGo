// Package scheduler provides internal scheduler implementations for RxGo
package scheduler

import "time"

// Adapter adapts internal scheduler to pkg/observable Scheduler interface
type Adapter struct {
	s Scheduler
}

// NewAdapter creates a new scheduler adapter
func NewAdapter(s Scheduler) *Adapter {
	return &Adapter{s: s}
}

// Schedule schedules a task to run
func (a *Adapter) Schedule(run func()) {
	a.s.Schedule(func() {
		run()
	})
}

// ScheduleAt schedules a task to run after a delay
func (a *Adapter) ScheduleAt(run func(), delay time.Duration) {
	a.s.ScheduleAt(func() {
		run()
	}, delay)
}

// ToObservableScheduler adapts internal scheduler to observable scheduler
func ToObservableScheduler(s Scheduler) *Adapter {
	return NewAdapter(s)
}
