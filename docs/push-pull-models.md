# Push vs Pull Models & Backpressure

Understanding the fundamental differences between push and pull models in RxGo and how backpressure works.

## Push Model (Observable API)
The Observable API in `pkg/observable` implements a push-based model:

```go
import "github.com/droxer/RxGo/pkg/observable"

// Producer pushes data as fast as it can generate it
obs := observable.Just(1, 2, 3, 4, 5)
if err := obs.Subscribe(context.Background(), observable.NewSubscriber(
    func(v int) { fmt.Printf("Received: %d\n", v) },
    func() { fmt.Println("Completed") },
    func(err error) { fmt.Printf("Error: %v\n", err) },
)); err != nil {
    fmt.Printf("Subscription failed: %v\n", err)
}
```

### Characteristics of Push Model:
- **Producer Controlled**: The producer determines when data is emitted
- **Immediate Delivery**: Data is pushed to subscribers immediately
- **No Backpressure**: No mechanism for subscribers to control data flow
- **Simple API**: Easy to understand and use for basic scenarios
- **Potential Issues**: Can overwhelm slow consumers, leading to memory issues

## Pull Model (Reactive Streams)

The Reactive Streams API in `pkg/streams` implements a pull-based model with backpressure:

```go
import "github.com/droxer/RxGo/pkg/streams"

// Subscriber requests data in controlled amounts
publisher := streams.NewCompliantRangePublisher(1, 1000)
publisher.Subscribe(context.Background(), &MyReactiveSubscriber{})
```

With a custom subscriber that implements proper backpressure:

```go
type MyReactiveSubscriber struct {
    subscription streams.Subscription
}

func (s *MyReactiveSubscriber) OnSubscribe(sub streams.Subscription) {
    s.subscription = sub
    // Request only 10 items initially
    sub.Request(10)
}

func (s *MyReactiveSubscriber) OnNext(value int) {
    fmt.Printf("Received: %d
", value)
    // Request one more item after processing
    s.subscription.Request(1)
}

func (s *MyReactiveSubscriber) OnError(err error) {
    fmt.Printf("Error: %v
", err)
}

func (s *MyReactiveSubscriber) OnComplete() {
    fmt.Println("Completed")
}
```

### Characteristics of Pull Model:
- **Subscriber Controlled**: Subscribers request data using `Request(n)`
- **Backpressure Support**: Publishers must respect subscriber demand
- **Reactive Streams Compliant**: Full specification compliance
- **Production Ready**: Suitable for systems with unbounded data streams
- **Resource Control**: Prevents memory exhaustion and thread starvation

## Backpressure Strategies

The pull model supports multiple backpressure strategies:

### Buffer Strategy
Keeps all items in a bounded buffer:
```go
config := streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 100,
}
publisher := streams.NewBufferedPublisher[int](config, func(ctx context.Context, sub streams.Subscriber[int]) {
    // Implementation here
})
```

### Drop Strategy
Discards new items when buffer is full:
```go
config := streams.BackpressureConfig{
    Strategy:   streams.Drop,
    BufferSize: 50,
}
publisher := streams.NewBufferedPublisher[int](config, func(ctx context.Context, sub streams.Subscriber[int]) {
    // Implementation here
})
```

### Latest Strategy
Keeps only the latest item:
```go
config := streams.BackpressureConfig{
    Strategy:   streams.Latest,
    BufferSize: 1,
}
publisher := streams.NewBufferedPublisher[int](config, func(ctx context.Context, sub streams.Subscriber[int]) {
    // Implementation here
})
```

### Error Strategy
Signals an error when buffer overflows:
```go
config := streams.BackpressureConfig{
    Strategy:   streams.Error,
    BufferSize: 10,
}
publisher := streams.NewBufferedPublisher[int](config, func(ctx context.Context, sub streams.Subscriber[int]) {
    // Implementation here
})
```

## When to Use Each Model

### Use Push Model When:
- Building simple applications with predictable data rates
- Prototyping or learning reactive programming concepts
- Working with data sources that naturally emit at controlled rates
- Consumer can handle data at the producer's rate

### Use Pull Model When:
- Building production systems with potentially unbounded data streams
- Need to handle producer/consumer speed mismatches
- Working with external data sources (network, file I/O, database)
- Requiring Reactive Streams 1.0.4 compliance
- Need for resource control to prevent memory exhaustion

## Key Differences Summary

| Aspect | Push Model (Observable) | Pull Model (Reactive Streams) |
|--------|-------------------------|-------------------------------|
| Control Flow | Producer controlled | Subscriber controlled |
| Backpressure | Not supported | Fully supported |
| Complexity | Simple API | More complex but powerful |
| Use Cases | Simple scenarios | Production systems |
| Compliance | None | Reactive Streams 1.0.4 |
| Resource Safety | Potential issues | Built-in protection |

## Best Practices

1. **Start Simple**: Use the Observable API for learning and simple use cases
2. **Move to Streams**: Switch to Reactive Streams API for production systems
3. **Implement Proper Backpressure**: Always request data in controlled amounts
4. **Handle Errors**: Properly handle backpressure-related errors
5. **Monitor Performance**: Watch for memory and CPU usage in high-volume scenarios