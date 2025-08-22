package streams

import (
	"context"
	"errors"
	"math"
	"time"
)

// RetryConfig configures retry and backoff behavior
//
//	config := RetryConfig{
//	    MaxRetries:     3,
//	    InitialDelay:   100 * time.Millisecond,
//	    MaxDelay:       5 * time.Second,
//	    BackoffFactor:  2.0,
//	    BackoffPolicy:  RetryExponential,
//	}
type RetryConfig struct {
	MaxRetries     int                // Maximum number of retry attempts (0 = infinite)
	InitialDelay   time.Duration      // Initial delay before first retry
	MaxDelay       time.Duration      // Maximum delay between retries
	BackoffFactor  float64            // Multiplier for backoff (must be >= 1.0)
	BackoffPolicy  RetryBackoffPolicy // Type of backoff strategy
	RetryCondition RetryCondition     // Function to determine if retry should happen
}

// RetryBackoffPolicy defines the backoff strategy
type RetryBackoffPolicy int

const (
	RetryFixed RetryBackoffPolicy = iota
	RetryLinear
	RetryExponential
)

func (r RetryBackoffPolicy) String() string {
	switch r {
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
//
//	condition := func(err error, attempt int) bool {
//	    if err == context.DeadlineExceeded {
//	        return false
//	    }
//	    return attempt < 3
//	}
type RetryCondition func(err error, attempt int) bool

var DefaultRetryCondition = func(err error, attempt int) bool {
	return true
}

type RetryPublisher[T any] struct {
	source Publisher[T]
	config RetryConfig
}

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

func (r *RetryPublisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) error {
	if sub == nil {
		return errors.New("subscriber cannot be nil")
	}

	retrySub := &retrySubscriber[T]{
		subscriber: sub,
		config:     r.config,
		source:     r.source,
		ctx:        ctx,
	}
	retrySub.start()

	return nil
}

type retrySubscriber[T any] struct {
	subscriber   Subscriber[T]
	config       RetryConfig
	source       Publisher[T]
	ctx          context.Context
	attempt      int
	subscription Subscription
}

func (r *retrySubscriber[T]) start() {
	r.subscriber.OnSubscribe(&retrySubscription[T]{
		cancelled: false,
		parent:    r,
	})
}

func (r *retrySubscriber[T]) subscribeToSource() {
	sourceSub := &retrySourceSubscriber[T]{
		parent: r,
	}
	err := r.source.Subscribe(r.ctx, sourceSub)
	if err != nil {
		sourceSub.OnError(err)
	}
}

func (r *retrySubscriber[T]) OnNext(value T) {
	r.subscriber.OnNext(value)
}

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

	select {
	case <-r.ctx.Done():
		r.subscriber.OnError(r.ctx.Err())
	case <-time.After(delay):
		if r.ctx.Err() == nil {
			r.subscribeToSource()
		}
	}
}

func (r *retrySubscriber[T]) OnComplete() {
	r.subscriber.OnComplete()
}

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

	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}

	return delay
}

type retrySourceSubscriber[T any] struct {
	parent *retrySubscriber[T]
}

func (r *retrySourceSubscriber[T]) OnSubscribe(sub Subscription) {
	r.parent.subscription = sub
	sub.Request(math.MaxInt64)
}

func (r *retrySourceSubscriber[T]) OnNext(value T) {
	r.parent.OnNext(value)
}

func (r *retrySourceSubscriber[T]) OnError(err error) {
	r.parent.OnError(err)
}

func (r *retrySourceSubscriber[T]) OnComplete() {
	r.parent.OnComplete()
}

type retrySubscription[T any] struct {
	cancelled bool
	parent    *retrySubscriber[T]
}

func (r *retrySubscription[T]) Request(n int64) {
	if !r.cancelled && r.parent.subscription != nil {
		r.parent.subscription.Request(n)
	}
}

func (r *retrySubscription[T]) Cancel() {
	r.cancelled = true
	if r.parent.subscription != nil {
		r.parent.subscription.Cancel()
	}
}

type RetryBuilder[T any] struct {
	source Publisher[T]
	config RetryConfig
}

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

func (b *RetryBuilder[T]) MaxRetries(max int) *RetryBuilder[T] {
	b.config.MaxRetries = max
	return b
}

func (b *RetryBuilder[T]) InitialDelay(delay time.Duration) *RetryBuilder[T] {
	b.config.InitialDelay = delay
	return b
}

func (b *RetryBuilder[T]) MaxDelay(max time.Duration) *RetryBuilder[T] {
	b.config.MaxDelay = max
	return b
}

func (b *RetryBuilder[T]) FixedBackoff() *RetryBuilder[T] {
	b.config.BackoffPolicy = RetryFixed
	return b
}

func (b *RetryBuilder[T]) LinearBackoff() *RetryBuilder[T] {
	b.config.BackoffPolicy = RetryLinear
	return b
}

func (b *RetryBuilder[T]) ExponentialBackoff(factor float64) *RetryBuilder[T] {
	b.config.BackoffPolicy = RetryExponential
	b.config.BackoffFactor = factor
	return b
}

func (b *RetryBuilder[T]) RetryCondition(condition RetryCondition) *RetryBuilder[T] {
	b.config.RetryCondition = condition
	return b
}

func (b *RetryBuilder[T]) Build() Publisher[T] {
	return NewRetryPublisher(b.source, b.config)
}

func WithFixedRetry[T any](source Publisher[T], maxRetries int, delay time.Duration) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    maxRetries,
		InitialDelay:  delay,
		BackoffPolicy: RetryFixed,
	})
}

func WithLinearRetry[T any](source Publisher[T], maxRetries int, initialDelay time.Duration) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    maxRetries,
		InitialDelay:  initialDelay,
		BackoffPolicy: RetryLinear,
	})
}

func WithExponentialRetry[T any](source Publisher[T], maxRetries int, initialDelay time.Duration, factor float64) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    maxRetries,
		InitialDelay:  initialDelay,
		BackoffPolicy: RetryExponential,
		BackoffFactor: factor,
	})
}

func WithInfiniteRetry[T any](source Publisher[T], initialDelay time.Duration, policy RetryBackoffPolicy) Publisher[T] {
	return NewRetryPublisher(source, RetryConfig{
		MaxRetries:    0,
		InitialDelay:  initialDelay,
		BackoffPolicy: policy,
	})
}
