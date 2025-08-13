# API Reference

## Core Interfaces

### Unified Rx API (Recommended)

#### Observable[T] struct
```go
type Observable[T any] struct {
    // Contains filtered or unexported fields
}

func (o *Observable[T]) Subscribe(ctx context.Context, sub Subscriber[T])
```

#### Subscriber[T] Interface
```go
type Subscriber[T any] interface {
    Start()
    OnNext(next T)
    OnCompleted()
    OnError(e error)
}
```

## Constants

```go
const (
    rxgo.Unlimited = 1 << 62 // Maximum request value
)
```

## Utility Functions

| Function | Description | Example |
|----------|-------------|---------|
| `rx.Just[T](values...)` | Create from literal values | `rx.Just(1, 2, 3)` |
| `rx.Range(start, count)` | Integer sequence observable | `rx.Range(1, 10)` |
| `rx.Create[T](fn)` | Custom observable creation | `rx.Create(customProducer)` |

## Package Structure

For information about the package layout and API choices, see [architecture.md](./architecture.md).

## Performance Considerations

### Optimization Features
- **Zero-allocation** signal delivery for common cases
- **Lock-free** data structures where possible
- **Context-aware** cancellation to prevent goroutine leaks
- **Bounded buffers** to prevent memory issues
- **Backpressure** to handle producer/consumer speed mismatches

### Best Practices
- Always use context cancellation for long-running streams
- Implement proper backpressure with `Request(n)` calls
- Use bounded buffers to prevent memory exhaustion
- Prefer the Reactive Streams API for production use