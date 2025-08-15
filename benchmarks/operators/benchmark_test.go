// Package operators provides performance benchmarks for Observable operators.
package operators

import (
	"context"
	"fmt"
	"sync"
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

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		source := observable.Range(0, 1000)
		mapped := observable.Map(source, func(v int) int { return v * 2 })
		sub := &subscriber[int]{}
		mapped.Subscribe(context.Background(), sub)
	}
}

func BenchmarkFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		source := observable.Range(0, 1000)
		filtered := observable.Filter(source, func(v int) bool { return v%2 == 0 })
		sub := &subscriber[int]{}
		filtered.Subscribe(context.Background(), sub)
	}
}

func BenchmarkMapFilterChain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		source := observable.Range(0, 1000)
		mapped := observable.Map(source, func(v int) int { return v * 2 })
		filtered := observable.Filter(mapped, func(v int) bool { return v > 500 })
		sub := &subscriber[int]{}
		filtered.Subscribe(context.Background(), sub)
	}
}

func BenchmarkDataTypes(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := observable.Range(0, 100)
			mapped := observable.Map(source, func(v int) int { return v * 2 })
			sub := &subscriber[int]{}
			mapped.Subscribe(context.Background(), sub)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := observable.Just(1.0, 2.0, 3.0, 4.0, 5.0)
			mapped := observable.Map(source, func(v float64) float64 { return v * 2.0 })
			sub := &subscriber[float64]{}
			mapped.Subscribe(context.Background(), sub)
		}
	})

	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := observable.Just("a", "b", "c", "d", "e")
			mapped := observable.Map(source, func(v string) string { return v + "x" })
			sub := &subscriber[string]{}
			mapped.Subscribe(context.Background(), sub)
		}
	})
}

func BenchmarkDatasetSizes(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				source := observable.Range(0, size)
				mapped := observable.Map(source, func(v int) int { return v * 2 })
				sub := &subscriber[int]{}
				mapped.Subscribe(context.Background(), sub)
			}
		})
	}
}

func BenchmarkComplexChains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		source := observable.Range(0, 1000)
		mapped1 := observable.Map(source, func(v int) int { return v * 2 })
		filtered := observable.Filter(mapped1, func(v int) bool { return v > 500 })
		mapped2 := observable.Map(filtered, func(v int) int { return v + 1 })
		sub := &subscriber[int]{}
		mapped2.Subscribe(context.Background(), sub)
	}
}

func BenchmarkConcurrentChains(b *testing.B) {
	source := observable.Range(0, 1000)
	mapped := observable.Map(source, func(v int) int { return v * 2 })
	filtered := observable.Filter(mapped, func(v int) bool { return v > 500 })

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sub := &subscriber[int]{}
			filtered.Subscribe(context.Background(), sub)
		}
	})
}

func BenchmarkMemoryAllocations(b *testing.B) {
	b.ReportAllocs()
	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := observable.Range(0, 1000)
			mapped := observable.Map(source, func(v int) int { return v * 2 })
			sub := &subscriber[int]{}
			mapped.Subscribe(context.Background(), sub)
		}
	})

	b.Run("Filter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			source := observable.Range(0, 1000)
			filtered := observable.Filter(source, func(v int) bool { return v%2 == 0 })
			sub := &subscriber[int]{}
			filtered.Subscribe(context.Background(), sub)
		}
	})
}
