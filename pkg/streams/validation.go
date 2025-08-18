package streams

import "sync"

type SignalSerializer struct {
	mu       sync.Mutex
	terminal bool
}

func (s *SignalSerializer) Serialize(signal func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.terminal {
		signal()
	}
}

func (s *SignalSerializer) MarkTerminal() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.terminal = true
}

func (s *SignalSerializer) IsTerminal() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.terminal
}

type ConcurrentAccessValidator[T any] struct {
	mu       sync.Mutex
	accesses []string
}

func (v *ConcurrentAccessValidator[T]) RecordAccess(access string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.accesses = append(v.accesses, access)
}

func (v *ConcurrentAccessValidator[T]) GetAccesses() []string {
	v.mu.Lock()
	defer v.mu.Unlock()
	return append([]string{}, v.accesses...)
}

func (v *ConcurrentAccessValidator[T]) ResetAccesses() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.accesses = nil
}

type ReactiveStreamsValidator struct{}

func (v *ReactiveStreamsValidator) ValidatePublisher(p Publisher[any]) error {
	// Implement publisher validation logic
	return nil
}

func (v *ReactiveStreamsValidator) ValidateSubscriber(s Subscriber[any]) error {
	// Implement subscriber validation logic
	return nil
}

func (v *ReactiveStreamsValidator) ValidateProcessor(p Processor[any, any]) error {
	// Implement processor validation logic
	return nil
}
