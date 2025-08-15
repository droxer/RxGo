// Package observable provides performance benchmarks for the Observable pattern.
package observable

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/droxer/RxGo/pkg/observable"
)

type subscriber[T any] struct {
	received  []T
	completed bool
	errors    []error
	mu        sync.Mutex
}

func (s *subscriber[T]) Start() {}

func (s *subscriber[T]) OnNext(next T) {
	s.mu.Lock()
	s.received = append(s.received, next)
	s.mu.Unlock()
}

func (s *subscriber[T]) OnCompleted() {
	s.mu.Lock()
	s.completed = true
	s.mu.Unlock()
}

func (s *subscriber[T]) OnError(e error) {
	s.mu.Lock()
	s.errors = append(s.errors, e)
	s.mu.Unlock()
}

func BenchmarkCreation(b *testing.B) {
	b.Run("Just", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obs := observable.Just(1, 2, 3, 4, 5)
			_ = obs
		}
	})

	b.Run("Range", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obs := observable.Range(0, 100)
			_ = obs
		}
	})

	b.Run("Empty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obs := observable.Empty[int]()
			_ = obs
		}
	})

	b.Run("Error", func(b *testing.B) {
		err := fmt.Errorf("test error")
		for i := 0; i < b.N; i++ {
			obs := observable.Error[int](err)
			_ = obs
		}
	})
}

func BenchmarkSubscription(b *testing.B) {
	b.Run("Just with subscriber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obs := observable.Just(1, 2, 3, 4, 5)
			sub := &subscriber[int]{}
			obs.Subscribe(context.Background(), sub)
		}
	})

	b.Run("Range with subscriber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obs := observable.Range(0, 100)
			sub := &subscriber[int]{}
			obs.Subscribe(context.Background(), sub)
		}
	})
}

func BenchmarkDataTypes(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		obs := observable.Just(1, 2, 3, 4, 5)
		for i := 0; i < b.N; i++ {
			sub := &subscriber[int]{}
			obs.Subscribe(context.Background(), sub)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		obs := observable.Just(1.0, 2.0, 3.0, 4.0, 5.0)
		for i := 0; i < b.N; i++ {
			sub := &subscriber[float64]{}
			obs.Subscribe(context.Background(), sub)
		}
	})

	b.Run("String", func(b *testing.B) {
		obs := observable.Just("a", "b", "c", "d", "e")
		for i := 0; i < b.N; i++ {
			sub := &subscriber[string]{}
			obs.Subscribe(context.Background(), sub)
		}
	})
}

func BenchmarkDatasetSizes(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				obs := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
					for _, v := range data {
						sub.OnNext(v)
					}
					sub.OnCompleted()
				})
				sub := &subscriber[int]{}
				obs.Subscribe(context.Background(), sub)
			}
		})
	}
}

func BenchmarkFromSlice(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obs := observable.FromSlice(data)
		sub := &subscriber[int]{}
		obs.Subscribe(context.Background(), sub)
	}
}

func BenchmarkConcurrentSubscribers(b *testing.B) {
	obs := observable.Just(1, 2, 3, 4, 5)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sub := &subscriber[int]{}
			obs.Subscribe(context.Background(), sub)
		}
	})
}

func BenchmarkErrorHandling(b *testing.B) {
	err := fmt.Errorf("test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obs := observable.Error[int](err)
		sub := &subscriber[int]{}
		obs.Subscribe(context.Background(), sub)
	}
}

func BenchmarkContextCancellation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		obs := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
			for j := 0; j < 100; j++ {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(j)
				}
			}
			sub.OnCompleted()
		})

		sub := &subscriber[int]{}
		go obs.Subscribe(ctx, sub)
		cancel()
	}
}

func BenchmarkMemoryAllocations(b *testing.B) {
	b.ReportAllocs()

	b.Run("Just", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obs := observable.Just(1, 2, 3, 4, 5)
			sub := &subscriber[int]{}
			obs.Subscribe(context.Background(), sub)
		}
	})

	b.Run("Range", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			obs := observable.Range(0, 100)
			sub := &subscriber[int]{}
			obs.Subscribe(context.Background(), sub)
		}
	})
}

func BenchmarkStructData(b *testing.B) {
	type testStruct struct {
		ID    int
		Name  string
		Value float64
	}

	data := []testStruct{
		{ID: 1, Name: "test1", Value: 1.0},
		{ID: 2, Name: "test2", Value: 2.0},
		{ID: 3, Name: "test3", Value: 3.0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obs := observable.Create(func(ctx context.Context, sub observable.Subscriber[testStruct]) {
			for _, v := range data {
				sub.OnNext(v)
			}
			sub.OnCompleted()
		})
		sub := &subscriber[testStruct]{}
		obs.Subscribe(context.Background(), sub)
	}
}

func BenchmarkAtomicSubscriber(b *testing.B) {
	var counter atomic.Int64
	obs := observable.Range(0, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Store(0)
		sub := &atomicSubscriber{
			doNext: func(p int) {
				counter.Add(int64(p))
			},
		}
		obs.Subscribe(context.Background(), sub)
	}
}

type atomicSubscriber struct {
	doNext func(next int)
}

func (s *atomicSubscriber) Start()          {}
func (s *atomicSubscriber) OnNext(next int) { s.doNext(next) }
func (s *atomicSubscriber) OnCompleted()    {}
func (s *atomicSubscriber) OnError(e error) {}
