package scheduler

import (
	"sync"
	"testing"
)

func TestBasicScheduler(t *testing.T) {
	tests := []struct {
		name  string
		sched Scheduler
	}{
		{"Trampoline", Trampoline},
		{"NewThread", NewThread},
		{"SingleThread", SingleThread},
		{"Computation", Computation},
		{"IO", IO},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			result := 0
			tt.sched.Schedule(func() {
				result = 123
				wg.Done()
			})

			wg.Wait()

			if result != 123 {
				t.Errorf("%s: expected 123, got %d", tt.name, result)
			}
		})
	}
}
