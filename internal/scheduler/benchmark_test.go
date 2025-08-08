package scheduler

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Benchmark for scheduler creation and basic scheduling
func BenchmarkSchedulerSchedule(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			// Simple work
		})
	}
}

// Benchmark for thread pool scheduler
func BenchmarkThreadPoolScheduler(b *testing.B) {
	scheduler := newThreadPoolScheduler(time.Second)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			// Simple work
		})
	}
}

// Benchmark for scheduled tasks with delay
func BenchmarkSchedulerScheduleAt(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.ScheduleAt(func() {
			// Simple work
		}, time.Millisecond)
	}
}

// Benchmark for concurrent scheduling
func BenchmarkSchedulerConcurrent(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			scheduler.Schedule(func() {
				// Simple work
			})
		}
	})
}

// Benchmark for thread pool concurrent scheduling
func BenchmarkThreadPoolConcurrent(b *testing.B) {
	scheduler := newThreadPoolScheduler(time.Second)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			scheduler.Schedule(func() {
				// Simple work
			})
		}
	})
}

// Benchmark for worker creation (simplified)
func BenchmarkWorkerCreation(b *testing.B) {
	jobChanQueue := make(chan chan job, 10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker := newThreadWorker(time.Second, jobChanQueue)
		worker.start()
		// Skip stop to avoid deadlock in benchmark
	}
}

// Benchmark for pool worker creation
func BenchmarkPoolWorkerCreation(b *testing.B) {
	fixedWorkerPool := make(chan chan job, 4)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker := newPoolWorker(fixedWorkerPool)
		worker.start()
		worker.stop()
	}
}

// Benchmark for actual work scheduling
func BenchmarkSchedulerWorkThroughput(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	var counter atomic.Int64
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			counter.Add(1)
		})
	}
	
	// Wait for all work to complete
	time.Sleep(100 * time.Millisecond)
}

// Benchmark for thread pool work throughput
func BenchmarkThreadPoolWorkThroughput(b *testing.B) {
	scheduler := newThreadPoolScheduler(time.Second)
	scheduler.Start()
	defer scheduler.Stop()
	
	var counter atomic.Int64
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			counter.Add(1)
		})
	}
	
	// Wait for all work to complete
	time.Sleep(100 * time.Millisecond)
}

// Benchmark for memory allocations in scheduling
func BenchmarkSchedulerMemoryAllocations(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			// Simple work
		})
	}
}

// Benchmark for complex work scheduling
func BenchmarkSchedulerComplexWork(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	var counter atomic.Int64
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			// Simulate some CPU work
			sum := 0
			for j := 0; j < 1000; j++ {
				sum += j
			}
			counter.Add(int64(sum))
		})
	}
	
	// Wait for all work to complete
	time.Sleep(200 * time.Millisecond)
}

// Benchmark for worker timeout handling
func BenchmarkWorkerTimeout(b *testing.B) {
	pool := newCachedThreadPool(time.Millisecond)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker := newThreadWorker(time.Millisecond, pool.jobChanQueue)
		worker.start()
		// Let worker timeout naturally
		time.Sleep(2 * time.Millisecond)
	}
}

// Benchmark for event loop scheduler throughput
func BenchmarkEventLoopThroughput(b *testing.B) {
	eventLoop := newEventLoopScheduler(8)
	eventLoop.Start()
	defer eventLoop.Stop()
	
	var counter atomic.Int64
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventLoop.Schedule(func() {
			counter.Add(1)
		})
	}
	
	// Wait for completion
	time.Sleep(100 * time.Millisecond)
}

// Benchmark for context switching overhead
func BenchmarkContextSwitching(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		scheduler.Schedule(func() {
			wg.Done()
		})
	}
	wg.Wait()
}

// Benchmark for job channel operations
func BenchmarkJobChannelOperations(b *testing.B) {
	jobQueue := make(chan job, 1000)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			select {
			case jobQueue <- job{run: func() {}}:
			default:
				// Drop if full
			}
		}
	})
}

// Benchmark for fixed worker pool
func BenchmarkFixedWorkerPool(b *testing.B) {
	fixedWorkerPool := make(chan chan job, 4)
	workers := make([]*poolWorker, 4)
	
	// Start workers
	for i := 0; i < 4; i++ {
		workers[i] = newPoolWorker(fixedWorkerPool)
		workers[i].start()
	}
	
	scheduler := &eventLoopScheduler{
		workers:         workers,
		fixedWorkerPool: fixedWorkerPool,
		jobQueue:        make(chan job, 1000),
		quit:            make(chan bool),
	}
	go scheduler.dispatch()
	
	var counter atomic.Int64
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			counter.Add(1)
		})
	}
	
	// Cleanup
	close(scheduler.quit)
}

// Benchmark for channel communication overhead
func BenchmarkChannelCommunication(b *testing.B) {
	ch := make(chan int, 100)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			select {
			case ch <- 1:
			case <-ch:
			}
		}
	})
}

// Benchmark for worker lifecycle
func BenchmarkWorkerLifecycle(b *testing.B) {
	jobChanQueue := make(chan chan job, 10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker := newThreadWorker(time.Second, jobChanQueue)
		worker.start()
		worker.stop()
	}
}

// Benchmark for panic recovery
func BenchmarkPanicRecovery(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			defer func() {
				if r := recover(); r != nil {
					// Simulate panic recovery
					_ = r
				}
			}()
			// Simulate work
		})
	}
}

// Benchmark for memory allocations in scheduling
func BenchmarkSchedulerAllocations(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	scheduler.Start()
	defer scheduler.Stop()
	
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {
			// Simple work
		})
	}
}

// Benchmark for timer operations
func BenchmarkTimerOperations(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timer := time.NewTimer(time.Millisecond)
		timer.Stop()
	}
}

// Benchmark for atomic operations
func BenchmarkAtomicOperations(b *testing.B) {
	var counter atomic.Int64
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}

// Benchmark for multi-threaded scheduling
func BenchmarkMultiThreadedScheduling(b *testing.B) {
	computations := newEventLoopScheduler(runtime.NumCPU())
	computations.Start()
	defer computations.Stop()
	
	ios := newThreadPoolScheduler(time.Second)
	ios.Start()
	defer ios.Stop()
	
	var counter atomic.Int64
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Alternate between schedulers
			if counter.Load()%2 == 0 {
				computations.Schedule(func() {
					counter.Add(1)
				})
			} else {
				ios.Schedule(func() {
					counter.Add(1)
				})
			}
		}
	})
	
	// Wait for completion
	time.Sleep(200 * time.Millisecond)
}