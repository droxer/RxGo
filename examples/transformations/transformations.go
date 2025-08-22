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

func (s *ObservableSubscriber) OnComplete() {
	fmt.Printf("[Observable %s] Completed\n", s.name)
}

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

	words := observable.Just("hello", "world", "rxgo")
	uppercased := observable.Map(words, strings.ToUpper)
	if err := uppercased.Subscribe(context.Background(), &ObservableSubscriber{name: "Mapper"}); err != nil {
		fmt.Printf("Error subscribing to uppercased observable: %v\n", err)
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== Reactive Streams API Example ===")

	wordPublisher := streams.FromSlicePublisher([]string{"hello", "world", "rxgo"})
	wordPublisher.Subscribe(context.Background(), &StreamsSubscriber{name: "Direct"})

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== New Observable Operators Example ===")

	numbers := observable.Range(1, 10)

	firstFive := observable.Take(numbers, 5)
	fmt.Println("\n--- Take Operator ---")
	if err := firstFive.Subscribe(context.Background(), observable.NewSubscriber(
		func(v int) { fmt.Printf("Taken: %d\n", v) },
		func() { fmt.Println("Take completed") },
		func(err error) { fmt.Printf("Take error: %v\n", err) },
	)); err != nil {
		fmt.Printf("Error subscribing to firstFive observable: %v\n", err)
	}

	skipped := observable.Skip(numbers, 3)
	fmt.Println("\n--- Skip Operator ---")
	if err := skipped.Subscribe(context.Background(), observable.NewSubscriber(
		func(v int) { fmt.Printf("Skipped to: %d\n", v) },
		func() { fmt.Println("Skip completed") },
		func(err error) { fmt.Printf("Skip error: %v\n", err) },
	)); err != nil {
		fmt.Printf("Error subscribing to skipped observable: %v\n", err)
	}

	evens := observable.Filter(observable.Range(1, 10), func(x int) bool { return x%2 == 0 })
	odds := observable.Filter(observable.Range(1, 10), func(x int) bool { return x%2 != 0 })
	merged := observable.Merge(evens, odds)
	fmt.Println("\n--- Merge Operator ---")
	if err := merged.Subscribe(context.Background(), observable.NewSubscriber(
		func(v int) { fmt.Printf("Merged: %d\n", v) },
		func() { fmt.Println("Merge completed") },
		func(err error) { fmt.Printf("Merge error: %v\n", err) },
	)); err != nil {
		fmt.Printf("Error subscribing to merged observable: %v\n", err)
	}

	withDuplicates := observable.FromSlice([]int{1, 2, 2, 3, 3, 3, 4, 5, 5})
	distinct := observable.Distinct(withDuplicates)
	fmt.Println("\n--- Distinct Operator ---")
	if err := distinct.Subscribe(context.Background(), observable.NewSubscriber(
		func(v int) { fmt.Printf("Distinct: %d\n", v) },
		func() { fmt.Println("Distinct completed") },
		func(err error) { fmt.Printf("Distinct error: %v\n", err) },
	)); err != nil {
		fmt.Printf("Error subscribing to distinct observable: %v\n", err)
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== New Stream Processors Example ===")

	numberPublisher := streams.NewCompliantRangePublisher(1, 10)

	takeProcessor := streams.NewTakeProcessor[int](5)
	numberPublisher.Subscribe(context.Background(), takeProcessor)
	fmt.Println("\n--- Take Processor ---")
	takeProcessor.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Taken: %d\n", v) },
		func(err error) { fmt.Printf("Take error: %v\n", err) },
		func() { fmt.Println("Take completed") },
	))

	numberPublisher = streams.NewCompliantRangePublisher(1, 10)

	skipProcessor := streams.NewSkipProcessor[int](3)
	numberPublisher.Subscribe(context.Background(), skipProcessor)
	fmt.Println("\n--- Skip Processor ---")
	skipProcessor.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Skipped to: %d\n", v) },
		func(err error) { fmt.Printf("Skip error: %v\n", err) },
		func() { fmt.Println("Skip completed") },
	))

	numberPublisher = streams.NewCompliantRangePublisher(1, 10)

	duplicatePublisher := streams.NewCompliantFromSlicePublisher([]int{1, 2, 2, 3, 3, 3, 4, 5, 5})
	distinctProcessor := streams.NewDistinctProcessor[int]()
	duplicatePublisher.Subscribe(context.Background(), distinctProcessor)
	fmt.Println("\n--- Distinct Processor ---")
	distinctProcessor.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Distinct: %d\n", v) },
		func(err error) { fmt.Printf("Distinct error: %v\n", err) },
		func() { fmt.Println("Distinct completed") },
	))

	// Additional stream processors examples
	fmt.Println("\n--- Merge Processor ---")
	pub1 := streams.NewCompliantFromSlicePublisher([]int{1, 2, 3})
	pub2 := streams.NewCompliantFromSlicePublisher([]int{4, 5, 6})
	mergeProcessor := streams.NewMergeProcessor[int](pub1, pub2)
	mergeProcessor.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Merged processor value: %d\n", v) },
		func(err error) { fmt.Printf("Merge processor error: %v\n", err) },
		func() { fmt.Println("Merge processor completed") },
	))

	fmt.Println("\n--- Concat Processor ---")
	pub3 := streams.NewCompliantFromSlicePublisher([]int{7, 8, 9})
	pub4 := streams.NewCompliantFromSlicePublisher([]int{10, 11, 12})
	concatProcessor := streams.NewConcatProcessor[int](pub3, pub4)
	concatProcessor.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Concat processor value: %d\n", v) },
		func(err error) { fmt.Printf("Concat processor error: %v\n", err) },
		func() { fmt.Println("Concat processor completed") },
	))

	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n=== All examples completed ===")
}
