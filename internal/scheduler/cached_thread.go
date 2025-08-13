package scheduler

import (
	"fmt"
	"time"
)

type CachedThreadScheduler struct {
	workerPool *cachedThreadPool
	jobQueue   chan job
	quit       chan bool
}

func NewCachedThreadScheduler(ttl time.Duration) *CachedThreadScheduler {
	return &CachedThreadScheduler{
		workerPool: newCachedThreadPool(ttl),
		jobQueue:   make(chan job),
		quit:       make(chan bool),
	}
}

func (cts *CachedThreadScheduler) Start() {
	go cts.run()
}

func (cts *CachedThreadScheduler) Stop() {
}

func (cts *CachedThreadScheduler) Schedule(run Runnable) {
	job := job{
		run: run,
	}

	cts.jobQueue <- job
}

func (cts *CachedThreadScheduler) ScheduleAt(run Runnable, delay time.Duration) {
}

func (cts *CachedThreadScheduler) run() {
	for job := range cts.jobQueue {
		cts.workerPool.get() <- job
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
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("Thread worker recovered from panic: %v\n", r)
						}
					}()
					job.run()
				}()
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
