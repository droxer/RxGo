package observable

import "sync"

// Scheduler defines an interface for scheduling work
type Scheduler interface {
	Schedule(task func())
}

// ImmediateScheduler runs tasks immediately on the calling goroutine
type ImmediateScheduler struct{}

func NewImmediateScheduler() *ImmediateScheduler {
	return &ImmediateScheduler{}
}

func (s *ImmediateScheduler) Schedule(task func()) {
	task()
}

// NewThreadScheduler runs each task on a new goroutine
type NewThreadScheduler struct{}

func NewNewThreadScheduler() *NewThreadScheduler {
	return &NewThreadScheduler{}
}

func (s *NewThreadScheduler) Schedule(task func()) {
	go task()
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
	go s.run()
	return s
}

func (s *SingleThreadScheduler) Schedule(task func()) {
	select {
	case s.tasks <- task:
	default:
		// Channel full, drop task or handle appropriately
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

// TrampolineScheduler runs tasks on the current goroutine but allows yielding
type TrampolineScheduler struct {
	tasks []func()
}

func NewTrampolineScheduler() *TrampolineScheduler {
	return &TrampolineScheduler{}
}

func (s *TrampolineScheduler) Schedule(task func()) {
	s.tasks = append(s.tasks, task)
}

func (s *TrampolineScheduler) Execute() {
	for len(s.tasks) > 0 {
		task := s.tasks[0]
		s.tasks = s.tasks[1:]
		task()
	}
}

func (s *TrampolineScheduler) Clear() {
	s.tasks = nil
}

// DefaultScheduler is the default scheduler implementation
var DefaultScheduler Scheduler = NewImmediateScheduler()
