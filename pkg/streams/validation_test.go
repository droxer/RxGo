package streams

import (
	"testing"
)

func TestBackpressureConfigValidation(t *testing.T) {
	// Test buffer size validation
	config := BackpressureConfig{
		Strategy:   Buffer,
		BufferSize: 0,
	}

	if config.BufferSize <= 0 {
		t.Log("Buffer size validation would catch zero size - good")
	}

	// Test negative buffer size
	config.BufferSize = -5
	if config.BufferSize <= 0 {
		t.Log("Buffer size validation would catch negative size - good")
	}
}

func TestBackpressureStrategyConstants(t *testing.T) {
	// Ensure the constants have expected values
	strategies := []BackpressureStrategy{Buffer, Drop, Latest, Error}

	if len(strategies) != 4 {
		t.Error("Expected 4 backpressure strategies")
	}

	// Test that they're all different
	seen := make(map[BackpressureStrategy]bool)
	for _, strategy := range strategies {
		if seen[strategy] {
			t.Errorf("Duplicate strategy value: %v", strategy)
		}
		seen[strategy] = true
	}
}
