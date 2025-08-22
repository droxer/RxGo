package streams

import (
	"context"
	"testing"
)

func TestSubscribeErrorHandling(t *testing.T) {
	t.Run("ReactivePublisher nil subscriber", func(t *testing.T) {
		publisher := NewPublisher(func(ctx context.Context, sub Subscriber[int]) {
			sub.OnNext(1)
			sub.OnComplete()
		})

		err := publisher.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("BufferedPublisher nil subscriber", func(t *testing.T) {
		publisher := NewBufferedPublisher(BackpressureConfig{
			Strategy:   Buffer,
			BufferSize: 5,
		}, func(ctx context.Context, sub Subscriber[int]) {
			sub.OnNext(1)
			sub.OnComplete()
		})

		err := publisher.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("RetryPublisher nil subscriber", func(t *testing.T) {
		source := NewPublisher(func(ctx context.Context, sub Subscriber[int]) {
			sub.OnNext(1)
			sub.OnComplete()
		})

		publisher := NewRetryPublisher(source, RetryConfig{
			MaxRetries:    3,
			InitialDelay:  100,
			BackoffPolicy: RetryFixed,
		})

		err := publisher.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("MapProcessor nil subscriber", func(t *testing.T) {
		processor := NewMapProcessor(func(x int) string {
			return "test"
		})

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("FilterProcessor nil subscriber", func(t *testing.T) {
		processor := NewFilterProcessor(func(x int) bool {
			return x > 0
		})

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("FlatMapProcessor nil subscriber", func(t *testing.T) {
		processor := NewFlatMapProcessor(func(x int) Publisher[string] {
			return NewPublisher(func(ctx context.Context, sub Subscriber[string]) {
				sub.OnNext("test")
				sub.OnComplete()
			})
		})

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("MergeProcessor nil subscriber", func(t *testing.T) {
		source := NewPublisher(func(ctx context.Context, sub Subscriber[int]) {
			sub.OnNext(1)
			sub.OnComplete()
		})

		processor := NewMergeProcessor(source)

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("ConcatProcessor nil subscriber", func(t *testing.T) {
		source := NewPublisher(func(ctx context.Context, sub Subscriber[int]) {
			sub.OnNext(1)
			sub.OnComplete()
		})

		processor := NewConcatProcessor(source)

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("TakeProcessor nil subscriber", func(t *testing.T) {
		processor := NewTakeProcessor[int](5)

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("SkipProcessor nil subscriber", func(t *testing.T) {
		processor := NewSkipProcessor[int](2)

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("DistinctProcessor nil subscriber", func(t *testing.T) {
		processor := NewDistinctProcessor[int]()

		err := processor.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("CompliantRangePublisher nil subscriber", func(t *testing.T) {
		publisher := NewCompliantRangePublisher(1, 5)

		err := publisher.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("CompliantFromSlicePublisher nil subscriber", func(t *testing.T) {
		publisher := NewCompliantFromSlicePublisher([]int{1, 2, 3, 4, 5})

		err := publisher.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})
}
