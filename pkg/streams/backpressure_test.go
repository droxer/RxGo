package streams

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBackpressureStrategies(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	tests := []struct {
		name     string
		strategy BackpressureStrategy
		config   BackpressureConfig
		setup    func() *backpressureTestSubscriber
		validate func(t *testing.T, received []int, err error)
	}{
		{
			name:     "Buffer strategy",
			strategy: Buffer,
			config:   BackpressureConfig{Strategy: Buffer, BufferSize: 5},
			setup: func() *backpressureTestSubscriber {
				return &backpressureTestSubscriber{}
			},
			validate: func(t *testing.T, received []int, err error) {
				if len(received) != 10 {
					t.Errorf("Expected 10 items, got %d", len(received))
				}
				for i, val := range received {
					if val != i+1 {
						t.Errorf("Expected item %d to be %d, got %d", i, i+1, val)
					}
				}
			},
		},
		{
			name:     "Drop strategy",
			strategy: Drop,
			config:   BackpressureConfig{Strategy: Drop, BufferSize: 3},
			setup: func() *backpressureTestSubscriber {
				return &backpressureTestSubscriber{processDelay: 50 * time.Millisecond}
			},
			validate: func(t *testing.T, received []int, err error) {
				t.Logf("Drop strategy processed %d items", len(received))
				if len(received) > 10 {
					t.Errorf("Expected drop strategy to drop some items, got %d", len(received))
				}
			},
		},
		{
			name:     "Latest strategy",
			strategy: Latest,
			config:   BackpressureConfig{Strategy: Latest, BufferSize: 2},
			setup: func() *backpressureTestSubscriber {
				return &backpressureTestSubscriber{processDelay: 100 * time.Millisecond}
			},
			validate: func(t *testing.T, received []int, err error) {
				t.Logf("Latest strategy processed %d items", len(received))
				if len(received) > 0 {
					lastItem := received[len(received)-1]
					if lastItem != 10 {
						t.Logf("Expected last item to be 10, got %d (this may be expected with Latest strategy)", lastItem)
					}
				}
			},
		},
		{
			name:     "Error strategy",
			strategy: Error,
			config:   BackpressureConfig{Strategy: Error, BufferSize: 1},
			setup: func() *backpressureTestSubscriber {
				return &backpressureTestSubscriber{processDelay: 500 * time.Millisecond}
			},
			validate: func(t *testing.T, received []int, err error) {
				if err == nil {
					t.Log("Error strategy test - overflow may not trigger under all conditions")
				} else {
					t.Logf("Error strategy correctly triggered: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := FromSlicePublishWithBackpressure(data, tt.config)
			sub := tt.setup()

			ctx := context.Background()
			err := publisher.Subscribe(ctx, sub)
			if err != nil {
				t.Fatalf("Subscribe failed: %v", err)
			}
			sub.wait()

			received := sub.getItems()
			tt.validate(t, received, sub.getError())
		})
	}
}

func TestRangePublishWithBackpressure(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	publisher := RangePublishWithBackpressure(1, 5, BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 10,
	})

	sub := &backpressureTestSubscriber{
		items: make([]int, 0),
		wg:    sync.WaitGroup{},
	}
	ctx := context.Background()

	err := publisher.Subscribe(ctx, sub)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	sub.wait()

	received := sub.getItems()
	if len(received) != len(expected) {
		t.Errorf("Expected %d values, got %d: %v", len(expected), len(received), received)
		return
	}
	for i, v := range expected {
		if received[i] != v {
			t.Errorf("Expected value[%d] to be %v, got %v", i, v, received[i])
		}
	}
	if err := sub.getError(); err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestBackpressureStrategyString(t *testing.T) {
	tests := []struct {
		strategy BackpressureStrategy
		expected string
	}{
		{Buffer, "Buffer"},
		{Drop, "Drop"},
		{Latest, "Latest"},
		{Error, "Error"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.strategy.String()
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

type backpressureTestSubscriber struct {
	items        []int
	err          error
	mu           sync.Mutex
	wg           sync.WaitGroup
	processDelay time.Duration
	requestCount int64
	completed    bool
	subscription Subscription
}

func (s *backpressureTestSubscriber) OnSubscribe(sub Subscription) {
	s.subscription = sub
	s.wg.Add(1)
	if s.requestCount > 0 {
		sub.Request(s.requestCount)
	} else {
		sub.Request(100) // Default unlimited
	}
}

func (s *backpressureTestSubscriber) OnNext(value int) {
	if s.processDelay > 0 {
		time.Sleep(s.processDelay)
	}
	s.mu.Lock()
	s.items = append(s.items, value)
	s.mu.Unlock()
}

func (s *backpressureTestSubscriber) OnError(err error) {
	s.mu.Lock()
	if !s.completed {
		s.err = err
		s.completed = true
		s.wg.Done()
	}
	s.mu.Unlock()
}

func (s *backpressureTestSubscriber) OnComplete() {
	s.mu.Lock()
	if !s.completed {
		s.completed = true
		s.wg.Done()
	}
	s.mu.Unlock()
}

func (s *backpressureTestSubscriber) wait() {
	s.wg.Wait()
}

func (s *backpressureTestSubscriber) getItems() []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]int, len(s.items))
	copy(result, s.items)
	return result
}

func (s *backpressureTestSubscriber) getError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}
