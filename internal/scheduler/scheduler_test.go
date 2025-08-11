package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestEventLoopScheduler(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"Schedule", testEventLoopSchedulerSchedule},
		{"ScheduleAt", testEventLoopSchedulerScheduleAt},
		{"StartStop", testEventLoopSchedulerStartStop},
		{"Concurrent", testEventLoopSchedulerConcurrent},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}

func testEventLoopSchedulerSchedule(t *testing.T) {
	scheduler := newEventLoopScheduler(2)
	defer scheduler.Stop()

	var result int
	var wg sync.WaitGroup
	wg.Add(1)

	scheduler.Schedule(func() {
		result = 42
		wg.Done()
	})

	wg.Wait()

	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func testEventLoopSchedulerScheduleAt(t *testing.T) {
	scheduler := newEventLoopScheduler(2)
	defer scheduler.Stop()

	start := time.Now()
	var executed bool
	var wg sync.WaitGroup
	wg.Add(1)

	scheduler.ScheduleAt(func() {
		executed = true
		wg.Done()
	}, 50*time.Millisecond)

	wg.Wait()

	elapsed := time.Since(start)
	if !executed {
		t.Error("task was not executed")
	}
	if elapsed < 45*time.Millisecond {
		t.Errorf("task executed too early: %v", elapsed)
	}
}

func testEventLoopSchedulerStartStop(t *testing.T) {
	scheduler := newEventLoopScheduler(2)

	// Should start successfully
	scheduler.Start()

	// Schedule should work after start
	var executed bool
	var wg sync.WaitGroup
	wg.Add(1)

	scheduler.Schedule(func() {
		executed = true
		wg.Done()
	})

	wg.Wait()

	if !executed {
		t.Error("task was not executed")
	}

	// Stop should not panic
	scheduler.Stop()
}

func testEventLoopSchedulerConcurrent(t *testing.T) {
	scheduler := newEventLoopScheduler(4)
	defer scheduler.Stop()

	const numTasks = 100
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		scheduler.Schedule(func() {
			mu.Lock()
			counter++
			mu.Unlock()
			wg.Done()
		})
	}

	wg.Wait()

	if counter != numTasks {
		t.Errorf("expected %d, got %d", numTasks, counter)
	}
}

func TestThreadPoolScheduler(t *testing.T) {
	scheduler := newThreadPoolScheduler(100 * time.Millisecond)
	defer scheduler.Stop()

	var result int
	var wg sync.WaitGroup
	wg.Add(1)

	scheduler.Schedule(func() {
		result = 100
		wg.Done()
	})

	wg.Wait()

	if result != 100 {
		t.Errorf("expected 100, got %d", result)
	}
}

func TestThreadPoolSchedulerStop(t *testing.T) {
	scheduler := newThreadPoolScheduler(100 * time.Millisecond)

	// Start and schedule a task
	scheduler.Start()

	var executed bool
	var wg sync.WaitGroup
	wg.Add(1)

	scheduler.Schedule(func() {
		executed = true
		wg.Done()
	})

	wg.Wait()

	if !executed {
		t.Error("task was not executed")
	}

	// Stop should not panic
	scheduler.Stop()
}

func TestMaxParallelism(t *testing.T) {
	cores := maxParallelism()
	if cores <= 0 {
		t.Errorf("maxParallelism should be > 0, got %d", cores)
	}
	if cores > 1000 {
		t.Errorf("maxParallelism seems too high: %d", cores)
	}
}

func TestGlobalSchedulers(t *testing.T) {
	if Computation == nil {
		t.Error("Computation scheduler should not be nil")
	}
	if IO == nil {
		t.Error("IO scheduler should not be nil")
	}

	// Test that global schedulers can schedule work
	var result1, result2 int
	var wg sync.WaitGroup
	wg.Add(2)

	Computation.Schedule(func() {
		result1 = 1
		wg.Done()
	})

	IO.Schedule(func() {
		result2 = 2
		wg.Done()
	})

	wg.Wait()

	if result1 != 1 {
		t.Error("Computation scheduler failed")
	}
	if result2 != 2 {
		t.Error("IO scheduler failed")
	}
}

func TestThreadPoolSchedulerTimeout(t *testing.T) {
	scheduler := newThreadPoolScheduler(10 * time.Millisecond)
	defer scheduler.Stop()

	var executed bool
	var wg sync.WaitGroup
	wg.Add(1)

	scheduler.Schedule(func() {
		executed = true
		wg.Done()
	})

	wg.Wait()

	if !executed {
		t.Error("task was not executed")
	}

	// Allow timeout to trigger
	time.Sleep(20 * time.Millisecond)
}

func BenchmarkEventLoopScheduler(b *testing.B) {
	scheduler := newEventLoopScheduler(4)
	defer scheduler.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {})
	}
}

func BenchmarkThreadPoolScheduler(b *testing.B) {
	scheduler := newThreadPoolScheduler(100 * time.Millisecond)
	defer scheduler.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Schedule(func() {})
	}
}
