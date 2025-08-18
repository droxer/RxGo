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

func RunCompliantRangePublisher() {
	fmt.Println("\n--- Compliant Range Publisher ---")

	publisher := streams.NewCompliantRangePublisher(1, 10)
	subscriber := &intSubscriber{name: "CompliantSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()

	fmt.Printf("Compliant publisher processed %d items\n", subscriber.processed)
}

func RunCompliantProcessorChain() {
	fmt.Println("\n--- Compliant Processor Chain ---")

	publisher := streams.NewCompliantRangePublisher(1, 10)
	mapProcessor := streams.NewMapProcessor(func(i int) string {
		return fmt.Sprintf("Value-%d", i*2)
	})
	filterProcessor := streams.NewFilterProcessor(func(s string) bool {
		return len(s) > 7
	})
	subscriber := &stringSubscriber{name: "ProcessorChainSubscriber"}

	ctx := context.Background()

	publisher.Subscribe(ctx, mapProcessor)
	mapProcessor.Subscribe(ctx, filterProcessor)
	filterProcessor.Subscribe(ctx, subscriber)

	subscriber.wait()
	fmt.Printf("Processor chain processed %d items\n", subscriber.processed)
}

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

func RunThreadSafetyDemo() {
	fmt.Println("\n--- Thread Safety Demo ---")

	publisher := streams.NewCompliantRangePublisher(1, 50)
	var wg sync.WaitGroup

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

func RunProcessorExamples() {
	fmt.Println("\n--- Processor Examples ---")

	fmt.Println("Map Processor:")
	mapProc := streams.NewMapProcessor(func(i int) string {
		return fmt.Sprintf("mapped-%d", i)
	})
	sub1 := &stringSubscriber{name: "MapSubscriber"}
	ctx := context.Background()

	source := streams.NewCompliantFromSlicePublisher[int]([]int{1, 2, 3})
	source.Subscribe(ctx, mapProc)
	mapProc.Subscribe(ctx, sub1)
	sub1.wait()

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

type intSubscriber struct {
	name      string
	processed int
	wg        sync.WaitGroup
}

func (s *intSubscriber) OnSubscribe(sub streams.Subscription) {
	s.wg.Add(1)
	sub.Request(100)
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

type stringSubscriber struct {
	name      string
	processed int
	wg        sync.WaitGroup
}

func (s *stringSubscriber) OnSubscribe(sub streams.Subscription) {
	s.wg.Add(1)
	sub.Request(100)
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
	sub.Request(int64(s.maxItems))
}

func (s *demandSubscriber) OnNext(value int) {
	fmt.Printf("[%s] Processing: %v\n", s.name, value)
	s.received++
	time.Sleep(s.delay)
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

func RunBuilderExample() {
	fmt.Println("\n--- Builder Example ---")

	builder := streams.NewCompliantBuilder[int]()
	publisher := builder.Range(1, 5)

	subscriber := &intSubscriber{name: "BuilderSubscriber"}
	ctx := context.Background()
	publisher.Subscribe(ctx, subscriber)
	subscriber.wait()

	fmt.Printf("Builder processed %d items\n", subscriber.processed)
}

func RunProcessorBuilderExample() {
	fmt.Println("\n--- Processor Builder Example ---")

	publisher := streams.NewCompliantRangePublisher(1, 10)

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

func init() {
	RunBuilderExample()
	RunProcessorBuilderExample()
}

func RunMemorySafetyTest() {
	fmt.Println("\n--- Memory Safety Test ---")

	for i := 0; i < 100; i++ {
		publisher := streams.NewCompliantRangePublisher(1, 10)
		subscriber := &intSubscriber{name: fmt.Sprintf("Memory-%d", i)}
		ctx := context.Background()
		publisher.Subscribe(ctx, subscriber)
		subscriber.wait()
	}

	fmt.Println("Memory safety test completed")
}

func RunStressTest() {
	fmt.Println("\n--- Stress Test ---")

	publisher := streams.NewCompliantRangePublisher(1, 1000)
	var wg sync.WaitGroup

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

func TestCompliantAPIs() {
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

var (
	NewCompliantRangePublisher     = streams.NewCompliantRangePublisher
	NewCompliantFromSlicePublisher = streams.NewCompliantFromSlicePublisher[string]
	NewMapProcessor                = streams.NewMapProcessor[int, string]
	NewFilterProcessor             = streams.NewFilterProcessor[int]
	NewFlatMapProcessor            = streams.NewFlatMapProcessor[int, string]
	NewCompliantBuilder            = streams.NewCompliantBuilder[int]
	NewProcessorBuilder            = streams.NewProcessorBuilder[int, string]
)

func init() {
	fmt.Println("Reactive Streams 1.0.4 Compliant Examples Initialized")
}

var _ = RunAllCompliantExamples
var _ = TestCompliantAPIs

func NewCompliantIntPublisher() streams.Publisher[int] {
	return streams.NewCompliantRangePublisher(1, 100)
}

func NewCompliantStringPublisher() streams.Publisher[string] {
	return streams.NewCompliantFromSlicePublisher([]string{"a", "b", "c"})
}
