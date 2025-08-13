# Backpressure Control

Handle producer/consumer speed mismatches with Reactive Streams 1.0.4 compliant backpressure strategies.

## Overview

RxGo provides full Reactive Streams 1.0.4 compliance with four backpressure strategies to control flow when producers outpace consumers.

## Four Backpressure Strategies

| Strategy | When to Use | Memory | Behavior |
|----------|-------------|--------|----------|
| **Buffer** | Keep all data | High | Queue items |
| **Drop** | High volume | Medium | Drop new items |
| **Latest** | Real-time | Low | Keep newest |
| **Error** | Strict requirements | Low | Signal error |

## Quick Start

```go
import "github.com/droxer/RxGo/pkg/rx/streams"

// Publisher with backpressure
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 100,
})

// Subscriber with demand control
subscriber := &mySubscriber{}
publisher.Subscribe(context.Background(), subscriber)
```

## Usage Examples

### Buffer Strategy
```go
publisher := streams.RangePublisherWithConfig(1, 100, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 50,
})
```

### Drop Strategy
```go
publisher := streams.FromSlicePublisherWithConfig(data, streams.BackpressureConfig{
    Strategy:   streams.Drop,
    BufferSize: 20,
})
```

### Latest Strategy
```go
publisher := streams.RangePublisherWithConfig(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Latest,
    BufferSize: 1,
})
```

### Error Strategy
```go
publisher := streams.RangePublisherWithConfig(1, 100, streams.BackpressureConfig{
    Strategy:   streams.Error,
    BufferSize: 10,
})
```

## Processors for Data Flow

Chain processors to transform data:

```go
// Map integers to strings
mapProcessor := streams.NewMapProcessor(func(i int) string {
    return fmt.Sprintf("item-%d", i)
})

// Filter even numbers
filterProcessor := streams.NewFilterProcessor(func(i int) bool {
    return i%2 == 0
})
```

## Running Examples

### Complete Examples

**Buffer Strategy** - Keeps all items in bounded buffer:
```go
publisher := streams.RangePublishWithBackpressure(1, 10, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 5,
})
// Processes: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
```

**Drop Strategy** - Discards new items when buffer is full:
```go
publisher := streams.RangePublishWithBackpressure(1, 10, streams.BackpressureConfig{
    Strategy:   streams.Drop,
    BufferSize: 3,
})
// Processes: 1, 2, 3 (rest dropped due to slow consumer)
```

**Latest Strategy** - Keeps only the latest item:
```go
publisher := streams.RangePublishWithBackpressure(1, 10, streams.BackpressureConfig{
    Strategy:   streams.Latest,
    BufferSize: 2,
})
// Processes: 9, 10 (only latest items kept)
```

**Error Strategy** - Signals error on overflow:
```go
publisher := streams.RangePublishWithBackpressure(1, 10, streams.BackpressureConfig{
    Strategy:   streams.Error,
    BufferSize: 2,
})
// Processes: 1, 2 then signals error due to overflow
```

## API Reference

### Types
```go
type BackpressureConfig struct {
    Strategy   BackpressureStrategy
    BufferSize int64
}

type BackpressureStrategy int
const (
    Buffer BackpressureStrategy = iota
    Drop
    Latest
    Error
)
```

### Builders
- `RangePublishWithBackpressure(start, count int, config BackpressureConfig) Publisher[int]`
- `FromSlicePublishWithBackpressure[T any](items []T, config BackpressureConfig) Publisher[T]`

### Processors
- `NewMapProcessor[T, R](transform func(T) R) Processor[T, R]`
- `NewFilterProcessor[T](predicate func(T) bool) Processor[T, T]`
- `NewFlatMapProcessor[T, R](transform func(T) Publisher[R]) Processor[T, R]`