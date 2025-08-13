package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/droxer/RxGo/pkg/rx/scheduler"
)

func main() {
	fmt.Println("=== RxGo Scheduler Examples ===")

	// Quick demo of all schedulers
	work := func(id int) { fmt.Printf("  Task %d on goroutine %d\n", id, getGID()) }

	// 1. Trampoline - immediate execution
	fmt.Println("1. Trampoline (immediate execution):")
	scheduler.Trampoline.Schedule(func() { work(1) })

	// 2. NewThread - new goroutine for each task
	fmt.Println("\n2. NewThread (new goroutine for each task):")
	scheduler.NewThread.Schedule(func() { work(2) })

	// 3. SingleThread - sequential on dedicated goroutine
	fmt.Println("\n3. SingleThread (sequential):")
	st := scheduler.SingleThread
	st.Schedule(func() { work(3) })
	st.Schedule(func() { work(4) })

	// 4. Computation - fixed thread pool
	fmt.Println("\n4. Computation (fixed thread pool):")
	scheduler.Computation.Schedule(func() { work(5) })

	// 5. IO - cached thread pool
	fmt.Println("\n5. IO (cached thread pool):")
	scheduler.IO.Schedule(func() { work(6) })

	// 6. Performance comparison
	fmt.Println("\n6. Performance (100 tasks) - using Computation scheduler:")
	task := func() {
		for i := 0; i < 1000; i++ {
			// Simulate work
		}
	}

	benchmark := func(name string, s scheduler.Scheduler, n int) time.Duration {
		start := time.Now()
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			s.Schedule(func() { defer wg.Done(); task() })
		}
		wg.Wait()
		return time.Since(start)
	}

	fmt.Printf("  Computation: %v\n", benchmark("Computation", scheduler.Computation, 100))
	fmt.Printf("  IO: %v\n", benchmark("IO", scheduler.IO, 100))

	// 7. Decision guide
	fmt.Println("\n7. When to use:")
	fmt.Println("  Trampoline - UI updates, lightweight ops, immediate execution")
	fmt.Println("  NewThread - CPU/I/O intensive, parallel processing")
	fmt.Println("  SingleThread - ordered processing, transactions")
	fmt.Println("  Computation - CPU-bound work with optimal thread count")
	fmt.Println("  IO - I/O-bound work with dynamic thread creation")
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	var id uint64
	_, _ = fmt.Sscanf(string(b), "goroutine %d", &id)
	return id
}
