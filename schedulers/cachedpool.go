package schedulers

import (
	"time"
)

type cachedPoolScheduler struct {
	timer      *time.Timer
	jobQueue   chan job
	workerPool chan chan job
	quit       chan bool
}

func newCachedPoolScheduler(maxWorkers int, keepAliveTime time.Duration) *cachedPoolScheduler {
	return &cachedPoolScheduler{
		timer:      time.NewTimer(keepAliveTime),
		workerPool: make(chan chan job, maxWorkers),
		jobQueue:   make(chan job),
		quit:       make(chan bool),
	}
}

func (cps *cachedPoolScheduler) Start() {

}

func (cps *cachedPoolScheduler) Stop() {

}

func (cps *cachedPoolScheduler) Schedule(run Runnable) {

}

func (cps *cachedPoolScheduler) ScheduleAt(run Runnable, delay time.Duration) {

}
