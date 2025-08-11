# Basic Usage

This guide demonstrates the simple Observable API using the unified `rxgo` package.

## Creating Observables

### Using Just

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

// Simple subscriber implementation
type IntSubscriber struct{}

func (s *IntSubscriber) Start() {}
func (s *IntSubscriber) OnNext(value int) { fmt.Println(value) }
func (s *IntSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *IntSubscriber) OnCompleted() { fmt.Println("Completed!") }

func main() {
    // Create observable from literal values
    obs := rxgo.Just(1, 2, 3, 4, 5)
    obs.Subscribe(context.Background(), &IntSubscriber{})
}
```

### Using Range

```go
// Create observable from range of integers
obs := rxgo.Range(1, 10) // Emits 1, 2, 3, ..., 10
obs.Subscribe(context.Background(), &IntSubscriber{})
```

### Using Create

```go
// Create custom observable
obs := rxgo.Create(func(ctx context.Context, sub rxgo.Subscriber[string]) {
    messages := []string{"hello", "world", "from", "rxgo"}
    
    for _, msg := range messages {
        select {
        case <-ctx.Done():
            sub.OnError(ctx.Err())
            return
        default:
            sub.OnNext(msg)
        }
    }
    sub.OnCompleted()
})
```

### From Slice

```go
// Create observable from slice
numbers := []int{10, 20, 30, 40, 50}
obs := rxgo.FromSlice(numbers)
obs.Subscribe(context.Background(), &IntSubscriber{})
```

## Data Transformations

### Map Operation

Transform each value in the stream:

```go
// Map integers to their squares
numbers := rxgo.Range(1, 5)
squares := rxgo.Map(numbers, func(x int) int {
    return x * x
})

// Output: 1, 4, 9, 16, 25
squares.Subscribe(context.Background(), &IntSubscriber{})
```

### Filter Operation

Filter values based on predicate:

```go
// Filter even numbers
evens := rxgo.Filter(rxgo.Range(1, 10), func(x int) bool {
    return x%2 == 0
})

// Output: 2, 4, 6, 8, 10
evens.Subscribe(context.Background(), &IntSubscriber{})
```

### Chain Operations

Combine multiple operations:

```go
// Chain map and filter operations
result := rxgo.Filter(
    rxgo.Map(rxgo.Range(1, 10), func(x int) int { return x * 2 }),
    func(x int) bool { return x > 10 },
)

// Output: 12, 14, 16, 18, 20
result.Subscribe(context.Background(), &IntSubscriber{})
```

## Error Handling

### Error Observable

```go
// Create observable that immediately emits error
obs := rxgo.Error[string](errors.New("something went wrong"))
obs.Subscribe(context.Background(), &StringSubscriber{})
```

### Empty Observable

```go
// Create observable that completes without emitting values
obs := rxgo.Empty[string]()
obs.Subscribe(context.Background(), &StringSubscriber{})
```

## Context Cancellation

Use context for graceful cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// Infinite stream that respects context cancellation
obs := rxgo.Create(func(ctx context.Context, sub rxgo.Subscriber[int]) {
    for i := 0; ; i++ {
        select {
        case <-ctx.Done():
            sub.OnError(ctx.Err())
            return
        default:
            sub.OnNext(i)
            time.Sleep(100 * time.Millisecond)
        }
    }
})

obs.Subscribe(ctx, &IntSubscriber{})
```

## Complete Example

Here's a complete example combining multiple concepts:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/droxer/RxGo/pkg/rxgo"
)

type User struct {
    ID   int
    Name string
    Age  int
}

type UserSubscriber struct{}

func (s *UserSubscriber) Start() {}
func (s *UserSubscriber) OnNext(user User) {
    fmt.Printf("Processing user: %s (ID: %d, Age: %d)\n", user.Name, user.ID, user.Age)
}
func (s *UserSubscriber) OnError(err error) { fmt.Printf("Error: %v\n", err) }
func (s *UserSubscriber) OnCompleted() { fmt.Println("User processing completed!") }

func main() {
    // Create users observable
    users := rxgo.Just(
        User{ID: 1, Name: "Alice", Age: 25},
        User{ID: 2, Name: "Bob", Age: 30},
        User{ID: 3, Name: "Charlie", Age: 35},
        User{ID: 4, Name: "Diana", Age: 28},
    )

    // Filter adults (age >= 30)
    adults := rxgo.Filter(users, func(u User) bool {
        return u.Age >= 30
    })

    // Transform to user names
    names := rxgo.Map(adults, func(u User) string {
        return u.Name
    })

    // Subscribe and process
    names.Subscribe(context.Background(), &UserSubscriber{})
}
```

## Output

```
Processing user: Bob (ID: 2, Age: 30)
Processing user: Charlie (ID: 3, Age: 35)
User processing completed!
```