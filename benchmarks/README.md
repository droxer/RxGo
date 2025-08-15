# Benchmarks

This directory contains performance benchmarks for the RxGo library, organized by core features:

## Structure

- **`observable/`** - Benchmarks for the Observable pattern (push-based)
- **`streams/`** - Benchmarks for Reactive Streams (pull-based with backpressure)
- **`operators/`** - Comparative benchmarks for operators across both patterns
- **`scheduler/`** - Benchmarks for scheduler implementations

## Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./benchmarks/...

# Run specific package benchmarks
go test -bench=. ./benchmarks/observable
go test -bench=. ./benchmarks/streams
go test -bench=. ./benchmarks/operators
go test -bench=. ./benchmarks/scheduler

# Run with memory allocations
go test -bench=. -benchmem ./benchmarks/...
```

## Benchmark Categories

### Observable Pattern (`observable/`)
- **Creation**: `Just`, `Range`, `Empty`, `Error`, `Create`
- **Subscription**: Direct subscription performance
- **Data Types**: `int`, `float64`, `string`, `struct`
- **Dataset Sizes**: 10, 100, 1000, 10000 items
- **Concurrency**: Parallel subscribers
- **Memory**: Allocation profiling
- **Error Handling**: Error propagation
- **Cancellation**: Context cancellation

### Reactive Streams (`streams/`)
- **Publisher Creation**: `RangePublisher`, `FromSlicePublisher`
- **Subscription**: Reactive streams subscription
- **Backpressure**: Different request patterns
- **Processors**: Map, Filter, and chain processors
- **Dataset Sizes**: Various input sizes
- **Concurrency**: Parallel subscribers
- **Memory**: Allocation profiling
- **Error Handling**: Publisher error scenarios

### Operators (`operators/`)
Comparative benchmarks between Observable and Reactive Streams for:
- **Map**: Transformation performance
- **Filter**: Filtering performance  
- **Chain Operations**: Multi-operator chains
- **Data Types**: Performance across different types
- **Complex Chains**: 3+ operator pipelines
- **Concurrent Usage**: Parallel execution
- **Memory Usage**: Allocation comparisons

### Schedulers (`scheduler/`)
- **Creation**: Scheduler instantiation
- **Execution**: Task scheduling performance
- **Concurrent**: Parallel scheduling
- **Throughput**: Batch processing
- **Latency**: Scheduling latency measurements
- **Memory**: Allocation profiling
- **Scheduler Types**: Trampoline, NewThread, SingleThread, Computation, IO

## Key Metrics

- **ns/op**: Nanoseconds per operation
- **B/op**: Bytes allocated per operation
- **allocs/op**: Number of allocations per operation

## Usage Examples

```bash
# Compare Observable vs Reactive Streams Map performance
go test -bench=BenchmarkMap -benchmem ./benchmarks/operators

# Test scheduler performance
go test -bench=BenchmarkSchedulerExecution ./benchmarks/scheduler

# Profile memory usage for large datasets
go test -bench=BenchmarkDatasetSizes -benchmem ./benchmarks/observable
```
