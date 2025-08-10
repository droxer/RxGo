# RxGo Usage Examples

This directory contains comprehensive usage examples for RxGo, organized by topic.

## 📁 Example Categories

### [Basic Usage](./basic-usage.md)
- Observable API (Simple)
- Range Observable
- Create with Custom Logic

### [Reactive Streams API](./reactive-streams.md)
- Reactive Streams Publisher
- Custom Publisher Creation
- FromSlice Publisher

### [Backpressure Control](./backpressure.md)
- Slow Subscriber with Backpressure
- Buffered Backpressure Strategy
- Demand-Based Flow Control

### [Context Cancellation](./context-cancellation.md)
- Timeout-Based Cancellation
- Manual Cancellation
- Context with Deadline
- Parent Context Cancellation

### [Data Transformation](./data-transformation.md)
- Basic Data Transformation
- Filter and Map Operations
- String Processing Pipeline
- Aggregation Operations

## 🚀 Quick Start

Choose your preferred API and start with the corresponding example:

- **Beginners**: Start with [Basic Usage](./basic-usage.md)
- **Advanced Users**: Use [Reactive Streams API](./reactive-streams.md)
- **Production**: Review [Backpressure Control](./backpressure.md) and [Context Cancellation](./context-cancellation.md)

## 📋 Running Examples

Each example can be run independently:

```bash
# Basic examples
go run examples/basic/basic.go

# Advanced examples
go run examples/backpressure/backpressure.go
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