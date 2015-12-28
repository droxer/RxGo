package schedulers

import (
	"time"
)

type threadPoolScheduler struct {
	workers    []*threadWorker
	jobQueue   chan job
	workerPool chan chan job
	quit       chan bool
	ttl        time.Duration
}

func newThreadPoolScheduler(ttl time.Duration) *threadPoolScheduler {
	return &threadPoolScheduler{
		workers:    make([]*threadWorker, 0),
		jobQueue:   make(chan job),
		workerPool: make(chan chan job),
		quit:       make(chan bool),
		ttl:        ttl,
	}
}

func (tps *threadPoolScheduler) Start() {
	go tps.run()
}

func (tps *threadPoolScheduler) Stop() {
	for _, worker := range tps.workers {
		worker.stop()
	}
}

func (tps *threadPoolScheduler) Schedule(run Runnable) {
	job := job{
		run: run,
	}

	tps.jobQueue <- job
}

func (tps *threadPoolScheduler) ScheduleAt(run Runnable, delay time.Duration) {
}

func (tps *threadPoolScheduler) run() {
	for {
		for job := range tps.jobQueue {
			select {
			case jobChan := <-tps.workerPool:
				jobChan <- job
			default:
				worker := newThreadWorker(tps.ttl, tps.workerPool)
				tps.workers = append(tps.workers, worker)
				worker.start()
				worker.jobChan <- job
			}
		}
	}
}

type threadWorker struct {
	ttl        time.Duration
	timer      *time.Timer
	workerPool chan chan job
	jobChan    chan job
	quit       chan bool
}

func newThreadWorker(ttl time.Duration, workerPool chan chan job) *threadWorker {
	return &threadWorker{
		ttl:        ttl,
		timer:      time.NewTimer(ttl),
		workerPool: workerPool,
		jobChan:    make(chan job),
		quit:       make(chan bool),
	}
}

func (t *threadWorker) start() {
	go func() {
		for {
			t.workerPool <- t.jobChan

			select {
			case job := <-t.jobChan:
				time.Sleep(job.delay)
				job.run()
				t.timer.Reset(t.ttl)
			case <-t.timer.C:
				return
			case <-t.quit:
				return
			}
		}
	}()
}

func (t *threadWorker) stop() {
	t.quit <- true
}
