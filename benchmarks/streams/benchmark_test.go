// Package streams provides performance benchmarks for the Reactive Streams implementation.
package streams

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/droxer/RxGo/pkg/streams"
)

type testSubscriber[T any] struct {
	received  []T
	completed bool
	errors    []error
	mu        sync.Mutex
}

func (s *testSubscriber[T]) OnSubscribe(sub streams.Subscription) {
	sub.Request(1000) // Request all items
}

func (s *testSubscriber[T]) OnNext(next T) {
	s.mu.Lock()
	s.received = append(s.received, next)
	s.mu.Unlock()
}

func (s *testSubscriber[T]) OnComplete() {
	s.mu.Lock()
	s.completed = true
	s.mu.Unlock()
}

func (s *testSubscriber[T]) OnError(e error) {
	s.mu.Lock()
	s.errors = append(s.errors, e)
	s.mu.Unlock()
}

type controlledSubscriber[T any] struct {
	received  []T
	completed bool
	errors    []error
	mu        sync.Mutex
	request   int64
}

func (s *controlledSubscriber[T]) OnSubscribe(sub streams.Subscription) {
	sub.Request(s.request)
}

func (s *controlledSubscriber[T]) OnNext(next T) {
	s.mu.Lock()
	s.received = append(s.received, next)
	s.mu.Unlock()
}

func (s *controlledSubscriber[T]) OnComplete() {
	s.mu.Lock()
	s.completed = true
	s.mu.Unlock()
}

func (s *controlledSubscriber[T]) OnError(e error) {
	s.mu.Lock()
	s.errors = append(s.errors, e)
	s.mu.Unlock()
}

func BenchmarkPublisherCreation(b *testing.B) {
	b.Run("RangePublisher", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := streams.RangePublisher(0, 100)
			_ = publisher
		}
	})

	b.Run("FromSlicePublisher", func(b *testing.B) {
		data := make([]int, 100)
		for i := 0; i < 100; i++ {
			data[i] = i
		}
		for i := 0; i < b.N; i++ {
			publisher := streams.FromSlicePublisher(data)
			_ = publisher
		}
	})
}

func BenchmarkSubscription(b *testing.B) {
	ctx := context.Background()
	b.Run("RangePublisher with subscriber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := streams.RangePublisher(0, 100)
			sub := &testSubscriber[int]{}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})

	b.Run("FromSlicePublisher with subscriber", func(b *testing.B) {
		data := make([]int, 100)
		for i := 0; i < 100; i++ {
			data[i] = i
		}
		for i := 0; i < b.N; i++ {
			publisher := streams.FromSlicePublisher(data)
			sub := &testSubscriber[int]{}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})
}

func BenchmarkBackpressureStrategies(b *testing.B) {
	ctx := context.Background()
	b.Run("Unbounded", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := streams.RangePublisher(0, 1000)
			sub := &testSubscriber[int]{}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})

	b.Run("ControlledRequest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := streams.RangePublisher(0, 1000)
			sub := &controlledSubscriber[int]{request: 10}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})

	b.Run("SingleRequest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := streams.RangePublisher(0, 1000)
			sub := &controlledSubscriber[int]{request: 1}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})
}

func BenchmarkProcessors(b *testing.B) {
	ctx := context.Background()
	b.Run("MapProcessor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := streams.RangePublisher(0, 100)
			processor := streams.NewMapProcessor(func(v int) int { return v * 2 })
			err := source.Subscribe(ctx, processor)
			if err != nil {
				b.Fatalf("Source subscribe failed: %v", err)
			}
			sub := &testSubscriber[int]{}
			err = processor.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Processor subscribe failed: %v", err)
			}
		}
	})

	b.Run("FilterProcessor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := streams.RangePublisher(0, 100)
			processor := streams.NewFilterProcessor(func(v int) bool { return v%2 == 0 })
			err := source.Subscribe(ctx, processor)
			if err != nil {
				b.Fatalf("Source subscribe failed: %v", err)
			}
			sub := &testSubscriber[int]{}
			err = processor.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Processor subscribe failed: %v", err)
			}
		}
	})

	b.Run("MapFilterChain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := streams.RangePublisher(0, 100)
			mapProc := streams.NewMapProcessor(func(v int) int { return v * 2 })
			filterProc := streams.NewFilterProcessor(func(v int) bool { return v > 50 })

			err := source.Subscribe(ctx, mapProc)
			if err != nil {
				b.Fatalf("Source subscribe failed: %v", err)
			}
			err = mapProc.Subscribe(ctx, filterProc)
			if err != nil {
				b.Fatalf("Map processor subscribe failed: %v", err)
			}

			sub := &testSubscriber[int]{}
			err = filterProc.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Filter processor subscribe failed: %v", err)
			}
		}
	})
}

func BenchmarkDatasetSizes(b *testing.B) {
	ctx := context.Background()
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				publisher := streams.RangePublisher(0, size)
				sub := &testSubscriber[int]{}
				err := publisher.Subscribe(ctx, sub)
				if err != nil {
					b.Fatalf("Subscribe failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkConcurrentSubscribers(b *testing.B) {
	ctx := context.Background()
	publisher := streams.RangePublisher(0, 100)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sub := &testSubscriber[int]{}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})
}

func BenchmarkErrorHandling(b *testing.B) {
	ctx := context.Background()
	b.Run("PublisherError", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := streams.NewPublisher(func(ctx context.Context, sub streams.Subscriber[int]) {
				sub.OnError(fmt.Errorf("test error"))
			})
			sub := &testSubscriber[int]{}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})
}

func BenchmarkMemoryAllocations(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()

	b.Run("RangePublisher", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := streams.RangePublisher(0, 100)
			sub := &testSubscriber[int]{}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})

	b.Run("FromSlicePublisher", func(b *testing.B) {
		data := make([]int, 100)
		for i := 0; i < 100; i++ {
			data[i] = i
		}
		for i := 0; i < b.N; i++ {
			publisher := streams.FromSlicePublisher(data)
			sub := &testSubscriber[int]{}
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Subscribe failed: %v", err)
			}
		}
	})

	b.Run("MapProcessor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := streams.RangePublisher(0, 100)
			processor := streams.NewMapProcessor(func(v int) int { return v * 2 })
			err := source.Subscribe(ctx, processor)
			if err != nil {
				b.Fatalf("Source subscribe failed: %v", err)
			}
			sub := &testSubscriber[int]{}
			err = processor.Subscribe(ctx, sub)
			if err != nil {
				b.Fatalf("Processor subscribe failed: %v", err)
			}
		}
	})
}
