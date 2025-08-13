// Package rx provides the main RxGo API for reactive programming.
//
// This package offers a clean and simplified interface for working with
// reactive streams in Go. It includes all necessary components:
//
// - Observable: The main reactive type for handling streams of data
// - Scheduler: Threading and execution control
// - Operators: Map, Filter, and other transformation utilities
//
// Usage:
//
//  // Create an observable
//  obs := rx.Just(1, 2, 3, 4, 5)
//
//  // Use with scheduler
//  obs.ObserveOn(rx.Computation).Subscribe(ctx, subscriber)
package rx