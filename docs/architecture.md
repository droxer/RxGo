# Package Structure

This document describes the new clean package structure for RxGo v0.1.1.

## Overview

RxGo has been refactored into a clean, modular structure with no backward compatibility constraints. The new structure provides clear separation between different APIs and internal implementations.

## Package Layout

```
RxGo/
├── pkg/
│   └── rx/             # Unified API (recommended)
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

### 1. Unified Rx API (Recommended)
**Import:** `github.com/droxer/RxGo/pkg/rx`

Provides a clean, unified interface for both Observable and Reactive Streams patterns:

```go
// Observable operations
obs := rx.Just(1, 2, 3)
obs := rx.Range(1, 10)
obs := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) { ... })
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

**New (recommended):**
```go
import "github.com/droxer/RxGo/pkg/rx"

obs := rx.Just(1, 2, 3)  // Clean unified API
```

## Internal Architecture

The internal packages provide the actual implementations but are not exposed in the public API. This allows for:

- **Implementation flexibility** - Internal changes don't break public API
- **Performance optimization** - Can optimize without API changes
- **Clean abstractions** - Public APIs are clean and focused
- **Testing** - Internal packages can be tested independently

## Best Practices

1. **Use `pkg/rx`** for all applications - provides the cleanest unified experience
2. **Never import from `internal/`** - these are implementation details
3. **Use context cancellation** for all long-running operations