// Package rx provides the main RxGo API for reactive programming.
//
// This package offers a comprehensive reactive programming library for Go
// with full Reactive Streams 1.0.4 compliance and backpressure support.
//
// Package Structure:
//
// - rx/          : Core Observable API and main entry points
// - rx/scheduler : Advanced threading and execution control
// - rx/streams   : Reactive Streams 1.0.4 compliant Publisher/Subscriber interfaces
// - rx/operators : Data transformation operators like Map, Filter, ObserveOn
//
// Features:
// - Reactive Streams 1.0.4 compliance with backpressure
// - Context-based cancellation support
// - Multiple scheduler types: Computation, IO, NewThread, SingleThread, Trampoline
// - Fixed and cached thread pool implementations
// - Type-safe generics throughout
//
// Quick Start:
//
//	import "github.com/droxer/RxGo/pkg/rx"
//	import "github.com/droxer/RxGo/pkg/rx/operators"
//	import "github.com/droxer/RxGo/pkg/rx/scheduler"
//
//	// Basic observable
//	obs := rx.Just(1, 2, 3, 4, 5)
//	obs.ObserveOn(scheduler.Computation).Subscribe(ctx, subscriber)
//
//	// Using operators
//	transformed := operators.Map(obs, func(x int) int { return x * 2 })
//	filtered := operators.Filter(transformed, func(x int) bool { return x > 5 })
//
//	// Reactive Streams Publisher
//	import "github.com/droxer/RxGo/pkg/rx/streams"
//	publisher := streams.RangePublisher(1, 10)
//	publisher.Subscribe(ctx, subscriber)
package rx
