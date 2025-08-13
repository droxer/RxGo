package rx

import (
	"runtime"
	"sync"
)

// Scheduler defines the interface for scheduling work
type Scheduler interface {
	Schedule(task func())
}

// ComputationScheduler provides a fixed thread pool based on CPU cores
type ComputationScheduler struct {
	workers int
}

func NewComputationScheduler() *ComputationScheduler {
	return &ComputationScheduler{
		workers: maxParallelism(),
	}
}

func (c *ComputationScheduler) Schedule(task func()) {
	go task() // Simplified for now
}

// IOScheduler provides a cached thread pool for I/O operations
type IOScheduler struct{}

func NewIOScheduler() *IOScheduler {
	return &IOScheduler{}
}

func (i *IOScheduler) Schedule(task func()) {
	go task()
}

// NewThreadScheduler creates a new goroutine for each task
type NewThreadScheduler struct{}

func NewNewThreadScheduler() *NewThreadScheduler {
	return &NewThreadScheduler{}
}

func (n *NewThreadScheduler) Schedule(task func()) {
	go task()
}

// ImmediateScheduler runs tasks immediately on current goroutine
type ImmediateScheduler struct{}

func NewImmediateScheduler() *ImmediateScheduler {
	return &ImmediateScheduler{}
}

func (i *ImmediateScheduler) Schedule(task func()) {
	task()
}

// SingleThreadScheduler runs tasks sequentially on a single goroutine
type SingleThreadScheduler struct {
	tasks chan func()
	wg    sync.WaitGroup
}

func NewSingleThreadScheduler() *SingleThreadScheduler {
	s := &SingleThreadScheduler{
		tasks: make(chan func(), 100),
	}
	s.wg.Add(1)
	go s.run()
	return s
}

func (s *SingleThreadScheduler) Schedule(task func()) {
	select {
	case s.tasks <- task:
	default:
		// Channel full, skip task
	}
}

func (s *SingleThreadScheduler) run() {
	defer s.wg.Done()
	for task := range s.tasks {
		task()
	}
}

func (s *SingleThreadScheduler) Close() {
	close(s.tasks)
	s.wg.Wait()
}

// TrampolineScheduler batches tasks for later execution
type TrampolineScheduler struct {
	tasks []func()
}

func NewTrampolineScheduler() *TrampolineScheduler {
	return &TrampolineScheduler{}
}

func (t *TrampolineScheduler) Schedule(task func()) {
	t.tasks = append(t.tasks, task)
}

func (t *TrampolineScheduler) Execute() {
	for len(t.tasks) > 0 {
		task := t.tasks[0]
		t.tasks = t.tasks[1:]
		task()
	}
}

func (t *TrampolineScheduler) Clear() {
	t.tasks = nil
}

// Predefined schedulers
var (
	Computation = NewComputationScheduler()
	IO          = NewIOScheduler()
)

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCpu := runtime.NumCPU()
	if maxProcs < numCpu {
		return maxProcs
	}
	return numCpu
}