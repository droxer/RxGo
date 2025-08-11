// Package reactive provides full Reactive Streams 1.0.3 compliance for RxGo.
//
// This package implements the complete Reactive Streams specification with
// backpressure support, demand control, and full interoperability between
// publishers and subscribers.
//
// Basic usage:
//
//	publisher := reactive.Range(1, 10)
//	publisher.Subscribe(context.Background(), subscriber)
//
// Features:
// - Full Reactive Streams 1.0.3 compliance
// - Backpressure support with demand control
// - Type-safe generics throughout API
// - Publisher lifecycle management
// - Subscription-based cancellation
package reactive
