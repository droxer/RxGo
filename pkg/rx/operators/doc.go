// Package operators provides transformation operators for RxGo Observables.
//
// This package contains operators that transform the values emitted by Observables.
// Operators are pure functions that take an Observable as input and return a new
// Observable as output, allowing for chaining operations.
//
// Available Operators:
//
// Map: Transforms each value emitted by an Observable by applying a function
//
//	import "github.com/droxer/RxGo/pkg/rx/operators"
//
//	// Transform each value by doubling it
//	source := rx.Just(1, 2, 3, 4, 5)
//	doubled := operators.Map(source, func(x int) int { return x * 2 })
//
// Filter: Selectively emits values from an Observable based on a predicate function
//
//	// Only emit even numbers
//	source := rx.Just(1, 2, 3, 4, 5, 6)
//	evens := operators.Filter(source, func(x int) bool { return x%2 == 0 })
//
// ObserveOn: Specifies the Scheduler on which an Observable will operate
//
//	// Emit values on a computation scheduler
//	scheduled := operators.ObserveOn(source, scheduler.Computation)
//
// Chaining Operators:
//
// Operators can be chained together to create complex data transformations:
//
//	result := operators.Filter(
//		operators.Map(source, func(x int) int { return x * 2 }),
//		func(x int) bool { return x > 5 },
//	)
//
// This will first double each value, then only emit values greater than 5.
package operators
