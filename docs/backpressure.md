# Backpressure Control

This document demonstrates how to handle producer/consumer speed mismatches using backpressure strategies with RxGo's Reactive Streams implementation.

## Overview

RxGo provides full Reactive Streams 1.0.4 compliance with backpressure support through the `pkg/rx/streams` package. This enables demand-based flow control to prevent memory issues when producers outpace consumers.

## Key Components

- **Publisher[T]** - Type-safe data source with demand control
- **Subscriber[T]** - Complete subscriber interface with lifecycle
- **Subscription** - Request/cancel control with backpressure
- **ReactivePublisher[T]** - Full implementation with backpressure

## Quick Start

```go
import (
    "context"
    "github.com/droxer/RxGo/pkg/rx/streams"
)

// Create publisher with backpressure support
publisher := streams.RangePublisher(1, 1000)

// Subscribe with controlled demand
consumer := &controlledConsumer{demand: 10}
publisher.Subscribe(ctx, consumer)
```

## Examples

### 1. Basic Backpressure

**File**: `examples/backpressure/basic_backpressure.go`

Controlled consumption with explicit demand:

```go
// Controlled consumer with backpressure
type controlledConsumer struct {
    demand    int64
    processed []int
    sub       streams.Subscription
}

func (c *controlledConsumer) OnSubscribe(sub streams.Subscription) {
    c.sub = sub
    sub.Request(c.demand) // Request initial batch
}

func (c *controlledConsumer) OnNext(value int) {
    c.processed = append(c.processed, value)
    
    // Request next batch when current one is processed
    if len(c.processed)%int(c.demand) == 0 {
        c.sub.Request(c.demand)
    }
}
```

**Run**: `go run examples/backpressure/basic_backpressure.go`

### 2. Advanced Patterns

**File**: `examples/backpressure/advanced_backpressure.go`

#### Dynamic Demand Adjustment
```go
// Adjust demand based on processing speed
func calculateDemand() int64 {
    if averageProcessingTime < 100*time.Millisecond {
        return min(currentDemand+5, maxDemand)
    } else if averageProcessingTime > 300*time.Millisecond {
        return max(currentDemand-2, minDemand)
    }
    return currentDemand
}
```

#### Buffer Overflow Protection
```go
// Monitor memory usage and apply backpressure
if currentMemoryUsage >= memoryLimit {
    flushBuffer()
    sub.Request(smallBatchSize)
}
```

**Run**: `go run examples/backpressure/advanced_backpressure.go`

### 3. Usage Patterns

#### Simple Backpressure
```go
// Create publisher with backpressure support
publisher := streams.RangePublisher(1, 100)

// Subscribe with controlled demand
consumer := &controlledConsumer{demand: 10}
publisher.Subscribe(ctx, consumer)
```

#### Batch Processing
```go
// Process in batches of fixed size
publisher := streams.RangePublisher(1, 1000)
consumer := &batchConsumer{batchSize: 50}
publisher.Subscribe(ctx, consumer)
```

#### Rate Limiting
```go
// Process at controlled rate
consumer := &rateLimitedConsumer{
    rate: time.Second,
    sub:  sub,
}
```

## Running All Examples

```bash
# Run all backpressure examples
./examples/backpressure/run_all.sh

# Or run individually
go run examples/backpressure/basic_backpressure.go
go run examples/backpressure/advanced_backpressure.go
```

## Key Concepts Demonstrated

1. **Demand-based flow control** via `Subscription.Request(n int64)`
2. **Memory pressure handling** with buffer management
3. **Dynamic demand adjustment** based on processing metrics
4. **Error recovery** while maintaining flow control
5. **Rate limiting** with controlled consumption

## Best Practices

- **Start with moderate demand** (10-50 items) and adjust based on performance
- **Monitor memory usage** to prevent buffer overflow
- **Implement error handling** with retry mechanisms
- **Use context cancellation** for graceful shutdown
- **Test with different load patterns** to find optimal batch sizes

## Performance Guidelines

| Use Case | Recommended Demand | Memory Impact |
|----------|-------------------|---------------|
| CPU-intensive | 10-20 items | Low |
| I/O-bound | 5-10 items | Medium |
| Network calls | 1-5 items | High |
| Batch processing | 50-100 items | Very High |

## Integration with Schedulers

Combine backpressure with schedulers for optimal performance:

```go
import (
    "github.com/droxer/RxGo/pkg/rx/streams"
    "github.com/droxer/RxGo/pkg/rx/scheduler"
)

// Use Computation scheduler for CPU-bound processing
scheduler := scheduler.Computation
consumer := &cpuBoundConsumer{scheduler: scheduler}
```