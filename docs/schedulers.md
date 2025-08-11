# RxGo Schedulers Guide

Schedulers in RxGo control **where** and **when** your code executes. They provide fine-grained control over concurrency, thread management, and execution context, making them essential for building responsive and efficient reactive applications.

## Overview

RxGo provides four built-in schedulers, each optimized for specific use cases:

- **Immediate**: Executes work synchronously on the calling goroutine
- **NewThread**: Creates a new goroutine for each unit of work
- **SingleThread**: Uses a single dedicated goroutine for sequential execution
- **Trampoline**: Queues work for batch execution on the current goroutine

## Scheduler Types

### Immediate Scheduler

The `ImmediateScheduler` executes work synchronously on the calling goroutine. This is the most lightweight scheduler with zero overhead.

**Use Cases:**
- Lightweight operations (logging, simple calculations)
- Real-time UI updates where latency is critical
- Testing and debugging scenarios
- When work must complete before continuing

**Characteristics:**
- Zero context switching overhead
- Blocks the calling goroutine
- Executes immediately when scheduled
- No concurrency benefits

```go
scheduler := observable.NewImmediateScheduler()
scheduler.Schedule(func() {
    fmt.Println("Executed immediately")
})
```

### NewThread Scheduler

The `NewThreadScheduler` creates a new goroutine for each unit of work, providing maximum parallelism.

**Use Cases:**
- CPU-intensive operations (image processing, calculations)
- I/O-bound operations (database queries, HTTP requests)
- Parallel processing of independent tasks
- Background processing

**Characteristics:**
- Maximum CPU utilization
- Unbounded goroutine creation
- Parallel execution
- Higher memory overhead

```go
scheduler := observable.NewNewThreadScheduler()
scheduler.Schedule(func() {
    fmt.Println("Executed on new goroutine")
})
```

### SingleThread Scheduler

The `SingleThreadScheduler` uses a single dedicated goroutine to execute all scheduled work sequentially.

**Use Cases:**
- Database transactions requiring consistency
- Sequential processing pipelines
- State machines or ordered operations
- Thread-safe operations without locks

**Characteristics:**
- Sequential execution guarantees order
- Single goroutine for all work
- Thread-safe by design
- Predictable execution order

```go
scheduler := observable.NewSingleThreadScheduler()
defer scheduler.Close() // Important: always close when done

scheduler.Schedule(func() {
    fmt.Println("Task 1")
})
scheduler.Schedule(func() {
    fmt.Println("Task 2") // Always executes after Task 1
})
```

### Trampoline Scheduler

The `TrampolineScheduler` queues work for later execution on the current goroutine, allowing for batch processing and stack unwinding.

**Use Cases:**
- Batch processing of events
- Recursive operations
- Tail-call optimization scenarios
- When you need to defer work but stay on same goroutine

**Characteristics:**
- Queued execution
- Same goroutine as caller
- Batch processing capability
- Non-blocking scheduling

```go
scheduler := observable.NewTrampolineScheduler()
scheduler.Schedule(func() {
    fmt.Println("Queued task 1")
})
scheduler.Schedule(func() {
    fmt.Println("Queued task 2")
})
scheduler.Execute() // Executes all queued tasks
```

## Performance Comparison

| Scheduler | Latency | Throughput | Memory Usage | Use Case |
|-----------|---------|------------|--------------|----------|
| Immediate | Lowest | Single-threaded | Minimal | Lightweight ops |
| NewThread | Medium | Highest | High | CPU/I/O intensive |
| SingleThread | Low | Medium | Low | Ordered processing |
| Trampoline | Low | Medium | Minimal | Batch operations |

## Practical Examples

### CPU-Intensive Workload

```go
func processImages(imageIDs []int) {
    scheduler := observable.NewNewThreadScheduler()
    
    for _, id := range imageIDs {
        id := id // capture variable
        scheduler.Schedule(func() {
            processImage(id) // Runs in parallel
        })
    }
}
```

### Sequential Database Operations

```go
func sequentialDatabaseOperations() {
    scheduler := observable.NewSingleThreadScheduler()
    defer scheduler.Close()
    
    scheduler.Schedule(func() { insertUser() })
    scheduler.Schedule(func() { insertProfile() })
    scheduler.Schedule(func() { commitTransaction() })
    // All operations execute sequentially
}
```

### Real-time UI Updates

```go
func updateUI(data int) {
    scheduler := observable.NewImmediateScheduler()
    
    // Update happens immediately on main thread
    scheduler.Schedule(func() {
        updateProgressBar(data)
    })
}
```

### Batch Event Processing

```go
func processEvents(events []Event) {
    scheduler := observable.NewTrampolineScheduler()
    
    for _, event := range events {
        event := event // capture variable
        scheduler.Schedule(func() {
            processEvent(event)
        })
    }
    
    scheduler.Execute() // Process all events at once
}
```

## Best Practices

### 1. Choose Based on Workload Type

- **CPU-intensive**: Use `NewThread` for maximum parallelism
- **I/O-bound**: Use `NewThread` to hide latency
- **Ordered operations**: Use `SingleThread` for consistency
- **Lightweight work**: Use `Immediate` to avoid overhead
- **Batch processing**: Use `Trampoline` for efficiency

### 2. Resource Management

Always close `SingleThreadScheduler` when done:

```go
scheduler := observable.NewSingleThreadScheduler()
defer scheduler.Close()
```

### 3. Avoid Goroutine Leaks

Be cautious with `NewThreadScheduler` in high-frequency scenarios:

```go
// Good: Use for bounded workloads
for i := 0; i < 10; i++ {
    scheduler.Schedule(doWork)
}

// Bad: Can create unlimited goroutines
for range unboundedStream {
    scheduler.Schedule(doWork)
}
```

### 4. Testing Considerations

Use `ImmediateScheduler` for deterministic tests:

```go
func TestMyFunction(t *testing.T) {
    scheduler := observable.NewImmediateScheduler()
    result := myFunction(scheduler)
    // Result is immediately available
    assert.Equal(t, expected, result)
}
```

## Advanced Patterns

### Custom Pool Implementation

When the built-in `NewThreadScheduler` creates too many goroutines, implement a custom pool:

```go
type PoolScheduler struct {
    work chan func()
    wg   sync.WaitGroup
}

func NewPoolScheduler(size int) *PoolScheduler {
    p := &PoolScheduler{
        work: make(chan func(), size*2),
    }
    
    for i := 0; i < size; i++ {
        p.wg.Add(1)
        go func() {
            defer p.wg.Done()
            for task := range p.work {
                task()
            }
        }()
    }
    
    return p
}

func (p *PoolScheduler) Schedule(task func()) {
    p.work <- task
}

func (p *PoolScheduler) Close() {
    close(p.work)
    p.wg.Wait()
}
```

### Work Stealing Pattern

For dynamic load balancing:

```go
type WorkStealingScheduler struct {
    workers int
    tasks   chan func()
    results chan Result
}

func (w *WorkStealingScheduler) Schedule(task func()) {
    w.tasks <- task
}
```

### Scheduler Selection Guide

Create a decision tree for your application:

```go
func selectScheduler(workloadType string) observable.Scheduler {
    switch workloadType {
    case "CPU-intensive":
        return observable.NewNewThreadScheduler()
    case "IO-bound":
        return observable.NewNewThreadScheduler()
    case "ordered":
        return observable.NewSingleThreadScheduler()
    case "lightweight":
        return observable.NewImmediateScheduler()
    case "batch":
        return observable.NewTrampolineScheduler()
    default:
        return observable.NewImmediateScheduler()
    }
}
```

## Monitoring and Debugging

### Goroutine Count Monitoring

```go
func monitorGoroutines(scheduler observable.Scheduler) {
    initial := runtime.NumGoroutine()
    
    for i := 0; i < 100; i++ {
        scheduler.Schedule(func() {
            time.Sleep(time.Second)
        })
    }
    
    final := runtime.NumGoroutine()
    fmt.Printf("Goroutines created: %d\n", final-initial)
}
```

### Performance Profiling

```go
func profileScheduler(scheduler observable.Scheduler, tasks int) time.Duration {
    start := time.Now()
    var wg sync.WaitGroup
    
    for i := 0; i < tasks; i++ {
        wg.Add(1)
        scheduler.Schedule(func() {
            defer wg.Done()
            // Your workload here
        })
    }
    
    wg.Wait()
    return time.Since(start)
}
```

## Common Pitfalls

1. **Resource Leaks**: Not closing `SingleThreadScheduler`
2. **Goroutine Explosion**: Using `NewThreadScheduler` in unbounded loops
3. **Blocking Operations**: Using `ImmediateScheduler` for long-running tasks
4. **Order Dependencies**: Using parallel schedulers for ordered operations
5. **Memory Pressure**: Creating too many goroutines with `NewThreadScheduler`

## Examples and Usage

For comprehensive examples demonstrating all schedulers in action, see the `examples/schedulers/` directory which contains:

- **`scheduler_comparison.go`** - Basic scheduler comparison with timing measurements
- **`goroutine_patterns.go`** - Advanced goroutine usage patterns and benchmarks  
- **`computation_examples.go`** - Workload-specific examples for CPU, I/O, and memory-intensive tasks

Run the examples to see the differences:

```bash
# Comprehensive scheduler examples
go run examples/schedulers/scheduler_examples.go
```

## Examples Overview

The `examples/schedulers/scheduler_examples.go` file provides a concise, comprehensive demonstration:

- **Quick Demo** - Shows all 4 schedulers in action (50 lines of code)
- **Performance Benchmark** - Simple timing comparison for 100 tasks
- **Decision Guide** - Straightforward when-to-use recommendations

## Summary

Schedulers are a powerful feature in RxGo that provide fine-grained control over execution context. Choose the right scheduler based on your specific use case:

- **Immediate**: When you need immediate, synchronous execution
- **NewThread**: When you need maximum parallelism for CPU or I/O work
- **SingleThread**: When order matters and you need sequential processing
- **Trampoline**: When you want to batch work on the current goroutine

Understanding these schedulers and their appropriate use cases is crucial for building efficient, responsive reactive applications in Go.