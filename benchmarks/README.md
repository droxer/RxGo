# Benchmarks

This directory contains comprehensive benchmark tests for the RxGo library. 

## Structure

- `observable/benchmark_test.go` - Benchmarks for the observable package
- `rxgo/benchmark_test.go` - Benchmarks for the rxgo package
- `operators/benchmark_test.go` - Benchmarks for operators package
- `backpressure/benchmark_test.go` - Benchmarks for backpressure strategies

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
- **From Slice**: `BenchmarkObservableFromSlice`
- **Empty Observable**: `BenchmarkEmptyObservable`
- **Never Observable**: `BenchmarkNeverObservable`
- **Error Observable**: `BenchmarkErrorObservable`

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
- **From Slice**: `BenchmarkRxGoObservableFromSlice`
- **Empty Observable**: `BenchmarkRxGoEmptyObservable`
- **Never Observable**: `BenchmarkRxGoNeverObservable`
- **Error Observable**: `BenchmarkRxGoErrorObservable`

### Operators Package Benchmarks
- **Map Operator**: `BenchmarkOperatorMap`
- **Filter Operator**: `BenchmarkOperatorFilter`
- **Operator Chains**: `BenchmarkOperatorMapFilterChain`
- **ObserveOn Operator**: `BenchmarkOperatorObserveOn`
- **Memory Allocations**: `BenchmarkOperatorMapMemoryAllocations`, `BenchmarkOperatorFilterMemoryAllocations`, `BenchmarkOperatorChainMemoryAllocations`
- **Data Types**: `BenchmarkOperatorMapDifferentDataTypes`, `BenchmarkOperatorFilterDifferentDataTypes`

### Backpressure Package Benchmarks
- **Buffer Strategy**: `BenchmarkBackpressureBufferStrategy`
- **Drop Strategy**: `BenchmarkBackpressureDropStrategy`
- **Latest Strategy**: `BenchmarkBackpressureLatestStrategy`
- **Error Strategy**: `BenchmarkBackpressureErrorStrategy`
- **Strategy Comparison**: `BenchmarkBackpressureStrategiesComparison`
- **Buffer Sizes**: `BenchmarkBackpressureWithDifferentBufferSizes`
- **Memory Allocations**: `BenchmarkBackpressureMemoryAllocations`

## Example Usage

```bash
# Run all benchmarks with verbose output
go test -v -bench=. ./benchmarks/...

# Run benchmarks with CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./benchmarks/observable
go test -bench=. -cpuprofile=cpu.prof ./benchmarks/rxgo
go test -bench=. -cpuprofile=cpu.prof ./benchmarks/operators
go test -bench=. -cpuprofile=cpu.prof ./benchmarks/backpressure

# Run benchmarks with memory profiling
go test -bench=. -memprofile=mem.prof ./benchmarks/observable
go test -bench=. -memprofile=mem.prof ./benchmarks/rxgo
go test -bench=. -memprofile=mem.prof ./benchmarks/operators
go test -bench=. -memprofile=mem.prof ./benchmarks/backpressure
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
- **Operator chaining** performance
- **Scheduler integration** performance
