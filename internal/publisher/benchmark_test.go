package publisher

import (
	"context"
	"testing"
	"time"
)

// Benchmark for basic publisher creation and subscription
func BenchmarkPublisherCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publisher := FromSlice([]int{1, 2, 3, 4, 5})
		_ = publisher
	}
}

// Benchmark for publisher with simple subscriber
func BenchmarkPublisherWithSubscriber(b *testing.B) {
	subscriber := NewTestReactiveSubscriber[int]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice([]int{1, 2, 3, 4, 5})
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for publisher with large dataset
func BenchmarkPublisherLargeDataset(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice(data)
		subscriber := NewTestReactiveSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for range publisher
func BenchmarkRangePublisher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publisher := RangePublisher(0, 100)
		subscriber := NewTestReactiveSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for backpressure handling
func BenchmarkPublisherWithBackpressure(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice(data)
		subscriber := NewTestReactiveSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
		
		// Simulate backpressure
		sub := subscriber.getSubscription()
		if sub != nil {
			sub.Request(10)
		}
	}
}

// Benchmark for context cancellation
func BenchmarkPublisherWithCancellation(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		publisher := FromSlice(data)
		subscriber := NewTestReactiveSubscriber[int]()
		publisher.Subscribe(ctx, subscriber)
		cancel()
	}
}

// Benchmark for custom publisher creation
func BenchmarkCustomPublisher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publisher := NewReactivePublisher(func(ctx context.Context, sub ReactiveSubscriber[int]) {
			for j := 0; j < 100; j++ {
				sub.OnNext(j)
			}
			sub.OnComplete()
		})
		_ = publisher
	}
}

// Benchmark for concurrent subscribers
func BenchmarkPublisherConcurrentSubscribers(b *testing.B) {
	publisher := FromSlice([]int{1, 2, 3, 4, 5})
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			subscriber := NewTestReactiveSubscriber[int]()
			publisher.Subscribe(context.Background(), subscriber)
		}
	})
}

// Benchmark memory allocations for publisher creation
func BenchmarkPublisherMemoryAllocations(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice([]int{1, 2, 3, 4, 5})
		subscriber := NewTestReactiveSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for string vs int performance
func BenchmarkPublisherStringData(b *testing.B) {
	data := []string{"hello", "world", "benchmark", "test", "performance"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice(data)
		subscriber := NewTestReactiveSubscriber[string]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for struct data handling
func BenchmarkPublisherStructData(b *testing.B) {
	type TestStruct struct {
		ID    int
		Name  string
		Value float64
	}
	
	data := []TestStruct{
		{ID: 1, Name: "test1", Value: 1.0},
		{ID: 2, Name: "test2", Value: 2.0},
		{ID: 3, Name: "test3", Value: 3.0},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice(data)
		subscriber := NewTestReactiveSubscriber[TestStruct]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}