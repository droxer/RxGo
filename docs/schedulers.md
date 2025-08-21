# RxGo Schedulers Guide

Schedulers in RxGo control **where** and **when** your code executes. They provide fine-grained control over concurrency, thread management, and execution context, making them essential for building responsive and efficient reactive applications.

## Overview

RxGo provides five built-in schedulers, each optimized for specific use cases:

- **Computation**: Fixed thread pool for CPU-bound work
- **IO**: Cached thread pool for I/O-bound work
- **NewThread**: Creates a new goroutine for each unit of work
- **SingleThread**: Uses a single dedicated goroutine for sequential execution
- **Trampoline**: Executes work immediately on the current goroutine

## Scheduler Types

### Computation Scheduler

The `Computation` scheduler uses a fixed-size thread pool optimized for CPU-bound work. It automatically sizes the pool based on available CPU cores.

**Use Cases:**
- CPU-intensive operations (mathematical calculations, data processing)
- Parallel algorithms and computations
- Image processing and compression
- Cryptographic operations

**Characteristics:**
- Fixed thread pool size (number of CPU cores)
- Optimized for CPU-bound operations
- Prevents goroutine explosion
- Efficient resource utilization

```go
import "github.com/droxer/RxGo/pkg/scheduler"

scheduler := scheduler.Computation
scheduler.Schedule(func() {
    fmt.Println("Executed on computation thread pool")
})
```

### IO Scheduler

The `IO` scheduler uses a cached thread pool that grows and shrinks based on demand, optimized for I/O-bound work.

**Use Cases:**
- Network requests and HTTP calls
- Database operations and queries
- File system operations
- Time-consuming I/O operations

**Characteristics:**
- Cached thread pool with automatic sizing
- Optimized for blocking I/O operations
- Handles bursty workloads efficiently
- Unbounded thread creation (within limits)

```go
import "github.com/droxer/RxGo/pkg/scheduler"

scheduler := scheduler.IO
scheduler.Schedule(func() {
    fmt.Println("Executed on IO thread pool")
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
import "github.com/droxer/RxGo/pkg/scheduler"

scheduler := scheduler.NewThread
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
import "github.com/droxer/RxGo/pkg/scheduler"

scheduler := scheduler.SingleThread
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
import "github.com/droxer/RxGo/pkg/scheduler"

scheduler := scheduler.Trampoline
scheduler.Schedule(func() {
    fmt.Println("Executed immediately")
})
```

## Performance Comparison

| Scheduler | Latency | Throughput | Memory Usage | Use Case |
|-----------|---------|------------|--------------|----------|
| Computation | Low | High CPU | Medium | CPU-bound work |
| IO | Low | High I/O | Low-Medium | I/O operations |
| NewThread | Medium | Highest | High | Parallel processing |
| SingleThread | Low | Medium | Low | Ordered processing |
| Trampoline | Lowest | Single-threaded | Minimal | Immediate execution |

## Practical Examples

### CPU-Intensive Workload

```go
import "github.com/droxer/RxGo/pkg/scheduler"

func processImages(imageIDs []int) {
    sched := scheduler.Computation // Use Computation for CPU work
    
    for _, id := range imageIDs {
        id := id // capture variable
        sched.Schedule(func() {
            processImage(id) // Runs in parallel
        })
    }
}
```

### Sequential Database Operations

```go
import "github.com/droxer/RxGo/pkg/scheduler"

func sequentialDatabaseOperations() {
    sched := scheduler.SingleThread
    
    sched.Schedule(func() { insertUser() })
    sched.Schedule(func() { insertProfile() })
    sched.Schedule(func() { commitTransaction() })
    // All operations execute sequentially
}
```

### Real-time UI Updates

```go
import "github.com/droxer/RxGo/pkg/scheduler"

func updateUI(data int) {
    sched := scheduler.Trampoline
    
    // Update happens immediately on current goroutine
    sched.Schedule(func() {
        updateProgressBar(data)
    })
}
```

### Batch Event Processing

```go
import "github.com/droxer/RxGo/pkg/scheduler"

func processEvents(events []Event) {
    sched := scheduler.Trampoline
    
    for _, event := range events {
        event := event // capture variable
        sched.Schedule(func() {
            processEvent(event)
        })
    }
    
    // All events processed immediately on current goroutine
}
```

## Best Practices

### 1. Choose Based on Workload Type

- **CPU-intensive**: Use `Computation` for efficient CPU utilization
- **I/O-bound**: Use `IO` for blocking operations
- **Maximum parallelism**: Use `NewThread` for unlimited goroutines
- **Ordered operations**: Use `SingleThread` for consistency
- **Immediate execution**: Use `Trampoline` for current goroutine


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

Use `Trampoline` for deterministic tests:

```go
import "github.com/droxer/RxGo/pkg/scheduler"

func TestMyFunction(t *testing.T) {
    sched := scheduler.Trampoline
    result := myFunction(sched)
    // Result is immediately available
    assert.Equal(t, expected, result)
}
```

## Advanced Patterns

### Custom Pool Implementation

When the built-in `NewThread` creates too many goroutines, implement a custom pool:

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
func selectScheduler(workloadType string) Scheduler {
    switch workloadType {
    case "CPU-intensive":
        return scheduler.Computation
    case "IO-bound":
        return scheduler.IO
    case "ordered":
        return scheduler.SingleThread
    case "maximum-parallelism":
        return scheduler.NewThread
    case "immediate":
        return scheduler.Trampoline
    default:
        return scheduler.Trampoline
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

## Complete Examples

### Quick Demo - All Schedulers

```go
import "github.com/droxer/RxGo/pkg/scheduler"

// Demo all 5 schedulers
work := func(id int) { fmt.Printf("Task %d on goroutine %d\n", id, getGID()) }

// 1. Trampoline - immediate execution
scheduler.Trampoline.Schedule(func() { work(1) })

// 2. NewThread - new goroutine for each task
scheduler.NewThread.Schedule(func() { work(2) })

// 3. SingleThread - sequential on dedicated goroutine
st := scheduler.SingleThread
st.Schedule(func() { work(3) })
st.Schedule(func() { work(4) }) // Always executes after task 3

// 4. Computation - fixed thread pool
scheduler.Computation.Schedule(func() { work(5) })

// 5. IO - cached thread pool
scheduler.IO.Schedule(func() { work(6) })
```

### Performance Comparison

```go
import "time"

func benchmarkScheduler() {
    task := func() {
        // Simulate CPU work
        for i := 0; i < 1000; i++ {}
    }
    
    benchmark := func(name string, s scheduler.Scheduler) time.Duration {
        start := time.Now()
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            s.Schedule(func() { defer wg.Done(); task() })
        }
        wg.Wait()
        return time.Since(start)
    }
    
    fmt.Printf("Computation: %v\n", benchmark("Computation", scheduler.Computation))
    fmt.Printf("IO: %v\n", benchmark("IO", scheduler.IO))
}
```

### Decision Guide

| Scheduler | When to Use | Memory | Use Case |
|-----------|-------------|--------|----------|
| **Computation** | CPU-intensive work | Medium | Math, compression |
| **IO** | I/O-bound work | Low-Medium | HTTP, database |
| **NewThread** | Maximum parallelism | High | Parallel processing |
| **SingleThread** | Ordered operations | Low | Transactions |
| **Trampoline** | Immediate execution | Minimal | UI updates |

## Summary

Schedulers are a powerful feature in RxGo that provide fine-grained control over execution context. Choose the right scheduler based on your specific use case:

- **Computation**: When you need efficient CPU utilization for CPU-bound work
- **IO**: When you need efficient handling of I/O-bound operations
- **NewThread**: When you need maximum parallelism with new goroutines
- **SingleThread**: When order matters and you need sequential processing
- **Trampoline**: When you need immediate execution on the current goroutine

Understanding these schedulers and their appropriate use cases is crucial for building efficient, responsive reactive applications in Go.