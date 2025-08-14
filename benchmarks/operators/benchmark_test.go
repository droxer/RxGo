// Package benchmarks provides performance benchmarks for the RxGo operators.
package benchmarks

import (
	"context"
	"sync"
	"testing"

	"github.com/droxer/RxGo/pkg/rx"
	"github.com/droxer/RxGo/pkg/rx/operators"
	"github.com/droxer/RxGo/pkg/rx/scheduler"
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

func (t *TestSubscriber[T]) OnComplete() {
	t.mu.Lock()
	t.completed = true
	t.mu.Unlock()
}

func (t *TestSubscriber[T]) OnError(e error) {
	t.mu.Lock()
	t.errors = append(t.errors, e)
	t.mu.Unlock()
}

func BenchmarkOperatorMap(b *testing.B) {
	observable := rx.Range(0, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapped := operators.Map(observable, func(v int) int { return v * 2 })
		subscriber := &TestSubscriber[int]{}
		mapped.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkOperatorFilter(b *testing.B) {
	observable := rx.Range(0, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := operators.Filter(observable, func(v int) bool { return v%2 == 0 })
		subscriber := &TestSubscriber[int]{}
		filtered.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkOperatorMapFilterChain(b *testing.B) {
	observable := rx.Range(0, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Chain Map and Filter operators
		mapped := operators.Map(observable, func(v int) int { return v * 2 })
		filtered := operators.Filter(mapped, func(v int) bool { return v > 500 })
		subscriber := &TestSubscriber[int]{}
		filtered.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkOperatorObserveOn(b *testing.B) {
	observable := rx.Range(0, 1000)
	sched := scheduler.Trampoline

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		observed := operators.ObserveOn(observable, sched)
		subscriber := &TestSubscriber[int]{}
		observed.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkOperatorMapMemoryAllocations(b *testing.B) {
	observable := rx.Range(0, 1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapped := operators.Map(observable, func(v int) int { return v * 2 })
		subscriber := &TestSubscriber[int]{}
		mapped.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkOperatorFilterMemoryAllocations(b *testing.B) {
	observable := rx.Range(0, 1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := operators.Filter(observable, func(v int) bool { return v%2 == 0 })
		subscriber := &TestSubscriber[int]{}
		filtered.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkOperatorChainMemoryAllocations(b *testing.B) {
	observable := rx.Range(0, 1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Chain multiple operators
		mapped := operators.Map(observable, func(v int) int { return v * 2 })
		filtered := operators.Filter(mapped, func(v int) bool { return v > 500 })
		mapped2 := operators.Map(filtered, func(v int) int { return v + 1 })
		subscriber := &TestSubscriber[int]{}
		mapped2.Subscribe(context.Background(), subscriber)
	}
}

func BenchmarkOperatorMapDifferentDataTypes(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		observable := rx.Range(0, 100)
		for i := 0; i < b.N; i++ {
			mapped := operators.Map(observable, func(v int) int { return v * 2 })
			subscriber := &TestSubscriber[int]{}
			mapped.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		floatObservable := rx.Create(func(ctx context.Context, sub rx.Subscriber[float64]) {
			defer sub.OnComplete()
			for i := 0; i < 100; i++ {
				sub.OnNext(float64(i))
			}
		})

		for i := 0; i < b.N; i++ {
			mapped := operators.Map(floatObservable, func(v float64) float64 { return v * 2.0 })
			subscriber := &TestSubscriber[float64]{}
			mapped.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("String", func(b *testing.B) {
		stringObservable := rx.Just("a", "b", "c", "d", "e")
		for i := 0; i < b.N; i++ {
			mapped := operators.Map(stringObservable, func(v string) string { return v + "x" })
			subscriber := &TestSubscriber[string]{}
			mapped.Subscribe(context.Background(), subscriber)
		}
	})
}

func BenchmarkOperatorFilterDifferentDataTypes(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		observable := rx.Range(0, 100)
		for i := 0; i < b.N; i++ {
			filtered := operators.Filter(observable, func(v int) bool { return v%2 == 0 })
			subscriber := &TestSubscriber[int]{}
			filtered.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("Float64", func(b *testing.B) {
		floatObservable := rx.Create(func(ctx context.Context, sub rx.Subscriber[float64]) {
			defer sub.OnComplete()
			for i := 0; i < 100; i++ {
				sub.OnNext(float64(i))
			}
		})

		for i := 0; i < b.N; i++ {
			filtered := operators.Filter(floatObservable, func(v float64) bool { return v > 50.0 })
			subscriber := &TestSubscriber[float64]{}
			filtered.Subscribe(context.Background(), subscriber)
		}
	})

	b.Run("String", func(b *testing.B) {
		stringObservable := rx.Just("apple", "banana", "cherry", "date", "elderberry")
		for i := 0; i < b.N; i++ {
			filtered := operators.Filter(stringObservable, func(v string) bool { return len(v) > 5 })
			subscriber := &TestSubscriber[string]{}
			filtered.Subscribe(context.Background(), subscriber)
		}
	})
}
