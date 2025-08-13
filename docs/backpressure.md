# Backpressure Control

This document demonstrates how to handle producer/consumer speed mismatches using backpressure strategies with RxGo's Reactive Streams implementation.

## Overview

RxGo provides full Reactive Streams 1.0.4 compliance with backpressure support through the `pkg/rx/streams` package. This enables demand-based flow control to prevent memory issues when producers outpace consumers.

## Backpressure Strategies

RxGo implements four reactive backpressure strategies for handling overflow when producers are faster than consumers:

### 1. Buffer Strategy
Maintains a bounded buffer to store items when the consumer is slow.

```go
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 100,
})
```

### 2. Drop Strategy
Discards new items when the buffer is full.

```go
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Drop,
    BufferSize: 50,
})
```

### 3. Latest Strategy
Keeps only the latest item, discarding older ones.

```go
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Latest,
    BufferSize: 1,
})
```

### 4. Error Strategy
Signals an error when the buffer overflows.

```go
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Error,
    BufferSize: 10,
})
```

## Key Components

- **Publisher[T]** - Type-safe data source with demand control
- **Subscriber[T]** - Complete subscriber interface with lifecycle
- **Subscription** - Request/cancel control with backpressure
- **ReactivePublisher[T]** - Full implementation with backpressure
- **BufferedPublisher[T]** - Publisher with configurable backpressure strategies

## Quick Start

```go
import (
    "context"
    "github.com/droxer/RxGo/pkg/rx/streams"
)

// Create publisher with backpressure strategy
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 100,
})

// Subscribe with controlled demand
consumer := &controlledConsumer{demand: 10}
publisher.Subscribe(ctx, consumer)
```

## Examples

### Basic Backpressure with Buffer Strategy

```go
type controlledConsumer struct {
    demand    int64
    processed []int
    sub       streams.Subscription
}

func (c *controlledConsumer) OnSubscribe(sub streams.Subscription) {
    c.sub = sub
    sub.Request(c.demand)
}

func (c *controlledConsumer) OnNext(value int) {
    c.processed = append(c.processed, value)
    
    if len(c.processed)%int(c.demand) == 0 {
        c.sub.Request(c.demand)
    }
}
```

### Drop Strategy for High-Volume Data

```go
// Fast producer, slow consumer - drops excess items
publisher := streams.FromSlicePublisherWithConfig(largeDataSet, streams.BackpressureConfig{
    Strategy:   streams.Drop,
    BufferSize: 20,
})
```

### Latest Strategy for Real-Time Data

```go
// Keep only the most recent data point
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Latest,
    BufferSize: 1,
})
```

### Error Strategy for Strict Requirements

```go
// Fail fast when overflow occurs
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Error,
    BufferSize: 10,
})
```

## Running Examples

```bash
# Run comprehensive backpressure examples
go run examples/backpressure/backpressure_examples.go
```

## Strategy Selection Guide

| Strategy | Use Case | Memory Impact | Data Loss | Error Handling |
|----------|----------|---------------|-----------|----------------|
| **Buffer** | Preserve all data | High | None | Automatic |
| **Drop** | High-volume streams | Medium | New items | Silent |
| **Latest** | Real-time updates | Low | Old items | Silent |
| **Error** | Strict requirements | Low | None | Explicit error |

## Best Practices

- **Buffer**: Use when data completeness is critical
- **Drop**: Use when data volume exceeds processing capacity
- **Latest**: Use for real-time applications where only current state matters
- **Error**: Use when overflow conditions must be explicitly handled

## Integration with Schedulers

Combine backpressure strategies with schedulers for optimal performance:

```go
import (
    "github.com/droxer/RxGo/pkg/rx/streams"
    "github.com/droxer/RxGo/pkg/rx/scheduler"
)

// CPU-intensive processing with buffer strategy
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 50,
})

// Process with computation scheduler
scheduler := scheduler.Computation
// ... integrate with scheduler
```

## API Reference

### BackpressureConfig

```go
type BackpressureConfig struct {
    Strategy   BackpressureStrategy
    BufferSize int64
}
```

### BackpressureStrategy

```go
const (
    Buffer BackpressureStrategy = iota  // Keep all items in bounded buffer
    Drop                               // Discard new items when full
    Latest                             // Keep only latest item
    Error                              // Signal error on overflow
)
```

### Publisher Builders

- `RangePublisherWithConfig(start, count int, config BackpressureConfig) Publisher[int]`
- `FromSlicePublisherWithConfig[T any](items []T, config BackpressureConfig) Publisher[T]`
- `NewBufferedPublisher[T any](config BackpressureConfig, source func(ctx context.Context, sub Subscriber[T])) Publisher[T]`

### Complete Examples

See the `examples/backpressure/` directory for complete working examples of all four strategies:

- `backpressure_examples.go` - Comprehensive demonstration of all strategies