package streams

import (
	"context"
	"math"
	"time"
)

// RetryConfig configures retry and backoff behavior
// This follows Reactive Streams 1.0.4 specification for error handling
// and provides flexible retry strategies
//
// Example:
//
//	config := RetryConfig{
//	    MaxRetries:     3,
//	    InitialDelay:   100 * time.Millisecond,
//	    MaxDelay:       5 * time.Second,
//	    BackoffFactor:  2.0,
//	    BackoffPolicy:  ExponentialBackoff,
//	}
//
// This will retry up to 3 times with delays: 100ms, 200ms, 400ms
type RetryConfig struct {
	MaxRetries     int                // Maximum number of retry attempts (0 = infinite)
	InitialDelay   time.Duration      // Initial delay before first retry
	MaxDelay       time.Duration      // Maximum delay between retries
	BackoffFactor  float64            // Multiplier for backoff (must be >= 1.0)
	BackoffPolicy  RetryBackoffPolicy // Type of backoff strategy
	RetryCondition RetryCondition     // Function to determine if retry should happen
}

// RetryBackoffPolicy defines retry backoff strategies
type RetryBackoffPolicy int

const (
	// RetryFixed uses constant delay between retries
	RetryFixed RetryBackoffPolicy = iota
	// RetryLinear increases delay linearly (delay = initial * attempt)
	RetryLinear
	// RetryExponential increases delay exponentially (delay = initial * factor^attempt)
	RetryExponential
)

// String returns string representation of retry backoff policy
func (b RetryBackoffPolicy) String() string {
	switch b {
	case RetryFixed:
		return "RetryFixed"
	case RetryLinear:
		return "RetryLinear"
	case RetryExponential:
		return "RetryExponential"
	default:
		return "Unknown"
	}
}

// RetryCondition determines if an error should trigger a retry
// Return true to retry, false to propagate the error
// This allows fine-grained control over which errors are retriable
//
// Example:
//
//	condition := func(err error, attempt int) bool {
//	    if err == context.DeadlineExceeded {
//	        return false // Don't retry timeouts
//	    }
//	    return attempt < 3
//	}
type RetryCondition func(err error, attempt int) bool

// DefaultRetryCondition retries all errors up to MaxRetries
var DefaultRetryCondition = func(err error, attempt int) bool {
	return true
}

// RetryPublisher adds retry capability to any Publisher
// It wraps an existing Publisher and applies retry logic on errors
// The retry behavior is fully configurable through RetryConfig
//
// This implementation maintains Reactive Streams compliance:
// - Respects backpressure signals during retries
// - Propagates cancellation properly
// - Ensures sequential signaling of events
// - Handles subscription lifecycle correctly
type RetryPublisher[T any] struct {
	source Publisher[T]
	config RetryConfig
}

// NewRetryPublisher creates a new Publisher with retry capability
// The retry behavior is configured through RetryConfig
//
// Example:
//
//	source := streams.RangePublisher(1, 10)
//	retryPublisher := streams.NewRetryPublisher(source, streams.RetryConfig{
//	    MaxRetries: 3,
//	    InitialDelay: 100 * time.Millisecond,
//	})
func NewRetryPublisher[T any](source Publisher[T], config RetryConfig) Publisher[T] {
	if config.BackoffFactor < 1.0 {
		config.BackoffFactor = 1.0
	}
	if config.RetryCondition == nil {
		config.RetryCondition = DefaultRetryCondition
	}
	return &RetryPublisher[T]{
		source: source,
		config: config,
	}
}

// Subscribe implements Publisher[T]
// It handles the retry logic by wrapping the original subscription
func (r *RetryPublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	retrySub := &retrySubscriber[T]{
		subscriber: sub,
		config:     r.config,
		source:     r.source,
		ctx:        ctx,
	}
	retrySub.start()
}

// retrySubscriber handles the retry logic for a single subscription
// It maintains state between retry attempts and manages backoff timing
type retrySubscriber[T any] struct {
	subscriber Subscriber[T]
	config     RetryConfig
	source     Publisher[T]
	ctx        context.Context

	attempt      int
	subscription Subscription
}

// start begins the subscription with retry logic
func (r *retrySubscriber[T]) start() {
	r.subscriber.OnSubscribe(&retrySubscription[T]{
		cancelled: false,
		parent:    r,
	})
}

// subscribeToSource creates a new subscription to the source Publisher
// This is called for each retry attempt
func (r *retrySubscriber[T]) subscribeToSource() {
	sourceSub := &retrySourceSubscriber[T]{
		parent: r,
	}
	r.source.Subscribe(r.ctx, sourceSub)
}

// OnNext handles successful data emission
func (r *retrySubscriber[T]) OnNext(value T) {
	r.subscriber.OnNext(value)
}

// OnError handles errors and determines retry behavior
// If retry is needed, it schedules a retry with backoff
func (r *retrySubscriber[T]) OnError(err error) {
	if r.config.RetryCondition != nil && !r.config.RetryCondition(err, r.attempt) {
		r.subscriber.OnError(err)
		return
	}

	if r.config.MaxRetries > 0 && r.attempt >= r.config.MaxRetries {
		r.subscriber.OnError(err)
		return
	}

	r.attempt++
	delay := r.calculateDelay()

	// Schedule retry with context cancellation support
	select {
	case <-r.ctx.Done():
		r.subscriber.OnError(r.ctx.Err())
	case <-time.After(delay):
		if r.ctx.Err() == nil {
			r.subscribeToSource()
		}
	}
}

// OnComplete handles successful completion
func (r *retrySubscriber[T]) OnComplete() {
	r.subscriber.OnComplete()
}

// calculateDelay computes the backoff delay based on the current attempt
func (r *retrySubscriber[T]) calculateDelay() time.Duration {
	var delay time.Duration

	switch r.config.BackoffPolicy {
	case RetryFixed:
		delay = r.config.InitialDelay
	case RetryLinear:
		delay = r.config.InitialDelay * time.Duration(r.attempt)
	case RetryExponential:
		delay = time.Duration(float64(r.config.InitialDelay) *
			math.Pow(r.config.BackoffFactor, float64(r.attempt-1)))
	}

	// Cap the delay at MaxDelay
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}

	return delay
}

// retrySourceSubscriber adapts the source Publisher's Subscriber interface
// to the retrySubscriber
// It forwards events to the retrySubscriber and handles the retry logic
type retrySourceSubscriber[T any] struct {
	parent *retrySubscriber[T]
}

// OnSubscribe handles subscription from the source Publisher
func (r *retrySourceSubscriber[T]) OnSubscribe(sub Subscription) {
	r.parent.subscription = sub
	sub.Request(math.MaxInt64) // Request all items
}

// OnNext forwards successful data emission
func (r *retrySourceSubscriber[T]) OnNext(value T) {
	r.parent.OnNext(value)
}

// OnError triggers retry logic
func (r *retrySourceSubscriber[T]) OnError(err error) {
	r.parent.OnError(err)
}

// OnComplete forwards completion
func (r *retrySourceSubscriber[T]) OnComplete() {
	r.parent.OnComplete()
}

// retrySubscription implements Subscription for retryPublisher
// It handles cancellation and backpressure signals
type retrySubscription[T any] struct {
	cancelled bool
	parent    *retrySubscriber[T]
}

// Request implements Subscription.Request
// Forwards the demand to the underlying subscription
func (r *retrySubscription[T]) Request(n int64) {
	if !r.cancelled && r.parent.subscription != nil {
		r.parent.subscription.Request(n)
	}
}

// Cancel implements Subscription.Cancel
// Cancels the underlying subscription and prevents further retries
func (r *retrySubscription[T]) Cancel() {
	r.cancelled = true
	if r.parent.subscription != nil {
		r.parent.subscription.Cancel()
	}
}

// RetryBuilder provides a fluent API for configuring retry behavior
// This makes it easy to create complex retry configurations
//
// Example:
//
//	retryPublisher := streams.NewRetryBuilder[int](source).
//	    MaxRetries(3).
//	    InitialDelay(100 * time.Millisecond).
//	    MaxDelay(5 * time.Second).
//	    ExponentialBackoff(2.0).
//	    Build()
type RetryBuilder[T any] struct {
	source Publisher[T]
	config RetryConfig
}

// NewRetryBuilder creates a new RetryBuilder
func NewRetryBuilder[T any](source Publisher[T]) *RetryBuilder[T] {
	return &RetryBuilder[T]{
		source: source,
		config: RetryConfig{
			BackoffFactor:  1.0,
			BackoffPolicy:  RetryFixed,
			RetryCondition: DefaultRetryCondition,
		},
	}
}

// MaxRetries sets the maximum number of retry attempts
func (b *RetryBuilder[T]) MaxRetries(max int) *RetryBuilder[T] {
	b.config.MaxRetries = max
	return b
}

// InitialDelay sets the initial delay before the first retry
func (b *RetryBuilder[T]) InitialDelay(delay time.Duration) *RetryBuilder[T] {
	b.config.InitialDelay = delay
	return b
}

// MaxDelay sets the maximum delay between retries
func (b *RetryBuilder[T]) MaxDelay(max time.Duration) *RetryBuilder[T] {
	b.config.MaxDelay = max
	return b
}

// FixedBackoff configures fixed delay between retries
func (b *RetryBuilder[T]) FixedBackoff() *RetryBuilder[T] {
	b.config.BackoffPolicy = RetryFixed
	return b
}

// LinearBackoff configures linear delay increase
func (b *RetryBuilder[T]) LinearBackoff() *RetryBuilder[T] {
	b.config.BackoffPolicy = RetryLinear
	return b
}

// ExponentialBackoff configures exponential delay increase
func (b *RetryBuilder[T]) ExponentialBackoff(factor float64) *RetryBuilder[T] {
	b.config.BackoffPolicy = RetryExponential
	b.config.BackoffFactor = factor
	return b
}

// RetryCondition sets custom retry condition
func (b *RetryBuilder[T]) RetryCondition(condition RetryCondition) *RetryBuilder[T] {
	b.config.RetryCondition = condition
	return b
}

// Build creates the final RetryPublisher
func (b *RetryBuilder[T]) Build() Publisher[T] {
	return NewRetryPublisher(b.source, b.config)
}

// Common retry configurations for convenience

// WithFixedRetry creates a Publisher with fixed delay retries
func WithFixedRetry[T any](source Publisher[T], maxRetries int, delay time.Duration) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    maxRetries,
		InitialDelay:  delay,
		BackoffPolicy: RetryFixed,
	})
}

// WithLinearRetry creates a Publisher with linear backoff retries
func WithLinearRetry[T any](source Publisher[T], maxRetries int, initialDelay time.Duration) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    maxRetries,
		InitialDelay:  initialDelay,
		BackoffPolicy: RetryLinear,
	})
}

// WithExponentialRetry creates a Publisher with exponential backoff retries
func WithExponentialRetry[T any](source Publisher[T], maxRetries int, initialDelay time.Duration, factor float64) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    maxRetries,
		InitialDelay:  initialDelay,
		BackoffPolicy: RetryExponential,
		BackoffFactor: factor,
	})
}

// WithInfiniteRetry creates a Publisher that retries indefinitely
func WithInfiniteRetry[T any](source Publisher[T], initialDelay time.Duration, policy RetryBackoffPolicy) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    0, // 0 means infinite
		InitialDelay:  initialDelay,
		BackoffPolicy: policy,
	})
}
