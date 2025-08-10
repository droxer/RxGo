# Package Structure

This document describes the new clean package structure for RxGo v0.1.1.

## Overview

RxGo has been refactored into a clean, modular structure with no backward compatibility constraints. The new structure provides clear separation between different APIs and internal implementations.

## Package Layout

```
RxGo/
├── pkg/
│   ├── rxgo/           # Unified API (recommended)
│   ├── observable/     # Observable API (simple)
│   └── reactive/       # Reactive Streams API (full spec)
├── internal/
│   ├── publisher/      # Internal publisher implementations
│   ├── scheduler/      # Scheduler implementations  
│   ├── subscriber/     # Subscriber implementations
│   └── backpressure/   # Backpressure strategies
├── examples/
├── benchmarks/
└── docs/
```

## API Choices

### 1. Unified RxGo API (Recommended)
**Import:** `github.com/droxer/RxGo/pkg/rxgo`

Provides a clean, unified interface for both Observable and Reactive Streams patterns:

```go
// Observable operations
obs := rxgo.Just(1, 2, 3)
obs := rxgo.Range(1, 10)
obs := rxgo.Create(func(ctx context.Context, sub rxgo.Subscriber[int]) { ... })

// Reactive Streams operations  
publisher := rxgo.RangePublisher(1, 10)
publisher := rxgo.JustPublisher("hello", "world")
publisher := rxgo.NewPublisher(func(ctx context.Context, sub rxgo.SubscriberReactive[string]) { ... })
```

### 2. Observable API (Simple)
**Import:** `github.com/droxer/RxGo/pkg/observable`

Simple, callback-based API for basic reactive programming:

```go
obs := observable.Just(1, 2, 3)
obs := observable.Range(1, 10)
obs := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) { ... })
```

### 3. Reactive Streams API (Full Spec)
**Import:** `github.com/droxer/RxGo/pkg/reactive`

Full Reactive Streams 1.0.3 compliance with backpressure support:

```go
publisher := reactive.Just(1, 2, 3)
publisher := reactive.Range(1, 10)
publisher := reactive.NewPublisher(func(ctx context.Context, sub reactive.Subscriber[int]) { ... })
```

## Key Changes from Legacy Structure

### Removed Packages
- `pkg/scheduler` → moved to `internal/scheduler`
- `pkg/subscriber` → moved to `internal/subscriber` 
- `pkg/publisher` → moved to `internal/publisher`
- All legacy internal structures

### New Structure Benefits
1. **Clean API boundaries** - No internal package dependencies in public API
2. **Type safety** - Full generics support throughout
3. **Modular design** - Each package has a single, clear purpose
4. **No backward compatibility** constraints - Clean slate for future development
5. **Unified vs specialized** - Choose the right API for your use case

## Migration Guide

### From Legacy to New Structure

**Old (backward compatibility not maintained):**
```go
import "github.com/droxer/RxGo/pkg/observable"

obs := observable.Just(1, 2, 3)
```

**New (recommended):**
```go
import "github.com/droxer/RxGo/pkg/rxgo"

obs := rxgo.Just(1, 2, 3)  // Unified API
// OR
obs := observable.Just(1, 2, 3)  // Direct Observable API
```

## Internal Architecture

The internal packages provide the actual implementations but are not exposed in the public API. This allows for:

- **Implementation flexibility** - Internal changes don't break public API
- **Performance optimization** - Can optimize without API changes
- **Clean abstractions** - Public APIs are clean and focused
- **Testing** - Internal packages can be tested independently

## Best Practices

1. **Use `pkg/rxgo`** for most applications - provides the cleanest experience
2. **Use `pkg/observable`** when you need only the simple Observable API
3. **Use `pkg/reactive`** when you need full Reactive Streams compliance
4. **Never import from `internal/`** - these are implementation details
5. **Use context cancellation** for all long-running operations