package schedulers

import (
	"time"
)

type job struct {
	run   Runnable
	delay time.Duration
}

type poolWorker struct {
	workerPool chan chan job
	jobChan    chan job
	quit       chan bool
}

func newPoolWorker(workerPool chan chan job) *poolWorker {
	return &poolWorker{
		workerPool: workerPool,
		jobChan:    make(chan job),
		quit:       make(chan bool),
	}
}

func (p *poolWorker) start() {
	go func() {
		for {
			p.workerPool <- p.jobChan

			select {
			case job := <-p.jobChan:
				time.Sleep(job.delay)
				job.run()
			case <-p.quit:
				return
			}
		}
	}()
}

func (p *poolWorker) stop() {
	p.quit <- true
}
