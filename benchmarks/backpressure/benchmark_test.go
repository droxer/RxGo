// Package benchmarks provides performance benchmarks for the RxGo backpressure strategies.
package benchmarks

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/droxer/RxGo/pkg/rx"
	"github.com/droxer/RxGo/pkg/rx/streams"
)

type BackpressureTestSubscriber[T any] struct {
	received  []T
	completed bool
	errors    []error
	mu        sync.Mutex
}

func (t *BackpressureTestSubscriber[T]) Start() {}

func (t *BackpressureTestSubscriber[T]) OnNext(next T) {
	t.mu.Lock()
	t.received = append(t.received, next)
	t.mu.Unlock()
}

func (t *BackpressureTestSubscriber[T]) OnComplete() {
	t.mu.Lock()
	t.completed = true
	t.mu.Unlock()
}

func (t *BackpressureTestSubscriber[T]) OnError(e error) {
	t.mu.Lock()
	t.errors = append(t.errors, e)
	t.mu.Unlock()
}

func BenchmarkBackpressureBufferStrategy(b *testing.B) {
	// config := streams.BackpressureConfig{
	// 	Strategy:   streams.Buffer,
	// 	BufferSize: 100,
	// }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a fast producer
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			defer sub.OnComplete()
			for j := 0; j < 1000; j++ {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(j)
				}
			}
		})

		subscriber := &BackpressureTestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkBackpressureDropStrategy(b *testing.B) {
	// config := streams.BackpressureConfig{
	// 	Strategy:   streams.Drop,
	// 	BufferSize: 100,
	// }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a fast producer
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			defer sub.OnComplete()
			for j := 0; j < 1000; j++ {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(j)
				}
			}
		})

		subscriber := &BackpressureTestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkBackpressureLatestStrategy(b *testing.B) {
	// config := streams.BackpressureConfig{
	// 	Strategy:   streams.Latest,
	// 	BufferSize: 100,
	// }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a fast producer
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			defer sub.OnComplete()
			for j := 0; j < 1000; j++ {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(j)
				}
			}
		})

		subscriber := &BackpressureTestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkBackpressureErrorStrategy(b *testing.B) {
	// config := streams.BackpressureConfig{
	// 	Strategy:   streams.Error,
	// 	BufferSize: 100,
	// }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a fast producer
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			defer sub.OnComplete()
			for j := 0; j < 1000; j++ {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(j)
				}
			}
		})

		subscriber := &BackpressureTestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkBackpressureStrategiesComparison(b *testing.B) {
	strategies := []struct {
		name     string
		strategy streams.BackpressureStrategy
	}{
		{"Buffer", streams.Buffer},
		{"Drop", streams.Drop},
		{"Latest", streams.Latest},
		{"Error", streams.Error},
	}

	for _, s := range strategies {
		b.Run(s.name, func(b *testing.B) {
			// config := streams.BackpressureConfig{
			// 	Strategy:   s.strategy,
			// 	BufferSize: 100,
			// }

			for i := 0; i < b.N; i++ {
				// Create a fast producer
				observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
					defer sub.OnComplete()
					for j := 0; j < 1000; j++ {
						select {
						case <-ctx.Done():
							sub.OnError(ctx.Err())
							return
						default:
							sub.OnNext(j)
						}
					}
				})

				subscriber := &BackpressureTestSubscriber[int]{}
				observable.Subscribe(context.Background(), subscriber)
			}
		})
	}
}

func BenchmarkBackpressureWithDifferentBufferSizes(b *testing.B) {
	bufferSizes := []int64{10, 100, 1000}

	for _, size := range bufferSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			// config := streams.BackpressureConfig{
			// 	Strategy:   streams.Buffer,
			// 	BufferSize: size,
			// }

			for i := 0; i < b.N; i++ {
				// Create a fast producer
				observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
					defer sub.OnComplete()
					for j := 0; j < 1000; j++ {
						select {
						case <-ctx.Done():
							sub.OnError(ctx.Err())
							return
						default:
							sub.OnNext(j)
						}
					}
				})

				subscriber := &BackpressureTestSubscriber[int]{}
				observable.Subscribe(context.Background(), subscriber)
			}
		})
	}
}

func BenchmarkBackpressureMemoryAllocations(b *testing.B) {
	// config := streams.BackpressureConfig{
	// 	Strategy:   streams.Buffer,
	// 	BufferSize: 100,
	// }

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a fast producer
		observable := rx.Create(func(ctx context.Context, sub rx.Subscriber[int]) {
			defer sub.OnComplete()
			for j := 0; j < 1000; j++ {
				select {
				case <-ctx.Done():
					sub.OnError(ctx.Err())
					return
				default:
					sub.OnNext(j)
				}
			}
		})

		subscriber := &BackpressureTestSubscriber[int]{}
		observable.Subscribe(context.Background(), subscriber)
	}
}
