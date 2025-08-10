// Package RxGo provides modern Reactive Extensions for Go with full Reactive Streams 1.0.3 compliance.
//
// RxGo offers two compatible APIs:
//
// 1. Observable API - Simple, callback-based reactive programming
// 2. Reactive Streams API - Full specification compliance with backpressure support
//
// Basic usage with the simple Observable API:
//
//	observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
//	    for i := 0; i < 5; i++ {
//	        sub.OnNext(i)
//	    }
//	    sub.OnCompleted()
//	}).Subscribe(context.Background(), subscriber)
//
// Advanced usage with Reactive Streams API:
//
//	publisher := publisher.NewRangePublisher(1, 10)
//	subscriber := &CustomSubscriber[int]{}
//	publisher.Subscribe(context.Background(), subscriber)
//
// Key Features:
//   - Type-safe generics throughout the API
//   - Context-based cancellation support
//   - Backpressure strategies (Buffer, Drop, Latest, Error)
//   - Thread-safe signal delivery
//   - Memory efficient with bounded buffers
//   - Full Reactive Streams 1.0.3 compliance
//
// Backpressure Support:
//
// RxGo implements the Reactive Streams specification with full backpressure support,
// allowing producers and consumers to handle speed mismatches gracefully through
// demand-based flow control.
//
// For more information and examples, see the README.md file and the examples directory.
package RxGo
