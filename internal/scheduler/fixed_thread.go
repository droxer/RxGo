package scheduler

import (
	"fmt"
	"time"
)

type FixedThreadScheduler struct {
	workers         []*poolWorker
	jobQueue        chan job
	fixedWorkerPool chan chan job
	quit            chan bool
}

func NewFixedThreadScheduler(maxWorkers int) Scheduler {
	return &FixedThreadScheduler{
		workers:         make([]*poolWorker, maxWorkers),
		fixedWorkerPool: make(chan chan job, maxWorkers),
		jobQueue:        make(chan job),
		quit:            make(chan bool),
	}
}

func (fts *FixedThreadScheduler) Start() {
	for i := 0; i < len(fts.workers); i++ {
		fts.workers[i] = newPoolWorker(fts.fixedWorkerPool)
		fts.workers[i].start()
	}

	go fts.dispatch()
}

func (fts *FixedThreadScheduler) Stop() {
	fts.quit <- true
	for _, worker := range fts.workers {
		worker.stop()
	}
}

func (fts *FixedThreadScheduler) Schedule(run Runnable) {
	job := job{
		run: run,
	}
	fts.jobQueue <- job
}

func (fts *FixedThreadScheduler) ScheduleAt(run Runnable, delay time.Duration) {
	job := job{
		run:   run,
		delay: delay,
	}
	fts.jobQueue <- job
}

func (fts *FixedThreadScheduler) dispatch() {
	for {
		select {
		case job := <-fts.jobQueue:
			jobChan := <-fts.fixedWorkerPool
			jobChan <- job
		case <-fts.quit:
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
