# RxGo

Reactive Extensions for Go with full Reactive Streams 1.0.4 compliance.

[![Go Report Card](https://goreportcard.com/badge/github.com/droxer/RxGo)](https://goreportcard.com/report/github.com/droxer/RxGo)
[![GoDoc](https://godoc.org/github.com/droxer/RxGo?status.svg)](https://godoc.org/github.com/droxer/RxGo)

## Features

- **Type-safe generics** - Full Go generics support
- **Reactive Streams 1.0.4** - Complete specification compliance
- **Backpressure strategies** - Buffer, Drop, Latest, Error
- **Retry and backoff** - Fixed, Linear, Exponential backoff with configurable retry limits
- **Context cancellation** - Graceful shutdown
- **Thread-safe** - Safe concurrent access

## Quick Start

Install the library:

```bash
go get github.com/droxer/RxGo@latest
```

### Usage Example

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/rx"
)

func main() {
    obs := rx.Just(1, 2, 3, 4, 5)
    obs.Subscribe(context.Background(), rx.NewSubscriber(
        func(v int) { fmt.Printf("Got %d\n", v) },
        func() { fmt.Println("Done") },
        func(err error) { fmt.Printf("Error: %v\n", err) },
    ))
}
```

**Output:**
```
Got 1
Got 2
Got 3
Got 4
Got 5
Done
```

### Backpressure Strategies

Handle producer/consumer speed mismatches with four strategies:

```go
import "github.com/droxer/RxGo/pkg/rx/streams"

// Buffer - keep all items in bounded buffer
publisher := streams.RangePublishWithBackpressure(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Buffer,
    BufferSize: 100,
})

// Drop - discard new items when full
publisher := streams.RangePublishWithBackpressure(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Drop,
    BufferSize: 50,
})

// Latest - keep only latest item
publisher := streams.RangePublishWithBackpressure(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Latest,
    BufferSize: 1,
})

// Error - signal error on overflow
publisher := streams.RangePublishWithBackpressure(1, 1000, streams.BackpressureConfig{
    Strategy:   streams.Error,
    BufferSize: 10,
})
```

## Documentation

- [Basic Usage](./docs/basic-usage.md) - Simple Observable API examples
- [Reactive Streams](./docs/reactive-streams.md) - Full Reactive Streams 1.0.4 compliance
- [Backpressure](./docs/backpressure.md) - Handle producer/consumer speed mismatches
- [Retry and Backoff](./docs/retry-backoff.md) - Configurable retry with backoff strategies
- [Data Transformation](./docs/data-transformation.md) - Transform and process data streams
- [Context Cancellation](./docs/context-cancellation.md) - Graceful cancellation using Go context
- [Schedulers](./docs/schedulers.md) - Execution context control
- [Architecture](./docs/architecture.md) - Package structure and design decisions
- [API Reference](./docs/api-reference.md) - Complete API documentation

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

[MIT License](./LICENSE)
