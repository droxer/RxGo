// Package scheduler provides performance benchmarks for scheduler implementations.
package scheduler

import (
	"sync"
	"testing"
	"time"

	"github.com/droxer/RxGo/pkg/scheduler"
)

func BenchmarkSchedulerCreation(b *testing.B) {
	b.Run("Trampoline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = scheduler.Trampoline
		}
	})

	b.Run("NewThread", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = scheduler.NewThread
		}
	})

	b.Run("SingleThread", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = scheduler.SingleThread
		}
	})

	b.Run("Computation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = scheduler.Computation
		}
	})

	b.Run("IO", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = scheduler.IO
		}
	})
}

func BenchmarkSchedulerExecution(b *testing.B) {
	schedulers := map[string]scheduler.Scheduler{
		"Trampoline":   scheduler.Trampoline,
		"NewThread":    scheduler.NewThread,
		"SingleThread": scheduler.SingleThread,
		"Computation":  scheduler.Computation,
		"IO":           scheduler.IO,
	}

	for name, sched := range schedulers {
		b.Run(name, func(b *testing.B) {
			var wg sync.WaitGroup
			wg.Add(b.N)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sched.Schedule(func() {
					wg.Done()
				})
			}
			wg.Wait()
		})
	}
}

func BenchmarkConcurrentScheduling(b *testing.B) {
	schedulers := map[string]scheduler.Scheduler{
		"Trampoline":   scheduler.Trampoline,
		"NewThread":    scheduler.NewThread,
		"SingleThread": scheduler.SingleThread,
		"Computation":  scheduler.Computation,
		"IO":           scheduler.IO,
	}

	for name, sched := range schedulers {
		b.Run(name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				var wg sync.WaitGroup
				for pb.Next() {
					wg.Add(1)
					sched.Schedule(func() {
						wg.Done()
					})
				}
				wg.Wait()
			})
		})
	}
}

func BenchmarkSchedulerThroughput(b *testing.B) {
	b.Run("SingleThread Sequential", func(b *testing.B) {
		sched := scheduler.SingleThread
		var wg sync.WaitGroup
		wg.Add(b.N)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sched.Schedule(func() {
				// Simple work
				_ = 42 * 42
				wg.Done()
			})
		}
		wg.Wait()
	})

	b.Run("NewThread Parallel", func(b *testing.B) {
		sched := scheduler.NewThread
		var wg sync.WaitGroup
		wg.Add(b.N)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sched.Schedule(func() {
				// Simple work
				_ = 42 * 42
				wg.Done()
			})
		}
		wg.Wait()
	})

	b.Run("Computation Parallel", func(b *testing.B) {
		sched := scheduler.Computation
		var wg sync.WaitGroup
		wg.Add(b.N)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sched.Schedule(func() {
				// Simple work
				_ = 42 * 42
				wg.Done()
			})
		}
		wg.Wait()
	})
}

func BenchmarkMemoryAllocations(b *testing.B) {
	b.ReportAllocs()

	schedulers := map[string]scheduler.Scheduler{
		"Trampoline":   scheduler.Trampoline,
		"NewThread":    scheduler.NewThread,
		"SingleThread": scheduler.SingleThread,
		"Computation":  scheduler.Computation,
		"IO":           scheduler.IO,
	}

	for name, sched := range schedulers {
		b.Run(name, func(b *testing.B) {
			var wg sync.WaitGroup
			wg.Add(b.N)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sched.Schedule(func() {
					wg.Done()
				})
			}
			wg.Wait()
		})
	}
}

func BenchmarkSchedulerLatency(b *testing.B) {
	schedulers := map[string]scheduler.Scheduler{
		"Trampoline":   scheduler.Trampoline,
		"NewThread":    scheduler.NewThread,
		"SingleThread": scheduler.SingleThread,
		"Computation":  scheduler.Computation,
		"IO":           scheduler.IO,
	}

	for name, sched := range schedulers {
		b.Run(name, func(b *testing.B) {
			start := time.Now()
			var wg sync.WaitGroup
			wg.Add(b.N)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sched.Schedule(func() {
					wg.Done()
				})
			}
			wg.Wait()

			elapsed := time.Since(start)
			b.ReportMetric(float64(elapsed.Nanoseconds())/float64(b.N), "ns/op")
		})
	}
}

func BenchmarkBatchScheduling(b *testing.B) {
	b.Run("BatchSize10", func(b *testing.B) {
		sched := scheduler.Trampoline
		var wg sync.WaitGroup
		wg.Add(b.N)

		b.ResetTimer()
		for i := 0; i < b.N; i += 10 {
			remaining := 10
			if b.N-i < 10 {
				remaining = b.N - i
			}
			for j := 0; j < remaining; j++ {
				sched.Schedule(func() {
					wg.Done()
				})
			}
		}
		wg.Wait()
	})

	b.Run("BatchSize100", func(b *testing.B) {
		sched := scheduler.Trampoline
		var wg sync.WaitGroup
		wg.Add(b.N)

		b.ResetTimer()
		for i := 0; i < b.N; i += 100 {
			remaining := 100
			if b.N-i < 100 {
				remaining = b.N - i
			}
			for j := 0; j < remaining; j++ {
				sched.Schedule(func() {
					wg.Done()
				})
			}
		}
		wg.Wait()
	})
}
