package streams

// RetryBackoffPolicy defines retry backoff strategies
type RetryBackoffPolicy int

const (
	// RetryFixed uses constant delay between retries
	RetryFixed RetryBackoffPolicy = iota
	// RetryLinear increases delay linearly (delay = initial * attempt)
	RetryLinear
	// RetryExponential increases delay exponentially (delay = initial * factor^(attempt-1))
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
