package schedulers

import (
	"time"
)

type threadPoolScheduler struct {
	workerPool *cachedThreadPool
	jobQueue   chan job
	quit       chan bool
}

func newThreadPoolScheduler(ttl time.Duration) *threadPoolScheduler {
	return &threadPoolScheduler{
		workerPool: newCachedThreadPool(ttl),
		jobQueue:   make(chan job),
		quit:       make(chan bool),
	}
}

func (tps *threadPoolScheduler) Start() {
	go tps.run()
}

func (tps *threadPoolScheduler) Stop() {
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
		select {
		case job := <-tps.jobQueue:
			tps.workerPool.get() <- job
		}
	}
}

type cachedThreadPool struct {
	ttl          time.Duration
	jobChanQueue chan chan job
}

func newCachedThreadPool(ttl time.Duration) *cachedThreadPool {
	return &cachedThreadPool{
		ttl:          ttl,
		jobChanQueue: make(chan chan job, 10),
	}
}

func (c *cachedThreadPool) get() chan job {
	select {
	case jobChan := <-c.jobChanQueue:
		return jobChan
	default:
		worker := newThreadWorker(c.ttl, c.jobChanQueue)
		worker.start()
		return worker.jobChan
	}
}

type threadWorker struct {
	ttl          time.Duration
	timer        *time.Timer
	jobChanQueue chan chan job
	jobChan      chan job
	quit         chan bool
}

func newThreadWorker(ttl time.Duration, jobChanQueue chan chan job) *threadWorker {
	return &threadWorker{
		ttl:          ttl,
		timer:        time.NewTimer(ttl),
		jobChanQueue: jobChanQueue,
		jobChan:      make(chan job),
		quit:         make(chan bool),
	}
}

func (t *threadWorker) start() {
	go func() {
		for {
			t.jobChanQueue <- t.jobChan

			select {
			case job := <-t.jobChan:
				time.Sleep(job.delay)
				job.run()
				t.timer.Reset(t.ttl)
			case <-t.timer.C:
				t.stop()
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
