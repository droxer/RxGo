package scheduler

import (
	"fmt"
	"time"
)

type eventLoopScheduler struct {
	workers         []*poolWorker
	jobQueue        chan job
	fixedWorkerPool chan chan job
	quit            chan bool
}

func newEventLoopScheduler(maxWorkers int) Scheduler {
	return &eventLoopScheduler{
		workers:         make([]*poolWorker, maxWorkers),
		fixedWorkerPool: make(chan chan job, maxWorkers),
		jobQueue:        make(chan job),
		quit:            make(chan bool),
	}
}

func (s *eventLoopScheduler) Start() {
	for i := 0; i < len(s.workers); i++ {
		s.workers[i] = newPoolWorker(s.fixedWorkerPool)
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

func (s *eventLoopScheduler) Schedule(run Runnable) {
	job := job{
		run: run,
	}
	s.jobQueue <- job
}

func (s *eventLoopScheduler) ScheduleAt(run Runnable, delay time.Duration) {
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
			jobChan := <-s.fixedWorkerPool
			jobChan <- job
		case <-s.quit:
			return
		}
	}
}

type poolWorker struct {
	fixedWorkerPool chan chan job
	jobChan         chan job
	quit            chan bool
}

func newPoolWorker(fixedWorkerPool chan chan job) *poolWorker {
	return &poolWorker{
		fixedWorkerPool: fixedWorkerPool,
		jobChan:         make(chan job),
		quit:            make(chan bool),
	}
}

func (p *poolWorker) start() {
	go func() {
		for {
			p.fixedWorkerPool <- p.jobChan

			select {
			case job := <-p.jobChan:
				time.Sleep(job.delay)
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("Worker recovered from panic: %v\n", r)
						}
					}()
					job.run()
				}()
			case <-p.quit:
				return
			}
		}
	}()
}

func (p *poolWorker) stop() {
	p.quit <- true
}
