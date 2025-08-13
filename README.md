# RxGo

Reactive Extensions for Go with full Reactive Streams 1.0.4 compliance.

[![Go Report Card](https://goreportcard.com/badge/github.com/droxer/RxGo)](https://goreportcard.com/report/github.com/droxer/RxGo)
[![GoDoc](https://godoc.org/github.com/droxer/RxGo?status.svg)](https://godoc.org/github.com/droxer/RxGo)

## Features

- **Type-safe generics** - Full Go generics support
- **Reactive Streams 1.0.4** - Complete specification compliance
- **Backpressure strategies** - Buffer, Drop, Latest, Error
- **Context cancellation** - Graceful shutdown
- **Thread-safe** - Safe concurrent access

## Quick Start

```bash
go get github.com/droxer/RxGo@latest
```

```go
import "github.com/droxer/RxGo/pkg/rx"

// Simple usage
rx.Just(1, 2, 3).Subscribe(context.Background(), 
    rx.NewSubscriber(
        func(v int) { fmt.Printf("Got %d\n", v) },
        func() { fmt.Println("Done") },
        func(err error) { fmt.Printf("Error: %v\n", err) },
    ))
```

## Quick Start

```bash
go get github.com/droxer/RxGo@latest
```

```go
import "github.com/droxer/RxGo/pkg/rx"

// Simple usage
rx.Just(1, 2, 3).Subscribe(context.Background(), 
    rx.NewSubscriber(
        func(v int) { fmt.Printf("Got %d\n", v) },
        func() { fmt.Println("Done") },
        func(err error) { fmt.Printf("Error: %v\n", err) },
    ))
```

## Documentation

- [Quick Start](./docs/quick-start.md)
- [Reactive Streams](./docs/reactive-streams.md)
- [Backpressure](./docs/backpressure.md)
- [API Reference](./docs/api-reference.md)

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

[MIT License](./LICENSE)
