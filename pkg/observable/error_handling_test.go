package observable

import (
	"context"
	"testing"
)

func TestObservableSubscribeErrorHandling(t *testing.T) {
	t.Run("Observable nil subscriber", func(t *testing.T) {
		obs := Just(1, 2, 3)

		err := obs.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("Map operator error propagation", func(t *testing.T) {
		// Test that Map operator properly propagates Subscribe errors
		obs := Just(1, 2, 3)
		mapped := Map(obs, func(x int) string { return "test" })

		err := mapped.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("Filter operator error propagation", func(t *testing.T) {
		obs := Just(1, 2, 3)
		filtered := Filter(obs, func(x int) bool { return x > 1 })

		err := filtered.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("Merge operator error propagation", func(t *testing.T) {
		obs1 := Just(1, 2, 3)
		obs2 := Just(4, 5, 6)
		merged := Merge(obs1, obs2)

		err := merged.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("Concat operator error propagation", func(t *testing.T) {
		obs1 := Just(1, 2, 3)
		obs2 := Just(4, 5, 6)
		concatenated := Concat(obs1, obs2)

		err := concatenated.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("Take operator error propagation", func(t *testing.T) {
		obs := Just(1, 2, 3, 4, 5)
		taken := Take(obs, 3)

		err := taken.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("Skip operator error propagation", func(t *testing.T) {
		obs := Just(1, 2, 3, 4, 5)
		skipped := Skip(obs, 2)

		err := skipped.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})

	t.Run("Distinct operator error propagation", func(t *testing.T) {
		obs := Just(1, 2, 2, 3, 3, 3)
		distinct := Distinct(obs)

		err := distinct.Subscribe(context.Background(), nil)
		if err == nil {
			t.Error("Expected error for nil subscriber")
		}
		if err.Error() != "subscriber cannot be nil" {
			t.Errorf("Expected 'subscriber cannot be nil', got: %v", err)
		}
	})
}
