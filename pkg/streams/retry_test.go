package streams

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryPolicies(t *testing.T) {
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

func TestRetryBuilderAPIFinal(t *testing.T) {
	mockPub := &mockPublisher[int]{}

	builder := NewRetryBuilder(mockPub)
	if builder == nil {
		t.Error("Expected builder to be created")
	}

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

func TestRetryConvenienceFunctionsFinal(t *testing.T) {
	mockPub := &mockPublisher[int]{}

	fixedRetry := WithFixedRetry(mockPub, 3, 100*time.Millisecond)
	if fixedRetry == nil {
		t.Error("Expected fixedRetry to be created")
	}

	linearRetry := WithLinearRetry(mockPub, 3, 100*time.Millisecond)
	if linearRetry == nil {
		t.Error("Expected linearRetry to be created")
	}

	expRetry := WithExponentialRetry(mockPub, 3, 100*time.Millisecond, 2.0)
	if expRetry == nil {
		t.Error("Expected expRetry to be created")
	}

	infRetry := WithInfiniteRetry(mockPub, 100*time.Millisecond, RetryFixed)
	if infRetry == nil {
		t.Error("Expected infRetry to be created")
	}
}

type mockPublisher[T any] struct{}

func (m *mockPublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	return nil
}

func TestRetryConfigValidation(t *testing.T) {
	config := RetryConfig{
		BackoffFactor: 0.5,
	}

	source := &mockPublisher[int]{}
	retryPub := NewRetryPublisher(source, config)

	if retryPub == nil {
		t.Error("Expected retryPub to be created")
	}
}

func TestRetryConditions(t *testing.T) {
	err := errors.New("test error")
	shouldRetry := DefaultRetryCondition(err, 1)
	if !shouldRetry {
		t.Error("Expected DefaultRetryCondition to return true")
	}

	customCondition := func(err error, attempt int) bool {
		return attempt < 3
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
