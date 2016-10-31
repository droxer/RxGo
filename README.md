# RxGo

[![Build Status](https://travis-ci.org/droxer/RxGo.svg?branch=develop)](https://travis-ci.org/droxer/RxGo)

Experimental Reactive Extensions implementation for the Go Programming Language

## Example

```go
package main

import (
    rx "github.com/droxer/RxGo"
    "fmt"
)

type SampleSubscriber struct {
}

func (s *SampleSubscriber) Start() {
}

func (s *SampleSubscriber) OnNext(next interface{}) {
    if v, ok := next.(int); ok {
        fmt.Println(v)
    }
}

func (s *SampleSubscriber) OnCompleted() {
    fmt.Println("Completed!")
}

func (s *SampleSubscriber) OnError(e error) {
}

func main() {
    rx.Create(func(sub rx.Subscriber) {
        for i := 0; i < 10; i++ {
            sub.OnNext(i)
        }
        sub.OnCompleted()
    }).Subscribe(&SampleSubscriber{})
}


```

## License

Copyright 2015

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
