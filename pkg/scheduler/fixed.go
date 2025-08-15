package scheduler

import (
	"time"
)

// job represents a unit of work to be scheduled
type job struct {
	run   Runnable
	delay time.Duration
}

// FixedThreadScheduler provides a fixed-size thread pool for computation tasks
type FixedThreadScheduler struct {
	workers         []*poolWorker
	jobQueue        chan job
	fixedWorkerPool chan chan job
	quit            chan bool
	started         bool
}

// NewFixedThreadScheduler creates a new fixed-size thread scheduler
func NewFixedThreadScheduler(maxWorkers int) Scheduler {
	return &FixedThreadScheduler{
		workers:         make([]*poolWorker, maxWorkers),
		fixedWorkerPool: make(chan chan job, maxWorkers),
		jobQueue:        make(chan job),
		quit:            make(chan bool),
		started:         false,
	}
}

func (fts *FixedThreadScheduler) Start() {
	if fts.started {
		return
	}
	fts.started = true

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

	// Auto-start if not already started
	if !fts.started {
		fts.Start()
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

// ComputationScheduler returns a FixedThreadScheduler optimized for CPU-bound tasks
func ComputationScheduler() Scheduler {
	return NewFixedThreadScheduler(maxParallelism())
}

// SingleThreadScheduler uses a single dedicated thread
func SingleThreadScheduler() Scheduler {
	return NewFixedThreadScheduler(1)
}

// poolWorker for FixedThreadScheduler
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
		defer func() {
			// Clean up when worker exits
			_ = recover()
		}()

		for {
			// Register this worker
			select {
			case p.fixedWorkerPool <- p.jobChan:
				// Successfully registered, wait for job
				select {
				case job := <-p.jobChan:
					time.Sleep(job.delay)
					job.run()
				case <-p.quit:
					return
				}
			case <-p.quit:
				return
			}
		}
	}()
}

func (p *poolWorker) stop() {
	p.quit <- true
}
