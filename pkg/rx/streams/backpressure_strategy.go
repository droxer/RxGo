package streams

// BackpressureStrategy defines how to handle buffer overflow scenarios

type BackpressureStrategy int

const (
	// Buffer keeps all items in a bounded buffer, blocking when full
	Buffer BackpressureStrategy = iota

	// Drop discards new items when buffer is full
	Drop

	// Latest keeps only the latest item, discarding older ones
	Latest

	// Error signals an error when buffer overflows
	Error
)

// BackoffPolicy defines retry backoff strategies
type BackoffPolicy int

const (
	// FixedBackoff uses constant delay between retries
	FixedBackoff BackoffPolicy = iota
	// LinearBackoff increases delay linearly
	LinearBackoff
	// ExponentialBackoff increases delay exponentially
	ExponentialBackoff
)

// String returns string representation of backpressure strategy
func (bs BackpressureStrategy) String() string {
	switch bs {
	case Buffer:
		return "Buffer"
	case Drop:
		return "Drop"
	case Latest:
		return "Latest"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}

// String returns string representation of backoff policy
func (b BackoffPolicy) String() string {
	switch b {
	case FixedBackoff:
		return "FixedBackoff"
	case LinearBackoff:
		return "LinearBackoff"
	case ExponentialBackoff:
		return "ExponentialBackoff"
	default:
		return "Unknown"
	}
}

// BackpressureConfig configures backpressure behavior for publishers
type BackpressureConfig struct {
	Strategy   BackpressureStrategy
	BufferSize int64
}
