// Package observable provides the Observable API for RxGo.
//
// This package offers a simple, callback-based reactive programming interface
// with full type safety using Go generics. The Observable API is designed for
// ease of use while providing powerful reactive capabilities.
//
// Basic usage:
//
//	obs := observable.Just(1, 2, 3)
//	obs.Subscribe(context.Background(), subscriber)
//
// Features:
// - Type-safe generics throughout API
// - Context-based cancellation support
// - Simple callback-based interface
// - Chainable operations (Map, Filter, etc.)
package observable
