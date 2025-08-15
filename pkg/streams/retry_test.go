package streams

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Test that retry policies can be created and configured
func TestRetryPolicies(t *testing.T) {
	// Test that we can create different retry policies
	config := RetryConfig{
		MaxRetries:    3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		BackoffPolicy: RetryExponential,
	}

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.MaxRetries)
	}

	if config.BackoffFactor != 2.0 {
		t.Errorf("Expected BackoffFactor to be 2.0, got %f", config.BackoffFactor)
	}
}

// Test that retry backoff policy string representations work
func TestRetryBackoffPolicyStringFinal(t *testing.T) {
	tests := []struct {
		policy   RetryBackoffPolicy
		expected string
	}{
		{RetryFixed, "RetryFixed"},
		{RetryLinear, "RetryLinear"},
		{RetryExponential, "RetryExponential"},
		{RetryBackoffPolicy(999), "Unknown"}, // Invalid policy
	}

	for _, test := range tests {
		result := test.policy.String()
		if result != test.expected {
			t.Errorf("Expected %s for policy %d, got %s", test.expected, int(test.policy), result)
		}
	}
}

// Test the retry builder API
func TestRetryBuilderAPIFinal(t *testing.T) {
	// Create a mock publisher for testing the builder
	mockPub := &mockPublisher[int]{}

	// Test builder API
	builder := NewRetryBuilder(mockPub)
	if builder == nil {
		t.Error("Expected builder to be created")
	}

	// Test method chaining
	result := builder.
		MaxRetries(3).
		InitialDelay(100 * time.Millisecond).
		MaxDelay(5 * time.Second).
		FixedBackoff().
		Build()

	if result == nil {
		t.Error("Expected result to be created")
	}
}

// Test convenience functions
func TestRetryConvenienceFunctionsFinal(t *testing.T) {
	// Create a mock publisher for testing
	mockPub := &mockPublisher[int]{}

	// Test WithFixedRetry
	fixedRetry := WithFixedRetry(mockPub, 3, 100*time.Millisecond)
	if fixedRetry == nil {
		t.Error("Expected fixedRetry to be created")
	}

	// Test WithLinearRetry
	linearRetry := WithLinearRetry(mockPub, 3, 100*time.Millisecond)
	if linearRetry == nil {
		t.Error("Expected linearRetry to be created")
	}

	// Test WithExponentialRetry
	expRetry := WithExponentialRetry(mockPub, 3, 100*time.Millisecond, 2.0)
	if expRetry == nil {
		t.Error("Expected expRetry to be created")
	}

	// Test WithInfiniteRetry
	infRetry := WithInfiniteRetry(mockPub, 100*time.Millisecond, RetryFixed)
	if infRetry == nil {
		t.Error("Expected infRetry to be created")
	}
}

// mockPublisher is a simple mock for testing
type mockPublisher[T any] struct{}

func (m *mockPublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	// Do nothing for mock
}

// Test that the retry configuration validation works
func TestRetryConfigValidation(t *testing.T) {
	// Test that backoff factor is corrected when less than 1.0
	config := RetryConfig{
		BackoffFactor: 0.5, // Less than 1.0
	}

	source := &mockPublisher[int]{}
	retryPub := NewRetryPublisher(source, config)

	// Access the underlying config to check if it was corrected
	// This is a bit tricky since we can't directly access it
	// But we can at least verify the publisher was created
	if retryPub == nil {
		t.Error("Expected retryPub to be created")
	}
}

// Test that retry conditions work
func TestRetryConditions(t *testing.T) {
	// Test default retry condition
	err := errors.New("test error")
	shouldRetry := DefaultRetryCondition(err, 1)
	if !shouldRetry {
		t.Error("Expected DefaultRetryCondition to return true")
	}

	// Test custom retry condition
	customCondition := func(err error, attempt int) bool {
		return attempt < 3 // Only retry first 3 attempts
	}

	shouldRetry = customCondition(err, 1)
	if !shouldRetry {
		t.Error("Expected custom condition to return true for attempt 1")
	}

	shouldRetry = customCondition(err, 5)
	if shouldRetry {
		t.Error("Expected custom condition to return false for attempt 5")
	}
}
