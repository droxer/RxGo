# API Reference

## Core Interfaces

### Observable API (Simple)

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

### Reactive Streams API (Full Spec)

#### Publisher[T] Interface
```go
type Publisher[T any] interface {
    Subscribe(ctx context.Context, s SubscriberReactive[T])
}
```

#### SubscriberReactive[T] Interface
```go
type SubscriberReactive[T any] interface {
    OnSubscribe(s Subscription)
    OnNext(t T)
    OnError(err error)
    OnComplete()
}
```

#### Subscription Interface
```go
type Subscription interface {
    Request(n int64)
    Cancel()
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
| `rxgo.Just[T](values...)` | Create from literal values | `rxgo.Just(1, 2, 3)` |
| `rxgo.Range(start, count)` | Integer sequence observable | `rxgo.Range(1, 10)` |
| `rxgo.Create[T](fn)` | Custom observable creation | `rxgo.Create(customProducer)` |
| `rxgo.RangePublisher(start, count)` | Integer sequence publisher | `rxgo.RangePublisher(1, 10)` |
| `rxgo.FromSlicePublisher[T](slice)` | Create publisher from slice | `rxgo.FromSlicePublisher([]int{1, 2, 3})` |
| `rxgo.NewPublisher[T](fn)` | Custom publisher creation | `rxgo.NewPublisher(customProducer)` |

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