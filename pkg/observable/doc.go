// Package observable provides the main RxGo API for reactive programming.
//
// This package offers a comprehensive reactive programming library for Go
// with type-safe generics and context-based cancellation support.
//
// Features:
// - Type-safe generics throughout
// - Context-based cancellation support
// - Fluent operator chaining
// - Built-in operators: Map, Filter, ObserveOn
// - Convenient creation methods: Just, Range, FromSlice, Empty, Error
//
// Quick Start:
//
//	import "github.com/droxer/RxGo/pkg/observable"
//	import "github.com/droxer/RxGo/pkg/scheduler"
//
//	// Basic observable
//	obs := observable.Just(1, 2, 3, 4, 5)
//	obs.Subscribe(ctx, subscriber)
//
//	// Using operators
//	transformed := observable.Map(obs, func(x int) int { return x * 2 })
//	filtered := observable.Filter(transformed, func(x int) bool { return x > 5 })
//
//	// Chaining operators
//	result := observable.Chain(obs,
//		func(o *observable.Observable[int]) *observable.Observable[int] {
//			return observable.Map(o, func(x int) int { return x * 2 })
//		},
//		func(o *observable.Observable[int]) *observable.Observable[int] {
//			return observable.Filter(o, func(x int) bool { return x > 5 })
//		},
//	)
package observable
