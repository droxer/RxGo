// Package observable provides the Observable API for backward compatibility.
//
// This package offers a simple, callback-based reactive programming interface
// that predates the Reactive Streams specification. While fully functional,
// new applications should consider using the publisher package for full
// Reactive Streams compliance and backpressure support.
//
// Basic usage:
//
//	// Create an observable
//	obs := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
//	    sub.OnNext(1)
//	    sub.OnNext(2)
//	    sub.OnNext(3)
//	    sub.OnCompleted()
//	})
//
//	// Subscribe with a custom subscriber
//	obs.Subscribe(context.Background(), &MySubscriber[int]{)
//
// Utility functions:
//
//	// From literal values
//	obs := observable.Just(1, 2, 3, 4, 5)
//
//	// From a range
//	obs := observable.Range(1, 10)
//
// Subscriber interface:
//
//	type MySubscriber[T any] struct{}
//
//	func (s *MySubscriber[T]) Start() {}
//	func (s *MySubscriber[T]) OnNext(value T) { fmt.Println(value) }
//	func (s *MySubscriber[T]) OnCompleted() { fmt.Println("Done") }
//	func (s *MySubscriber[T]) OnError(err error) { fmt.Printf("Error: %v\n", err) }
package observable
