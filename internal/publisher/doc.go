// Package publisher provides the Reactive Streams Publisher implementation with full backpressure support.
//
// This package implements the Reactive Streams 1.0.3 specification for Go, providing:
//
// • Publisher[T]: A type-safe data source that can be subscribed to
// • ReactiveSubscriber[T]: Complete subscriber interface with lifecycle methods
// • Subscription: Request/cancel control with backpressure
// • Processor[T,R]: Transforming publisher
//
// Basic usage:
//
//	// Create a publisher that emits integers 1-10
//	publisher := publisher.NewRangePublisher(1, 10)
//	
//	// Create a subscriber
//	subscriber := &MySubscriber[int]{}
//	
//	// Subscribe with context
//	publisher.Subscribe(context.Background(), subscriber)
//
// Backpressure control:
//
//	subscriber := &ControlledSubscriber[int]{
//	    onSubscribe: func(sub publisher.Subscription) {
//	        // Request 3 items initially
//	        sub.Request(3)
//	    },
//	}
//
// Utility publishers:
//
//	// From slice
//	publisher := publisher.FromSlice([]int{1, 2, 3, 4, 5})
//	
//	// Custom publisher
//	publisher := publisher.NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[string]) {
//	    sub.OnNext("Hello")
//	    sub.OnNext("World")
//	    sub.OnComplete()
//	})
package publisher