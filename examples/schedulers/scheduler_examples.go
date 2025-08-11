package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/droxer/RxGo/pkg/observable"
)

func main() {
	fmt.Println("=== RxGo Scheduler Examples ===")

	// Quick demo of all schedulers
	work := func(id int) { fmt.Printf("  Task %d on goroutine %d\n", id, getGID()) }

	// 1. Immediate - same goroutine
	fmt.Println("1. Immediate (same goroutine):")
	observable.NewImmediateScheduler().Schedule(func() { work(1) })

	// 2. NewThread - new goroutine
	fmt.Println("\n2. NewThread (new goroutine):")
	observable.NewNewThreadScheduler().Schedule(func() { work(2) })

	// 3. SingleThread - sequential on dedicated goroutine
	fmt.Println("\n3. SingleThread (sequential):")
	st := observable.NewSingleThreadScheduler()
	defer st.Close()
	st.Schedule(func() { work(3) })
	st.Schedule(func() { work(4) })

	// 4. Trampoline - queued batch
	fmt.Println("\n4. Trampoline (queued):")
	t := observable.NewTrampolineScheduler()
	t.Schedule(func() { work(5) })
	t.Schedule(func() { work(6) })
	t.Execute()

	// 5. Performance comparison
	fmt.Println("\n5. Performance (100 tasks):")
	task := func() {
		for i := 0; i < 1000; i++ {
		}
	}

	benchmark := func(name string, s observable.Scheduler, n int) time.Duration {
		start := time.Now()
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			s.Schedule(func() { defer wg.Done(); task() })
		}
		if name == "SingleThread" {
			defer s.(*observable.SingleThreadScheduler).Close()
		}
		wg.Wait()
		return time.Since(start)
	}

	fmt.Printf("  Immediate: %v\n", benchmark("Immediate", observable.NewImmediateScheduler(), 100))
	fmt.Printf("  NewThread: %v\n", benchmark("NewThread", observable.NewNewThreadScheduler(), 100))

	// 6. Decision guide
	fmt.Println("\n6. When to use:")
	fmt.Println("  Immediate - UI updates, lightweight ops")
	fmt.Println("  NewThread - CPU/I/O intensive, parallel")
	fmt.Println("  SingleThread - ordered processing, transactions")
	fmt.Println("  Trampoline - batch processing, queueing")
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	var id uint64
	fmt.Sscanf(string(b), "goroutine %d", &id)
	return id
}
