// Package rxgo provides the unified API for RxGo library.
//
// This package offers a clean, unified interface for both Observable and
// Reactive Streams patterns. It serves as the main entry point for users
// who want a consistent API across different reactive programming styles.
//
// Basic usage:
//
//	// Observable API
//	obs := rxgo.Just(1, 2, 3)
//	obs.Subscribe(context.Background(), subscriber)
//
//	// Reactive Streams API
//	publisher := rxgo.RangePublisher(1, 10)
//	publisher.Subscribe(context.Background(), subscriber)
//
// Features:
// - Unified API for both Observable and Reactive Streams
// - Type-safe generics throughout API
// - Context-based cancellation support
// - Seamless integration between API styles
package rxgo
