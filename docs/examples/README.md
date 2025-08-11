# RxGo Usage Examples

This directory contains comprehensive usage examples for RxGo, organized by topic and consistent with the actual examples in the codebase.

## 📁 Example Categories

### [Quick Start](./quick-start.md)
- Minimal working example
- Just and Range observables
- Basic subscriber implementation

### [Basic Usage](./basic-usage.md)
- Observable API with rxgo package
- Range and Just observables
- Create with custom logic
- Context cancellation integration

### [Backpressure Control](./backpressure.md)
- Slow consumer processing
- Thread-safe subscriber implementation
- WaitGroup for completion handling
- Real-world backpressure patterns

### [Context Cancellation](./context-cancellation.md)
- Timeout-based cancellation
- Context-aware observables
- Graceful shutdown patterns
- Production-ready cancellation

### [Reactive Streams API](./reactive-streams.md)
- Observable as publisher pattern
- Range and Just publishers
- Context support in streams
- Synchronous processing examples

### [Data Transformation](./data-transformation.md)
- Range-based transformations
- Just value processing
- Custom create logic
- Context-aware transformations

## 🚀 Quick Start

Choose your preferred API and start with the corresponding example:

- **Beginners**: Start with [Quick Start](./quick-start.md)
- **Basic Usage**: Use [Basic Usage](./basic-usage.md)
- **Production**: Review [Backpressure Control](./backpressure.md) and [Context Cancellation](./context-cancellation.md)

## 📋 Running Examples

Each example can be run independently and matches the actual code:

```bash
# Basic examples
go run examples/basic/basic.go

# Backpressure example
go run examples/backpressure/backpressure.go

# Context cancellation
go run examples/context/context.go

# Or use compiled binaries
./bin/basic
./bin/backpressure  
./bin/context
```

## 🔗 API Reference

For complete API documentation, see:
- [Observable API](../pkg/observable/doc.go)
- [Reactive Streams API](../pkg/rxgo/doc.go)
- [GoDoc](https://godoc.org/github.com/droxer/RxGo)