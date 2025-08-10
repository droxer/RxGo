# Data Transformation

This document demonstrates how to transform and process data streams using RxGo.

## 1. Basic Data Transformation

Transform data as it flows through the stream:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/observable"
)

func main() {
    // Create source observable
    source := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        defer sub.OnCompleted()
        for i := 1; i <= 5; i++ {
            sub.OnNext(i)
        }
    })
    
    // Transform the stream by doubling values
    transformed := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        source.Subscribe(ctx, observable.Subscriber[int]{
            Start: func() {
                sub.Start()
            },
            OnNext: func(value int) {
                sub.OnNext(value * 2) // Double the value
            },
            OnError: func(err error) {
                sub.OnError(err)
            },
            OnCompleted: func() {
                sub.OnCompleted()
            },
        })
    })
    
    // Subscribe to the transformed stream
    type IntSubscriber struct{}
    
    func (s *IntSubscriber) Start() {
        fmt.Println("Starting transformation subscriber")
    }
    
    func (s *IntSubscriber) OnNext(value int) {
        fmt.Printf("Received: %d\n", value)
    }
    
    func (s *IntSubscriber) OnError(err error) {
        fmt.Printf("Error: %v\n", err)
    }
    
    func (s *IntSubscriber) OnCompleted() {
        fmt.Println("Transformation completed")
    }
    
    subscriber := &IntSubscriber{}
    transformed.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```

## 2. Filter and Map Operations

Filter and map data in reactive streams:

```go
package main

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/droxer/RxGo/pkg/observable"
)

// FilterSubscriber implements filtering logic
type FilterSubscriber struct {
    predicate func(int) bool
    next      observable.Subscriber[int]
}

func (s *FilterSubscriber) Start() {
    s.next.Start()
}

func (s *FilterSubscriber) OnNext(value int) {
    if s.predicate(value) {
        s.next.OnNext(value)
    }
}

func (s *FilterSubscriber) OnError(err error) {
    s.next.OnError(err)
}

func (s *FilterSubscriber) OnCompleted() {
    s.next.OnCompleted()
}

// MapSubscriber implements mapping logic
type MapSubscriber struct {
    mapper func(int) string
    next   observable.Subscriber[string]
}

func (s *MapSubscriber) Start() {
    s.next.Start()
}

func (s *MapSubscriber) OnNext(value int) {
    s.next.OnNext(s.mapper(value))
}

func (s *MapSubscriber) OnError(err error) {
    s.next.OnError(err)
}

func (s *MapSubscriber) OnCompleted() {
    s.next.OnCompleted()
}

func main() {
    // Create source with numbers 1-10
    source := observable.Range(1, 10)
    
    // Filter even numbers
    filtered := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        source.Subscribe(ctx, &FilterSubscriber{
            predicate: func(n int) bool { return n%2 == 0 },
            next:      sub,
        })
    })
    
    // Map to strings
    mapped := observable.Create(func(ctx context.Context, sub observable.Subscriber[string]) {
        filtered.Subscribe(ctx, &MapSubscriber{
            mapper: func(n int) string { return fmt.Sprintf("even-%d", n) },
            next:   sub,
        })
    })
    
    // Subscribe to final result
    type StringSubscriber struct{}
    
    func (s *StringSubscriber) Start() {
        fmt.Println("Starting filter-map pipeline")
    }
    
    func (s *StringSubscriber) OnNext(value string) {
        fmt.Printf("Result: %s\n", value)
    }
    
    func (s *StringSubscriber) OnError(err error) {
        fmt.Printf("Pipeline error: %v\n", err)
    }
    
    func (s *StringSubscriber) OnCompleted() {
        fmt.Println("Filter-map pipeline completed")
    }
    
    subscriber := &StringSubscriber{}
    mapped.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```

## 3. String Processing Pipeline

Create complex string processing pipelines:

```go
package main

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/droxer/RxGo/pkg/observable"
)

func main() {
    // Create string source
    source := observable.Just(
        "hello world",
        "reactive programming",
        "go generics",
        "stream processing",
        "data transformation",
    )
    
    // Split into words
    words := observable.Create(func(ctx context.Context, sub observable.Subscriber[string]) {
        source.Subscribe(ctx, observable.Subscriber[string]{
            Start: func() {
                sub.Start()
            },
            OnNext: func(value string) {
                for _, word := range strings.Fields(value) {
                    sub.OnNext(strings.ToLower(word))
                }
            },
            OnError: func(err error) {
                sub.OnError(err)
            },
            OnCompleted: func() {
                sub.OnCompleted()
            },
        })
    })
    
    // Filter words longer than 3 characters
    filtered := observable.Create(func(ctx context.Context, sub observable.Subscriber[string]) {
        words.Subscribe(ctx, observable.Subscriber[string]{
            Start: func() {
                sub.Start()
            },
            OnNext: func(value string) {
                if len(value) > 3 {
                    sub.OnNext(value)
                }
            },
            OnError: func(err error) {
                sub.OnError(err)
            },
            OnCompleted: func() {
                sub.OnCompleted()
            },
        })
    })
    
    // Transform to uppercase
    transformed := observable.Create(func(ctx context.Context, sub observable.Subscriber[string]) {
        filtered.Subscribe(ctx, observable.Subscriber[string]{
            Start: func() {
                sub.Start()
            },
            OnNext: func(value string) {
                sub.OnNext(strings.ToUpper(value))
            },
            OnError: func(err error) {
                sub.OnError(err)
            },
            OnCompleted: func() {
                sub.OnCompleted()
            },
        })
    })
    
    // Subscribe to final result
    type StringSubscriber struct{}
    
    func (s *StringSubscriber) Start() {
        fmt.Println("Starting string processing pipeline")
    }
    
    func (s *StringSubscriber) OnNext(value string) {
        fmt.Printf("Processed word: %s\n", value)
    }
    
    func (s *StringSubscriber) OnError(err error) {
        fmt.Printf("Pipeline error: %v\n", err)
    }
    
    func (s *StringSubscriber) OnCompleted() {
        fmt.Println("String processing pipeline completed")
    }
    
    subscriber := &StringSubscriber{}
    transformed.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```

## 4. Aggregation Operations

Perform aggregation operations on streams:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/observable"
)

// SumAggregator accumulates values
type SumAggregator struct {
    sum   int
    count int
    next  observable.Subscriber[int]
}

func (s *SumAggregator) Start() {
    s.next.Start()
}

func (s *SumAggregator) OnNext(value int) {
    s.sum += value
    s.count++
}

func (s *SumAggregator) OnError(err error) {
    s.next.OnError(err)
}

func (s *SumAggregator) OnCompleted() {
    fmt.Printf("Sum: %d, Count: %d, Average: %.2f\n", s.sum, s.count, float64(s.sum)/float64(s.count))
    s.next.OnNext(s.sum)
    s.next.OnCompleted()
}

// AverageSubscriber receives aggregated results
type AverageSubscriber struct{}

func (s *AverageSubscriber) Start() {
    fmt.Println("Starting aggregation")
}

func (s *AverageSubscriber) OnNext(value int) {
    fmt.Printf("Final sum: %d\n", value)
}

func (s *AverageSubscriber) OnError(err error) {
    fmt.Printf("Aggregation error: %v\n", err)
}

func (s *AverageSubscriber) OnCompleted() {
    fmt.Println("Aggregation completed")
}

func main() {
    // Create source with numbers 1-5
    source := observable.Range(1, 5)
    
    // Aggregate and sum
    aggregated := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
        source.Subscribe(ctx, &SumAggregator{
            next: sub,
        })
    })
    
    // Subscribe to final result
    subscriber := &AverageSubscriber{}
    aggregated.Subscribe(context.Background(), subscriber)
    
    time.Sleep(100 * time.Millisecond)
}
```