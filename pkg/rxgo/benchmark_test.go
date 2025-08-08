package rxgo

import (
	"context"
	"sync"
	"testing"

	"github.com/droxer/RxGo/internal/publisher"
)

// BenchmarkSubscriber implements ReactiveSubscriber for benchmarks
type BenchmarkSubscriber[T any] struct {
	received  []T
	errors    []error
	completed bool
	mu        sync.Mutex
}

func NewBenchmarkSubscriber[T any]() *BenchmarkSubscriber[T] {
	return &BenchmarkSubscriber[T]{}
}

func (b *BenchmarkSubscriber[T]) OnSubscribe(s publisher.Subscription) {
	s.Request(1 << 62) // Request all items
}

func (b *BenchmarkSubscriber[T]) OnNext(next T) {
	b.mu.Lock()
	b.received = append(b.received, next)
	b.mu.Unlock()
}

func (b *BenchmarkSubscriber[T]) OnError(err error) {
	b.mu.Lock()
	b.errors = append(b.errors, err)
	b.mu.Unlock()
}

func (b *BenchmarkSubscriber[T]) OnComplete() {
	b.mu.Lock()
	b.completed = true
	b.mu.Unlock()
}

// BenchmarkSubscriberOld implements Subscriber for legacy observable benchmarks
type BenchmarkSubscriberOld[T any] struct {
	received  []T
	errors    []error
	completed bool
	mu        sync.Mutex
}

func NewBenchmarkSubscriberOld[T any]() *BenchmarkSubscriberOld[T] {
	return &BenchmarkSubscriberOld[T]{}
}

func (b *BenchmarkSubscriberOld[T]) Start() {}

func (b *BenchmarkSubscriberOld[T]) OnNext(next T) {
	b.mu.Lock()
	b.received = append(b.received, next)
	b.mu.Unlock()
}

func (b *BenchmarkSubscriberOld[T]) OnError(err error) {
	b.mu.Lock()
	b.errors = append(b.errors, err)
	b.mu.Unlock()
}

func (b *BenchmarkSubscriberOld[T]) OnCompleted() {
	b.mu.Lock()
	b.completed = true
	b.mu.Unlock()
}

// Benchmark for reactive publisher creation
func BenchmarkRxGoPublisherCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publisher := FromSlice([]int{1, 2, 3, 4, 5})
		_ = publisher
	}
}

// Benchmark for reactive publisher with subscriber
func BenchmarkRxGoPublisherWithSubscriber(b *testing.B) {
	subscriber := NewBenchmarkSubscriber[int]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice([]int{1, 2, 3, 4, 5})
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for observable creation
func BenchmarkRxGoObservableCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := Just(1, 2, 3, 4, 5)
		_ = observable
	}
}

// Benchmark for observable with subscriber
func BenchmarkRxGoObservableWithSubscriber(b *testing.B) {
	subscriber := NewBenchmarkSubscriberOld[int]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observable := Just(1, 2, 3, 4, 5)
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for large dataset handling
func BenchmarkRxGoLargeDataset(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice(data)
		subscriber := NewBenchmarkSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for range operations
func BenchmarkRxGoRangeOperations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := Range(0, 1000)
		subscriber := NewBenchmarkSubscriberOld[int]()
		observable.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for range publisher operations
func BenchmarkRxGoRangePublisher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publisher := RangePublisher(0, 1000)
		subscriber := NewBenchmarkSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for custom publisher creation
func BenchmarkRxGoCustomPublisher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publisher := NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[int]) {
			for j := 0; j < 100; j++ {
				sub.OnNext(j)
			}
			sub.OnComplete()
		})
		_ = publisher
	}
}

// Benchmark for custom observable creation
func BenchmarkRxGoCustomObservable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		observable := Create(func(ctx context.Context, sub Subscriber[int]) {
			for j := 0; j < 100; j++ {
				sub.OnNext(j)
			}
			sub.OnCompleted()
		})
		_ = observable
	}
}

// Benchmark for concurrent operations
func BenchmarkRxGoConcurrentPublishers(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			publisher := FromSlice([]int{1, 2, 3, 4, 5})
			subscriber := NewBenchmarkSubscriber[int]()
			publisher.Subscribe(context.Background(), subscriber)
		}
	})
}

// Benchmark for concurrent observables
func BenchmarkRxGoConcurrentObservables(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			observable := Just(1, 2, 3, 4, 5)
			subscriber := NewBenchmarkSubscriberOld[int]()
			observable.Subscribe(context.Background(), subscriber)
		}
	})
}

// Benchmark for memory allocations
func BenchmarkRxGoMemoryAllocations(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice([]int{1, 2, 3, 4, 5})
		subscriber := NewBenchmarkSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for string data handling
func BenchmarkRxGoStringData(b *testing.B) {
	data := []string{"hello", "world", "benchmark", "test", "performance"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice(data)
		subscriber := NewBenchmarkSubscriber[string]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for struct data handling
type TestStruct struct {
	ID    int
	Name  string
	Value float64
}

func BenchmarkRxGoStructData(b *testing.B) {
	data := []TestStruct{
		{ID: 1, Name: "test1", Value: 1.0},
		{ID: 2, Name: "test2", Value: 2.0},
		{ID: 3, Name: "test3", Value: 3.0},
		{ID: 4, Name: "test4", Value: 4.0},
		{ID: 5, Name: "test5", Value: 5.0},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher := FromSlice(data)
		subscriber := NewBenchmarkSubscriber[TestStruct]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for error handling
func BenchmarkRxGoErrorHandling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		publisher := NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[int]) {
			sub.OnError(context.Canceled)
		})
		subscriber := NewBenchmarkSubscriber[int]()
		publisher.Subscribe(context.Background(), subscriber)
	}
}

// Benchmark for context cancellation
func BenchmarkRxGoContextCancellation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		
		publisher := NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[int]) {
			for j := 0; j < 100; j++ {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(j)
				}
			}
			sub.OnComplete()
		})
		
		subscriber := NewBenchmarkSubscriber[int]()
		go publisher.Subscribe(ctx, subscriber)
		cancel()
	}
}

// Benchmark for performance comparison across data types
func BenchmarkRxGoDataTypes(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := FromSlice([]int{1, 2, 3, 4, 5})
			subscriber := NewBenchmarkSubscriber[int]()
			publisher.Subscribe(context.Background(), subscriber)
		}
	})
	
	b.Run("Float64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := FromSlice([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
			subscriber := NewBenchmarkSubscriber[float64]()
			publisher.Subscribe(context.Background(), subscriber)
		}
	})
	
	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := FromSlice([]string{"a", "b", "c", "d", "e"})
			subscriber := NewBenchmarkSubscriber[string]()
			publisher.Subscribe(context.Background(), subscriber)
		}
	})
}

// Benchmark for different dataset sizes
func BenchmarkRxGoDatasetSizes(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(string(rune('0'+rune(size))), func(b *testing.B) {
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				publisher := FromSlice(data)
				subscriber := NewBenchmarkSubscriber[int]()
				publisher.Subscribe(context.Background(), subscriber)
			}
		})
	}
}

// Benchmark for concurrent operations
func BenchmarkRxGoConcurrentOperations(b *testing.B) {
	b.Run("Publishers", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				publisher := FromSlice([]int{1, 2, 3, 4, 5})
				subscriber := NewBenchmarkSubscriber[int]()
				publisher.Subscribe(context.Background(), subscriber)
			}
		})
	})
	
	b.Run("Observables", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				observable := Just(1, 2, 3, 4, 5)
				subscriber := NewBenchmarkSubscriberOld[int]()
				observable.Subscribe(context.Background(), subscriber)
			}
		})
	})
}


// Benchmark for error propagation
func BenchmarkRxGoErrorPropagation(b *testing.B) {
	b.Run("Publisher", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			publisher := NewReactivePublisher(func(ctx context.Context, sub publisher.ReactiveSubscriber[int]) {
				sub.OnError(context.Canceled)
			})
			subscriber := NewBenchmarkSubscriber[int]()
			publisher.Subscribe(context.Background(), subscriber)
		}
	})
	
	b.Run("Observable", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			observable := Create(func(ctx context.Context, sub Subscriber[int]) {
				sub.OnError(context.Canceled)
			})
			subscriber := NewBenchmarkSubscriberOld[int]()
			observable.Subscribe(context.Background(), subscriber)
		}
	})
}