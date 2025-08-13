package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestTrampolineScheduler(t *testing.T) {
	result := 0
	var wg sync.WaitGroup
	wg.Add(1)

	Trampoline.Schedule(func() {
		result = 42
		wg.Done()
	})

	wg.Wait()

	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestNewThreadScheduler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	executed := false
	NewThread.Schedule(func() {
		executed = true
		wg.Done()
	})

	wg.Wait()

	if !executed {
		t.Error("Task was not executed")
	}
}

func TestSingleThreadScheduler(t *testing.T) {
	scheduler := SingleThread
	var wg sync.WaitGroup
	wg.Add(3)

	results := make([]int, 0, 3)
	mutex := sync.Mutex{}

	// Schedule multiple tasks that should run sequentially
	for i := 0; i < 3; i++ {
		val := i
		scheduler.Schedule(func() {
			mutex.Lock()
			results = append(results, val)
			mutex.Unlock()
			wg.Done()
		})
	}

	wg.Wait()

	// Should run in order due to single thread
	expected := []int{0, 1, 2}
	for i, v := range results {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestComputationScheduler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	result := 0
	Computation.Schedule(func() {
		result = 100
		wg.Done()
	})

	wg.Wait()

	if result != 100 {
		t.Errorf("Expected 100, got %d", result)
	}
}

func TestIOScheduler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	result := 0
	IO.Schedule(func() {
		result = 200
		wg.Done()
	})

	wg.Wait()

	if result != 200 {
		t.Errorf("Expected 200, got %d", result)
	}
}

func TestSchedulersAreThreadSafe(t *testing.T) {
	schedulers := []Scheduler{Trampoline, NewThread, SingleThread, Computation, IO}

	for _, sched := range schedulers {
		t.Run("concurrent scheduling", func(t *testing.T) {
			var wg sync.WaitGroup
			const numTasks = 10
			results := make([]int, numTasks)
			mutex := sync.Mutex{}

			wg.Add(numTasks)
			for i := 0; i < numTasks; i++ {
				val := i
				go func() {
					sched.Schedule(func() {
						mutex.Lock()
						results[val] = val * 2
						mutex.Unlock()
						wg.Done()
					})
				}()
			}

			wg.Wait()

			for i, v := range results {
				if v != i*2 {
					t.Errorf("Expected %d at index %d, got %d", i*2, i, v)
				}
			}
		})
	}
}

func TestSingleThreadSchedulerClose(t *testing.T) {
	// Since we're using the global SingleThread, we can't close it
	// This test verifies the SingleThreadScheduler type exists and works
	scheduler := SingleThread

	var wg sync.WaitGroup
	wg.Add(1)

	result := false
	scheduler.Schedule(func() {
		result = true
		wg.Done()
	})

	wg.Wait()

	if !result {
		t.Error("SingleThread scheduler did not execute task")
	}
}

func TestSchedulerPerformance(t *testing.T) {
	schedulers := map[string]Scheduler{
		"Trampoline":   Trampoline,
		"NewThread":    NewThread,
		"SingleThread": SingleThread,
		"Computation":  Computation,
		"IO":           IO,
	}

	for name, sched := range schedulers {
		t.Run(name+" performance", func(t *testing.T) {
			const numTasks = 10
			var wg sync.WaitGroup
			wg.Add(numTasks)

			start := time.Now()
			for i := 0; i < numTasks; i++ {
				sched.Schedule(func() {
					wg.Done()
				})
			}
			wg.Wait()
			elapsed := time.Since(start)

			// All should complete within reasonable time
			if elapsed > 1*time.Second {
				t.Errorf("%s took too long: %v", name, elapsed)
			}
		})
	}
}
