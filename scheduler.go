package RxGo

import (
	"time"
)

type runnable func()

type Scheduler interface {
	Start()
	Stop()
	Schedule(run runnable)
	ScheduleAt(run runnable, delay time.Duration)
}

type eventLoopScheduler struct {
	workers    []*poolWorker
	jobQueue   chan job
	workerPool chan chan job
	quit       chan bool
}

func newEventLoopScheduler(maxWorkers int) Scheduler {
	return &eventLoopScheduler{
		workers:    make([]*poolWorker, maxWorkers),
		workerPool: make(chan chan job, maxWorkers),
		jobQueue:   make(chan job),
		quit:       make(chan bool),
	}
}

func (s *eventLoopScheduler) Start() {
	for i := 0; i < len(s.workers); i++ {
		s.workers[i] = newPoolWorker(s.workerPool)
		s.workers[i].start()
	}

	go s.dispatch()
}

func (s *eventLoopScheduler) Stop() {
	s.quit <- true
	for _, worker := range s.workers {
		worker.stop()
	}
}

func (s *eventLoopScheduler) Schedule(run runnable) {
	job := job{
		run: run,
	}
	s.jobQueue <- job
}

func (s *eventLoopScheduler) ScheduleAt(run runnable, delay time.Duration) {
	job := job{
		run:   run,
		delay: delay,
	}
	s.jobQueue <- job
}

func (s *eventLoopScheduler) dispatch() {
	for {
		select {
		case job := <-s.jobQueue:
			jobChan := <-s.workerPool
			jobChan <- job
		case <-s.quit:
			return
		}
	}
}
