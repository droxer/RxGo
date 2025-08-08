# RxGo


Experimental Reactive Extensions implementation for the Go Programming Language

## Modern Example (Go 1.22+)

```go
package main

import (
    "context"
    "fmt"
    
    rx "github.com/droxer/RxGo"
)

// Type-safe subscriber with generics
type IntSubscriber struct{}

func (s *IntSubscriber) Start() {}

func (s *IntSubscriber) OnNext(value int) {
    fmt.Println(value)
}

func (s *IntSubscriber) OnError(err error) {
    fmt.Printf("Error: %v\n", err)
}

func (s *IntSubscriber) OnCompleted() {
    fmt.Println("Completed!")
}

func main() {
    // Create an observable with type safety
    rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
        for i := 0; i < 10; i++ {
            sub.OnNext(i)
        }
        sub.OnCompleted()
    }).Subscribe(context.Background(), &IntSubscriber{})
}

// Additional examples:
// rx.Just(1, 2, 3, 4, 5).Subscribe(context.Background(), subscriber)
// rx.Range(0, 10).Subscribe(context.Background(), subscriber)
```

## License

Copyright 2015

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
