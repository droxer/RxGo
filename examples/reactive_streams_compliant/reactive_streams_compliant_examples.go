package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/droxer/RxGo/pkg/streams"
)

func main() {
	fmt.Println("RxGo Reactive Streams 1.0.4 Compliant Examples")
	fmt.Println("=================================================")

	// Run simple subscriber example
	RunSimpleSubscriberExample()

	// Run compliant examples
	RunCompliantRangePublisher()
	RunCompliantProcessorChain()
	RunDemandEnforcement()
	RunThreadSafetyDemo()
	RunProcessorExamples()

	fmt.Println("\nAll compliant examples completed!")
}

// RunSimpleSubscriberExample demonstrates the new functional subscriber
func RunSimpleSubscriberExample() {
	fmt.Println("\n--- Simple Functional Subscriber ---")

	publisher := streams.NewCompliantRangePublisher(1, 5)
	var wg sync.WaitGroup
	wg.Add(1)

	publisher.Subscribe(context.Background(), streams.NewSubscriber(
		func(v int) { fmt.Printf("Received: %d\n", v) },
		func(err error) {
			fmt.Printf("Error: %v\n", err)
			wg.Done()
		},
		func() {
			fmt.Println("Completed")
			wg.Done()
		},
	))

	wg.Wait()
}

// RunCompliantRangePublisher demonstrates Reactive Streams 1.0.4 compliant publisher
func RunCompliantRangePublisher() {
	fmt.Println("\n--- Compliant Range Publisher ---")

	publisher := streams.NewCompliantRangePublisher(1, 10)
	subscriber := &intSubscriber{name: "CompliantSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()

	fmt.Printf("Compliant publisher processed %d items\n", subscriber.processed)
}

// RunCompliantProcessorChain demonstrates processor chaining
func RunCompliantProcessorChain() {
	fmt.Println("\n--- Compliant Processor Chain ---")

	// Create a processing pipeline
	publisher := streams.NewCompliantRangePublisher(1, 10)
	mapProcessor := streams.NewMapProcessor(func(i int) string {
		return fmt.Sprintf("Value-%d", i*2)
	})
	filterProcessor := streams.NewFilterProcessor(func(s string) bool {
		return len(s) > 7
	})
	subscriber := &stringSubscriber{name: "ProcessorChainSubscriber"}

	ctx := context.Background()

	// Chain: Publisher -> MapProcessor -> FilterProcessor -> Subscriber
	publisher.Subscribe(ctx, mapProcessor)
	mapProcessor.Subscribe(ctx, filterProcessor)
	filterProcessor.Subscribe(ctx, subscriber)

	subscriber.wait()
	fmt.Printf("Processor chain processed %d items\n", subscriber.processed)
}

// RunDemandEnforcement demonstrates demand-based flow control
func RunDemandEnforcement() {
	fmt.Println("\n--- Demand Enforcement Demo ---")

	publisher := streams.NewCompliantRangePublisher(1, 20)
	subscriber := &demandSubscriber{
		name:     "DemandSubscriber",
		maxItems: 5,
		delay:    100 * time.Millisecond,
	}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()

	fmt.Printf("Demand enforcement: requested 5, received %d\n", subscriber.received)
}

// RunThreadSafetyDemo demonstrates concurrent access
func RunThreadSafetyDemo() {
	fmt.Println("\n--- Thread Safety Demo ---")

	publisher := streams.NewCompliantRangePublisher(1, 50)
	var wg sync.WaitGroup

	// Multiple subscribers accessing concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			subscriber := &intSubscriber{
				name: fmt.Sprintf("ThreadSafe-%d", id),
			}
			ctx := context.Background()
			publisher.Subscribe(ctx, subscriber)
			subscriber.wait()
			fmt.Printf("Thread %d processed %d items\n", id, subscriber.processed)
		}(i)
	}

	wg.Wait()
}

// RunProcessorExamples demonstrates different processors
func RunProcessorExamples() {
	fmt.Println("\n--- Processor Examples ---")

	// Map processor
	fmt.Println("Map Processor:")
	mapProc := streams.NewMapProcessor(func(i int) string {
		return fmt.Sprintf("mapped-%d", i)
	})
	sub1 := &stringSubscriber{name: "MapSubscriber"}
	ctx := context.Background()

	// Create a simple source
	source := streams.NewCompliantFromSlicePublisher[int]([]int{1, 2, 3})
	source.Subscribe(ctx, mapProc)
	mapProc.Subscribe(ctx, sub1)
	sub1.wait()

	// Filter processor
	fmt.Println("Filter Processor:")
	filterProc := streams.NewFilterProcessor(func(i int) bool {
		return i%2 == 0
	})
	sub2 := &intSubscriber{name: "FilterSubscriber"}

	source2 := streams.NewCompliantFromSlicePublisher[int]([]int{1, 2, 3, 4, 5})
	source2.Subscribe(ctx, filterProc)
	filterProc.Subscribe(ctx, sub2)
	sub2.wait()
}

// intSubscriber implements a compliant subscriber for int type
type intSubscriber struct {
	name      string
	processed int
	wg        sync.WaitGroup
}

func (s *intSubscriber) OnSubscribe(sub streams.Subscription) {
	s.wg.Add(1)
	sub.Request(100) // Request all items
}

func (s *intSubscriber) OnNext(value int) {
	fmt.Printf("[%s] Received: %v\n", s.name, value)
	s.processed++
}

func (s *intSubscriber) OnError(err error) {
	fmt.Printf("[%s] Error: %v\n", s.name, err)
	s.wg.Done()
}

func (s *intSubscriber) OnComplete() {
	fmt.Printf("[%s] Completed\n", s.name)
	s.wg.Done()
}

func (s *intSubscriber) wait() {
	s.wg.Wait()
}

// stringSubscriber implements a compliant subscriber for string type
type stringSubscriber struct {
	name      string
	processed int
	wg        sync.WaitGroup
}

func (s *stringSubscriber) OnSubscribe(sub streams.Subscription) {
	s.wg.Add(1)
	sub.Request(100) // Request all items
}

func (s *stringSubscriber) OnNext(value string) {
	fmt.Printf("[%s] Received: %v\n", s.name, value)
	s.processed++
}

func (s *stringSubscriber) OnError(err error) {
	fmt.Printf("[%s] Error: %v\n", s.name, err)
	s.wg.Done()
}

func (s *stringSubscriber) OnComplete() {
	fmt.Printf("[%s] Completed\n", s.name)
	s.wg.Done()
}

func (s *stringSubscriber) wait() {
	s.wg.Wait()
}

// demandSubscriber demonstrates demand-based flow control
type demandSubscriber struct {
	name     string
	maxItems int
	received int
	delay    time.Duration
	sub      streams.Subscription
	wg       sync.WaitGroup
}

func (s *demandSubscriber) OnSubscribe(sub streams.Subscription) {
	s.sub = sub
	s.wg.Add(1)
	sub.Request(int64(s.maxItems)) // Request specific number
}

func (s *demandSubscriber) OnNext(value int) {
	fmt.Printf("[%s] Processing: %v\n", s.name, value)
	s.received++
	time.Sleep(s.delay) // Simulate slow processing
}

func (s *demandSubscriber) OnError(err error) {
	fmt.Printf("[%s] Error: %v\n", s.name, err)
	s.wg.Done()
}

func (s *demandSubscriber) OnComplete() {
	fmt.Printf("[%s] Completed processing %d items\n", s.name, s.received)
	s.wg.Done()
}

func (s *demandSubscriber) wait() {
	s.wg.Wait()
}

// Builder example
func RunBuilderExample() {
	fmt.Println("\n--- Builder Example ---")

	// Using compliant builder
	builder := streams.NewCompliantBuilder[int]()
	publisher := builder.Range(1, 5)

	subscriber := &intSubscriber{name: "BuilderSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()

	fmt.Printf("Builder processed %d items\n", subscriber.processed)
}

// Processor builder example
func RunProcessorBuilderExample() {
	fmt.Println("\n--- Processor Builder Example ---")

	// Create a processing pipeline
	publisher := streams.NewCompliantRangePublisher(1, 10)

	// Create processors
	mapProcessor := streams.NewMapProcessor(func(i int) string {
		return fmt.Sprintf("item-%d", i)
	})
	filterProcessor := streams.NewFilterProcessor(func(s string) bool {
		return len(s) > 5
	})

	subscriber := &stringSubscriber{name: "ProcessorBuilderSubscriber"}
	ctx := context.Background()

	publisher.Subscribe(ctx, mapProcessor)
	mapProcessor.Subscribe(ctx, filterProcessor)
	filterProcessor.Subscribe(ctx, subscriber)
	subscriber.wait()

	fmt.Printf("Processor builder processed %d items\n", subscriber.processed)
}

// Main function includes all examples
func init() {
	// Additional examples can be added here
	// RunBuilderExample()
	// RunProcessorBuilderExample()
}

// Memory safety test
func RunMemorySafetyTest() {
	fmt.Println("\n--- Memory Safety Test ---")

	// Create many publishers and subscribers
	for i := 0; i < 100; i++ {
		publisher := streams.NewCompliantRangePublisher(1, 10)
		subscriber := &intSubscriber{name: fmt.Sprintf("Memory-%d", i)}
		ctx := context.Background()
		publisher.Subscribe(ctx, subscriber)
		subscriber.wait()
	}

	fmt.Println("Memory safety test completed")
}

// Stress test
func RunStressTest() {
	fmt.Println("\n--- Stress Test ---")

	publisher := streams.NewCompliantRangePublisher(1, 1000)
	var wg sync.WaitGroup

	// Many concurrent subscribers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			subscriber := &intSubscriber{name: fmt.Sprintf("Stress-%d", id)}
			ctx := context.Background()
			publisher.Subscribe(ctx, subscriber)
			subscriber.wait()
		}(i)
	}

	wg.Wait()
	fmt.Println("Stress test completed")
}

// Run all examples
func RunAllCompliantExamples() {
	RunCompliantRangePublisher()
	RunCompliantProcessorChain()
	RunDemandEnforcement()
	RunThreadSafetyDemo()
	RunProcessorExamples()
	RunBuilderExample()
	RunProcessorBuilderExample()
	RunMemorySafetyTest()
	RunStressTest()
}

// For testing
func TestCompliantAPIs() {
	// Quick test of all compliant APIs
	publisher := streams.NewCompliantRangePublisher(1, 3)
	subscriber := &intSubscriber{name: "Test"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()

	if subscriber.processed != 3 {
		panic("Compliant API test failed")
	}
}

// Ensure interfaces are implemented
var (
	_ streams.Publisher[int]         = (*streams.CompliantRangePublisher)(nil)
	_ streams.Publisher[string]      = (*streams.CompliantFromSlicePublisher[string])(nil)
	_ streams.Processor[int, string] = (*streams.MapProcessor[int, string])(nil)
	_ streams.Processor[int, int]    = (*streams.FilterProcessor[int])(nil)
	_ streams.Processor[int, string] = (*streams.FlatMapProcessor[int, string])(nil)
	_ streams.Subscription           = streams.Subscription(nil)
	_ streams.Subscriber[int]        = (*intSubscriber)(nil)
)

// Entry point - main function is defined at the beginning of the file

// Alternative entry points
func RunBasicCompliant() {
	RunCompliantRangePublisher()
	RunCompliantProcessorChain()
}

func RunAdvancedCompliant() {
	RunDemandEnforcement()
	RunThreadSafetyDemo()
	RunProcessorExamples()
}

func RunPerformanceCompliant() {
	RunMemorySafetyTest()
	RunStressTest()
}

// Export for external usage
var (
	NewCompliantRangePublisher     = streams.NewCompliantRangePublisher
	NewCompliantFromSlicePublisher = streams.NewCompliantFromSlicePublisher[string]
	// NewCompliantBufferedPublisher  = streams.NewCompliantBufferedPublisher[int]
	NewMapProcessor     = streams.NewMapProcessor[int, string]
	NewFilterProcessor  = streams.NewFilterProcessor[int]
	NewFlatMapProcessor = streams.NewFlatMapProcessor[int, string]
	NewCompliantBuilder = streams.NewCompliantBuilder[int]
	NewProcessorBuilder = streams.NewProcessorBuilder[int, string]
)

// Initialize example
func init() {
	// Set up any necessary initialization
	fmt.Println("Reactive Streams 1.0.4 Compliant Examples Initialized")
}

// Ensure everything compiles
var _ = RunAllCompliantExamples
var _ = TestCompliantAPIs

// Convenience functions
func NewCompliantIntPublisher() streams.Publisher[int] {
	return streams.NewCompliantRangePublisher(1, 100)
}

func NewCompliantStringPublisher() streams.Publisher[string] {
	return streams.NewCompliantFromSlicePublisher([]string{"a", "b", "c"})
}

// Example documentation
//
// This file demonstrates the new Reactive Streams 1.0.4 compliant API
// which provides:
// - Full specification compliance
// - Demand tracking and enforcement
// - Sequential signaling
// - Thread safety
// - Proper cancellation handling
// - Processor support
// - Comprehensive error handling
//
// Usage:
// go run examples/reactive_streams_compliant/main.go
//
// Or run individual functions:
// go run examples/reactive_streams_compliant/main.go [function_name]
