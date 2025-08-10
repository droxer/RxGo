package benchmarks

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/droxer/RxGo/pkg/observable"
)

// TestSubscriber for benchmarks
type TestSubscriber[T any] struct {
	received  []T
	completed bool
	errors    []error
	mu        sync.Mutex
}

func (t *TestSubscriber[T]) Start() {}

func (t *TestSubscriber[T]) OnNext(next T) {
	t.mu.Lock()
	t.received = append(t.received, next)
	t.mu.Unlock()
}

func (t *TestSubscriber[T]) OnCompleted() {
	t.mu.Lock()
	t.completed = true
	t.mu.Unlock()
}

func (t *TestSubscriber[T]) OnError(e error) {
	t.mu.Lock()
	t.errors = append(t.errors, e)
	t.mu.Unlock()
}

// Benchmark for observable creation
func BenchmarkObservableCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := observable.Just(1, 2, 3, 4, 5)
		_ = observable
	}
}

// Benchmark for observable with subscriber
func BenchmarkObservableWithSubscriber(b *testing.B) {
	subscriber := &TestSubscriber[int]{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := observable.Just(1, 2, 3, 4, 5)
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for observable with large dataset
func BenchmarkObservableLargeDataset(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
			for _, v := range data {
				sub.OnNext(v)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for range observable
func BenchmarkRangeObservable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := observable.Range(0, 100)
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for just observable
func BenchmarkJustObservable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := observable.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for create observable
func BenchmarkCreateObservable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
			for j := 0; j < 100; j++ {
				sub.OnNext(j)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for concurrent subscribers
func BenchmarkObservableConcurrentSubscribers(b *testing.B) {
	observable := observable.Just(1, 2, 3, 4, 5)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			subscriber := &TestSubscriber[int]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})
}

// Benchmark memory allocations for observable
func BenchmarkObservableMemoryAllocations(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		observable := observable.Just(1, 2, 3, 4, 5)
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for string vs int performance
func BenchmarkObservableStringData(b *testing.B) {
	data := []string{"hello", "world", "benchmark", "test", "performance", "observable", "reactive", "streams"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[string]) {
			for _, v := range data {
				sub.OnNext(v)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[string]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for struct data handling
func BenchmarkObservableStructData(b *testing.B) {
	type TestStruct struct {
		ID    int
		Name  string
		Value float64
	}

	data := []TestStruct{
		{ID: 1, Name: "test1", Value: 1.0},
		{ID: 2, Name: "test2", Value: 2.0},
		{ID: 3, Name: "test3", Value: 3.0},
		{ID: 4, Name: "test4", Value: 4.0},
		{ID: 5, Name: "test5", Value: 5.0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[TestStruct]) {
			for _, v := range data {
				sub.OnNext(v)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[TestStruct]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for error handling
func BenchmarkObservableErrorHandling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
			sub.OnError(fmt.Errorf("test error"))
		})
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for context cancellation
func BenchmarkObservableContextCancellation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
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

		subscriber := &TestSubscriber[int]{}
		go observable.Subscribe(ctx, subscriber)
		cancel()
	}
}

// AtomicSubscriber for benchmarks
type AtomicSubscriber struct {
	doNext func(next int)
}

func (s *AtomicSubscriber) Start() {}

func (s *AtomicSubscriber) OnNext(next int) {
	s.doNext(next)
}

func (s *AtomicSubscriber) OnCompleted() {}

func (s *AtomicSubscriber) OnError(e error) {}

// Benchmark for atomic operations in subscribers
func BenchmarkAtomicSubscriber(b *testing.B) {
	var counter atomic.Int64

	observable := observable.Range(0, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Store(0)
		subscriber := &AtomicSubscriber{
			doNext: func(p int) {
				counter.Add(int64(p))
			},
		}
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for data types comparison
func BenchmarkObservableDataTypes(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := observable.Just(1, 2, 3, 4, 5)
			subscriber := &TestSubscriber[int]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := observable.Just(1.0, 2.0, 3.0, 4.0, 5.0)
			subscriber := &TestSubscriber[float64]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := observable.Just("a", "b", "c", "d", "e")
			subscriber := &TestSubscriber[string]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})
}

// Benchmark for different dataset sizes
func BenchmarkObservableDatasetSizes(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				observable := observable.Create(func(ctx context.Context, sub observable.Subscriber[int]) {
					for _, v := range data {
						sub.OnNext(v)
					}
					sub.OnCompleted()
				})
				subscriber := &TestSubscriber[int]{}
				observable.Subscribe(context.Background(), subscriber)
			}
		})
	}
}
