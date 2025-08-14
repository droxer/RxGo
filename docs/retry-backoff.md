# Retry and Backoff Strategies

Handle failures gracefully with configurable retry and backoff strategies that maintain Reactive Streams 1.0.4 compliance.

## Overview

RxGo provides flexible retry mechanisms with three backoff strategies to handle transient failures in reactive applications. All retry implementations maintain full Reactive Streams 1.0.4 compliance including backpressure support and proper cancellation handling.

## Retry Features

| Feature | Description |
|-----------|-------------|
| **Max Retries** | Configurable retry limit (0 = infinite) |
| **Backoff Strategies** | Fixed, Linear, Exponential |
| **Retry Conditions** | Custom retry logic per error type |
| **Context Support** | Respects context cancellation |
| **Backpressure** | Maintains demand control |

## Quick Start

```go
import "github.com/droxer/RxGo/pkg/rx/streams"

// Basic retry with exponential backoff
publisher := streams.NewRetryPublisher(source, streams.RetryConfig{
    MaxRetries:    3,
    InitialDelay:  100 * time.Millisecond,
    BackoffFactor: 2.0,
    BackoffPolicy: streams.RetryExponential,
})
```

## Retry Strategies

### Fixed Backoff
```go
publisher := streams.NewRetryPublisher(source, streams.RetryConfig{
    MaxRetries:    5,
    InitialDelay:  200 * time.Millisecond,
    BackoffPolicy: streams.RetryFixed,
})
// Delays: 200ms, 200ms, 200ms, 200ms, 200ms
```

### Linear Backoff
```go
publisher := streams.NewRetryPublisher(source, streams.RetryConfig{
    MaxRetries:    3,
    InitialDelay:  100 * time.Millisecond,
    BackoffPolicy: streams.RetryLinear,
})
// Delays: 100ms, 200ms, 300ms
```

### Exponential Backoff
```go
publisher := streams.NewRetryPublisher(source, streams.RetryConfig{
    MaxRetries:    4,
    InitialDelay:  50 * time.Millisecond,
    BackoffFactor: 2.0,
    BackoffPolicy: streams.RetryExponential,
})
// Delays: 50ms, 100ms, 200ms, 400ms
```

## Custom Retry Conditions

```go
publisher := streams.NewRetryPublisher(source, streams.RetryConfig{
    MaxRetries:   3,
    InitialDelay: 100 * time.Millisecond,
    RetryCondition: func(err error, attempt int) bool {
        // Retry only for network errors
        if netErr, ok := err.(net.Error); ok {
            return netErr.Temporary()
        }
        return false
    },
})
```

## Fluent Builder API

```go
publisher := streams.NewRetryBuilder[int](source).
    MaxRetries(5).
    InitialDelay(100 * time.Millisecond).
    MaxDelay(5 * time.Second).
    ExponentialBackoff(2.0).
    RetryCondition(func(err error, attempt int) bool {
        return attempt < 3
    }).
    Build()
```

## Convenience Functions

```go
// Fixed delay retry
publisher := streams.WithFixedRetry(source, 3, 100*time.Millisecond)

// Linear backoff retry
publisher := streams.WithLinearRetry(source, 3, 100*time.Millisecond)

// Exponential backoff retry
publisher := streams.WithExponentialRetry(source, 3, 100*time.Millisecond, 2.0)

// Infinite retry
publisher := streams.WithInfiniteRetry(source, 100*time.Millisecond, streams.RetryFixed)
```

## Complete Examples

### Network Retry with Exponential Backoff

```go
// Simulate a network call that might fail
networkCall := streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[string]) {
    if rand.Float32() < 0.7 {
        sub.OnError(errors.New("network timeout"))
        return
    }
    sub.OnNext("API response data")
    sub.OnComplete()
})

// Configure retry with exponential backoff
retryPublisher := streams.NewRetryPublisher(networkCall, streams.RetryConfig{
    MaxRetries:    5,
    InitialDelay:  100 * time.Millisecond,
    MaxDelay:      5 * time.Second,
    BackoffFactor: 2.0,
    BackoffPolicy: streams.RetryExponential,
    RetryCondition: func(err error, attempt int) bool {
        return strings.Contains(err.Error(), "timeout")
    },
})

// Subscribe and process
retryPublisher.Subscribe(ctx, &MySubscriber{})
```

### Database Connection Retry

```go
dbQuery := streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[[]Row]) {
    if db.Ping() != nil {
        sub.OnError(errors.New("connection failed"))
        return
    }
    sub.OnNext(queryResults)
    sub.OnComplete()
})

// Retry with linear backoff for connection issues
retryPublisher := streams.WithLinearRetry(dbQuery, 10, 500*time.Millisecond)
retryPublisher.Subscribe(ctx, &MySubscriber{})
```

## API Reference

### Types

```go
type RetryConfig struct {
    MaxRetries     int
    InitialDelay   time.Duration
    MaxDelay       time.Duration
    BackoffFactor  float64
    BackoffPolicy  RetryBackoffPolicy
    RetryCondition RetryCondition
}

type RetryCondition func(err error, attempt int) bool
type RetryBackoffPolicy int

const (
    RetryFixed RetryBackoffPolicy = iota
    RetryLinear
    RetryExponential
)
```

### Builders

- `NewRetryPublisher[T any](source Publisher[T], config RetryConfig) Publisher[T]`
- `NewRetryBuilder[T any](source Publisher[T]) *RetryBuilder[T]`
- `WithFixedRetry[T any](source Publisher[T], maxRetries int, delay time.Duration) Publisher[T]`
- `WithLinearRetry[T any](source Publisher[T], maxRetries int, initialDelay time.Duration) Publisher[T]`
- `WithExponentialRetry[T any](source Publisher[T], maxRetries int, initialDelay time.Duration, factor float64) Publisher[T]`
- `WithInfiniteRetry[T any](source Publisher[T], initialDelay time.Duration, policy BackoffPolicy) Publisher[T]`

### Builder Methods

- `MaxRetries(max int) *RetryBuilder[T]`
- `InitialDelay(delay time.Duration) *RetryBuilder[T]`
- `MaxDelay(max time.Duration) *RetryBuilder[T]`
- `FixedBackoff() *RetryBuilder[T]`
- `LinearBackoff() *RetryBuilder[T]`
- `ExponentialBackoff(factor float64) *RetryBuilder[T]`
- `RetryCondition(condition RetryCondition) *RetryBuilder[T]`
- `Build() Publisher[T]`

## Best Practices

1. **Set reasonable retry limits** to avoid infinite loops
2. **Use exponential backoff** for network calls
3. **Implement retry conditions** to avoid retrying non-recoverable errors
4. **Respect context cancellation** for graceful shutdown
5. **Set maximum delays** to prevent excessive backoff
6. **Log retry attempts** for debugging and monitoring

## Integration with Backpressure

Retry strategies work seamlessly with backpressure control:

```go
// Combine retry with backpressure
publisher := streams.RangePublishWithBackpressure(1, 100, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 50,
})

retryPublisher := streams.NewRetryPublisher(publisher, streams.RetryConfig{
    MaxRetries:    3,
    InitialDelay:  100 * time.Millisecond,
    BackoffPolicy: streams.RetryExponential,
})
```

## Testing Retry Logic

```go
// Test retry behavior
testPublisher := streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[int]) {
    if testAttempts < 2 {
        sub.OnError(errors.New("simulated error"))
        testAttempts++
        return
    }
    sub.OnNext(42)
    sub.OnComplete()
})

retryPublisher := streams.WithFixedRetry(testPublisher, 3, 10*time.Millisecond)
```

## Error Handling

The retry system properly handles:
- **Context cancellation** during backoff
- **Subscription cancellation** preventing further retries
- **Backpressure signals** maintaining demand control
- **Terminal states** ensuring no further events after completion/error