package streams

import "sync"

// SignalSerializer ensures sequential signaling (Rule 1.3, 2.7)
type SignalSerializer struct {
	mu       sync.Mutex
	terminal bool
}

// Serialize executes the signal function sequentially
func (s *SignalSerializer) Serialize(signal func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.terminal {
		signal()
	}
}

// MarkTerminal marks the serializer as terminal
func (s *SignalSerializer) MarkTerminal() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.terminal = true
}

// IsTerminal checks if the serializer is in terminal state
func (s *SignalSerializer) IsTerminal() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.terminal
}

// ConcurrentAccessValidator validates thread safety
type ConcurrentAccessValidator[T any] struct {
	mu       sync.Mutex
	accesses []string
}

// RecordAccess records an access for validation
func (v *ConcurrentAccessValidator[T]) RecordAccess(access string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.accesses = append(v.accesses, access)
}

// GetAccesses returns all recorded accesses
func (v *ConcurrentAccessValidator[T]) GetAccesses() []string {
	v.mu.Lock()
	defer v.mu.Unlock()
	return append([]string{}, v.accesses...)
}

// ResetAccesses clears all recorded accesses
func (v *ConcurrentAccessValidator[T]) ResetAccesses() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.accesses = nil
}

// ReactiveStreamsValidator provides utilities for validating Reactive Streams compliance
type ReactiveStreamsValidator struct{}

// ValidatePublisher checks publisher compliance
func (v *ReactiveStreamsValidator) ValidatePublisher(p Publisher[any]) error {
	// Implement publisher validation logic
	return nil
}

// ValidateSubscriber checks subscriber compliance
func (v *ReactiveStreamsValidator) ValidateSubscriber(s Subscriber[any]) error {
	// Implement subscriber validation logic
	return nil
}

// ValidateProcessor checks processor compliance
func (v *ReactiveStreamsValidator) ValidateProcessor(p Processor[any, any]) error {
	// Implement processor validation logic
	return nil
}
