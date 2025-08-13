package benchmarks

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/droxer/RxGo/pkg/rx"
)

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

func BenchmarkRxGoObservableCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := rx.Just(1, 2, 3, 4, 5)
		_ = observable
	}
}

func BenchmarkRxGoObservableWithSubscriber(b *testing.B) {
	subscriber := &TestSubscriber[int]{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := rx.Just(1, 2, 3, 4, 5)
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoObservableLargeDataset(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			for _, v := range data {
				sub.OnNext(v)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoRangeObservable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := rx.Range(0, 100)
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoJustObservable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := rx.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoCreateObservable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			for j := 0; j < 100; j++ {
				sub.OnNext(j)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoObservableConcurrentSubscribers(b *testing.B) {
	observable := rx.Just(1, 2, 3, 4, 5)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			subscriber := &TestSubscriber[int]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})
}

func BenchmarkRxGoObservableMemoryAllocations(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		observable := rx.Just(1, 2, 3, 4, 5)
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoObservableStringData(b *testing.B) {
	data := []string{"hello", "world", "benchmark", "test", "performance", "observable", "reactive", "streams"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[string]) {
			for _, v := range data {
				sub.OnNext(v)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[string]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoObservableStructData(b *testing.B) {
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
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[TestStruct]) {
			for _, v := range data {
				sub.OnNext(v)
			}
			sub.OnCompleted()
		})
		subscriber := &TestSubscriber[TestStruct]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoObservableErrorHandling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			sub.OnError(fmt.Errorf("test error"))
		})
		subscriber := &TestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoObservableContextCancellation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
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

type RxGoAtomicSubscriber struct {
	doNext func(next int)
}

func (s *RxGoAtomicSubscriber) Start() {}

func (s *RxGoAtomicSubscriber) OnNext(next int) {
	s.doNext(next)
}

func (s *RxGoAtomicSubscriber) OnCompleted() {}

func (s *RxGoAtomicSubscriber) OnError(e error) {}

func BenchmarkRxGoAtomicSubscriber(b *testing.B) {
	var counter atomic.Int64

	observable := rx.Range(0, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Store(0)
		subscriber := &RxGoAtomicSubscriber{
			doNext: func(p int) {
				counter.Add(int64(p))
			},
		}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkRxGoObservableDataTypes(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := rx.Just(1, 2, 3, 4, 5)
			subscriber := &TestSubscriber[int]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := rx.Just(1.0, 2.0, 3.0, 4.0, 5.0)
			subscriber := &TestSubscriber[float64]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := rx.Just("a", "b", "c", "d", "e")
			subscriber := &TestSubscriber[string]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})
}

func BenchmarkRxGoObservableDatasetSizes(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
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

func BenchmarkRxGoBackpressureStrategies(b *testing.B) {
	b.Run("Buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := rx.Range(0, 1000)
			subscriber := &TestSubscriber[int]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("Drop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := rx.Range(0, 1000)
			subscriber := &TestSubscriber[int]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("Latest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := rx.Range(0, 1000)
			subscriber := &TestSubscriber[int]{}
			observable.Subscribe(context.Background(), subscriber)
		}
	})
}
