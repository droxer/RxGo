# Benchmarks

This directory contains comprehensive benchmark tests for the RxGo library. All benchmarks use only public APIs and do not depend on internal packages.

## Structure

- `observable/benchmark_test.go` - Benchmarks for the observable package
- `rxgo/benchmark_test.go` - Benchmarks for the rxgo package

## Running Benchmarks

### Run all benchmarks:
```bash
go test -bench=. ./benchmarks/...
```

### Run specific package benchmarks:
```bash
# Observable package benchmarks
go test -bench=. ./benchmarks/observable

# RxGo package benchmarks
go test -bench=. ./benchmarks/rxgo
```

### Run with memory allocation statistics:
```bash
go test -bench=. -benchmem ./benchmarks/...
```

### Run specific benchmark:
```bash
go test -bench=BenchmarkObservableCreation ./benchmarks/observable
```

## Benchmark Categories

### Observable Package Benchmarks
- **Creation**: `BenchmarkObservableCreation`
- **With Subscriber**: `BenchmarkObservableWithSubscriber`
- **Large Datasets**: `BenchmarkObservableLargeDataset`
- **Range Operations**: `BenchmarkRangeObservable`
- **Just Operations**: `BenchmarkJustObservable`
- **Create Operations**: `BenchmarkCreateObservable`
- **Concurrent Subscribers**: `BenchmarkObservableConcurrentSubscribers`
- **Memory Allocations**: `BenchmarkObservableMemoryAllocations`
- **Data Types**: `BenchmarkObservableDataTypes`
- **Dataset Sizes**: `BenchmarkObservableDatasetSizes`
- **Context Cancellation**: `BenchmarkObservableContextCancellation`
- **Error Handling**: `BenchmarkObservableErrorHandling`
- **Struct Handling**: `BenchmarkObservableStructData`
- **String Handling**: `BenchmarkObservableStringData`

### RxGo Package Benchmarks
- **Creation**: `BenchmarkRxGoObservableCreation`
- **With Subscriber**: `BenchmarkRxGoObservableWithSubscriber`
- **Large Datasets**: `BenchmarkRxGoObservableLargeDataset`
- **Range Operations**: `BenchmarkRxGoRangeObservable`
- **Just Operations**: `BenchmarkRxGoJustObservable`
- **Create Operations**: `BenchmarkRxGoCreateObservable`
- **Concurrent Subscribers**: `BenchmarkRxGoObservableConcurrentSubscribers`
- **Memory Allocations**: `BenchmarkRxGoObservableMemoryAllocations`
- **Data Types**: `BenchmarkRxGoObservableDataTypes`
- **Dataset Sizes**: `BenchmarkRxGoObservableDatasetSizes`
- **Context Cancellation**: `BenchmarkRxGoObservableContextCancellation`
- **Error Handling**: `BenchmarkRxGoObservableErrorHandling`
- **Struct Handling**: `BenchmarkRxGoObservableStructData`
- **String Handling**: `BenchmarkRxGoObservableStringData`
- **Backpressure Strategies**: `BenchmarkRxGoBackpressureStrategies`

## Example Usage

```bash
# Run all benchmarks with verbose output
go test -v -bench=. ./benchmarks/...

# Run benchmarks with CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./benchmarks/observable
go test -bench=. -cpuprofile=cpu.prof ./benchmarks/rxgo

# Run benchmarks with memory profiling
go test -bench=. -memprofile=mem.prof ./benchmarks/observable
go test -bench=. -memprofile=mem.prof ./benchmarks/rxgo
```

## Performance Analysis

The benchmarks cover:
- **Creation overhead** for different observable types
- **Memory allocation patterns** for various operations
- **Throughput** with different dataset sizes
- **Concurrency** performance with multiple subscribers
- **Data type handling** performance (int, string, struct)
- **Error handling** performance impact
- **Context cancellation** performance
- **Backpressure strategy** performance impact

## Notes

- All benchmarks use only public APIs from `pkg/observable` and `pkg/rxgo`
- No internal package dependencies are used
- Benchmarks are designed to be run independently
- Results may vary based on system specifications and load