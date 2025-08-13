package scheduler

import (
	"time"
)

// CachedThreadScheduler provides a dynamic thread pool with caching
// Workers are created on-demand and cached for reuse based on TTL
type CachedThreadScheduler struct {
	workerPool *cachedThreadPool
	jobQueue   chan job
	quit       chan bool
	started    bool
}

func NewCachedThreadScheduler(ttl time.Duration) Scheduler {
	return &CachedThreadScheduler{
		workerPool: newCachedThreadPool(ttl),
		jobQueue:   make(chan job),
		quit:       make(chan bool),
	}
}

func (cts *CachedThreadScheduler) Start() {
	if cts.started {
		return
	}
	cts.started = true
	go cts.run()
}

func (cts *CachedThreadScheduler) Stop() {
	if !cts.started {
		return
	}
	cts.started = false
	close(cts.quit)
}

func (cts *CachedThreadScheduler) Schedule(run Runnable) {
	job := job{
		run: run,
	}

	// Auto-start if not already started
	if !cts.started {
		cts.Start()
	}

	cts.jobQueue <- job
}

func (cts *CachedThreadScheduler) ScheduleAt(run Runnable, delay time.Duration) {
	job := job{
		run:   run,
		delay: delay,
	}

	// Auto-start if not already started
	if !cts.started {
		cts.Start()
	}

	cts.jobQueue <- job
}

func (cts *CachedThreadScheduler) run() {
	for {
		select {
		case job := <-cts.jobQueue:
			cts.workerPool.get() <- job
		case <-cts.quit:
			return
		}
	}
}

// IOScheduler returns a CachedThreadScheduler optimized for IO-bound tasks
func IOScheduler() Scheduler {
	return NewCachedThreadScheduler(120 * time.Second)
}

// NewThreadScheduler creates a new thread for each task
func NewThreadScheduler() Scheduler {
	return NewCachedThreadScheduler(0)
}

// cachedThreadPool for CachedThreadScheduler
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

// threadWorker for CachedThreadScheduler
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
		// Never expire if TTL is 0
		if t.ttl > 0 {
			t.timer.Reset(t.ttl)
		}

		for {
			t.jobChanQueue <- t.jobChan

			select {
			case job := <-t.jobChan:
				time.Sleep(job.delay)
				job.run()
				// Reset timer only if TTL > 0
				if t.ttl > 0 {
					t.timer.Reset(t.ttl)
				}
			case <-t.timer.C:
				if t.ttl > 0 {
					t.stop()
					return
				}
			case <-t.quit:
				return
			}
		}
	}()
}

func (t *threadWorker) stop() {
	t.quit <- true
}
