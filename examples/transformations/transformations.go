package main

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
	"github.com/droxer/RxGo/pkg/streams"
)

// Observable API Example

type ObservableSubscriber struct {
	name string
}

func (s *ObservableSubscriber) Start() {
	fmt.Printf("[Observable %s] Starting\n", s.name)
}

func (s *ObservableSubscriber) OnNext(value string) {
	fmt.Printf("[Observable %s] Received: %s\n", s.name, value)
}

func (s *ObservableSubscriber) OnError(err error) {
	fmt.Printf("[Observable %s] Error: %v\n", s.name, err)
}

func (s *ObservableSubscriber) OnCompleted() {
	fmt.Printf("[Observable %s] Completed\n", s.name)
}

// Reactive Streams API Example

type StreamsSubscriber struct {
	name string
}

func (s *StreamsSubscriber) OnSubscribe(sub streams.Subscription) {
	fmt.Printf("[Streams %s] Starting\n", s.name)
	sub.Request(math.MaxInt64) // Request all items
}

func (s *StreamsSubscriber) OnNext(value string) {
	fmt.Printf("[Streams %s] Received: %s\n", s.name, value)
}

func (s *StreamsSubscriber) OnError(err error) {
	fmt.Printf("[Streams %s] Error: %v\n", s.name, err)
}

func (s *StreamsSubscriber) OnComplete() {
	fmt.Printf("[Streams %s] Completed\n", s.name)
}

func main() {
	fmt.Println("=== Observable API Example ===")

	// Observable API transformation
	words := observable.Just("hello", "world", "rxgo")
	uppercased := observable.Map(words, strings.ToUpper)
	uppercased.Subscribe(context.Background(), &ObservableSubscriber{name: "Mapper"})

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== Reactive Streams API Example ===")

	// Reactive Streams API transformation
	wordPublisher := streams.FromSlicePublisher([]string{"hello", "world", "rxgo"})
	wordPublisher.Subscribe(context.Background(), &StreamsSubscriber{name: "Direct"})

	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n=== Both examples completed ===")
}
