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
